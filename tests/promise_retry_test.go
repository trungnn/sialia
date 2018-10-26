package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trungnn/sialia/promise"
	"net/http"
	"net/http/httptest"
	"testing"
)

var maxTries = 5

func newTestServer(readyC chan bool) *httptest.Server {
	var serverReady bool
	go func() {
		for {
			select {
			case serverReady = <-readyC:
			}
		}
	}()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !serverReady {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		name := r.URL.Query().Get("name")
		if score, err := getScore(name); err != nil {
			http.NotFound(w, r)
		} else {
			w.Write([]byte(fmt.Sprint(score)))
		}
	}))
}

func TestRetry(t *testing.T) {
	tries := 0
	p := promise.NewWithRetry(func() (interface{}, error) {
		tries++
		return SlowAdder(-1, 0)
	}, &promise.RetryOpts{
		MaxTries: maxTries,
	})

	_, err := p.Await()

	assert.Error(t, err)
	assert.Equal(t, tries, maxTries)
}

func TestRetryWithCondition(t *testing.T) {
	retryOpts := &promise.RetryOpts{
		MaxTries: maxTries,
		RetryCheck: func(err error) bool {
			if httperr, ok := err.(*httpErr); !ok {
				return false
			} else {
				return httperr.statusCode != http.StatusNotFound // don't retry a not found
			}
		},
	}

	readyC := make(chan bool)
	server := newTestServer(readyC)
	defer server.Close()

	tries := 0
	_, err := promise.NewWithRetry(func() (interface{}, error) {
		tries++
		return fetch(fmt.Sprintf("%s/%s", server.URL, "foo?name=alice"))
	}, retryOpts).Await()

	assert.Equal(t, tries, maxTries, "made use of all allowed tries")
	assert.IsType(t, &httpErr{}, err)
	assert.Equal(t, err.(*httpErr).statusCode, http.StatusServiceUnavailable)

	tries = 0
	_, err = promise.NewWithRetry(func() (interface{}, error) {
		if tries++; tries == 3 {
			readyC <- true
		}

		return fetch(fmt.Sprintf("%s/%s", server.URL, "foo?name=bob"))
	}, retryOpts).Await()

	assert.Equal(t, tries, 3, "should abort retry process with the right condition")
	assert.IsType(t, &httpErr{}, err)
	assert.Equal(t, err.(*httpErr).statusCode, http.StatusNotFound)

	aliceScore := fmt.Sprintf("%d", testScores["alice"])
	tries = 0
	res, err := promise.NewWithRetry(func() (interface{}, error) {
		tries++
		return fetch(fmt.Sprintf("%s/%s", server.URL, "foo?name=alice"))
	}, retryOpts).Await()

	assert.NoError(t, err)
	assert.Equal(t, tries, 1, "should not retry on success")
	assert.EqualValues(t, res, aliceScore, "should propagate correct result")

	tries = 0
	readyC <- false
	_, err = promise.NewWithRetry(func() (interface{}, error) {
		if tries++; tries == 3 {
			readyC <- true
		}

		return fetch(fmt.Sprintf("%s/%s", server.URL, "foo?name=alice"))
	}, retryOpts).Await()

	assert.Equal(t, tries, 3, "should abort retry process with the right condition")
	assert.NoError(t, err)
	assert.EqualValues(t, res, aliceScore, "should propagate correct result")
}
