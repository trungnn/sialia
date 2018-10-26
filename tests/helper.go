package test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

func toInterfaceSlice(orgSlice interface{}) []interface{} {
	v := reflect.ValueOf(orgSlice)
	if kind := v.Kind(); kind != reflect.Slice {
		panic(fmt.Sprintf("%s is not a slice", kind))
	}

	res := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		res[i] = v.Index(i).Interface()
	}

	return res
}

type httpErr struct {
	statusCode int
}

func (herr *httpErr) Error() string {
	return fmt.Sprintf("failed with code: %d", herr.statusCode)
}

func fetch(url string) (resp interface{}, err error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, &httpErr{statusCode: res.StatusCode}
	} else {
		bytes, err := ioutil.ReadAll(res.Body)
		return string(bytes), err
	}
}
