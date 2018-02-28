package parsers

import (
	"bytes"
	"encoding/csv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/aws"
	log "github.com/sirupsen/logrus"
)

// An ELB log parser
type Elb struct {
	Record    *events.S3EventRecord
	Event     [][]string
	IndexName string
}

// An ELB log
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

// Download records from S3
func (p *Elb) Download() error {
	s, err := aws.GetSession(cfg.Region)
	if err != nil {
		log.Fatalf("failed getting session: %v", err)
	}

	buff := awssdk.NewWriteAtBuffer([]byte{})
	s3api := s3manager.NewDownloader(s)

	bytes, err := s3api.Download(buff,
		&s3.GetObjectInput{
			Bucket: &p.Record.S3.Bucket.Name,
			Key:    &p.Record.S3.Object.Key,
		})
	log.Debug("got from S3: %d bytes", bytes)
	if err != nil {
		return err
	}

	split, err := splitRecord(buff.Bytes())
	p.Event = split

	return err
}

// Split an ELB log into its respective parts
func splitRecord(in []byte) ([][]string, error) {
	r := csv.NewReader(bytes.NewReader(in))
	r.Comma = ' '
	return r.ReadAll()
}

// Return a slice of logs with all necessary information included
func (p *Elb) GetSlice(idx int, end int) []interface{} {
	l := make([]interface{}, end-idx)

	for i, e := range p.Event[idx:end] {
		rs := strings.Fields(e[11])
		method := rs[0]
		url := rs[1][strings.LastIndex(rs[1], "/")-1 : len(rs[1])]
		agent := rs[2]

		l[i] = ElbLog{
			Timestamp:    e[0],
			Message:      strings.Join(e, " "),
			IndexName:    p.IndexName,
			Name:         e[1],
			ClientIp:     strings.Split(e[2], ":")[0],
			BackendIp:    strings.Split(e[3], ":")[0],
			RequestTime:  e[4],
			BackendTime:  e[5],
			ResponseTime: e[6],
			Code:         e[7],
			BackendCode:  e[8],
			Recieved:     e[9],
			Sent:         e[10],
			Method:       method,
			Url:          url,
			Agent:        agent,
			Cipher:       e[13],
			Protocol:     e[14],
		}
	}
	return l
}
