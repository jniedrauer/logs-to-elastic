package parsers

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type CloudwatchLogs struct {
	Events []BaseLog
}

func (c *CloudwatchLogs) ParseEvent(raw events.CloudwatchLogsRawData, indexName string) error {
	d, err := raw.Parse()
	if err != nil {
		return err
	}

	l := make([]BaseLog, len(d.LogEvents))
	for idx, evt := range d.LogEvents {
		l[idx] = BaseLog{
			Timestamp: unixToIso8601(evt.Timestamp),
			Message:   evt.Message,
			LogGroup:  d.LogGroup,
			IndexName: indexName,
		}
	}
	c.Events = l

	return nil
}

func unixToIso8601(unix int64) string {
	return time.Unix(unix, 0).UTC().Format("2006-01-02T15:04:05-0700")
}
