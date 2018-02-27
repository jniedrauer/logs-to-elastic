/*
HTTP client singleton
*/
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
