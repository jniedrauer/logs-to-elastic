package parsers

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/awsapi"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	logging.Init()
	rc := m.Run()
	os.Exit(rc)
}

func TestSplitRecord(t *testing.T) {
	tests := []struct {
		data   []byte
		err    error
		expect []string
	}{
		{
			data: []byte("2018-02-19T03:43:32.144378Z LoadBalancer01 123.123.123.123:15982 - -1 -1 -1 503 0 0 0 \"GET http://10.10.50.12:80/ HTTP/1.0\" \"scan/1.0\" - -"),
			expect: []string{
				"2018-02-19T03:43:32.144378Z",
				"LoadBalancer01",
				"123.123.123.123:15982",
				"-",
				"-1",
				"-1",
				"-1",
				"503",
				"0",
				"0",
				"0",
				"GET http://10.10.50.12:80/ HTTP/1.0",
				"scan/1.0",
				"-",
				"-",
			},
		},
	}

	for _, test := range tests {
		result, err := splitRecord(test.data)

		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, result)
	}
}

func TestElbGetChunk(t *testing.T) {
	tests := []struct {
		elb    *Elb
		start  int
		end    int
		offset int64
		data   []byte
		err    error
		expect []ElbLog
	}{
		{
			elb: &Elb{
				Config: &conf.Config{IndexName: "index", Delimiter: []byte("\n"), ChunkSize: 1},
			},
			start:  0,
			end:    1,
			offset: 0,
			data:   []byte("2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000086 0.001048 0.001337 200 200 0 57 \"GET https://www.example.com:443/ HTTP/1.1\" \"curl/7.38.0\" DHE-RSA-AES128-SHA TLSv1.2"),
			expect: []ElbLog{ElbLog{
				Timestamp:    "2015-05-13T23:39:43.945958Z",
				Message:      "my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000086 0.001048 0.001337 200 200 0 57 \"GET https://www.example.com:443/ HTTP/1.1\" \"curl/7.38.0\" DHE-RSA-AES128-SHA TLSv1.2",
				IndexName:    "index",
				Name:         "my-loadbalancer",
				ClientIp:     "192.168.131.39",
				BackendIp:    "10.0.0.1",
				RequestTime:  "0.000086",
				BackendTime:  "0.001048",
				ResponseTime: "0.001337",
				Code:         "200",
				BackendCode:  "200",
				Recieved:     "0",
				Sent:         "57",
				Method:       "GET",
				DomainName:   "www.example.com:443",
				Url:          "/",
				Agent:        "curl/7.38.0",
				Cipher:       "DHE-RSA-AES128-SHA",
				Protocol:     "TLSv1.2",
			}},
		},
	}

	for _, test := range tests {
		e := test.elb
		file, _ := ioutil.TempFile("", ".LogsToElasticTest")
		file.Write(test.data)
		file.Close()

		result, err := e.GetChunk(test.start, test.end, file.Name(), &test.offset)

		// We have to cast the expected result as a slice of interface{}
		expect := make([]interface{}, len(test.expect))
		for i, v := range test.expect {
			expect[i] = v
		}

		assert.IsType(t, test.err, err)
		assert.Equal(t, expect, result)

		os.Remove(file.Name())
	}
}

func TestElbGetChunks(t *testing.T) {
	tests := []struct {
		elb    *Elb
		src    string
		expect int
	}{
		{
			elb: &Elb{
				Records: []events.S3EventRecord{
					events.S3EventRecord{S3: events.S3Entity{
						Bucket: events.S3Bucket{Name: "notreal"},
						Object: events.S3Object{Key: "notreal"},
					}},
				},
				Config: &conf.Config{IndexName: "index", Delimiter: []byte("\n"), ChunkSize: 1},
			},
			src:    "testdata/issue_2_logs",
			expect: 16,
		},
	}

	for _, test := range tests {
		data, _ := ioutil.ReadFile(test.src)
		// This is a global variable. See also:
		// https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-handler-types.html#go-programming-model-handler-execution-environment-reuse
		S3Client = awsapi.S3Client{awsapi.MockS3Iface{TestData: data}}
		i := 0
		for v := range test.elb.GetChunks() {
			assert.True(t, strings.HasPrefix(string(v.Payload), "{\"timestamp\""))
			assert.True(t, strings.HasSuffix(string(v.Payload), "\"ssl_protocol\":\"TLSv1.2\"}"))
			i += 1
		}
		assert.Equal(t, test.expect, i)
	}
}
