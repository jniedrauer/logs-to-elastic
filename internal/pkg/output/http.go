/*
HTTP client singleton
*/
package output

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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
		client = &http.Client{Transport: tr}
	})

	return client
}

func Post(endpoint string, payload []byte, c *http.Client) ([]byte, error) {
	resp, err := c.Post(endpoint, "text/json", bytes.NewReader(payload))
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
