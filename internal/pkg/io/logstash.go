// HTTP functions
package io

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	log "github.com/sirupsen/logrus"
)

var client *http.Client
var once sync.Once

// Get an HTTP client using a singleton model
func GetClient() *http.Client {
	once.Do(func() {
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
func Consumer(in <-chan []byte, config *conf.Config) uint32 {
	c := GetClient()
	var oks uint32
	var wg sync.WaitGroup

	for p := range in {
		wg.Add(1)
		go func() {
			if Post(config.Logstash, p, c) {
				atomic.AddUint32(&oks, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return atomic.LoadUint32(&oks)
}
