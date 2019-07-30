// +build integration

package gocop_test

import (
	"fmt"
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
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASS", "testuser")
	db := gocop.ConnectDB(host, port, user, password)
	defer db.Close()

	runID, err := gocop.InsertRun(db, 2, "test_repo", "master", "1")
	if err != nil {
		t.Error(err)
	}

	testResults := make([]gocop.TestResult, 0)
	result := gocop.TestResult{Name: "test1", Result: "fail", Duration: time.Second}
	testResults = append(testResults, result)
	_, err = gocop.InsertTests(db, runID, testResults)
	if err != nil {
		t.Error(err)
		fmt.Println(err)
	}

	rows, err := gocop.GetTests(db, runID)
	var count int
	for rows.Next() {
		count = count + 1
	}
	if count != 1 {
		t.Fail()
	}
}
