package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trungnn/sialia/promise"
)

func TestCreatePromise(t *testing.T) {
	p := promise.New(func() (interface{}, error) {
		return nil, nil
	})

	assert.NotNil(t, p, "promise should be initialized")
}

func TestPromiseAwait(t *testing.T) {
	p := promise.New(func() (interface{}, error) {
		return SlowAdder(1, 2)
	})
	assert.False(t, p.IsSettled, "promise should not be settled yet")

	res, err := p.Await()
	assert.Nil(t, err)
	assert.Equal(t, res, 3)
}

func TestPromiseThen(t *testing.T) {
	pChain := promise.
		New(func() (interface{}, error) {
			return SlowAdder(1, 1)
		}).
		Then(func(res interface{}) (interface{}, error) {
			return SlowAdder(res.(int), 2)
		})

	res, err := pChain.Await()
	assert.Nil(t, err)
	assert.Equal(t, res, 4)

	pChain = promise.
		New(func() (interface{}, error) {
			return SlowAdder(1, 1)
		}).
		Then(func(res interface{}) (interface{}, error) {
			return SlowAdder(res.(int), 2)
		}).
		Then(func(res interface{}) (interface{}, error) {
			return SlowAdder(res.(int), -1)
		}).
		Then(func(res interface{}) (interface{}, error) {
			return SlowAdder(res.(int), 2)
		})

	res, err = pChain.Await()
	assert.Errorf(t, err, NegativeAdderErr.Error())
}

func TestPromiseCatch(t *testing.T) {
	pChain := promise.New(func() (interface{}, error) {
		return SlowAdder(-1, -1)
	}).Catch(func(err error) (interface{}, error) {
		if err == NegativeAdderErr {
			return nil, errors.New("transformed error")
		}

		return nil, err
	})

	res, err := pChain.Await()
	assert.Nil(t, res)
	assert.Errorf(t, err, "transformed error")

	pChain = promise.
		New(func() (interface{}, error) {
			return SlowAdder(1, 1)
		}).
		Then(func(res interface{}) (interface{}, error) {
			return SlowAdder(res.(int), 2)
		}).
		Catch(func(err error) (interface{}, error) {
			if err == NegativeAdderErr {
				return nil, errors.New("transformed error")
			}

			return nil, err
		})

	res, err = pChain.Await()
	assert.Nil(t, err, "should not reach catch's inner function")
	assert.Equal(t, res, 4)
}
