// Lambda handler for Cloudwatch Log events
package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/io"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

func CloudwatchHandler(event *events.CloudwatchLogsEvent) (Response, error) {
	log.Debug("got event: %v", event)

	cfg := conf.NewConfig()

	d, err := event.AWSLogs.Parse()
	if err != nil {
		log.Fatalf("failed to decode event")
	}

	p := parsers.Cloudwatch{&d, cfg}

	return NewResponse(io.Consumer(p.GetChunks(), cfg), len(d.LogEvents))
}
