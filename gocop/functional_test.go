package gocop_test

import (
	"testing"

	"github.com/digitalocean/gocop/gocop"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
)

func TestParseFileFailed(t *testing.T) {
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
		o.Spec(tt.name, func(expect expect.Expectation) {
			got := gocop.ParseFileFailed(tt.input)
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}

func TestFlakyFile(t *testing.T) {
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
			input: []string{"testdata/run0.txt", "testdata/run1.txt", "testdata/run3.txt"},
			want:  []string{},
		},
		{
			name:  "finds single flaky package",
			input: []string{"testdata/run0.txt", "testdata/run1.txt", "testdata/run2.txt", "testdata/run3.txt"},
			want:  []string{"github.com/digitalocean/gocop/sample/flaky"},
		},
	}

	for _, tt := range tests {
		o.Spec(tt.name, func(expect expect.Expectation) {
			got := gocop.FlakyFile(tt.input...)
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}
