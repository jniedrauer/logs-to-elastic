package net

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	"github.com/stretchr/testify/assert"
)

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
		c := Client()
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
