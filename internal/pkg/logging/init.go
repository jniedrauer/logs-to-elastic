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
	logLevel := log.InfoLevel
	lvl, set := os.LookupEnv("LOG_LEVEL")

	if !set {
		return logLevel, nil
	}

	logLevel, err := log.ParseLevel(lvl)
	return logLevel, err
}
