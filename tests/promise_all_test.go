package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trungnn/sialia/promise"
	"math"
	"testing"
	"time"
)

func genSlowAdderPromises(num int) []*promise.Promise {
	var ps []*promise.Promise

	for i := 0; i < num; i++ {
		x := i
		ps = append(ps, promise.New(func() (interface{}, error) {
			return SlowAdder(x, x)
		}))
	}

	return ps
}

func TestPromiseAll(t *testing.T) {
	res, err := promise.All(promise.PromiseAllOpts{
		Promises: genSlowAdderPromises(4),
	}).Await()

	assert.NoError(t, err)
	assert.EqualValues(t, res, promise.ResList{0, 2, 4, 6})
}

func TestPromiseAllWithErr(t *testing.T) {
	var promises []*promise.Promise
	for i := range names {
		name := names[i]
		promises = append(promises, promise.New(func() (interface{}, error) {
			return getScore(name)
		}))
	}

	res, err := promise.All(promise.PromiseAllOpts{
		Promises: promises,
	}).Await()

	assert.Nil(t, res)
	assert.EqualError(t, err, "bob not found")
}

func TestPromiseWaitAllSettled(t *testing.T) {
	var promises []*promise.Promise
	for i := range names {
		name := names[i]
		promises = append(promises, promise.New(func() (interface{}, error) {
			return getScore(name)
		}))
	}

	_, err := promise.All(promise.PromiseAllOpts{
		Promises:       promises,
		WaitAllSettled: true,
	}).Await()

	assert.Nil(t, err)
}

func TestPromiseAllEmpty(t *testing.T) {
	res, err := promise.All(promise.PromiseAllOpts{
		Promises: []*promise.Promise{},
	}).Await()

	assert.NoError(t, err)
	assert.EqualValues(t, res, promise.ResList{})

	res, err = promise.All(promise.PromiseAllOpts{
		Promises:       []*promise.Promise{},
		WaitAllSettled: true,
	}).Await()

	assert.NoError(t, err)
	assert.EqualValues(t, res, promise.PromiseList{})
}

// We use execution time as the proof of concurrency control
// TODO: think of a better way to test this
func TestPromiseAllConcurrency(t *testing.T) {
	start := time.Now()
	resC0, err := promise.All(promise.PromiseAllOpts{
		Promises: genSlowAdderPromises(4),
	}).Await()
	timeC0 := time.Since(start)
	assert.NoError(t, err)

	start = time.Now()
	resC1, err := promise.All(promise.PromiseAllOpts{
		Promises:       genSlowAdderPromises(4),
		MaxConcurrency: 1,
	}).Await()
	timeC1 := time.Since(start)
	assert.NoError(t, err)
	assert.EqualValues(t, resC0, resC1)

	start = time.Now()
	resC2, err := promise.All(promise.PromiseAllOpts{
		Promises:       genSlowAdderPromises(4),
		MaxConcurrency: 2,
	}).Await()
	timeC2 := time.Since(start)
	assert.NoError(t, err)
	assert.EqualValues(t, resC0, resC2)

	timeDiff1 := math.Round(float64(timeC1) / float64(timeC0))
	timeDiff2 := math.Round(float64(timeC2) / float64(timeC0))
	assert.EqualValues(t, timeDiff1, 4)
	assert.EqualValues(t, timeDiff2, 2)
}
