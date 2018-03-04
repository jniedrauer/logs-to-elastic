package net

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
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
		endpoint string
		payload  []byte
		code     int
		expect   bool
	}{
		{
			endpoint: "http://testurl",
			payload:  []byte("{\"key\":\"value\"}"),
			code:     200,
			expect:   true,
		},
		{
			endpoint: "localhost",
			payload:  []byte("{\"key\":\"value\"}"),
			code:     500,
			expect:   false,
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

		result := Post(ts.URL, test.payload, c)
		assert.Equal(t, test.expect, result)
	}
}

func TestLogstashConsumer(t *testing.T) {
	tests := []struct {
		input   int
		records uint32
		code    int
		expect  uint32
	}{
		// HTTP 200, expect 3 successes
		{
			input:   3,
			records: 1,
			code:    200,
			expect:  3,
		},
		// HTTP 500, expect 0 successes
		{
			input:   3,
			records: 1,
			code:    500,
			expect:  0,
		},
		// Multiple records per transmission
		{
			input:   10,
			records: 3,
			code:    200,
			expect:  30,
		},
	}

	for _, test := range tests {
		// Fake http endpoint
		client = &http.Client{}
		once = sync.Once{}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.code)
		}))
		defer ts.Close()

		// Fake config
		cfg := conf.Config{Logstash: ts.URL}

		// Fake channel
		out := make(chan *parsers.EncodedChunk, 10)
		go func() {
			for i := 0; i < test.input; i++ {
				out <- &parsers.EncodedChunk{Payload: []byte("f"), Records: test.records}
			}
			close(out)
		}()

		result := LogstashConsumer(out, &cfg)

		assert.Equal(t, int(test.expect), int(result))
	}
}
