package storer

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq" // import postgres driver
)

// PSQL is a PostgreSQL based storer.
type PSQL struct {
	db *sql.DB
}

// NewPSQL creates a new PostgresSQL storer instance.
func NewPSQL(host, port, user, password, dbname, sslmode string) (Storer, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("opening sql connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return &PSQL{db}, nil
}

// InsertRun inserts a new entry to the run table in the database
func (s *PSQL) InsertRun(ctx context.Context, run TestRun) error {
	sqlStr := `
		INSERT INTO run (created, build_id, repo, duration, branch, sha, cmd, benchmark, short, race, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	sort.Strings(run.Tags)
	tags := strings.Join(run.Tags, " ")

	_, err := s.db.ExecContext(
		ctx,
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

	return err
}

// GetRun retrieves information about a run
func (s *PSQL) GetRun(ctx context.Context, buildID int64) (*TestRun, error) {
	sqlStr := `SELECT build_id, created, duration, cmd, repo, branch, sha, benchmark, race, short, tags
		FROM run
		WHERE build_id=$1`

	var r TestRun
	err := s.db.QueryRowContext(ctx, sqlStr, buildID).Scan(
		&r.BuildID,
		&r.Created,
		&r.Duration,
		&r.Command,
		&r.Repo,
		&r.Branch,
		&r.Sha,
		&r.Benchmark,
		&r.Race,
		&r.Short,
		&r.Tags,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// InsertTests adds test results to database
func (s *PSQL) InsertTests(ctx context.Context, testResults []TestResult) error {
	sqlStr := "INSERT INTO test(created, package, result, duration, coverage) VALUES "
	vals := []interface{}{}

	for _, row := range testResults {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, row.Created, row.Package, row.Result, int(row.Duration/time.Millisecond), row.Coverage)
	}
	if len(vals) == 0 {
		return nil
	}

	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	// Replacing ? with $n for postgres
	sqlStr = replacePSQL(sqlStr, "?")
	stmt, err := s.db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, vals...)
	return err
}

// GetTests retrieves test results for a build
func (s *PSQL) GetTests(ctx context.Context, created time.Time) ([]*TestResult, error) {
	sqlStr := `
		SELECT created, package, result, duration, coverage
		FROM test
		WHERE created=$1
	`

	row, err := s.db.QueryContext(ctx, sqlStr, created)
	if err != nil {
		return nil, err
	}

	var results []*TestResult
	for row.Next() {
		var r TestResult
		err := row.Scan(
			&r.Created,
			&r.Package,
			&r.Result,
			&r.Duration,
			&r.Coverage,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Close closes the underlying SQL connection.
func (s *PSQL) Close() error {
	return s.db.Close()
}

// replacePSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func replacePSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
