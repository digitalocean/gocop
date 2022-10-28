package storer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/muesli/reflow/indent"
)

type Stdout struct{}

// NewStdout creates a new Stdout storer
func NewStdout() (Storer, error) {
	return &Stdout{}, nil
}

func (s *Stdout) InsertRun(ctx context.Context, run TestRun) error {
	t := table.NewWriter()
	t.SetOutputMirror(indent.NewWriterPipe(os.Stdout, 2, nil))
	t.SetStyle(table.StyleColoredDark)
	t.SetTitle("Test Run")
	t.AppendHeader(table.Row{
		"Job", "Repo", "Duration", "Flags",
	})
	var flags []string
	if run.Benchmark {
		flags = append(flags, "bench")
	}
	if run.Short {
		flags = append(flags, "short")
	}
	if run.Race {
		flags = append(flags, "race")
	}
	t.AppendRow(table.Row{
		fmt.Sprintf("%s#%d", run.JobName, run.BuildID),
		fmt.Sprintf("%s@%s", run.Repo, run.Branch),
		run.Duration.String(),
		fmt.Sprint(flags),
	})
	t.Render()
	fmt.Print("\n")

	return nil
}

func (s *Stdout) GetRun(ctx context.Context, buildID int64) (*TestRun, error) {
	return nil, fmt.Errorf("not supported by stdout storer")
}

func (s *Stdout) InsertTests(ctx context.Context, testResults []TestResult) error {
	t := table.NewWriter()
	t.SetOutputMirror(indent.NewWriterPipe(os.Stdout, 2, nil))
	t.SetStyle(table.StyleColoredDark)
	t.SetTitle("Test Results")
	t.AppendHeader(table.Row{
		"Package", "Test", "Result", "Duration", "Coverage",
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})

	for _, r := range testResults {
		var resultColor func(string, ...interface{}) string
		switch r.Result {
		case "pass":
			resultColor = color.GreenString
		case "fail":
			resultColor = color.RedString
		case "skip":
			resultColor = color.YellowString
		}
		result := resultColor(r.Result)

		t.AppendRow(table.Row{
			r.Package,
			r.Test,
			result,
			r.Duration.String(),
			fmt.Sprintf("%1.2f%%", r.Coverage/100),
		})
	}

	t.Render()
	fmt.Print("\n")

	return nil
}

func (s *Stdout) GetTests(ctx context.Context, created time.Time) ([]*TestResult, error) {
	return nil, fmt.Errorf("not supported by stdout storer")
}

func (s *Stdout) Close() error {
	return nil
}
