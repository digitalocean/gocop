package gocop

import (
	"io/ioutil"
	"log"
	"regexp"
)

// FailurePattern provides the REGEX pattern to find failed packages in test output
const FailurePattern = `FAIL\s+([A-Za-z\.\/]+)\s+[0-9s\.]+`

// Parse iterates over test output for failed packages
func Parse(output []byte) []string {
	re := regexp.MustCompile(FailurePattern)
	matches := re.FindAllSubmatch(output, -1)

	packages := make([]string, 0)
	for _, match := range matches {
		packages = append(packages, string(match[1]))
	}

	return packages
}

// ParseFile reads a file to Parse() failed packages
func ParseFile(path string) []string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return Parse(content)
}
