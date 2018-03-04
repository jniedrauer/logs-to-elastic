// HTTP functions
package net

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/parsers"
	log "github.com/sirupsen/logrus"
)

var client *http.Client
var once sync.Once

// Get an HTTP client using a singleton model
func GetClient() *http.Client {
	once.Do(func() {
		log.Debug("generating a new HTTP client")

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}

		client = &http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		}
	})

	return client
}

// Send a post request and only return HTTP status code pass/fail
func Post(endpoint string, payload []byte, c *http.Client) bool {
	resp, err := c.Post(endpoint, "text/json", bytes.NewReader(payload))
	if err != nil {
		log.Error(err.Error())
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("bad HTTP Status: %d", resp.StatusCode))
		return false
	}

	return true
}

// Asynchronous POST to endpoint
func LogstashConsumer(in <-chan *parsers.EncodedChunk, config *conf.Config) uint32 {
	c := GetClient()
	var oks uint32 = 0
	var wg sync.WaitGroup

	for p := range in {
		wg.Add(1)
		go func(p *parsers.EncodedChunk) {
			if Post(config.Logstash, p.Payload, c) {
				atomic.AddUint32(&oks, p.Records)
			}
			log.Debug(fmt.Sprintf("transmitted batch: %d, total: %d", p.Records, oks))
			wg.Done()
		}(p)
	}

	wg.Wait()

	return atomic.LoadUint32(&oks)
}
