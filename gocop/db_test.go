// +build integration

package gocop_test

import (
	"os"
	"testing"
	"time"

	"github.com/digitalocean/gocop/gocop"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func TestInsertResults(t *testing.T) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	name := getenv("DB_NAME", "postgres")
	ssl := getenv("DB_SSL", "disable")
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASS", "testuser")
	db := gocop.ConnectDB(host, port, user, password, name, ssl)
	defer db.Close()

	run := gocop.TestRun{
		BuildID: 2,
		Repo:    "test_repo",
		Branch:  "master",
		Created: time.Now().UTC(),
	}

	_, err := gocop.InsertRun(db, run)
	if err != nil {
		t.Error(err)
	}

	testResults := make([]gocop.TestResult, 0)
	result := gocop.TestResult{
		Package:  "test1",
		Result:   "fail",
		Created:  run.Created,
		Duration: time.Second,
		Coverage: 88.0,
	}

	testResults = append(testResults, result)
	_, err = gocop.InsertTests(db, run.Created, testResults)
	if err != nil {
		t.Error(err)
	}

	rows, err := gocop.GetTests(db, run.Created)
	var count int
	for rows.Next() {
		count = count + 1
	}
	if count != 1 {
		t.Fail()
	}
}
