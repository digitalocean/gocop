package gocop

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq" // import postgres driver
)

// TestRun contains data about a test run
type TestRun struct {
	Created   time.Time
	Repo      string
	Branch    string
	Sha       string
	BuildID   int64
	Config    string
	Command   string
	Benchmark bool
	Short     bool
	Race      bool
	Tags      []string
	Duration  time.Duration
	Results   []TestResult
}

// TestResult contains data about a test result
type TestResult struct {
	Created  time.Time
	Package  string
	Result   string
	Duration time.Duration
	Coverage float64
}

func (t *TestResult) setResult(v string) {
	switch v {
	case "ok":
		t.Result = "pass"
	case "FAIL":
		t.Result = "fail"
	case "?":
		t.Result = "skip"
	}
}

func (t *TestResult) setDuration(v string) error {
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Println(err)
		return err
	}
	t.Duration = d

	return nil
}

func (t *TestResult) setCoverage(v string) error {
	if v != "" && v != "[no statements]" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Println(err)
			return err
		}

		t.Coverage = f / 100
	}

	return nil
}

// ConvertResults is a helper for adding test results to a run from parsed output
func (run *TestRun) ConvertResults(pkgs [][]string) error {
	for _, entry := range pkgs {
		test := TestResult{Package: entry[1], Created: run.Created}
		test.setResult(entry[0])

		if r != "skip" {
			d, err := time.ParseDuration(entry[2])
			if err != nil {
				log.Println(err)
				continue
			}
			test.Duration = d
		}

		if entry[3] != "" {
			f, err := strconv.ParseFloat(entry[3], 64)
			if err != nil {
				log.Println(err)
				continue
			}

			test.Coverage = f / 100
		}

		run.AddResult(test)
	}

	return nil
}

// NewTestResult is a helper for creating a test result from parsed output
func NewTestResult(created time.Time, pkg []string) (TestResult, error) {
	test := TestResult{Package: pkg[1], Created: created}
	test.setResult(pkg[0])

	if pkg[2] == "[build failed]" || test.Result == "skip" {
		return test, nil
	}

	err := test.setDuration(pkg[2])
	if err != nil {
		return test, err
	}

	if pkg[3] == "" || pkg[3] == "[no statements]" {
		return test, nil
	}

	err = test.setCoverage(pkg[3])
	if err != nil {
		return test, err
	}

	return test, nil
}

// AddResult is a helper for adding a test result to a run
func (run *TestRun) AddResult(r TestResult) {
	run.Results = append(run.Results, r)
}

// ConnectDB connects to the database
func ConnectDB(host, port, user, password, dbname, sslmode string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// InsertRun inserts a new entry to the run table in the database
func InsertRun(db *sql.DB, run TestRun) (sql.Result, error) {
	sqlStr := `
    INSERT INTO run (created, build_id, repo, duration, branch, sha, cmd, benchmark, short, race, tags)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
  `

	sort.Strings(run.Tags)
	tags := strings.Join(run.Tags, " ")

	res, err := db.Exec(
		sqlStr,
		run.Created,
		run.BuildID,
		run.Repo,
		run.Duration,
		run.Branch,
		run.Sha,
		run.Command,
		run.Benchmark,
		run.Short,
		run.Race,
		tags,
	)

	return res, err
}

// GetRun retrieves information about a run
func GetRun(db *sql.DB, buildID int64) *sql.Row {
	sqlStr := `SELECT build_id, created, duration, cmd, repo, branch, sha, benchmark, race, short, tags
    FROM run
    WHERE build_id=$1`

	return db.QueryRow(sqlStr, buildID)
}

// InsertTests adds test results to database
func InsertTests(db *sql.DB, created time.Time, testResults []TestResult) (sql.Result, error) {
	sqlStr := "INSERT INTO test(created, package, result, duration, coverage) VALUES "
	vals := []interface{}{}

	for _, row := range testResults {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, row.Created, row.Package, row.Result, row.Duration/time.Millisecond, row.Coverage)
	}
	if len(vals) == 0 {
		return nil, errors.New("no test results found")
	}

	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	// Replacing ? with $n for postgres
	sqlStr = ReplaceSQL(sqlStr, "?")
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(vals...)
}

// GetTests retrieves test results for a build
func GetTests(db *sql.DB, created time.Time) (*sql.Rows, error) {
	sqlStr := `
    SELECT created, package, result, duration, coverage
    FROM test
    WHERE created=$1
  `

	return db.Query(sqlStr, created)
}

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
