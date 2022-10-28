package storer

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// TestRun contains data about a test run
type TestRun struct {
	ID        uuid.UUID
	Created   time.Time
	Team      string
	JobName   string
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
}

// TestResult contains data about a test result
type TestResult struct {
	RunID    uuid.UUID
	Created  time.Time
	Package  string
	Test     string
	Result   string
	Duration time.Duration
	Coverage float64
}

// Key represents a canonical key for a test result based on package and test name
func (r *TestResult) Key() string {
	return fmt.Sprintf("%s:%s", r.Package, r.Test)
}
