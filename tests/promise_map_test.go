package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trungnn/sialia/promise"
)

func TestMap(t *testing.T) {
	pMap := promise.Map(promise.PromiseMapOpts{
		Items: []interface{}{0, 1, 2, 3},
		MapFn: func(i interface{}) *promise.Promise {
			return promise.New(func() (interface{}, error) {
				num := i.(int)
				return SlowAdder(num, num)
			})
		},
	})

	res, err := pMap.Await()

	assert.NoError(t, err)
	assert.EqualValues(t, res, promise.ResList{0, 2, 4, 6})
}

func TestMapErr(t *testing.T) {
	res, err := promise.Map(promise.PromiseMapOpts{
		Items: toInterfaceSlice(names),
		MapFn: func(i interface{}) *promise.Promise {
			name := i.(string)
			return promise.New(func() (interface{}, error) {
				return getScore(name)
			})
		},
	}).Await()

	assert.Nil(t, res)
	assert.EqualError(t, err, "bob not found")
}

func TestMapAllSettled(t *testing.T) {
	res, err := promise.Map(promise.PromiseMapOpts{
		Items: toInterfaceSlice(names),
		MapFn: func(i interface{}) *promise.Promise {
			name := i.(string)
			return promise.New(func() (interface{}, error) {
				return getScore(name)
			})
		},
		WaitAllSettled: true,
	}).Await()

	assert.Nil(t, err)
	for i, p := range res.(promise.PromiseList) {
		assert.True(t, p.IsSettled)

		if i == 1 {
			assert.EqualError(t, p.Err, "bob not found")
		} else {
			assert.Nil(t, p.Err)
			assert.Equal(t, testScores[names[i]], p.Res)
		}
	}
}
