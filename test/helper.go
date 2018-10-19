package test

import (
	"errors"
	"fmt"
	"time"
)

var NegativeAdderErr = errors.New("can't handle negative number")

func SlowAdder(x, y int) (int, error) {
	time.Sleep(100 * time.Millisecond)

	if x < 0 || y < 0 {
		return 0, NegativeAdderErr
	}

	return x + y, nil
}

var names = []string{"alice", "bob", "charlie"}
var testScores = map[string]int{
	"alice":   80,
	"charlie": 75,
}

func getScore(name string) (int, error) {
	time.Sleep(10 * time.Millisecond)

	if content, ok := testScores[name]; !ok {
		return 0, errors.New(fmt.Sprintf("%s not found", name))
	} else {
		return content, nil
	}
}
