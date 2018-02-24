/*
TODO: Documentation
*/
package main

import (
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/logging"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.Init()

	log.Debug("Well hello there")
}
