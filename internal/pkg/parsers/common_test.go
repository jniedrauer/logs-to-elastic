package parsers

import (
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
	}

	for _, test := range tests {
		var result [][]int

		Chunk(test.sz, test.len, func(idx int, end int) {
			result = append(result, []int{idx, end})
		})

		assert.Equal(t, test.expect, result)
	}
}

func TestGetEncodedChunk(t *testing.T) {
	tests := []struct {
		data   events.CloudwatchLogsData
		delim  []byte
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

		result := GetEncodedChunk(0, 2, c.Config.Delimiter, c.GetChunk)

		assert.Equal(t, string(test.expect), string(result))
	}
}
