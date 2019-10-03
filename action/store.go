package action

import (
	"time"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

var host, port, dbName, user, password, sslMode, repo, branch, config, sha, start, runCommand string
var buildID int64
var short, race bool
var tags []string

var testResults []gocop.TestResult

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "stores test results to database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		db := gocop.ConnectDB(host, port, user, password, dbName, sslMode)
		defer db.Close()

		run := gocop.TestRun{
			BuildID: buildID,
			Repo:    repo,
			Branch:  branch,
			Sha:     sha,
			Command: runCommand,
			Short:   short,
			Race:    race,
			Tags:    tags,
		}
		var err error
		if len(start) != 0 {
			run.Created, err = time.Parse(time.RFC3339, start)
			if err != nil {
				panic(err)
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

				test := gocop.TestResult{Name: entry[1], Result: r, Created: run.Created}

				if r != "skip" {
					d, err := time.ParseDuration(entry[2])
					if err != nil {
						panic(err)
					}
					test.Duration = d
				}

				testResults = append(testResults, test)
			}
		}

		if len(retests) > 0 {
			pkgs := gocop.FlakyFile(retests...)
			for _, entry := range pkgs {
				testResults = append(testResults, gocop.TestResult{Name: entry, Result: "flaky"})
			}
		}

		id, err := gocop.InsertRun(db, run)
		if err != nil {
			panic(err)
		}
		_, err = gocop.InsertTests(db, id, testResults)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
	storeCmd.Flags().StringVarP(&host, "host", "a", "localhost", "database host")
	storeCmd.Flags().StringVarP(&port, "port", "t", "5432", "database port")
	storeCmd.Flags().StringVarP(&dbName, "database", "z", "postgres", "database name")
	storeCmd.Flags().StringVarP(&sslMode, "ssl", "y", "require", "database ssl mode")
	storeCmd.Flags().StringVarP(&password, "pass", "p", "", "database password")
	storeCmd.MarkFlagRequired("pass")
	storeCmd.Flags().StringVarP(&user, "user", "u", "postgres", "database username")
	storeCmd.Flags().StringVarP(&repo, "repo", "g", "", "repository name")
	storeCmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch name")
	storeCmd.Flags().Int64VarP(&buildID, "bld-id", "i", 0, "build id")
	storeCmd.MarkFlagRequired("build-id")
	storeCmd.Flags().StringVarP(&runCommand, "cmd", "c", "", "test execution command")
	storeCmd.Flags().StringVarP(&sha, "sha", "z", "", "git sha of test run")
	storeCmd.Flags().StringVarP(&start, "time", "m", "", "time of test run")
	storeCmd.Flags().StringVarP(&src, "src", "s", "", "source test output file")
	storeCmd.Flags().BoolVar(&short, "short", false, "indicate if test is run with -short flag")
	storeCmd.Flags().BoolVar(&race, "race", false, "indicate if test is run with -race flag")
	storeCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "comma-separated tags enabled for the run")
	storeCmd.Flags().StringSliceVarP(&retests, "rerun", "r", []string{}, "comma-separated source output for retests")

	RootCmd.AddCommand(storeCmd)
}
