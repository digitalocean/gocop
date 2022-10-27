package gocop

import (
	"log"
	"os"
)

// FlakyPackages reviews test output from multiple attempts and identifies potentially flaky packages
func FlakyPackages(runs ...[]byte) ([]string, error) {
	failCount := make(map[string]int)
	runCount := len(runs)
	flaky := make([]string, 0)

	for _, run := range runs {
		pkgs, err := ParseFailedPackages(run)
		if err != nil {
			return nil, err
		}
		for _, pkg := range pkgs {
			_, ok := failCount[pkg]
			if !ok {
				failCount[pkg] = 0
			}
			failCount[pkg] = failCount[pkg] + 1
		}
	}

	for k, v := range failCount {
		if v < runCount {
			flaky = append(flaky, k)
		}
	}

	return flaky, nil
}

// FlakyFilePackages reviews test output from multiple files to identify flaky packages
func FlakyFilePackages(files ...string) ([]string, error) {
	runs := make([][]byte, 0)
	for _, file := range files {
		run, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		runs = append(runs, run)
	}

	return FlakyPackages(runs...)
}
