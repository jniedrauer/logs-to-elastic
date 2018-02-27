package parsers

import (
	"github.com/aws/aws-lambda-go/events"
)

type Cloudwatch struct {
	Event     events.CloudwatchLogsData
	IndexName string
}

func (c *Cloudwatch) GetSlice(idx int, end int) []interface{} {
	l := make([]interface{}, end-idx)
	for i, e := range c.Event.LogEvents[idx:end] {
		l[i] = BaseLog{
			Timestamp: unixToIso8601(e.Timestamp),
			Message:   e.Message,
			LogGroup:  c.Event.LogGroup,
			IndexName: c.IndexName,
		}
	}
	return l
}
