package parsers

import (
	"encoding/json"
)

type BaseLog struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	LogGroup  string `json:"logGroup"`
	IndexName string `json:"indexname"`
}

func PayloadEncode(payload []interface{}, delim string) ([]byte, []error) {
	var logs []byte
	var errs []error
	bdelim := []byte(delim)

	for _, evt := range payload {
		// We have to do this because we might not use a json delimiter
		enc, err := json.Marshal(evt)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(logs) > 0 {
			logs = append(logs, bdelim...)
		}
		logs = append(logs, enc...)
	}

	return logs, errs
}
