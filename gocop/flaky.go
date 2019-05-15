package gocop

import (
	"io/ioutil"
	"log"
)

// Flaky reviews test output from multiple attempts and identifies potentially flaky packages
func Flaky(runs ...[]byte) []string {
	failCount := make(map[string]int, 0)
	runCount := len(runs)
	flaky := make([]string, 0)

	for _, run := range runs {
		pkgs := Parse(run)
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

	return flaky
}

// FlakyFile reviews test output from multiple files to identify flaky packages
func FlakyFile(files ...string) []string {
	runs := make([][]byte, 0)
	for _, file := range files {
		run, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		runs = append(runs, run)
	}

	return Flaky(runs...)
}
