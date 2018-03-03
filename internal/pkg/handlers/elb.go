// Lambda handler for Cloudwatch Log events
package elb

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/output"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(event events.S3Event) (Response, error) {
	log.Debug("Got event: %v", event)

	cfg := conf.Init()
	c := output.GetClient()

	errs := make([]error, len(event.Records))

	for i, record := range event.Records {
		errs[i] = recordHandler(cfg, c, &record)
	}

	if len(errs) > 0 {
		return Response{
			Message: fmt.Sprintf("sent %d records", len(event.Records)-len(errs)),
			Ok:      false,
		}, errors.New(fmt.Sprintf("failed to send %d records", len(errs)))
	} else {
		return Response{
			Message: fmt.Sprintf("sent %d records", len(event.Records)),
			Ok:      true,
		}, nil
	}
}

func recordHandler(cfg *conf.Conf, c *http.Client, record *events.S3EventRecord) error {
	stream := &parsers.Elb{Record: record, IndexName: cfg.IndexName}

	var errs []error

	stream.Download()

	chunk.Chunk(cfg.ChunkSize, len(stream.Event), func(idx int, end int) {
		payload := parsers.SliceEncode(stream, idx, end, "\n")
		err := output.Post(cfg.Logstash, payload, c)
		if err != nil {
            log.Error(err)
            log.Error("failed to send batch: %s", string(payload))
            errs = append(errs, err)
        }
		    })
			if len(errs) > 0 {
				return errors.New("failed to send %d batches",  len(errs)
			}
			return nil
}
