/*
TODO: Documentation
*/
package main

import (
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/handlers/cloudwatch"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/logging"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	logging.Init()
	lambda.Start(cloudwatch.Handler)
}
