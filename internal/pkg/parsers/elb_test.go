package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/stretchr/testify/assert"
)

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
		data   []byte
		err    error
		expect []ElbLog
	}{
		{
			elb: &Elb{
				Config:       &conf.Config{IndexName: "index", Delimiter: []byte("\n"), ChunkSize: 1},
				ReaderOffset: 0,
			},
			start: 0,
			end:   1,
			data:  []byte("2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000086 0.001048 0.001337 200 200 0 57 \"GET https://www.example.com:443/ HTTP/1.1\" \"curl/7.38.0\" DHE-RSA-AES128-SHA TLSv1.2"),
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
				Url:          "/",
				Agent:        "curl/7.38.0",
				Cipher:       "DHE-RSA-AES128-SHA",
				Protocol:     "TLSv1.2",
			}},
		},
	}

	for _, test := range tests {
		fmt.Println(test)
		e := test.elb
		e.BufferFile, _ = ioutil.TempFile("", ".LogsToElasticTest")

		e.BufferFile.Write(test.data)

		result, err := e.GetChunk(test.start, test.end)

		// We have to cast the expected result as a slice of interface{}
		expect := make([]interface{}, len(test.expect))
		for i, v := range test.expect {
			expect[i] = v
		}

		assert.IsType(t, test.err, err)
		assert.Equal(t, expect, result)

		e.BufferFile.Close()
		os.Remove(e.BufferFile.Name())
	}
}
