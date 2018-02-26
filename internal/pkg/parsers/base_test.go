package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadEncode(t *testing.T) {
	tests := []struct {
		payload []BaseLogEvent
		delim   string
		expect  []byte
		errs    []error
	}{
		{
			// Single record encoding
			payload: []BaseLogEvent{
				{Timestamp: "t", Message: "m", LogGroup: "l", IndexName: "i"},
			},
			delim: "\n",
			expect: []byte(
				"{\"timestamp\":\"t\",\"message\":\"m\",\"logGroup\":\"l\",\"indexname\":\"i\"}",
			),
			errs: []error{nil},
		},
		{
			// Multiple record encoding with delimiter
			payload: []BaseLogEvent{
				{Timestamp: "t1", Message: "m1", LogGroup: "l1", IndexName: "i1"},
				{Timestamp: "t2", Message: "m2", LogGroup: "l2", IndexName: "i2"},
			},
			delim: "\n",
			expect: []byte(
				"{\"timestamp\":\"t1\",\"message\":\"m1\",\"logGroup\":\"l1\",\"indexname\":\"i1\"}" +
					"\n" +
					"{\"timestamp\":\"t2\",\"message\":\"m2\",\"logGroup\":\"l2\",\"indexname\":\"i2\"}",
			),
			errs: []error{nil},
		},
	}

	for _, test := range tests {
		s := make([]interface{}, len(test.payload))
		for idx, val := range test.payload {
			s[idx] = val
		}
		result, errs := payloadEncode(s, test.delim)
		for idx, err := range errs {
			assert.IsType(t, test.errs[idx], err)
		}
		assert.Equal(t, string(test.expect), string(result))
	}
}
