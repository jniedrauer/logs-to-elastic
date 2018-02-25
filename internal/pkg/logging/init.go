package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var defaultLogLevel log.Level = log.DebugLevel

func Init() {
	logLevel, err := getLevel()
	log.SetLevel(logLevel)
	if err != nil {
		log.Error(err)
	}
}

func getLevel() (log.Level, error) {
	lvl, set := os.LookupEnv("LOG_LEVEL")
	if !set {
		return defaultLogLevel, nil
	}

	logLevel, err := log.ParseLevel(lvl)
	if err != nil {
		return defaultLogLevel, err
	}

	return logLevel, nil
}
