package action

import (
	"log"
	"strconv"
	"time"

	"github.com/digitalocean/gocop/gocop"
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

var testResults []gocop.TestResult

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "stores test results to database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		db := gocop.ConnectDB(
			storeFlags.host, storeFlags.port,
			storeFlags.user, storeFlags.password,
			storeFlags.dbName, storeFlags.sslMode,
		)
		defer func() {
			err = db.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}()

		run := gocop.TestRun{
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
			pkgs := gocop.ParseFile(storeFlags.src)
			for _, entry := range pkgs {
				var r string
				switch entry[0] {
				case "ok":
					r = "pass"
				case "FAIL":
					r = "fail"
				case "?":
					r = "skip"
				}

				test := gocop.TestResult{Package: entry[1], Result: r, Created: run.Created}

				if r != "skip" {
					d, err := time.ParseDuration(entry[2])
					if err != nil {
						log.Fatal(err)
					}
					test.Duration = d
				}

				if entry[3] != "" {
					f, err := strconv.ParseFloat(entry[3], 64)
					if err != nil {
						log.Fatal(err)
					}

					test.Coverage = f / 100
				}

				testResults = append(testResults, test)
			}
		}

		if len(storeFlags.retests) > 0 {
			pkgs := gocop.FlakyFile(storeFlags.retests...)
			for _, entry := range pkgs {
				testResults = append(testResults, gocop.TestResult{Package: entry, Result: "flaky"})
			}
		}

		_, err = gocop.InsertRun(db, run)
		if err != nil {
			log.Fatal(err)
		}

		_, err = gocop.InsertTests(db, run.Created, testResults)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
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
