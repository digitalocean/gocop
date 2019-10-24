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
	Created  time.Time
	Repo     string
	Branch   string
	Sha      string
	BuildID  int64
	Config   string
	Command  string
	Short    bool
	Race     bool
	Tags     []string
	Duration time.Duration
}

// TestResult contains data about a test result
type TestResult struct {
	Created  time.Time
	Package  string
	Result   string
	Duration time.Duration
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
		INSERT INTO run (created, build_id, repo, duration, branch, sha, cmd, short, race, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		run.Short,
		run.Race,
		tags,
	)

	return res, err
}

// GetRun retrieves information about a run
func GetRun(db *sql.DB, buildID int64) *sql.Row {
	sqlStr := `SELECT build_id, created, duration, cmd, repo, branch, sha, race, short, tags
		FROM run
		WHERE build_id=$1`

	return db.QueryRow(sqlStr, buildID)
}

// InsertTests adds test results to database
func InsertTests(db *sql.DB, created time.Time, testResults []TestResult) (sql.Result, error) {
	sqlStr := "INSERT INTO test(created, package, result, duration) VALUES "
	vals := []interface{}{}

	for _, row := range testResults {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, row.Created, row.Package, row.Result, row.Duration/time.Millisecond)
	}
	if len(vals) == 0 {
		return nil, errors.New("No test results found")
	}

	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	// Replacing ? with $n for postgres
	sqlStr = ReplaceSQL(sqlStr, "?")
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Fatal(err)
	}

	return stmt.Exec(vals...)
}

// GetTests retrieves test results for a build
// CURRENTLY BROKEN
func GetTests(db *sql.DB, buildID int64) (*sql.Rows, error) {
	sqlStr := `SELECT id, run_id, result, name, duration, created
		FROM test
		WHERE run_id=$1`

	return db.Query(sqlStr, buildID)
}

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
