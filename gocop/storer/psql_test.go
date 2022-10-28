//go:build integration
// +build integration

package storer

import (
	"context"
	"os"
	"testing"
	"time"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func TestInsertResults(t *testing.T) {
	ctx := context.Background()
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	name := getenv("DB_NAME", "postgres")
	ssl := getenv("DB_SSL", "disable")
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASS", "testuser")
	db, err := NewPSQL(host, port, user, password, name, ssl)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	run := TestRun{
		BuildID: 2,
		Repo:    "test_repo",
		Branch:  "master",
		Created: time.Now().UTC(),
	}

	err = db.InsertRun(ctx, run)
	if err != nil {
		t.Fatal(err)
	}

	testResults := make([]TestResult, 0)
	result := TestResult{
		Package:  "test1",
		Result:   "fail",
		Created:  run.Created,
		Duration: time.Second,
		Coverage: 0.834,
	}

	testResults = append(testResults, result)
	err = db.InsertTests(ctx, testResults)
	if err != nil {
		t.Fatal(err)
	}

	tests, err := db.GetTests(ctx, run.Created)
	if err != nil {
		t.Fatal(err)
	}
	if len(tests) != 1 {
		t.Fatalf("expected 1 test result, got %d", len(tests))
	}
}
