package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		aenv   string
		denv   string
		ienv   string
		lenv   string
		csenv  string
		expect Config
	}{
		{
			aenv:  "us-west-2",
			denv:  ",",
			ienv:  "indexname",
			lenv:  "logstash",
			csenv: "100",
			expect: Config{
				AwsRegion: "us-west-2",
				Delimiter: []byte(","),
				IndexName: "indexname",
				Logstash:  "logstash",
				ChunkSize: 100,
			},
		},
	}

	for _, test := range tests {
		os.Setenv("AWS_REGION", test.aenv)
		os.Setenv("DELIMITER", test.denv)
		os.Setenv("INDEXNAME", test.ienv)
		os.Setenv("LOGSTASH", test.lenv)
		os.Setenv("CHUNK_SIZE", test.csenv)

		result := NewConfig()

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
