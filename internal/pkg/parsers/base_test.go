package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadEncode(t *testing.T) {
	tests := []struct {
		idx    int
		end    int
		data   []BaseLog
		delim  string
		expect []byte
	}{
		{
			// Single record encoding
			idx:   0,
			end:   1,
			data:  []BaseLog{{Timestamp: "t", Message: "m", LogGroup: "l", IndexName: "i"}},
			delim: "\n",
			expect: []byte(
				"{\"timestamp\":\"t\",\"message\":\"m\",\"logGroup\":\"l\",\"indexname\":\"i\"}",
			),
		},
		{
			// Multiple record encoding with delimiter
			idx: 0,
			end: 2,
			data: []BaseLog{
				{Timestamp: "t1", Message: "m1", LogGroup: "l1", IndexName: "i1"},
				{Timestamp: "t2", Message: "m2", LogGroup: "l2", IndexName: "i2"},
			},
			delim: "\n",
			expect: []byte(
				"{\"timestamp\":\"t1\",\"message\":\"m1\",\"logGroup\":\"l1\",\"indexname\":\"i1\"}" +
					"\n" +
					"{\"timestamp\":\"t2\",\"message\":\"m2\",\"logGroup\":\"l2\",\"indexname\":\"i2\"}",
			),
		},
		{
			// Slice of possible values
			idx: 1,
			end: 3,
			data: []BaseLog{
				{},
				{Timestamp: "t1", Message: "m1", LogGroup: "l1", IndexName: "i1"},
				{Timestamp: "t2", Message: "m2", LogGroup: "l2", IndexName: "i2"},
				{},
				{},
			},
			delim: "\n",
			expect: []byte(
				"{\"timestamp\":\"t1\",\"message\":\"m1\",\"logGroup\":\"l1\",\"indexname\":\"i1\"}" +
					"\n" +
					"{\"timestamp\":\"t2\",\"message\":\"m2\",\"logGroup\":\"l2\",\"indexname\":\"i2\"}",
			),
		},
	}

	for _, test := range tests {
		m := MockPayload{Data: test.data}
		result := SliceEncode(m, test.idx, test.end, test.delim)
		assert.Equal(t, string(test.expect), string(result))
	}
}

type MockPayload struct {
	Data []BaseLog
}

func (m MockPayload) GetSlice(idx int, end int) []interface{} {
	l := make([]interface{}, end-idx)
	for i, e := range m.Data[idx:end] {
		l[i] = e
	}
	return l
}
