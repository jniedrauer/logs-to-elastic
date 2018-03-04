package parsers

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
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

func TestGetEncodedChunk(t *testing.T) {
	tests := []struct {
		data   events.CloudwatchLogsData
		delim  []byte
		err    error
		expect []byte
	}{
		// newline delimiter
		{
			data: events.CloudwatchLogsData{
				LogGroup: "g",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{Timestamp: 0, Message: "m1"},
					{Timestamp: 0, Message: "m2"},
				},
			},
			delim: []byte("\n"),
			expect: []byte(
				"{\"timestamp\":\"1970-01-01T00:00:00-0000\",\"message\":\"m1\",\"logGroup\":\"g\",\"indexname\":\"index\"}" +
					"\n" +
					"{\"timestamp\":\"1970-01-01T00:00:00-0000\",\"message\":\"m2\",\"logGroup\":\"g\",\"indexname\":\"index\"}",
			),
		},
		// comma delimiter
		{
			data: events.CloudwatchLogsData{
				LogGroup: "g",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{Timestamp: 0, Message: "m1"},
					{Timestamp: 0, Message: "m2"},
				},
			},
			delim: []byte(","),
			expect: []byte(
				"{\"timestamp\":\"1970-01-01T00:00:00-0000\",\"message\":\"m1\",\"logGroup\":\"g\",\"indexname\":\"index\"}" +
					"," +
					"{\"timestamp\":\"1970-01-01T00:00:00-0000\",\"message\":\"m2\",\"logGroup\":\"g\",\"indexname\":\"index\"}",
			),
		},
	}

	for _, test := range tests {
		c := Cloudwatch{Data: &test.data, Config: &conf.Config{IndexName: "index", Delimiter: test.delim}}

		result, err := GetEncodedChunk(0, 2, c.Config.Delimiter, c.GetChunk)

		assert.IsType(t, test.err, err)
		assert.Equal(t, string(test.expect), string(result))
	}
}

func TestLineCount(t *testing.T) {
	tests := []struct {
		data   io.Reader
		err    error
		expect int
	}{
		// Single line
		{
			data:   bytes.NewReader([]byte("test\n")),
			err:    nil,
			expect: 1,
		},
		// No newline
		{
			data:   bytes.NewReader([]byte("test")),
			err:    errors.New(""),
			expect: 1,
		},
		// Multiple lines
		{
			data:   bytes.NewReader([]byte("test1\ntest2\ntest3\n")),
			err:    nil,
			expect: 3,
		},
	}

	for _, test := range tests {
		result, err := LineCount(test.data)

		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, result)
	}
}

func TestGetLines(t *testing.T) {
	tests := []struct {
		start        int64
		lines        int
		data         io.ReadSeeker
		err          error
		expectData   [][]byte
		expectOffset int64
	}{
		// 1 line, no offset
		{
			start:        0,
			lines:        1,
			data:         bytes.NewReader([]byte("test\n")),
			err:          nil,
			expectData:   [][]byte{[]byte("test")},
			expectOffset: int64(len([]byte("test\n"))),
		},
		// 2 lines, no offset
		{
			start:        0,
			lines:        2,
			data:         bytes.NewReader([]byte("test1\ntest2\n")),
			err:          nil,
			expectData:   [][]byte{[]byte("test1"), []byte("test2")},
			expectOffset: int64(len([]byte("test1\ntest2\n"))),
		},
		// Multiline slice, with offset
		{
			start:        6,
			lines:        1,
			data:         bytes.NewReader([]byte("test1\ntest2\ntest2\n")),
			err:          nil,
			expectData:   [][]byte{[]byte("test2")},
			expectOffset: int64(len([]byte("test2\n"))),
		},
		// Ask for more lines than are left before EOF
		{
			start:        0,
			lines:        10,
			data:         bytes.NewReader([]byte("test1\ntest2\n")),
			err:          nil,
			expectData:   [][]byte{[]byte("test1"), []byte("test2")},
			expectOffset: int64(len([]byte("test1\ntest2\n"))),
		},
	}

	for _, test := range tests {
		result, offset, err := GetLines(test.start, test.lines, test.data)

		assert.IsType(t, test.err, err)
		for i, v := range test.expectData {
			assert.Equal(t, string(v), string(result[i]))
		}
		assert.Equal(t, test.expectOffset, offset)
	}
}
