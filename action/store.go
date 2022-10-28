package action

import (
	"context"
	"log"
	"time"

	"github.com/digitalocean/gocop/gocop"
	"github.com/digitalocean/gocop/gocop/storer"
	"github.com/spf13/cobra"
)

type storeCmdFlags struct {
	host       string
	port       string
	dbName     string
	user       string
	password   string
	sslMode    string
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
}

var storeFlags storeCmdFlags

var testResults []storer.TestResult

// storeCmd is the `store` subcommand
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "stores test results to database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var err error
		s, err := storer.NewSQL(
			storeFlags.host, storeFlags.port,
			storeFlags.user, storeFlags.password,
			storeFlags.dbName, storeFlags.sslMode,
		)
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
	storeCmd.Flags().StringVarP(&storeFlags.host, "host", "a", "localhost", "database host")
	storeCmd.Flags().StringVarP(&storeFlags.port, "port", "t", "5432", "database port")
	storeCmd.Flags().StringVarP(&storeFlags.dbName, "database", "x", "postgres", "database name")
	storeCmd.Flags().StringVarP(&storeFlags.sslMode, "ssl", "y", "require", "database ssl mode")
	storeCmd.Flags().StringVarP(&storeFlags.password, "pass", "p", "", "database password")
	err := storeCmd.MarkFlagRequired("pass")
	if err != nil {
		log.Fatal(err)
	}

	storeCmd.Flags().StringVarP(&storeFlags.user, "user", "u", "postgres", "database username")
	storeCmd.Flags().StringVarP(&storeFlags.repo, "repo", "g", "", "repository name")
	storeCmd.Flags().StringVarP(&storeFlags.branch, "branch", "b", "master", "branch name")
	storeCmd.Flags().Int64VarP(&storeFlags.buildID, "build-id", "i", 0, "build id")
	err = storeCmd.MarkFlagRequired("build-id")
	if err != nil {
		log.Fatal(err)
	}

	storeCmd.Flags().StringVarP(&storeFlags.runCommand, "cmd", "c", "", "test execution command")
	storeCmd.Flags().StringVarP(&storeFlags.sha, "sha", "z", "", "git sha of test run")
	storeCmd.Flags().StringVarP(&storeFlags.start, "time", "m", "", "time of test run")
	storeCmd.Flags().StringVarP(&storeFlags.src, "src", "s", "", "source test output file")
	storeCmd.Flags().BoolVar(&storeFlags.bench, "bench", false, "indicate if test ran benchmarks")
	storeCmd.Flags().BoolVar(&storeFlags.short, "short", false, "indicate if test is run with -short flag")
	storeCmd.Flags().BoolVar(&storeFlags.race, "race", false, "indicate if test is run with -race flag")
	storeCmd.Flags().StringSliceVar(&storeFlags.tags, "tags", []string{}, "comma-separated tags enabled for the run")
	storeCmd.Flags().StringSliceVarP(&storeFlags.retests, "rerun", "r", []string{}, "comma-separated source output for retests")

	RootCmd.AddCommand(storeCmd)
}
