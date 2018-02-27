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
