// +build sample

package fail

import (
	"testing"

	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestWillFail(t *testing.T) {
	if myNumber() == 11 {
		t.Fatal("number does equal eleven")
	}
}

func TestWillFailComplex(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{
			name:  "my simple test",
			input: 10,
		},
		{
			name:  "my simple test2",
			input: 12,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := myRandomNumber(8)
			Expect(t, out).To(Equal(out))
			Expect(t, out).To(Equal(tc.input))
		})
	}
}
