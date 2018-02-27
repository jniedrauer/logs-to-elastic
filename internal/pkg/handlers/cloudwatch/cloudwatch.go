package cloudwatch

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
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

	logs := parsers.CloudwatchLogs{}
	logs.ParseEvent(event.AWSLogs, cfg.IndexName)

	for i := 0; i < len(logs.Events); i += cfg.ChunkSize {
		end := i + cfg.ChunkSize

		if end > len(logs.Events) {
			end = len(logs.Events)
		}

		s := make([]interface{}, i-end)
		for idx, val := range logs.Events[i:end] {
			s[idx] = val
		}
		payload, errs := parsers.PayloadEncode(s, "\n")
		log.Debug(payload)

		c := output.GetClient()

		err := output.Post(cfg.Logstash, payload, c)
		if err != nil {
			log.Error(err)
		}

		for _, err := range errs {
			if err != nil {
				log.Error(err)
			}
		}
	}

	return Response{
		Message: fmt.Sprintf("sent %d records", len(logs.Events)),
		Ok:      true,
	}, nil
}
