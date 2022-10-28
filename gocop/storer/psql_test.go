//go:build integration
// +build integration

package storer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"
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

	t.Run("with only legacy fields", func(t *testing.T) {
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
	})

	t.Run("with new granular test fields", func(t *testing.T) {
		expectRun := TestRun{
			ID: uuid.NullUUID{
				Valid: true,
				UUID:  uuid.Must(uuid.NewV4()),
			},
			Team:    "my-team",
			JobName: "e2e",
			BuildID: 3,
			Repo:    "test_repo",
			Branch:  "master",
			Created: time.Now().UTC(),
		}

		err = db.InsertRun(ctx, expectRun)
		if err != nil {
			t.Fatal(err)
		}

		actualRun, err := db.GetRun(ctx, expectRun.BuildID)
		if err != nil {
			t.Fatal(err)
		}
		if expectRun.ID.UUID != actualRun.ID.UUID {
			t.Fatalf("expect run to have ID %s, got %s", expectRun.ID.UUID, actualRun.ID.UUID)
		}
		if expectRun.Team != actualRun.Team {
			t.Fatalf("expect run to have team %s, got %s", expectRun.Team, actualRun.Team)
		}
		if expectRun.JobName != actualRun.JobName {
			t.Fatalf("expect run to have job name %s, got %s", expectRun.JobName, actualRun.JobName)
		}

		testResults := make([]TestResult, 0)
		expectTest := TestResult{
			RunID:    expectRun.ID,
			Package:  "test1",
			Test:     "individual-test",
			Result:   "fail",
			Created:  expectRun.Created,
			Duration: time.Second,
			Coverage: 0.834,
		}

		testResults = append(testResults, expectTest)
		err = db.InsertTests(ctx, testResults)
		if err != nil {
			t.Fatal(err)
		}

		tests, err := db.GetTests(ctx, expectRun.Created)
		if err != nil {
			t.Fatal(err)
		}
		if len(tests) != 1 {
			t.Fatalf("expected 1 test result, got %d", len(tests))
		}
		actualTest := tests[0]
		if actualTest.RunID.UUID != expectTest.RunID.UUID {
			t.Fatalf("expect test result to have run ID %s, got %s", expectTest.RunID.UUID, actualTest.RunID.UUID)
		}
		if actualTest.Test != expectTest.Test {
			t.Fatalf("expect test result to have test %s, got %s", expectTest.Test, actualTest.Test)
		}
	})
}
