package logging

import (
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/env"
	log "github.com/sirupsen/logrus"
)

var defaultLogLevel log.Level = log.InfoLevel

func Init() {
	logLevel, err := getLevel()
	log.SetLevel(logLevel)
	if err != nil {
		log.Error(err)
	}
}

func getLevel() (log.Level, error) {
	lvl := env.GetEnvOrDefault("LOG_LEVEL", defaultLogLevel.String())

	logLevel, err := log.ParseLevel(lvl)
	if err != nil {
		return defaultLogLevel, err
	}

	return logLevel, nil
}
