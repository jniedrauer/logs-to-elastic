package aws

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegion(t *testing.T) {
	tests := []struct {
		env    string
		expect string
		err    error
	}{
		{
			// Test handling of no set region
			expect: "us-east-1",
		},
		{
			// Test handling of setting region
			env:    "us-west-2",
			expect: "us-west-2",
		},
	}

	for _, test := range tests {
		if len(test.env) <= 0 {
			os.Unsetenv("AWS_REGION")
		} else {
			os.Setenv("AWS_REGION", test.env)
		}
		response := getRegion()
		assert.Equal(t, test.expect, response)
	}
}
