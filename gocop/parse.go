package gocop

import (
	"io/ioutil"
	"log"
	"regexp"
)

const (
	// ResultsPattern provides the REGEX pattern to find all package output
	ResultsPattern = `((FAIL|ok|\?)\s+([A-Za-z\.\/\_]+)\s+([0-9s\.]+|\[build failed\]|\[no test files\]))`
)

// Parse iterates over test output for all packages
func Parse(output []byte) [][]string {
	re := regexp.MustCompile(ResultsPattern)
	matches := re.FindAllStringSubmatch(string(output), -1)

	packages := make([][]string, 0)
	for _, match := range matches {
		results := make([]string, 0)
		// outcome
		results = append(results, match[2])
		// package
		results = append(results, match[3])
		// duration
		results = append(results, match[4])
		packages = append(packages, results)
	}

	return packages
}

// ParseFailed iterates over test output for failed packages
func ParseFailed(output []byte) []string {
	re := regexp.MustCompile(ResultsPattern)
	matches := re.FindAllSubmatch(output, -1)

	packages := make([]string, 0)
	for _, match := range matches {
		if string(match[2]) == "FAIL" {
			packages = append(packages, string(match[3]))
		}
	}

	return packages
}

// ParseFileFailed reads a file to Parse() failed packages
func ParseFileFailed(path string) []string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return ParseFailed(content)
}

// ParseFile reads a file to Parse() results
func ParseFile(path string) [][]string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return Parse(content)
}
