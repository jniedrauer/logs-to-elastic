/*
An AWS Lambda function for shipping ELB logs from S3 to an
HTTP listener.

Required environment variables:
  CHUNK_SIZE:
    Number of logs to transmit in a single request.
  INDEXNAME:
	This will be passed as a json parameter `indexname`, to be used by the
	Logstash listener.
  LOGSTASH:
    HTTP endpoint to send logs to.

Optional environment variables:
  LOG_LEVEL:
    Minimum log level for this function's logging.
*/
package main

import (
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/handlers"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/logging"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	logging.Init()
	lambda.Start(handlers.CloudwatchHandler)
}
