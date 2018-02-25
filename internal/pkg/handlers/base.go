package handler

import (
	"encoding/json"
	"errors"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/config"
)

type Handler interface {
	LoadConfig(cfg *config.Config) error
	GetDataChunks() []*[]byte
	Decode()
}

type BaseHandle struct {
	LogGroup  string
	LogEvents []*LogEvent
	Config    *config.LogGroup
}

func (b *BaseHandle) LoadConfig(cfg *config.Config) error {
	for _, groupcfg := range cfg.LogGroups {
		if groupcfg.Name == b.LogGroup {
			b.Config = &groupcfg
			return nil
		}
	}
	return errors.New("no log group configuration found")
}

type LogEvent struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func payloadEncode(payload []*LogEvent, delim string) ([]byte, []error) {
	var logs []byte
	var errs []error
	bdelim := []byte(delim)

	for _, evt := range payload {
		// We have to do this because we're not using a json delimiter
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
