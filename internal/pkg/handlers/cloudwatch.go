// Lambda handler for Cloudwatch Log events
package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/net"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

func CloudwatchHandler(event *events.CloudwatchLogsEvent) (Response, error) {
	log.Debug("got event: ", event)

	cfg := conf.NewConfig()

	log.Debug("decoding event")
	d, err := event.AWSLogs.Parse()
	if err != nil {
		log.Fatalf("failed to decode event")
	}

	p := parsers.Cloudwatch{&d, cfg}

	log.Debug("encoding and transmitting logs")
	oks := net.LogstashConsumer(p.GetChunks(), cfg)

	return NewResponse(int(oks), len(d.LogEvents))
}
