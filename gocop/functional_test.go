package gocop_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/gocop/gocop"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
)

func TestParseFileFailedPackages(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "finds multiple failed packages",
			input: "testdata/run0.txt",
			want:  []string{"github.com/digitalocean/gocop/sample/fail", "github.com/digitalocean/gocop/sample/failbuild", "github.com/digitalocean/gocop/sample/flaky"},
		},
		{
			name:  "finds single failed packages",
			input: "testdata/run1.txt",
			want:  []string{"github.com/digitalocean/gocop/sample/fail", "github.com/digitalocean/gocop/sample/failbuild"},
		},
	}

	for _, tt := range tests {
		tt := tt
		o.Spec(tt.name, func(expect expect.Expectation) {
			got, err := gocop.ParseFileFailedPackages(&gocop.StandardParser{}, tt.input)
			expect(err).To(matchers.BeNil())
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}

func TestFlakyFilePackages(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "finds zero flaky packages",
			input: []string{"testdata/onlypass.txt"},
			want:  []string{},
		},
		{
			name:  "finds single flaky package",
			input: []string{"testdata/run0.txt", "testdata/run1.txt", "testdata/run2.txt", "testdata/run3.txt"},
			want:  []string{"github.com/digitalocean/gocop/sample/flaky"},
		},
	}

	for _, tt := range tests {
		tt := tt
		o.Spec(tt.name, func(expect expect.Expectation) {
			got, err := gocop.FlakyFilePackages(&gocop.StandardParser{}, tt.input...)
			expect(err).To(matchers.BeNil())
			fmt.Println("got: ", got)
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}

func TestMatchingParsedOutput(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name          string
		standardInput string
		jsonInput     string
	}{
		{
			name:          "run0",
			standardInput: "testdata/v2/run0.txt",
			jsonInput:     "testdata/v2/run0.json",
		},
		{
			name:          "run1",
			standardInput: "testdata/v2/run1.txt",
			jsonInput:     "testdata/v2/run1.json",
		},
		{
			name:          "run2",
			standardInput: "testdata/v2/run2.txt",
			jsonInput:     "testdata/v2/run2.json",
		},
		{
			name:          "run3",
			standardInput: "testdata/v2/run3.txt",
			jsonInput:     "testdata/v2/run3.json",
		},
	}

	for _, tt := range tests {
		tt := tt
		o.Spec(tt.name, func(expect expect.Expectation) {
			gotStandard, err := gocop.FlakyFilePackages(&gocop.StandardParser{}, tt.standardInput)
			expect(err).To(matchers.BeNil())
			gotJSON, err := gocop.FlakyFilePackages(&gocop.Test2JSONParser{}, tt.jsonInput)
			expect(err).To(matchers.BeNil())
			expect(gotStandard).To(matchers.Equal(gotJSON))
		})
	}
}
