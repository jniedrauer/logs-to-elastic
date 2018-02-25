package logging

import (
	"errors"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLevel(t *testing.T) {
	tests := []struct {
		env    string
		expect log.Level
		err    error
	}{
		{
			expect: defaultLogLevel,
			err:    nil,
		},
		{
			env:    "DEBUG",
			expect: log.DebugLevel,
			err:    nil,
		},
		{
			env:    "INFO",
			expect: log.InfoLevel,
			err:    nil,
		},
		{
			env:    "ERROR",
			expect: log.ErrorLevel,
			err:    nil,
		},
		{
			env:    "FATAL",
			expect: log.FatalLevel,
			err:    nil,
		},
		{
			env:    "fake",
			expect: defaultLogLevel,
			err:    errors.New(""),
		},
	}

	for _, test := range tests {
		if len(test.env) <= 0 {
			os.Unsetenv("LOG_LEVEL")
		} else {
			os.Setenv("LOG_LEVEL", test.env)
		}

		result, err := getLevel()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, result)
	}
}
