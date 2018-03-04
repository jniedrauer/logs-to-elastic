// Lambda handler for Cloudwatch Log events
package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/net"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

func ElbHandler(event *events.S3Event) (Response, error) {
	log.Debug("got event: ", event)

	cfg := conf.NewConfig()

	p := parsers.Elb{Records: event.Records, Config: cfg}

	log.Debug("encoding and transmitting logs")
	oks := net.LogstashConsumer(p.GetChunks(), cfg)

	return NewResponse(int(oks), p.LineCount)
}
