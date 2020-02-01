package action

import (
	"log"
	"strconv"
	"time"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

var host, port, dbName, user, password, sslMode, repo, branch, sha, start, runCommand string
var buildID int64
var bench, short, race bool
var tags []string

var testResults []gocop.TestResult

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "stores test results to database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		db := gocop.ConnectDB(host, port, user, password, dbName, sslMode)
		defer func() {
			err = db.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}()

		run := gocop.TestRun{
			BuildID:   buildID,
			Repo:      repo,
			Branch:    branch,
			Sha:       sha,
			Command:   runCommand,
			Benchmark: bench,
			Short:     short,
			Race:      race,
			Tags:      tags,
		}
		if len(start) != 0 {
			run.Created, err = time.Parse(time.RFC3339, start)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			run.Created = time.Now().UTC()
		}

		if len(src) > 0 {
			pkgs := gocop.ParseFile(src)
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

		if len(retests) > 0 {
			pkgs := gocop.FlakyFile(retests...)
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
	storeCmd.Flags().StringVarP(&host, "host", "a", "localhost", "database host")
	storeCmd.Flags().StringVarP(&port, "port", "t", "5432", "database port")
	storeCmd.Flags().StringVarP(&dbName, "database", "x", "postgres", "database name")
	storeCmd.Flags().StringVarP(&sslMode, "ssl", "y", "require", "database ssl mode")
	storeCmd.Flags().StringVarP(&password, "pass", "p", "", "database password")
	err := storeCmd.MarkFlagRequired("pass")
	if err != nil {
		log.Fatal(err)
	}

	storeCmd.Flags().StringVarP(&user, "user", "u", "postgres", "database username")
	storeCmd.Flags().StringVarP(&repo, "repo", "g", "", "repository name")
	storeCmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch name")
	storeCmd.Flags().Int64VarP(&buildID, "bld-id", "i", 0, "build id")
	err = storeCmd.MarkFlagRequired("build-id")
	if err != nil {
		log.Fatal(err)
	}

	storeCmd.Flags().StringVarP(&runCommand, "cmd", "c", "", "test execution command")
	storeCmd.Flags().StringVarP(&sha, "sha", "z", "", "git sha of test run")
	storeCmd.Flags().StringVarP(&start, "time", "m", "", "time of test run")
	storeCmd.Flags().StringVarP(&src, "src", "s", "", "source test output file")
	storeCmd.Flags().BoolVar(&bench, "bench", false, "indicate if test ran benchmarks")
	storeCmd.Flags().BoolVar(&short, "short", false, "indicate if test is run with -short flag")
	storeCmd.Flags().BoolVar(&race, "race", false, "indicate if test is run with -race flag")
	storeCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "comma-separated tags enabled for the run")
	storeCmd.Flags().StringSliceVarP(&retests, "rerun", "r", []string{}, "comma-separated source output for retests")

	RootCmd.AddCommand(storeCmd)
}
