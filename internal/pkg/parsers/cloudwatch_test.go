package parsers

import (
	"reflect"
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
		err    error
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
			err: nil,
			expect: []CloudwatchLog{
				{Timestamp: "1970-01-01T00:00:00Z", Message: "m", LogGroup: "g", IndexName: "i"},
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
					{Timestamp: 1519700693000, Message: "m2"},
					{Timestamp: 0, Message: "m3"},
				},
			},
			err: nil,
			expect: []CloudwatchLog{
				{Timestamp: "1970-01-01T00:00:00Z", Message: "m1", LogGroup: "g", IndexName: "i"},
				{Timestamp: "2018-02-27T03:04:53Z", Message: "m2", LogGroup: "g", IndexName: "i"},
			},
		},
	}

	for _, test := range tests {
		expect := make([]interface{}, len(test.expect))
		for i, v := range test.expect {
			expect[i] = v
		}
		c := Cloudwatch{Data: &test.data, Config: &conf.Config{IndexName: test.index}}

		result, err := c.GetChunk(test.start, test.end)

		assert.IsType(t, test.err, err)
		assert.Equal(t, expect, result)
	}
}

func TestGetChunks(t *testing.T) {
	tests := []struct {
		data        events.CloudwatchLogsData
		chunkSize   int
		expectCount uint32
		expect      [][]byte
	}{
		// Test output with single record per chunk
		{
			data: events.CloudwatchLogsData{
				LogGroup: "g",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{Timestamp: 0, Message: "m1"},
					{Timestamp: 0, Message: "m2"},
				},
			},
			chunkSize:   1,
			expectCount: 2,
			expect: [][]byte{
				[]byte("{\"timestamp\":\"1970-01-01T00:00:00Z\",\"message\":\"m1\",\"logGroup\":\"g\",\"indexname\":\"index\"}"),
				[]byte("{\"timestamp\":\"1970-01-01T00:00:00Z\",\"message\":\"m2\",\"logGroup\":\"g\",\"indexname\":\"index\"}"),
			},
		},
		// Test mismatched chunks and records to verify that count is correct
		{
			data: events.CloudwatchLogsData{
				LogGroup: "g",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{Timestamp: 0, Message: "m1"},
					{Timestamp: 0, Message: "m2"},
					{Timestamp: 0, Message: "m3"},
				},
			},
			chunkSize:   2,
			expectCount: 3,
		},
	}

	for _, test := range tests {
		config := conf.Config{IndexName: "index", Delimiter: []byte(","), ChunkSize: test.chunkSize}
		c := Cloudwatch{Data: &test.data, Config: &config}

		results := c.GetChunks()
		unmatched := test.expect

		var resultCount uint32

		/* The results come in asyncronously so we have to check all possible
		results against all expected results and remove them as they match.
		The end result should be no unmatched elements from the expected list */
		for result := range results {
			for i, expect := range test.expect {
				if reflect.DeepEqual(expect, result.Payload) {
					// Pop element off the unmatched slice
					unmatched[i] = unmatched[len(unmatched)-1]
					unmatched = unmatched[:len(unmatched)-1]
					break
				}
			}
			resultCount += result.Records
		}
		// Only run this test if we are asserting output
		if len(test.expect) > 0 {
			assert.Equal(t, 0, len(unmatched))
		}
		assert.Equal(t, int(test.expectCount), int(resultCount))
	}
}

func TestUnixToRfc3339(t *testing.T) {
	tests := []struct {
		unix   int64
		expect string
	}{
		{
			unix:   0,
			expect: "1970-01-01T00:00:00Z",
		},
		{
			unix:   1520270504080,
			expect: "2018-03-05T17:21:44.08Z",
		},
	}

	for _, test := range tests {
		result := unixToRfc3339(test.unix)

		assert.Equal(t, test.expect, result)
	}
}
