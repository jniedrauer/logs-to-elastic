// Log parsing utilities
package parsers

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

type Payloader interface {
	GetSlice(int, int) []interface{}
}

// The minimum information that an encoded log should include
type BaseLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	LogGroup  string `json:"logGroup"`
	IndexName string `json:"indexname"`
}

// Encode a slice of logs for transmission
func SliceEncode(p Payloader, idx int, end int, delim string) []byte {
	var encoded []byte
	bdelim := []byte(delim)

	for _, evt := range p.GetSlice(idx, end) {
		enc, err := json.Marshal(evt)
		if err != nil {
			log.Error("failed to encode: %v", evt)
			continue
		}

		if len(encoded) > 0 {
			encoded = append(encoded, bdelim...)
		}
		encoded = append(encoded, enc...)
	}

	return encoded
}

// Convert a unix timestamp to ISO 8601 format
func unixToIso8601(unix int64) string {
	return time.Unix(unix, 0).UTC().Format("2006-01-02T15:04:05-0000")
}
