package cloudwatch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(event events.CloudwatchLogsEvent) (Response, error) {
	log.Debug("Got event: %v", event)

	logstash := os.Getenv("LOGSTASH")
	indexName := os.Getenv("INDEXNAME")
	c64, err := strconv.ParseInt(os.Getenv("CHUNK_SIZE"), 10, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	chunkSize := int(c64)

	logs := parsers.CloudwatchLogs{}
	logs.ParseEvent(event.AWSLogs, indexName)

	for i := 0; i < len(logs.Events); i += chunkSize {
		end := i + chunkSize

		if end > len(logs.Events) {
			end = len(logs.Events)
		}

		s := make([]interface{}, i-end)
		for idx, val := range logs.Events[i:end] {
			s[idx] = val
		}
		payload, errs := parsers.PayloadEncode(s, "\n")
		log.Debug(payload)

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Post(logstash, "text/json", bytes.NewReader(payload))
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Debug("response: %v", body)

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
