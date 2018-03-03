// HTTP functions
package io

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
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
func Consumer(in <-chan []byte, config *conf.Config) int {
	c := GetClient()

	var wg sync.WaitGroup
	wg.Add(1)
	out := make(chan bool)

	go func() {
		for p := range in {
			if Post(config.Logstash, p, c) {
				out <- true
			}
		}
		wg.Done()
	}()

	wg.Wait()
	close(out)

	return len(out)
}
