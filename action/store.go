package action

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/gocop/gocop"
	"github.com/digitalocean/gocop/gocop/storer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type storeCmdFlags struct {
	repo       string
	branch     string
	sha        string
	start      string
	runCommand string
	src        string
	buildID    int64
	bench      bool
	short      bool
	race       bool
	tags       []string
	retests    []string
	storerName string
}

var storeFlags struct {
	storeCmdFlags
	storer Storer
}

var testResults []storer.TestResult

// storeCmd is the `store` subcommand
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "stores test results to database",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		storer := registeredStorers[storeFlags.storerName]
		if storer == nil {
			return fmt.Errorf("unrecognized storer %s", storeFlags.storerName)
		}

		for _, f := range storer.Required() {
			err := cmd.MarkFlagRequired(fmt.Sprintf("%s.%s", storer.Name(), f))
			if err != nil {
				return err
			}
		}

		storeFlags.storer = storer
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		s, err := storeFlags.storer.Storer()
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			err = s.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}()

		run := storer.TestRun{
			BuildID:   storeFlags.buildID,
			Repo:      storeFlags.repo,
			Branch:    storeFlags.branch,
			Sha:       storeFlags.sha,
			Command:   storeFlags.runCommand,
			Benchmark: storeFlags.bench,
			Short:     storeFlags.short,
			Race:      storeFlags.race,
			Tags:      storeFlags.tags,
		}
		if len(storeFlags.start) != 0 {
			run.Created, err = time.Parse(time.RFC3339, storeFlags.start)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			run.Created = time.Now().UTC()
		}

		if len(storeFlags.src) > 0 {
			tests, err := gocop.ParseFile(storeFlags.src)
			if err != nil {
				log.Fatal(err)
			}

			for _, test := range tests {
				test.Created = run.Created
				testResults = append(testResults, test)
			}
		}

		if len(storeFlags.retests) > 0 {
			pkgs, err := gocop.FlakyFilePackages(storeFlags.retests...)
			if err != nil {
				log.Fatal(err)
			}
			for _, entry := range pkgs {
				testResults = append(testResults, storer.TestResult{
					Created: run.Created,
					Package: entry,
					Result:  "flaky",
				})
			}
		}

		err = s.InsertRun(ctx, run)
		if err != nil {
			log.Fatal(err)
		}

		err = s.InsertTests(ctx, testResults)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	storeCmd.Flags().StringVarP(&storeFlags.storerName, "storer", "", "psql", fmt.Sprintf("storer name. supported options: %v", registeredStorerNames))
	storeCmd.Flags().StringVarP(&storeFlags.repo, "repo", "g", "", "repository name")
	storeCmd.Flags().StringVarP(&storeFlags.branch, "branch", "b", "master", "branch name")
	storeCmd.Flags().Int64VarP(&storeFlags.buildID, "build-id", "i", 0, "build id")
	err := storeCmd.MarkFlagRequired("build-id")
	if err != nil {
		log.Fatal(err)
	}

	storeCmd.Flags().StringVarP(&storeFlags.runCommand, "cmd", "c", "", "test execution command")
	storeCmd.Flags().StringVarP(&storeFlags.sha, "sha", "z", "", "git sha of test run")
	storeCmd.Flags().StringVarP(&storeFlags.start, "time", "m", "", "time of test run")
	storeCmd.Flags().StringVarP(&storeFlags.src, "src", "s", "", "source test output file")
	storeCmd.MarkFlagFilename("src")
	storeCmd.Flags().BoolVar(&storeFlags.bench, "bench", false, "indicate if test ran benchmarks")
	storeCmd.Flags().BoolVar(&storeFlags.short, "short", false, "indicate if test is run with -short flag")
	storeCmd.Flags().BoolVar(&storeFlags.race, "race", false, "indicate if test is run with -race flag")
	storeCmd.Flags().StringSliceVar(&storeFlags.tags, "tags", []string{}, "comma-separated tags enabled for the run")
	storeCmd.Flags().StringSliceVarP(&storeFlags.retests, "rerun", "r", []string{}, "comma-separated source output for retests")

	RootCmd.AddCommand(storeCmd)
	RegisterStorer(&PSQLStorer{})
}

var (
	registeredStorerNames []string
	registeredStorers     = map[string]Storer{}
)

// RegisterStorer registers a storer with the `store` command.
func RegisterStorer(storer Storer) {
	fs := storer.FlagSet()
	// add the storer name as a prefix to the flag names
	fs.VisitAll(func(f *pflag.Flag) {
		f.Name = fmt.Sprintf("%s.%s", storer.Name(), f.Name)
		// clear shorthands to avoid conflicts with other flags
		f.Shorthand = ""
	})

	registeredStorers[storer.Name()] = storer
	registeredStorerNames = append(registeredStorerNames, storer.Name())

	// add the flags to the store command
	storeCmd.Flags().AddFlagSet(&fs)
	storeCmd.Flags().Lookup("storer").Usage = fmt.Sprintf("storer name. supported options: %v", registeredStorerNames)
}

// Storer allows the `store` CLI to support different storers.
type Storer interface {
	// Name returns the storer's name.
	Name() string
	// FlagSet returns a flagset to be added to the command so that users can pass configuration options as flags.
	FlagSet() (flagSet pflag.FlagSet)
	// Required returns a list of required flags.
	Required() []string
	// Storer creates the storer instance.
	Storer() (storer.Storer, error)
}

type PSQLStorer struct {
	host     string
	port     string
	dbName   string
	user     string
	password string
	sslMode  string
}

// Name returns the storer's name.
func (s *PSQLStorer) Name() string {
	return "psql"
}

// FlagSet returns a flagset to be added to the command so that users can pass configuration options as flags.
func (s *PSQLStorer) FlagSet() pflag.FlagSet {
	var fs pflag.FlagSet
	fs.StringVarP(&s.host, "host", "a", "localhost", "database host")
	fs.StringVarP(&s.port, "port", "t", "5432", "database port")
	fs.StringVarP(&s.dbName, "database", "x", "postgres", "database name")
	fs.StringVarP(&s.sslMode, "ssl", "y", "require", "database ssl mode")
	fs.StringVarP(&s.password, "pass", "p", "", "database password")
	fs.StringVarP(&s.user, "user", "u", "postgres", "database username")

	return fs
}

// Storer creates the storer instance.
func (s *PSQLStorer) Storer() (storer.Storer, error) {
	return storer.NewPSQL(
		s.host, s.port,
		s.user, s.password,
		s.dbName, s.sslMode,
	)
}

// Required returns a list of required flags.
func (s *PSQLStorer) Required() []string {
	return []string{"pass"}
}
