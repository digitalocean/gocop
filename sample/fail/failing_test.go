// +build sample

package fail

import (
	"testing"
)

func TestWillFail(t *testing.T) {
	if myNumber() == 11 {
		t.Fatal("number does equal eleven")
	}
}
