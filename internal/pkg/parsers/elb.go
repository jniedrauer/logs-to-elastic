package parsers

import (
	"bytes"
	"encoding/csv"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/awsapi"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"

	log "github.com/sirupsen/logrus"
)

var S3Client = awsapi.S3Client{s3manager.NewDownloader(awsapi.Session)}

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
	DomainName   string `json:"domain_name"`
	Url          string `json:"url"`
	Agent        string `json:"user_agent"`
	Cipher       string `json:"ssl_cipher"`
	Protocol     string `json:"ssl_protocol"`
}

// An ELB log parser
type Elb struct {
	Records   []events.S3EventRecord
	Config    *conf.Config
	LineCount int
}

func (e *Elb) GetChunks() <-chan *EncodedChunk {
	var wg sync.WaitGroup
	out := make(chan *EncodedChunk, 10) // Buffer up to 10 records before transmitting

	for _, r := range e.Records {
		err := e.ParseRecord(&r, &wg, out)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Handle a single record
func (e *Elb) ParseRecord(record *events.S3EventRecord, wg *sync.WaitGroup, out chan<- *EncodedChunk) error {
	// Get the log file from S3
	rOffset := int64(0) // Reader offset for file
	fileName, err := S3Client.Get(record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		os.Remove(fileName)
		return err
	}

	// Count number of lines in the file
	lc, err := LineCount(fileName)
	if err != nil {
		log.Error(err.Error())
	}
	e.LineCount += lc
	log.Debug("found lines: ", e.LineCount)

	Chunk(e.Config.ChunkSize, e.LineCount, func(start int, end int) {
		wg.Add(1)
		go func(start int, end int, fileName string, rOffset *int64) {
			data, err := e.GetChunk(start, end, fileName, rOffset)
			if err != nil {
				wg.Done()
				log.Error(err.Error())
				return
			}

			payload, err := GetEncodedChunk(data, e.Config.Delimiter)
			if err != nil {
				wg.Done()
				log.Error(err.Error())
				return
			}

			out <- &EncodedChunk{
				Payload: payload,
				Records: uint32(end - start),
			}

			wg.Done()
		}(start, end, fileName, &rOffset)
	})

	go func() {
		wg.Wait()
		os.Remove(fileName)
	}()

	return err
}

// Return a slice of logs with logstash keys
func (e *Elb) GetChunk(start int, end int, fileName string, rOffset *int64) ([]interface{}, error) {
	lc := int(end - start)

	lines, offset, err := GetLines(atomic.LoadInt64(rOffset), lc, fileName)
	if err != nil {
		return make([]interface{}, 0), err
	}
	atomic.AddInt64(rOffset, offset)

	l := make([]interface{}, lc)
	for i, v := range lines {
		split, err := splitRecord(v)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		log.Debug("split line: ", split)
		if len(split) < 15 {
			log.Error("line less than 15 elements long")
			continue
		}
		request := strings.Fields(split[11])
		method := request[0]
		u, _ := url.Parse(request[1])
		domain := u.Host
		url := u.Path

		l[i] = ElbLog{
			IndexName:    e.Config.IndexName,
			Timestamp:    split[0],
			Message:      strings.SplitN(string(v[:]), " ", 2)[1],
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
			DomainName:   domain,
			Url:          url,
			Agent:        split[12],
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
