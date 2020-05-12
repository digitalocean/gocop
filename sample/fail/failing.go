package fail

import "github.com/digitalocean/gocop/sample/numbers"

func myNumber() int {
	return numbers.Eleven()
}

func myRandomNumber(loop int) int {
	sum := 0

	for i := 0; i < loop; i++ {
		sum += numbers.RandomInteger()
	}

	return sum
}
