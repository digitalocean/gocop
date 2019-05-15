package flaky

import (
	"testing"

	"github.com/digitalocean/gocop/sample/numbers"
)

func TestMightFail(t *testing.T) {
	if numbers.RandomInteger()%3 == 0 {
		t.Error("integer is factor of 3")
	}
}
