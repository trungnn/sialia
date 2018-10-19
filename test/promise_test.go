package test

import (
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
