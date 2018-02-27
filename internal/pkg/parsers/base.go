package parsers

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

type Payloader interface {
	GetSlice(int, int) []interface{}
}

type BaseLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	LogGroup  string `json:"logGroup"`
	IndexName string `json:"indexname"`
}

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

func unixToIso8601(unix int64) string {
	return time.Unix(unix, 0).UTC().Format("2006-01-02T15:04:05-0000")
}
