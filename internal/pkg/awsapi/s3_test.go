package awsapi

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tests := []struct {
		bucket string
		key    string
		iface  MockS3Iface
		expect string
		err    error
	}{
		// Test that downloaded data is written to tempfile
		{
			bucket: "foo",
			key:    "bar",
			iface:  MockS3Iface{TestData: []byte("foo")},
			expect: "foo",
			err:    nil,
		},
	}

	for _, test := range tests {
		c := S3Client{test.iface}
		result, err := c.Get(test.bucket, test.key)

		assert.IsType(t, test.err, err)

		content, _ := ioutil.ReadFile(result)
		assert.Equal(t, test.expect, string(content))
	}
}
