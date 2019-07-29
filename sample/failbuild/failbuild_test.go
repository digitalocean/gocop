package fail

import (
	"testing"

	"github.com/digitalocean/gocop/sample/numbers"
)

func TestWillFail(t *testing.T) {
	if numbers.Eleven() == 11 {
		t.Error("number does equal eleven")
	}
}
