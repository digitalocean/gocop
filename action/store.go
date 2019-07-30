package action

import (
	"time"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

var host string
var port string
var user string
var password string
var repo string
var branch string
var buildID int64
var config string
var sha string
var start string
var runCommand string

var testResults []gocop.TestResult

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "stores test results to database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		db := gocop.ConnectDB(host, port, user, password)
		defer db.Close()

		if len(src) > 0 {
			pkgs := gocop.ParseFile(src)
			for _, entry := range pkgs {
				testResults = append(testResults, gocop.TestResult{Name: entry, Result: "fail"})
			}
		}

		if len(retests) > 0 {
			pkgs := gocop.FlakyFile(retests...)
			for _, entry := range pkgs {
				testResults = append(testResults, gocop.TestResult{Name: entry, Result: "flaky"})
			}
		}

		var startTime time.Time
		var err error
		if len(start) != 0 {
			startTime, err = time.Parse(time.RFC3339, start)
			if err != nil {
				panic(err)
			}
		} else {
			startTime = time.Now().UTC()
		}

		id, err := gocop.InsertRun(db, buildID, repo, branch, sha, runCommand, startTime)
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

	storeCmd.Flags().StringSliceVarP(&retests, "rerun", "r", []string{}, "source output for retests")

	failedCmd.AddCommand(storeCmd)
}
