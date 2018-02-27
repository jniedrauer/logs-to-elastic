package parsers

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetSlice(t *testing.T) {
	tests := []struct {
		idx    int
		end    int
		event  events.CloudwatchLogsData
		index  string
		expect []BaseLog
	}{
		// Single event
		{
			idx:   0,
			end:   1,
			index: "i",
			event: events.CloudwatchLogsData{
				LogGroup:  "g",
				LogEvents: []events.CloudwatchLogsLogEvent{{Timestamp: 0, Message: "m"}},
			},
			expect: []BaseLog{
				{Timestamp: "1970-01-01T00:00:00-0000", Message: "m", LogGroup: "g", IndexName: "i"},
			},
		},
		// Multiple event slice
		{
			idx:   1,
			end:   3,
			index: "i",
			event: events.CloudwatchLogsData{
				LogGroup: "g",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{Timestamp: 0, Message: "m0"},
					{Timestamp: 0, Message: "m1"},
					{Timestamp: 1519700693, Message: "m2"},
					{Timestamp: 0, Message: "m3"},
				},
			},
			expect: []BaseLog{
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
		c := Cloudwatch{Event: test.event, IndexName: test.index}

		result := c.GetSlice(test.idx, test.end)

		assert.Equal(t, expect, result)
	}
}
