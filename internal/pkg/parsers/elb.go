package parsers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"net/url"
	"os"
	"strings"

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
	Received     string `json:"received_bytes"`
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
	out := make(chan *EncodedChunk, 10) // Buffer up to 10 records before transmitting

	// Put records into channel syncronously but don't block
	go func() {
		for _, r := range e.Records {
			err := e.ParseRecord(&r, out)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}
		close(out)
	}()

	return out
}

// Handle a single record
func (e *Elb) ParseRecord(record *events.S3EventRecord, out chan<- *EncodedChunk) error {
	// Get the log file from S3
	f, err := S3Client.Get(record.S3.Bucket.Name, record.S3.Object.Key)
	defer os.Remove(f)
	if err != nil {
		return err
	}

	fh, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)

	// Scan through the file by lines, transmitting in batches
	var chunk []interface{}
	i := 0
	eof := false
	for {
		if eof {
			break
		}

		if scanner.Scan() {
			chunk = append(chunk, e.ParseRow(scanner.Bytes()))
			i += 1
			e.LineCount += 1
		} else if err := scanner.Err(); err != nil {
			return err
		} else {
			eof = true
		}

		if i >= e.Config.ChunkSize || (eof && len(chunk) > 0) {
			i = 0
			payload, err := GetEncodedChunk(chunk, e.Config.Delimiter)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			out <- &EncodedChunk{
				Payload: payload,
				Records: uint32(len(chunk)),
			}

			chunk = nil
		}
	}

	return err
}

// Return a slice of logs with logstash keys
func (e *Elb) ParseRow(row []byte) interface{} {
	var result interface{}

	split, err := splitRow(row)
	if err != nil {
		log.Error(err.Error())
		return result
	}
	if len(split) < 15 {
		log.Error("line less than 15 elements long")
		return result
	}

	request := strings.Fields(split[11])
	method := request[0]
	u, _ := url.Parse(request[1])
	domain := u.Host
	url := u.Path

	parsed := ElbLog{
		IndexName:    e.Config.IndexName,
		Timestamp:    split[0],
		Message:      strings.SplitN(string(row[:]), " ", 2)[1],
		Name:         split[1],
		ClientIp:     strings.Split(split[2], ":")[0],
		BackendIp:    strings.Split(split[3], ":")[0],
		RequestTime:  split[4],
		BackendTime:  split[5],
		ResponseTime: split[6],
		Code:         split[7],
		BackendCode:  split[8],
		Received:     split[9],
		Sent:         split[10],
		Method:       method,
		DomainName:   domain,
		Url:          url,
		Agent:        split[12],
		Cipher:       split[13],
		Protocol:     split[14],
	}

	return parsed
}

// Split an ELB record into its parts
func splitRow(in []byte) ([]string, error) {
	//log.Debug("splitting line: ", string(in))
	r := csv.NewReader(bytes.NewReader(in))
	r.Comma = ' '
	return r.Read()
}
