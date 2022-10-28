package gocop

import (
	"testing"
	"time"

	"github.com/digitalocean/gocop/gocop/storer"
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

	o.Group("standard", func() {
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
			tt := tt
			o.Spec(tt.name, func(expect expect.Expectation) {
				got, err := ParseFailedPackages(&StandardParser{}, tt.input)
				expect(err).To(matchers.BeNil())
				expect(got).To(matchers.Equal(tt.want))
			})
		}
	})

	o.Group("test2json", func() {
		tests := []struct {
			name  string
			input []byte
			want  []string
		}{
			{
				name: "finds multiple failed packages",
				input: []byte(`
{"Time":"2022-10-28T08:25:04.156155695-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail"}
{"Time":"2022-10-28T08:25:04.156225976-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"=== RUN   TestWillFail\n"}
{"Time":"2022-10-28T08:25:04.156231824-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"    failing_test.go:13: number does equal eleven\n"}
{"Time":"2022-10-28T08:25:04.156237413-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"--- FAIL: TestWillFail (0.00s)\n"}
{"Time":"2022-10-28T08:25:04.15624131-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Elapsed":0}
{"Time":"2022-10-28T08:25:04.156246344-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\n"}
{"Time":"2022-10-28T08:25:04.156267072-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\tgithub.com/digitalocean/gocop/sample/fail\t0.003s\n"}
{"Time":"2022-10-28T08:25:04.15627312-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Elapsed":0.003}
{"Time":"2022-10-28T08:25:04.157186356-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail"}
{"Time":"2022-10-28T08:25:04.15720301-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"=== RUN   TestMightFail\n"}
{"Time":"2022-10-28T08:25:04.157208821-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"    flaky_test.go:13: integer is factor of 3\n"}
{"Time":"2022-10-28T08:25:04.157218117-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"--- FAIL: TestMightFail (0.00s)\n"}
{"Time":"2022-10-28T08:25:04.157221681-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Elapsed":0}
{"Time":"2022-10-28T08:25:04.157226381-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"FAIL\n"}
{"Time":"2022-10-28T08:25:04.157469606-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"FAIL\tgithub.com/digitalocean/gocop/sample/flaky\t0.003s\n"}
{"Time":"2022-10-28T08:25:04.157500533-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/flaky","Elapsed":0.003}
			`),
				want: []string{"github.com/digitalocean/gocop/sample/fail", "github.com/digitalocean/gocop/sample/flaky"},
			},
			{
				name: "finds single build failed",
				input: []byte(`
{"Time":"2022-10-28T08:25:04.156155695-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail"}
{"Time":"2022-10-28T08:25:04.156225976-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"=== RUN   TestWillFail\n"}
{"Time":"2022-10-28T08:25:04.156231824-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"    failing_test.go:13: number does equal eleven\n"}
{"Time":"2022-10-28T08:25:04.156237413-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"--- FAIL: TestWillFail (0.00s)\n"}
{"Time":"2022-10-28T08:25:04.15624131-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Elapsed":0}
{"Time":"2022-10-28T08:25:04.156246344-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\n"}
{"Time":"2022-10-28T08:25:04.156267072-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\tgithub.com/digitalocean/gocop/sample/fail\t0.003s\n"}
{"Time":"2022-10-28T08:25:04.15627312-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Elapsed":0.003}
			`),
				want: []string{"github.com/digitalocean/gocop/sample/fail"},
			},
			{
				name: "no failed builds",
				input: []byte(`
{"Time":"2022-10-28T08:25:04.758219281-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail"}
{"Time":"2022-10-28T08:25:04.758241033-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"=== RUN   TestMightFail\n"}
{"Time":"2022-10-28T08:25:04.758251604-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"--- PASS: TestMightFail (0.00s)\n"}
{"Time":"2022-10-28T08:25:04.75825788-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Elapsed":0}
{"Time":"2022-10-28T08:25:04.758263109-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"PASS\n"}
{"Time":"2022-10-28T08:25:04.758511484-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"ok  \tgithub.com/digitalocean/gocop/sample/flaky\t0.001s\n"}
{"Time":"2022-10-28T08:25:04.758523677-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/flaky","Elapsed":0.001}
{"Time":"2022-10-28T08:25:04.758835947-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/numbers","Output":"?   \tgithub.com/digitalocean/gocop/sample/numbers\t[no test files]\n"}
{"Time":"2022-10-28T08:25:04.75887762-07:00","Action":"skip","Package":"github.com/digitalocean/gocop/sample/numbers","Elapsed":0}
{"Time":"2022-10-28T08:25:04.761379566-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass"}
{"Time":"2022-10-28T08:25:04.761393919-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Output":"=== RUN   TestWillPass\n"}
{"Time":"2022-10-28T08:25:04.761404748-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Output":"--- PASS: TestWillPass (0.00s)\n"}
{"Time":"2022-10-28T08:25:04.761410316-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Elapsed":0}
{"Time":"2022-10-28T08:25:04.761415005-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"PASS\n"}
{"Time":"2022-10-28T08:25:04.761747146-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"ok  \tgithub.com/digitalocean/gocop/sample/pass\t0.001s\n"}
{"Time":"2022-10-28T08:25:04.761771584-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/pass","Elapsed":0.001}
			`),
				want: nil,
			},
		}

		for _, tt := range tests {
			tt := tt
			o.Spec(tt.name, func(expect expect.Expectation) {
				got, err := ParseFailedPackages(&Test2JSONParser{}, tt.input)
				expect(err).To(matchers.BeNil())
				expect(got).To(matchers.Equal(tt.want))
			})
		}
	})
}

