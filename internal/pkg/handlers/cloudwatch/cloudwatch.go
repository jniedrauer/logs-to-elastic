package cloudwatch

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(event events.CloudwatchLogsEvent) (Response, error) {
	log.Debug("Got event: %s", event.AWSLogs)

	// TODO: Call to another function for actual logic

	return Response{
		Message: fmt.Sprintf("did a thing"), // TODO: Do something here
		Ok:      true,
	}, nil
}
