package parsers

import (
	"encoding/json"
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

func (c *Cloudwatch) GetChunks() <-chan []byte {
	var wg sync.WaitGroup

	out := make(chan []byte, 10)

	Chunk(c.Config.ChunkSize, len(c.Data.LogEvents), func(start int, end int) {
		wg.Add(1)
		go func() {
			out <- c.GetEncodedChunk(start, end)
			wg.Done()
		}()
	})

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Return an encoded chunk of logs
func (c *Cloudwatch) GetEncodedChunk(start int, end int) []byte {
	var enc []byte

	for _, v := range c.GetChunk(start, end) {
		j, err := json.Marshal(v)
		if err != nil {
			log.Error("failed to encode: %v", v)
			continue
		}

		if len(enc) > 0 {
			enc = append(enc, c.Config.Delimiter...)
		}
		enc = append(enc, j...)
	}

	return enc
}

// Return a slice of logs with logstash keys
func (c *Cloudwatch) GetChunk(start int, end int) []interface{} {
	l := make([]interface{}, end-start)
	for i, v := range c.Data.LogEvents[start:end] {
		l[i] = CloudwatchLog{
			Timestamp: unixToIso8601(v.Timestamp),
			Message:   v.Message,
			LogGroup:  c.Data.LogGroup,
			IndexName: c.Config.IndexName,
		}
	}

	return l
}

// Convert a unix timestamp to ISO 8601 format
func unixToIso8601(unix int64) string {
	return time.Unix(unix, 0).UTC().Format("2006-01-02T15:04:05-0000")
}
