// Central logging configuration with logrus
package logging

import (
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/conf"
	log "github.com/sirupsen/logrus"
)

var defaultLogLevel log.Level = log.InfoLevel

// Create a new root logger with configuration
func Init() {
	logLevel, err := getLevel()
	log.SetLevel(logLevel)
	if err != nil {
		log.Error(err)
	}
}

// Get log level from an environment variable
func getLevel() (log.Level, error) {
	lvl := conf.GetEnvOrDefault("LOG_LEVEL", defaultLogLevel.String())

	logLevel, err := log.ParseLevel(lvl)
	if err != nil {
		return defaultLogLevel, err
	}

	return logLevel, nil
}
