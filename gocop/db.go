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
	Repo     string
	Branch   string
	Duration time.Duration
	BuildID  int64
	Config   string
	Sha      string
	Created  time.Time
	Command  string
	Short    bool
	Race     bool
	Tags     []string
}

// TestResult contains data about a test result
type TestResult struct {
	Name     string
	Result   string
	Created  time.Time
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
func InsertRun(db *sql.DB, run TestRun) (int64, error) {
	testRunInsert := `INSERT INTO run (build_id, repo, duration, branch, sha, run_cmd, created, short, race, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING build_id`
	id := int64(0)
	sort.Strings(run.Tags)
	tags := strings.Join(run.Tags, " ")
	err := db.QueryRow(testRunInsert,
		run.BuildID,
		run.Repo,
		run.Duration,
		run.Branch,
		run.Sha,
		run.Command,
		run.Created,
		run.Short,
		run.Race,
		tags).Scan(&id)
	return id, err
}

// GetRun retrieves information about a run
func GetRun(db *sql.DB, buildID int64) *sql.Row {
	runSelect := `SELECT build_id, created, duration, run_cmd, repo, branch, sha, race, short, tags
		FROM run
		WHERE build_id=$1`

	return db.QueryRow(runSelect, buildID)
}

// InsertTests adds test results to database
func InsertTests(db *sql.DB, runID int64, testResults []TestResult) (sql.Result, error) {
	sqlStr := "INSERT INTO test(run_id, result, name, duration, created) VALUES "
	vals := []interface{}{}

	for _, row := range testResults {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, runID, row.Result, row.Name, row.Duration/time.Millisecond, row.Created)
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
func GetTests(db *sql.DB, buildID int64) (*sql.Rows, error) {
	runSelect := `SELECT id, run_id, result, name, duration, created
		FROM test
		WHERE run_id=$1`

	return db.Query(runSelect, buildID)
}

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
