package parsers

import (
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	log "github.com/sirupsen/logrus"
)

// The format of a Cloudwatch log when sent to Logstash
type CloudwatchLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	LogGroup  string `json:"logGroup"`
	IndexName string `json:"indexname"`
}

// A Cloudwatch log parser
type Cloudwatch struct {
	Data   *events.CloudwatchLogsData
	Config *conf.Config
}

func (c *Cloudwatch) GetChunks() <-chan *EncodedChunk {
	var wg sync.WaitGroup

	out := make(chan *EncodedChunk, 10)

	Chunk(c.Config.ChunkSize, len(c.Data.LogEvents), func(start int, end int) {
		wg.Add(1)
		go func(start int, end int) {
			data, err := c.GetChunk(start, end)
			if err != nil {
				wg.Done()
				log.Error(err.Error())
				return
			}

			payload, err := GetEncodedChunk(data, c.Config.Delimiter)
			if err != nil {
				wg.Done()
				log.Error(err.Error())
				return
			}

			out <- &EncodedChunk{
				Payload: payload,
				Records: uint32(end - start),
			}

			wg.Done()
		}(start, end)
	})

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Return a slice of logs with logstash keys
func (c *Cloudwatch) GetChunk(start int, end int) ([]interface{}, error) {
	l := make([]interface{}, end-start)
	for i, v := range c.Data.LogEvents[start:end] {
		l[i] = CloudwatchLog{
			Timestamp: unixToIso8601(v.Timestamp),
			Message:   v.Message,
			LogGroup:  c.Data.LogGroup,
			IndexName: c.Config.IndexName,
		}
	}

	return l, nil
}

// Convert a unix timestamp to ISO 8601 format
func unixToIso8601(unix int64) string {
	return time.Unix(unix, 0).UTC().Format("2006-01-02T15:04:05-0000")
}
