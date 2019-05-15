package numbers

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Eleven returns the int 11
func Eleven() int {
	return 11
}

// RandomInteger returns a random int up to 100
func RandomInteger() int {
	return rand.Intn(100)
}
