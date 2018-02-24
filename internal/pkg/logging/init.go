package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	logLevel, err := getLevel()
	log.SetLevel(logLevel)
	if err != nil {
		log.Error(err)
	}
}

func getLevel() (log.Level, error) {
	defaultLevel := log.InfoLevel

	lvl, set := os.LookupEnv("LOG_LEVEL")
	if !set {
		return defaultLevel, nil
	}

	logLevel, err := log.ParseLevel(lvl)
	if err != nil {
		return defaultLevel, err
	}

	return logLevel, nil
}
