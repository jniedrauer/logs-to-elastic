package logging

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGetLevelDefault(t *testing.T) {
	// Verify that getLevel returns a default if environment level is unset
	os.Unsetenv("LOG_LEVEL")

	result, err := getLevel()

	if result != defaultLogLevel {
		t.Errorf("expected %v, got %v", defaultLogLevel, result)
	}

	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestGetLevelDebug(t *testing.T) {
	// Verify debug logging can be set
	os.Setenv("LOG_LEVEL", "DEBUG")

	result, err := getLevel()

	if result != log.DebugLevel {
		t.Errorf("expected %v, got %v", log.DebugLevel, result)
	}

	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestGetLevelInvalid(t *testing.T) {
	// Verify invalid level error handling
	os.Setenv("LOG_LEVEL", "fake")

	result, err := getLevel()

	if result != defaultLogLevel {
		t.Errorf("expected %v, got %v", defaultLogLevel, result)
	}

	if err == nil {
		t.Errorf("error not raised")
	}
}
