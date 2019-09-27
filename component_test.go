package main_test

import (
	"log"
	"os/exec"
	"testing"

	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestFailedPackages(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name   string
		action string
		input  string
		want   string
	}{
		{
			name:   "finds multiple failed packages",
			action: "failed",
			input:  "gocop/testdata/run0.txt",
			want:   "github.com/digitalocean/gocop/sample/fail\ngithub.com/digitalocean/gocop/sample/flaky",
		},
		{
			name:   "finds single failed packages",
			action: "failed",
			input:  "gocop/testdata/run1.txt",
			want:   "github.com/digitalocean/gocop/sample/fail",
		},
	}

	for _, tt := range tests {
		o.Spec(tt.name, func(expect expect.Expectation) {
			got, err := exec.Command("go", "run", "main.go", tt.action, "-s", tt.input).Output()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\n\n", got)

			expect(string(got)).To(Equal(tt.want))
		})
	}
}

func TestFlakyPackages(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "finds zero flaky packages",
			args: []string{"run", "main.go", "flaky", "-r", "gocop/testdata/run0.txt", "-r", "gocop/testdata/run1.txt", "-r", "gocop/testdata/run3.txt"},
			want: "",
		},
		{
			name: "finds single failed packages",
			args: []string{"run", "main.go", "flaky", "-r", "gocop/testdata/run0.txt", "-r", "gocop/testdata/run1.txt", "-r", "gocop/testdata/run2.txt", "-r", "gocop/testdata/run3.txt"},
			want: "github.com/digitalocean/gocop/sample/flaky",
		},
	}

	for _, tt := range tests {
		o.Spec(tt.name, func(expect expect.Expectation) {
			got, err := exec.Command("go", tt.args...).Output()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\n\n", got)
			expect(string(got)).To(Equal(tt.want))
		})
	}
}
