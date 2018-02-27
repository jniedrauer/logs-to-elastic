package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		lenv   string
		ienv   string
		csenv  string
		expect Conf
	}{
		{
			lenv:   "logstash",
			ienv:   "indexname",
			csenv:  "100",
			expect: Conf{Logstash: "logstash", IndexName: "indexname", ChunkSize: 100},
		},
	}

	for _, test := range tests {
		os.Setenv("LOGSTASH", test.lenv)
		os.Setenv("INDEXNAME", test.ienv)
		os.Setenv("CHUNK_SIZE", test.csenv)

		result := Init()

		assert.Equal(t, test.expect, *result)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		env    string
		def    string
		expect string
	}{
		{
			// Test handling of unset environment var
			def:    "foo",
			expect: "foo",
		},
		{
			// Test environment var set
			env:    "foo",
			def:    "bar",
			expect: "foo",
		},
	}
	for _, test := range tests {
		if len(test.env) <= 0 {
			os.Unsetenv("_TEST_ENV_VAR")
		} else {
			os.Setenv("_TEST_ENV_VAR", test.env)
		}
		response := GetEnvOrDefault("_TEST_ENV_VAR", test.def)
		assert.Equal(t, test.expect, response)
	}
}
