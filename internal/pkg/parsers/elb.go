package parsers

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/awsapi"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"

	log "github.com/sirupsen/logrus"
)

// The format of an ELB log when sent to Logstash
type ElbLog struct {
	Timestamp    string `json:"timestamp"`
	Message      string `json:"message"`
	IndexName    string `json:"indexname"`
	Name         string `json:"elb_name"`
	ClientIp     string `json:"client_ip"`
	BackendIp    string `json:"backend_ip"`
	RequestTime  string `json:"request_processing_time"`
	BackendTime  string `json:"backend_processing_time"`
	ResponseTime string `json:"response_processing_time"`
	Code         string `json:"elb_status_code"`
	BackendCode  string `json:"status_code"`
	Recieved     string `json:"recieved_bytes"`
	Sent         string `json:"sent_bytes"`
	Method       string `json:"method"`
	Url          string `json:"url"`
	Agent        string `json:"user_agent"`
	Cipher       string `json:"ssl_cipher"`
	Protocol     string `json:"ssl_protocol"`
}

// An ELB log parser
type Elb struct {
	Records      []events.S3EventRecord
	Config       *conf.Config
	BufferFile   *os.File
	ReaderOffset int64
	LineCount    int
}

func (e *Elb) GetChunks() <-chan *EncodedChunk {
	e.ReaderOffset = 0 // Start from the beginning of the data
	var err error
	var wg sync.WaitGroup
	out := make(chan *EncodedChunk, 10)

	for _, r := range e.Records {
		e.BufferFile, err = ioutil.TempFile("", "s3logs")
		if err != nil {
			e.BufferFile.Close()
			log.Fatalf(err.Error())
		}

		err = awsapi.GetFromS3(e.BufferFile, &r.S3, e.Config.AwsRegion)

		lc, lcerr := LineCount(e.BufferFile)
		if lcerr != nil {
			log.Error("line count may be incorrect")
		}
		e.LineCount += lc

		Chunk(e.Config.ChunkSize, e.LineCount, func(start int, end int) {
			wg.Add(1)
			go func() {
				log.Debug("encoding chunk from offset: ", e.ReaderOffset)
				payload, perr := GetEncodedChunk(start, end, e.Config.Delimiter, e.GetChunk)
				if perr != nil {
					log.Error(err.Error())
				}
				out <- &EncodedChunk{
					Payload: payload,
					Records: uint32(end - start),
				}
				wg.Done()
			}()
		})

		cerr := e.BufferFile.Close()
		if err == nil {
			err = cerr
		}
		if err != nil {
			log.Fatalf(err.Error())
		}

	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Return a slice of logs with logstash keys
func (e *Elb) GetChunk(start int, end int) ([]interface{}, error) {
	lc := int(end - start)
	lines, offset, err := GetLines(e.ReaderOffset, lc, e.BufferFile)
	if err != nil {
		return make([]interface{}, 0), err
	}
	atomic.AddInt64(&e.ReaderOffset, offset)

	l := make([]interface{}, lc)
	for i, v := range lines {
		split, err := splitRecord(v)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		request := strings.Fields(split[11])
		method := request[0]
		url := request[1][strings.LastIndex(request[1], "/")-1 : len(request[1])]
		agent := request[2]

		l[i] = ElbLog{
			IndexName:    e.Config.IndexName,
			Timestamp:    split[0],
			Message:      strings.Join(split[1:], " "),
			Name:         split[1],
			ClientIp:     strings.Split(split[2], ":")[0],
			BackendIp:    strings.Split(split[3], ":")[0],
			RequestTime:  split[4],
			BackendTime:  split[5],
			ResponseTime: split[6],
			Code:         split[7],
			BackendCode:  split[8],
			Recieved:     split[9],
			Sent:         split[10],
			Method:       method,
			Url:          url,
			Agent:        agent,
			Cipher:       split[13],
			Protocol:     split[14],
		}
	}

	return l, nil
}

// Split an ELB record into its parts
func splitRecord(in []byte) ([]string, error) {
	r := csv.NewReader(bytes.NewReader(in))
	r.Comma = ' '
	return r.Read()
}
