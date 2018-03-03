package parsers

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/stretchr/testify/assert"
)

func TestGetChunk(t *testing.T) {
	tests := []struct {
		start  int
		end    int
		data   events.CloudwatchLogsData
		index  string
		expect []CloudwatchLog
	}{
		// Single event
		{
			start: 0,
			end:   1,
			index: "i",
			data: events.CloudwatchLogsData{
				LogGroup:  "g",
				LogEvents: []events.CloudwatchLogsLogEvent{{Timestamp: 0, Message: "m"}},
			},
			expect: []CloudwatchLog{
				{Timestamp: "1970-01-01T00:00:00-0000", Message: "m", LogGroup: "g", IndexName: "i"},
			},
		},
		// Multiple event slice
		{
			start: 1,
			end:   3,
			index: "i",
			data: events.CloudwatchLogsData{
				LogGroup: "g",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{Timestamp: 0, Message: "m0"},
					{Timestamp: 0, Message: "m1"},
					{Timestamp: 1519700693, Message: "m2"},
					{Timestamp: 0, Message: "m3"},
				},
			},
			expect: []CloudwatchLog{
				{Timestamp: "1970-01-01T00:00:00-0000", Message: "m1", LogGroup: "g", IndexName: "i"},
				{Timestamp: "2018-02-27T03:04:53-0000", Message: "m2", LogGroup: "g", IndexName: "i"},
			},
		},
	}

	for _, test := range tests {
		expect := make([]interface{}, len(test.expect))
		for i, v := range test.expect {
			expect[i] = v
		}
		c := Cloudwatch{Data: &test.data, Config: &conf.Config{IndexName: test.index}}

		result := c.GetChunk(test.start, test.end)

		assert.Equal(t, expect, result)
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

		result := c.GetEncodedChunk(0, 2)

		assert.Equal(t, string(test.expect), string(result))
	}
}
