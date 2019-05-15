package gocop

import (
	"testing"

	"github.com/apoydence/onpar"
	"github.com/apoydence/onpar/expect"
	"github.com/apoydence/onpar/matchers"
)

func TestParse(t *testing.T) {
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
					failing_test.go:16: number does equal eleven
				FAIL
				FAIL	do/teams/cicd/fail	0.600s
				--- FAIL: TestMightFail (0.00s)
					flaky_test.go:16: integer is factor of 3
				FAIL
				FAIL	do/teams/cicd/flaky	1.685s
				ok  	do/teams/cicd/pass	1.129s
			`),
			want: []string{"do/teams/cicd/fail", "do/teams/cicd/flaky"},
		},
	}

	for _, tt := range tests {
		o.Spec(tt.name, func(expect expect.Expectation) {
			got := Parse(tt.input)
			expect(got).To(matchers.Equal(tt.want))
		})
	}
}
