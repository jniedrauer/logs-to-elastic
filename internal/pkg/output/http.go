// HTTP functions
package output

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
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
func Post(endpoint string, payload []byte, c *http.Client) error {
	resp, err := c.Post(endpoint, "text/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("bad HTTP Status: %d", resp.StatusCode))
	}

	return nil
}
