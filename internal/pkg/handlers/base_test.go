package handler

import (
	"errors"
	"testing"

	"github.com/jniedrauer/logs-to-elastic/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		cfg    *config.Config
		self   *BaseHandle
		err    error
		expect *config.LogGroup
	}{
		{
			// Single possible match
			cfg:    &config.Config{LogGroups: []config.LogGroup{{Name: "foo"}}},
			self:   &BaseHandle{LogGroup: "foo"},
			expect: &config.LogGroup{Name: "foo"},
			err:    nil,
		},
		{
			// Multiple possible matches
			cfg:    &config.Config{LogGroups: []config.LogGroup{{Name: "foo"}, {Name: "bar"}}},
			self:   &BaseHandle{LogGroup: "foo"},
			expect: &config.LogGroup{Name: "foo"},
			err:    nil,
		},
		{
			// Values are loaded
			cfg:    &config.Config{LogGroups: []config.LogGroup{{Name: "foo", IndexName: "bar"}}},
			self:   &BaseHandle{LogGroup: "foo"},
			expect: &config.LogGroup{Name: "foo", IndexName: "bar"},
			err:    nil,
		},
		{
			// No config found
			cfg:  &config.Config{LogGroups: []config.LogGroup{{Name: "foo"}}},
			self: &BaseHandle{LogGroup: "bar"},
			err:  errors.New(""),
		},
	}

	for _, test := range tests {
		b := test.self
		err := b.LoadConfig(test.cfg)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, b.Config)
	}
}

func TestPayloadEncode(t *testing.T) {
	tests := []struct {
		payload []*LogEvent
		delim   string
		expect  []byte
		errs    []error
	}{
		{
			// Single record encoding
			payload: []*LogEvent{{Timestamp: "foo", Message: "bar"}},
			delim:   "\n",
			expect:  []byte("{\"timestamp\":\"foo\",\"message\":\"bar\"}"),
			errs:    []error{nil},
		},
		{
			// Multiple record encoding with delimiter
			payload: []*LogEvent{{Timestamp: "foo", Message: "bar"}, {Timestamp: "bar", Message: "foo"}},
			delim:   "\n",
			expect:  []byte("{\"timestamp\":\"foo\",\"message\":\"bar\"}\n{\"timestamp\":\"bar\",\"message\":\"foo\"}"),
			errs:    []error{nil},
		},
	}

	for _, test := range tests {
		result, errs := payloadEncode(test.payload, test.delim)
		for idx, err := range errs {
			assert.IsType(t, test.errs[idx], err)
		}
		assert.Equal(t, string(test.expect), string(result))
	}
}
