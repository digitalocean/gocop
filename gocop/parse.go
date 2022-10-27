package gocop

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	// ResultsPattern provides the REGEX pattern to find all package output
	ResultsPattern = `((FAIL|ok|\?)\s+([\w\.\/\-]+)\s+([0-9s\.]+|\[build failed\]|\[no test files\])(\n|\s+coverage\:\s+([\d\.]+)\%\s+))`
)

// Parse iterates over test output for all packages
func Parse(output []byte) ([]TestResult, error) {
	var results []TestResult

	re := regexp.MustCompile(ResultsPattern)
	matches := re.FindAllStringSubmatch(string(output), -1)

	for _, match := range matches {
		event := TestResult{
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

// ParseFailedPackages iterates over test output for failed packages
func ParseFailedPackages(output []byte) ([]string, error) {
	tests, err := Parse(output)
	if err != nil {
		return nil, err
	}

	var packages []string
	for _, t := range tests {
		if t.Result == "fail" {
			packages = append(packages, t.Package)
		}
	}

	return packages, nil
}

// ParseFileFailedPackages reads a file to Parse() failed packages
func ParseFileFailedPackages(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return ParseFailedPackages(content)
}

// ParseFile reads a file to Parse() results
func ParseFile(path string) ([]TestResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return Parse(content)
}
