package test

import (
	"errors"
	"time"
)

var NegativeAdderErr = errors.New("can't handle negative number")

func SlowAdder(x, y int) (int, error) {
	time.Sleep(10 * time.Millisecond)

	if x < 0 || y < 0 {
		return 0, NegativeAdderErr
	}

	return x + y, nil
}
