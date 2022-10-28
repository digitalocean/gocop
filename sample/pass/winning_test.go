package pass

import (
	"testing"

	"github.com/digitalocean/gocop/sample/numbers"
)

func TestWillPass(t *testing.T) {
	err := Winning()
	if err != nil {
		t.Error("unexpected error while winning")
	}
	if numbers.Eleven() != 11 {
		t.Error("number does not equal eleven")
	}
}
