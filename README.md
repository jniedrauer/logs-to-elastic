# logs-to-elastic
AWS Lambda functions for shipping various AWS logs to an Elastic Stack

[![Go Report Card](https://goreportcard.com/badge/github.com/jniedrauer/logs-to-elastic?style=flat-square)](https://goreportcard.com/report/github.com/jniedrauer/logs-to-elastic)
[![Travis CI](https://img.shields.io/travis/jniedrauer/logs-to-elastic.svg?style=flat-square)](https://travis-ci.org/jniedrauer/logs-to-elastic)
[![Release](https://img.shields.io/github/release/jniedrauer/logs-to-elastic/all.svg?style=flat-square)](https://github.com/jniedrauer/logs-to-elastic/releases/latest)

## Building
Building requires the [dep](https://github.com/golang/dep) dependency manager
for Go. You can install it via `go get`:

```
go get -u github.com/golang/dep/cmd/dep
```

Default make target will download dependencies, run unit tests, and compile.

## Use
The binary artifacts are designed for use in AWS Lambda.

### Cloudwatch Logs
Set Cloudwatch log groups as a trigger for a Lambda function with the
artifact `cloudwatch.zip`. Every time the function triggers, the logs will
be split into keys, JSON encoded, and POSTed to the configured endpoint.

### ELB Logs
Set ELBs to ship access logs to S3. Then set the S3 bucket PUT as a trigger
for the Lambda function with a artifact `elb.zip`. The logs will be downloaded
from S3, split into keys, JSON encoded, and POSTed to the configured endpoint.

Note that encoding and transmission will occur in an indeterminate order and
with unlimited concurrency.

### Required environment variables
- `CHUNK_SIZE`

  Number of logs to transmit in a single request.

- `INDEXNAME`

  This will be passed as a json parameter `indexname`, to be used by the
  Logstash listener.

- `LOGSTASH`

  HTTP endpoint to send logs to.

### Optional environment variables
- `LOG_LEVEL`

  Minimum log level for this function's logging.

- `DELIMITER`

  Delimiter character to use between records when encoding.
