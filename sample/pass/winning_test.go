package pass

import (
	"testing"

	"github.com/digitalocean/gocop/sample/numbers"
)

func TestWillPass(t *testing.T) {
	if numbers.Eleven() != 11 {
		t.Error("number does not equal eleven")
	}
}
