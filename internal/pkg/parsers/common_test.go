package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		len    int
		sz     int
		expect [][]int
	}{
		{
			len:    4,
			sz:     2,
			expect: [][]int{{0, 2}, {2, 4}},
		},
		{
			len:    2,
			sz:     1,
			expect: [][]int{{0, 1}, {1, 2}},
		},
		{
			len:    5,
			sz:     2,
			expect: [][]int{{0, 2}, {2, 4}, {4, 5}},
		},
		{
			len:    8,
			sz:     5,
			expect: [][]int{{0, 5}, {5, 8}},
		},
		{
			len:    10,
			sz:     3,
			expect: [][]int{{0, 3}, {3, 6}, {6, 9}, {9, 10}},
		},
	}

	for _, test := range tests {
		var result [][]int

		Chunk(test.sz, test.len, func(idx int, end int) {
			result = append(result, []int{idx, end})
		})

		assert.Equal(t, test.expect, result)
		// All chunks except possibly the last should be the chunk size
		for i, v := range test.expect {
			if i+1 < len(test.expect) {
				assert.Equal(t, v[1]-v[0], test.sz)
			}
		}
	}
}

type mockEncode struct {
	TestString string `json:"test_string"`
	TestInt    int    `json:"test_int"`
}

func TestGetEncodedChunk(t *testing.T) {
	tests := []struct {
		data   []mockEncode
		delim  []byte
		err    error
		expect []byte
	}{
		// Single record
		{
			data:   []mockEncode{mockEncode{TestString: "ts", TestInt: 1}},
			delim:  []byte("\n"),
			err:    nil,
			expect: []byte("{\"test_string\":\"ts\",\"test_int\":1}"),
		},
		// Multirecord, newline delimiter
		{
			data: []mockEncode{
				mockEncode{TestString: "ts1", TestInt: 1},
				mockEncode{TestString: "ts2", TestInt: 2}},
			delim:  []byte("\n"),
			err:    nil,
			expect: []byte("{\"test_string\":\"ts1\",\"test_int\":1}\n{\"test_string\":\"ts2\",\"test_int\":2}"),
		},
		// Multirecord, comma delimiter
		{
			data: []mockEncode{
				mockEncode{TestString: "ts1", TestInt: 1},
				mockEncode{TestString: "ts2", TestInt: 2}},
			delim:  []byte(","),
			err:    nil,
			expect: []byte("{\"test_string\":\"ts1\",\"test_int\":1},{\"test_string\":\"ts2\",\"test_int\":2}"),
		},
	}

	for _, test := range tests {
		// Cast data as slice of interface[]
		data := make([]interface{}, len(test.data))
		for i, v := range test.data {
			data[i] = v
		}

		result, err := GetEncodedChunk(data, test.delim)

		assert.IsType(t, test.err, err)
		assert.Equal(t, string(test.expect), string(result))
	}
}
