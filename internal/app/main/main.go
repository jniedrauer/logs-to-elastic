/*
TODO: Documentation
*/
package main

import (
	"fmt"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/config"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/logging"

	"github.com/aws/aws-lambda-go/lambda"

	log "github.com/sirupsen/logrus"
)

type Request struct {
	ID    float64 `json:"id"`
	Value string  `json:"value"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(request Request) (Response, error) {
	log.Debug("Got event: %s", request.Value)

	cfg := config.Configuration{}
	cfg.LoadConfig()

	return Response{
		Message: fmt.Sprintf("Processed request ID %f", request.ID),
		Ok:      true,
	}, nil
}

func main() {
	logging.Init()
	lambda.Start(Handler)
}
