package logging

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

var defaultLevel log.Level = log.InfoLevel

func TestGetLevelDefault(t *testing.T) {
	// Verify that getLevel returns a default if environment level is unset
	os.Unsetenv("LOG_LEVEL")

	result, err := getLevel()

	if result != defaultLevel {
		t.Errorf("Expected %v, got %v", defaultLevel, result)
	}

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
}

func TestGetLevelDebug(t *testing.T) {
	// Verify debug logging can be set
	os.Setenv("LOG_LEVEL", "DEBUG")

	result, err := getLevel()

	if result != log.DebugLevel {
		t.Errorf("Expected %v, got %v", log.DebugLevel, result)
	}

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
}

func TestGetLevelInvalid(t *testing.T) {
	// Verify invalid level error handling
	os.Setenv("LOG_LEVEL", "fake")

	result, err := getLevel()

	if result != defaultLevel {
		t.Errorf("Expected %v, got %v", defaultLevel, result)
	}

	if err == nil {
		t.Errorf("Error not raised")
	}
}
