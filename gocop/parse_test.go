package gocop

import (
	"testing"
	"time"

	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
)

func TestParseFailed(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name  string
		input []byte
		want  []string
	}{
		{
			name: "finds multiple failed packages",
			input: []byte(`
				--- FAIL: TestWillFail (0.00s)
					failing_test.go:11: number does equal eleven
				FAIL
				FAIL	github.com/digitalocean/gocop/sample/fail	0.721s
				--- FAIL: TestMightFail (0.00s)
					flaky_test.go:11: integer is factor of 3
				FAIL
				coverage: 76.4% of statements
				ok  	github.com/digitalocean/gocop/sample/k8s	0.721s
				FAIL	github.com/digitalocean/gocop/sample/flaky	0.488s coverage: 50.0% of statements
				ok  	github.com/digitalocean/gocop/sample/pass	0.250s
			`),
			want: []string{"github.com/digitalocean/gocop/sample/fail", "github.com/digitalocean/gocop/sample/flaky"},
		},
		{
			name: "finds build failed package",
			input: []byte(`
				# github.com/digitalocean/gocop/sample/failbuild [github.com/digitalocean/gocop/sample/failbuild.test]
				sample\failbuild\failbuild.go:3:1: syntax error: non-declaration statement outside function body
				FAIL	github.com/digitalocean/gocop/sample/failbuild [build failed]
				?   	github.com/digitalocean/gocop/sample/numbers	[no test files]
				ok  	github.com/digitalocean/gocop/sample/pass	0.250s
			`),
			want: []string{"github.com/digitalocean/gocop/sample/failbuild"},
		},
		{
			name: "finds build failed package w/underscore",
			input: []byte(`
				# github.com/digitalocean/gocop/sample/failbuild [github.com/digitalocean/gocop/sample/failbuild.test]
				sample\failbuild\failbuild.go:3:1: syntax error: non-declaration statement outside function body
				FAIL	github.com/digitalocean/gocop/sample/fail_build [build failed]
				?   	github.com/digitalocean/gocop/sample/numbers	[no test files]
				ok  	github.com/digitalocean/gocop/sample/pass	0.250s
			`),
			want: []string{"github.com/digitalocean/gocop/sample/fail_build"},
		},
		{
			name: "finds build failed package w/0-9",
			input: []byte(`
				# github.com/digitalocean/gocop/sample/failbuild [github.com/digitalocean/gocop/sample/failbuild.test]
				sample\failbuild\failbuild.go:3:1: syntax error: non-declaration statement outside function body
				FAIL	github.com/digitalocean/gocop/sample/k8s [build failed]
				?   	github.com/digitalocean/gocop/sample/numbers	[no test files]
				ok  	github.com/digitalocean/gocop/sample/pass	0.250s coverage: 50.0% of statements
			`),
			want: []string{"github.com/digitalocean/gocop/sample/k8s"},
		},
		{
			name: "finds build failed package w/hyphen",
			input: []byte(`
				# github.com/digitalocean/gocop/sample/failbuild [github.com/digitalocean/gocop/sample/failbuild.test]
				sample\failbuild\failbuild.go:3:1: syntax error: non-declaration statement outside function body
				FAIL	github.com/digital-ocean/gocop/sample/hyphen [build failed]
				?   	github.com/digitalocean/gocop/sample/numbers	[no test files]
				ok  	github.com/digitalocean/gocop/sample/pass	0.250s coverage: 50.0% of statements
			`),
			want: []string{"github.com/digital-ocean/gocop/sample/hyphen"},
		},
	}

	for _, tt := range tests {
		o.Spec(tt.name, func(expect expect.Expectation) {
			got, err := ParseFailedPackages(tt.input)
			expect(err).To(matchers.BeNil())
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}

func TestParse(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	tests := []struct {
		name  string
		input []byte
		want  []TestResult
	}{
		{
			name: "finds multiple failed packages",
			input: []byte(`
				--- FAIL: TestWillFail (0.00s)
					failing_test.go:16: number does equal eleven
				FAIL
				FAIL	github.com/digitalocean/gocop/sample/fail	0.600s
				--- FAIL: TestMightFail (0.00s)
					flaky_test.go:16: integer is factor of 3
				FAIL
				FAIL	github.com/digitalocean/gocop/sample/flaky	1.685s
				ok  	github.com/digitalocean/gocop/sample/pass	1.129s coverage: 50.0% of statements
			`),
			want: []TestResult{
				// {"FAIL", "github.com/digitalocean/gocop/sample/fail", "0.600s", ""},
				{
					Result:   "fail",
					Package:  "github.com/digitalocean/gocop/sample/fail",
					Duration: time.Millisecond * 600,
				},
				// {"FAIL", "github.com/digitalocean/gocop/sample/flaky", "1.685s", ""},
				{
					Result:   "fail",
					Package:  "github.com/digitalocean/gocop/sample/flaky",
					Duration: time.Millisecond * 1685,
				},
				// {"ok", "github.com/digitalocean/gocop/sample/pass", "1.129s", "50.0"},
				{
					Result:   "pass",
					Package:  "github.com/digitalocean/gocop/sample/pass",
					Duration: time.Millisecond * 1129,
					Coverage: 0.5,
				},
			},
		},
	}

	for _, tt := range tests {
		o.Spec(tt.name, func(expect expect.Expectation) {
			got, err := Parse(tt.input)
			expect(err).To(matchers.BeNil())
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}
