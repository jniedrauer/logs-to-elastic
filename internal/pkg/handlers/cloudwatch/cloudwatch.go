package cloudwatch

import (
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/chunk"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/output"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(event events.CloudwatchLogsEvent) (Response, error) {
	log.Debug("Got event: %v", event)

	cfg := conf.Init()
	c := output.GetClient()

	d, err := event.AWSLogs.Parse()
	if err != nil {
		log.Fatalf("failed to parse event: %v", err)
	}
	stream := &parsers.Cloudwatch{Event: d, IndexName: cfg.IndexName}

	var errs []error

	chunk.Chunk(cfg.ChunkSize, len(stream.Event.LogEvents), func(idx int, end int) {
		payload := parsers.SliceEncode(stream, idx, end, "\n")
		err := output.Post(cfg.Logstash, payload, c)
		if err != nil {
			log.Error(err)
			log.Error("failed to transmit batch: %s", string(payload))
			errs = append(errs, err)
		}
	})

	if len(errs) > 0 {
		return Response{
			Message: fmt.Sprintf("sent %d records", len(stream.Event.LogEvents)-len(errs)),
			Ok:      false,
		}, errors.New(fmt.Sprintf("failed to send %d batches", len(errs)))
	} else {
		return Response{
			Message: fmt.Sprintf("sent %d records", len(stream.Event.LogEvents)),
			Ok:      true,
		}, nil
	}
}
