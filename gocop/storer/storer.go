package storer

import (
	"context"
	"time"
)

// Storer describes a data storage driver.
type Storer interface {
	InsertRun(ctx context.Context, run TestRun) error
	GetRun(ctx context.Context, buildID int64) (*TestRun, error)
	InsertTests(ctx context.Context, testResults []TestResult) error
	GetTests(ctx context.Context, created time.Time) ([]*TestResult, error)

	Close() error
}
