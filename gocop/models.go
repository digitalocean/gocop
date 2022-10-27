package gocop

import "time"

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
}

// TestResult contains data about a test result
type TestResult struct {
	Created  time.Time
	Package  string
	Result   string
	Duration time.Duration
	Coverage float64
}
