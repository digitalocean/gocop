package gocop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/digitalocean/gocop/gocop/storer"
)

const (
	// ResultsPattern provides the REGEX pattern to find all package output
	resultsPattern = `((FAIL|ok|\?)\s+([\w\.\/\-]+)\s+([0-9s\.]+|\[build failed\]|\[no test files\])(\n|\s+coverage\:\s+([\d\.]+)\%\s+))`
	// Test2JSONCoveragePattern parses the coverage percentage out of a test2json event
	test2JSONCoveragePattern = `"coverage\:\s+([\d\.]+)`
)

var (
	resultsPatternRE    = regexp.MustCompile(resultsPattern)
	test2JSONCoverageRE = regexp.MustCompile(test2JSONCoveragePattern)
)

// Parser represents a go test output parser
type Parser interface {
	Parse([]byte) ([]storer.TestResult, error)
}

// StandardParser will parse standard `go test` output.
type StandardParser struct{}

// Parse iterates over test output for all packages
func (p *StandardParser) Parse(output []byte) ([]storer.TestResult, error) {
	var results []storer.TestResult

	matches := resultsPatternRE.FindAllStringSubmatch(string(output), -1)

	for _, match := range matches {
		event := storer.TestResult{
			Package: match[3],
		}

		// result
		switch match[2] {
		case "ok":
			event.Result = "pass"
		case "FAIL":
			event.Result = "fail"
		case "?":
			event.Result = "skip"
		}

		// duration
		if event.Result != "skip" {
			// best effort, ignoring error
			d, _ := time.ParseDuration(match[4])
			event.Duration = d
		}

		// coverage
		if match[6] != "" {
			// best effort, ignoring error
			f, _ := strconv.ParseFloat(match[6], 64)
			event.Coverage = f / 100
		}

		results = append(results, event)
	}

	return results, nil
}

// Test2JSONParser will parse test2json `go test` output specified with `-json`
type Test2JSONParser struct {
	IncludeIndividualTests bool
}

// https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/cmd/internal/test2json/test2json.go
type test2JSONEvent struct {
	Time    *time.Time `json:",omitempty"`
	Action  string
	Package string   `json:",omitempty"`
	Test    string   `json:",omitempty"`
	Elapsed *float64 `json:",omitempty"`
	// varies from internal representation to support decoding
	Output json.RawMessage `json:",omitempty"`
}

func (e *test2JSONEvent) key() string {
	return fmt.Sprintf("%s:%s", e.Package, e.Test)
}

// Parse iterates over test output for all packages
func (p *Test2JSONParser) Parse(output []byte) ([]storer.TestResult, error) {
	resultMap := map[string]*storer.TestResult{}

	for _, l := range bytes.Split(output, []byte("\n")) {
		if len(bytes.TrimSpace(l)) == 0 {
			continue
		}

		t2j := test2JSONEvent{}
		err := json.Unmarshal(l, &t2j)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling test2json event: %w", err)
		}

		var event *storer.TestResult
		if _, ok := resultMap[t2j.key()]; ok {
			event = resultMap[t2j.key()]
		} else {
			event = &storer.TestResult{
				Package: t2j.Package,
				Test:    t2j.Test,
			}
			resultMap[t2j.key()] = event
		}

		switch t2j.Action {
		case "pass", "fail", "skip":
			event.Result = t2j.Action

			// duration
			if t2j.Elapsed != nil {
				event.Duration = time.Duration(*t2j.Elapsed * float64(time.Second))
			}
		case "output":
			matches := test2JSONCoverageRE.FindStringSubmatch(string(t2j.Output))
			if len(matches) == 2 {
				// best effort, ignoring error
				f, _ := strconv.ParseFloat(matches[1], 64)
				event.Coverage = f / 100
			}
		}
	}

	var results []storer.TestResult
	for _, r := range resultMap {
		if !p.IncludeIndividualTests && r.Test != "" {
			continue
		}
		results = append(results, *r)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key() < results[j].Key()
	})

	return results, nil
}

// ParseFailedPackages iterates over test output for failed packages
func ParseFailedPackages(p Parser, output []byte) ([]string, error) {
	tests, err := p.Parse(output)
	if err != nil {
		return nil, err
	}

	var packages []string
	for _, t := range tests {
		if t.Result == "fail" && t.Test == "" {
			packages = append(packages, t.Package)
		}
	}

	return packages, nil
}

// ParseFileFailedPackages reads a file to Parse() failed packages
func ParseFileFailedPackages(p Parser, path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseFailedPackages(p, content)
}

// ParseFile reads a file to Parse() results
func ParseFile(p Parser, path string) ([]storer.TestResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return p.Parse(content)
}