func TestParse(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) expect.Expectation {
		return expect.New(t)
	})

	o.Group("standard", func() {
		tests := []struct {
			name  string
			input []byte
			want  []storer.TestResult
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
				want: []storer.TestResult{
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
			tt := tt
			o.Spec(tt.name, func(expect expect.Expectation) {
				got, err := (&StandardParser{}).Parse(tt.input)
				expect(err).To(matchers.BeNil())
				expect(got).To(matchers.Equal(tt.want))
			})
		}
	})

	o.Group("test2json", func() {
		tests := []struct {
			name         string
			input        []byte
			includeTests bool
			want         []storer.TestResult
		}{
			{
				name:         "only package level results ",
				includeTests: false,
				input: []byte(`
{"Time":"2022-10-28T09:22:00.845059928-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail"}
{"Time":"2022-10-28T09:22:00.845150587-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"=== RUN   TestWillFail\n"}
{"Time":"2022-10-28T09:22:00.845206144-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"    failing_test.go:13: number does equal eleven\n"}
{"Time":"2022-10-28T09:22:00.84522339-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"--- FAIL: TestWillFail (0.00s)\n"}
{"Time":"2022-10-28T09:22:00.845236167-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Elapsed":1.5}
{"Time":"2022-10-28T09:22:00.84524241-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\n"}
{"Time":"2022-10-28T09:22:00.845251785-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"coverage: [no statements]\n"}
{"Time":"2022-10-28T09:22:00.845749934-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\tgithub.com/digitalocean/gocop/sample/fail\t0.003s\n"}
{"Time":"2022-10-28T09:22:00.845785227-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Elapsed":0.003}
{"Time":"2022-10-28T09:22:00.850030711-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail"}
{"Time":"2022-10-28T09:22:00.850041754-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"=== RUN   TestMightFail\n"}
{"Time":"2022-10-28T09:22:00.850047934-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"--- PASS: TestMightFail (0.00s)\n"}
{"Time":"2022-10-28T09:22:00.850063581-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Elapsed":0}
{"Time":"2022-10-28T09:22:00.850070891-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"PASS\n"}
{"Time":"2022-10-28T09:22:00.850073761-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"coverage: [no statements]\n"}
{"Time":"2022-10-28T09:22:00.85038881-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"ok  \tgithub.com/digitalocean/gocop/sample/flaky\t0.002s\tcoverage: [no statements]\n"}
{"Time":"2022-10-28T09:22:00.850407512-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/flaky","Elapsed":0.002}
{"Time":"2022-10-28T09:22:00.850769051-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/numbers","Output":"?   \tgithub.com/digitalocean/gocop/sample/numbers\t[no test files]\n"}
{"Time":"2022-10-28T09:22:00.85077946-07:00","Action":"skip","Package":"github.com/digitalocean/gocop/sample/numbers","Elapsed":0}
{"Time":"2022-10-28T09:22:00.863618934-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass"}
{"Time":"2022-10-28T09:22:00.863642487-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Output":"=== RUN   TestWillPass\n"}
{"Time":"2022-10-28T09:22:00.863657897-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Output":"--- PASS: TestWillPass (0.00s)\n"}
{"Time":"2022-10-28T09:22:00.863662676-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Elapsed":0}
{"Time":"2022-10-28T09:22:00.86366947-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"PASS\n"}
{"Time":"2022-10-28T09:22:00.863672863-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"coverage: 100.0% of statements\n"}
{"Time":"2022-10-28T09:22:00.863976847-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"ok  \tgithub.com/digitalocean/gocop/sample/pass\t0.001s\tcoverage: 100.0% of statements\n"}
{"Time":"2022-10-28T09:22:00.864000196-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/pass","Elapsed":0.001}
			`),
				want: []storer.TestResult{
					{
						Result:   "fail",
						Package:  "github.com/digitalocean/gocop/sample/fail",
						Duration: time.Millisecond * 3,
					},
					{
						Result:   "pass",
						Package:  "github.com/digitalocean/gocop/sample/flaky",
						Duration: time.Millisecond * 2,
					},
					{
						Result:  "skip",
						Package: "github.com/digitalocean/gocop/sample/numbers",
					},
					{
						Result:   "pass",
						Package:  "github.com/digitalocean/gocop/sample/pass",
						Duration: time.Millisecond * 1,
						Coverage: 1,
					},
				},
			},
			{
				name:         "include test level",
				includeTests: true,
				input: []byte(`
{"Time":"2022-10-28T09:22:00.845059928-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail"}
{"Time":"2022-10-28T09:22:00.845150587-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"=== RUN   TestWillFail\n"}
{"Time":"2022-10-28T09:22:00.845206144-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"    failing_test.go:13: number does equal eleven\n"}
{"Time":"2022-10-28T09:22:00.84522339-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Output":"--- FAIL: TestWillFail (0.00s)\n"}
{"Time":"2022-10-28T09:22:00.845236167-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Test":"TestWillFail","Elapsed":1.5}
{"Time":"2022-10-28T09:22:00.84524241-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\n"}
{"Time":"2022-10-28T09:22:00.845251785-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"coverage: [no statements]\n"}
{"Time":"2022-10-28T09:22:00.845749934-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/fail","Output":"FAIL\tgithub.com/digitalocean/gocop/sample/fail\t0.003s\n"}
{"Time":"2022-10-28T09:22:00.845785227-07:00","Action":"fail","Package":"github.com/digitalocean/gocop/sample/fail","Elapsed":0.003}
{"Time":"2022-10-28T09:22:00.850030711-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail"}
{"Time":"2022-10-28T09:22:00.850041754-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"=== RUN   TestMightFail\n"}
{"Time":"2022-10-28T09:22:00.850047934-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Output":"--- PASS: TestMightFail (0.00s)\n"}
{"Time":"2022-10-28T09:22:00.850063581-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/flaky","Test":"TestMightFail","Elapsed":0}
{"Time":"2022-10-28T09:22:00.850070891-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"PASS\n"}
{"Time":"2022-10-28T09:22:00.850073761-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"coverage: [no statements]\n"}
{"Time":"2022-10-28T09:22:00.85038881-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/flaky","Output":"ok  \tgithub.com/digitalocean/gocop/sample/flaky\t0.002s\tcoverage: [no statements]\n"}
{"Time":"2022-10-28T09:22:00.850407512-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/flaky","Elapsed":0.002}
{"Time":"2022-10-28T09:22:00.850769051-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/numbers","Output":"?   \tgithub.com/digitalocean/gocop/sample/numbers\t[no test files]\n"}
{"Time":"2022-10-28T09:22:00.85077946-07:00","Action":"skip","Package":"github.com/digitalocean/gocop/sample/numbers","Elapsed":0}
{"Time":"2022-10-28T09:22:00.863618934-07:00","Action":"run","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass"}
{"Time":"2022-10-28T09:22:00.863642487-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Output":"=== RUN   TestWillPass\n"}
{"Time":"2022-10-28T09:22:00.863657897-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Output":"--- PASS: TestWillPass (0.00s)\n"}
{"Time":"2022-10-28T09:22:00.863662676-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/pass","Test":"TestWillPass","Elapsed":0}
{"Time":"2022-10-28T09:22:00.86366947-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"PASS\n"}
{"Time":"2022-10-28T09:22:00.863672863-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"coverage: 100.0% of statements\n"}
{"Time":"2022-10-28T09:22:00.863976847-07:00","Action":"output","Package":"github.com/digitalocean/gocop/sample/pass","Output":"ok  \tgithub.com/digitalocean/gocop/sample/pass\t0.001s\tcoverage: 100.0% of statements\n"}
{"Time":"2022-10-28T09:22:00.864000196-07:00","Action":"pass","Package":"github.com/digitalocean/gocop/sample/pass","Elapsed":0.001}
			`),
				want: []storer.TestResult{
					{
						Result:   "fail",
						Package:  "github.com/digitalocean/gocop/sample/fail",
						Duration: time.Millisecond * 3,
					},
					{
						Result:   "fail",
						Package:  "github.com/digitalocean/gocop/sample/fail",
						Test:     "TestWillFail",
						Duration: time.Millisecond * 1500,
					},
					{
						Result:   "pass",
						Package:  "github.com/digitalocean/gocop/sample/flaky",
						Duration: time.Millisecond * 2,
					},
					{
						Result:   "pass",
						Package:  "github.com/digitalocean/gocop/sample/flaky",
						Test:     "TestMightFail",
						Duration: 0,
					},
					{
						Result:  "skip",
						Package: "github.com/digitalocean/gocop/sample/numbers",
					},
					{
						Result:   "pass",
						Package:  "github.com/digitalocean/gocop/sample/pass",
						Duration: time.Millisecond * 1,
						Coverage: 1,
					},
					{
						Result:  "pass",
						Package: "github.com/digitalocean/gocop/sample/pass",
						Test:    "TestWillPass",
					},
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			o.Spec(tt.name, func(expect expect.Expectation) {
				got, err := (&Test2JSONParser{
					IncludeIndividualTests: tt.includeTests,
				}).Parse(tt.input)
				expect(err).To(matchers.BeNil())
				expect(got).To(matchers.Equal(tt.want))
			})
		}
	})
}
