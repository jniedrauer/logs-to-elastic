package output

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientNew(t *testing.T) {
	// Get a new client

	// Clear previous values
	blank := &http.Client{}
	client = blank
	once = sync.Once{}

	c := GetClient()

	assert.IsType(t, c, &http.Client{})
	assert.NotEqual(t, c, blank)
}

func TestGetClientReuse(t *testing.T) {
	// Re-use client

	// Clear previous values
	client = &http.Client{}
	once = sync.Once{}

	c1 := GetClient()
	c2 := GetClient()

	assert.Equal(t, c1, c2)
}

func TestPost(t *testing.T) {
	tests := []struct {
		payload []byte
		code    int
		err     error
	}{
		{
			payload: []byte("{\"key\":\"value\"}"),
			code:    200,
			err:     nil,
		},
		{
			payload: []byte("{\"key\":\"value\"}"),
			code:    500,
			err:     errors.New(""),
		},
	}

	for _, test := range tests {
		// Use a new singleton for each test
		client = &http.Client{}
		once = sync.Once{}
		c := GetClient()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.code)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, "{\"rkey\":\"rvalue\"}")
		}))
		defer ts.Close()

		err := Post(ts.URL, test.payload, c)
		assert.IsType(t, test.err, err)
	}
}
