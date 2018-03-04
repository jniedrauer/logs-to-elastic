/*
An AWS Lambda function for shipping logs from S3 ELB access Logs events
to an HTTP listener.

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
  DELIMITER:
    Delimiter character to use between records when encoding.
*/
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/handlers"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/logging"
)

func main() {
	logging.Init()
	lambda.Start(handlers.ElbHandler)
}
