package handlers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	tests := []struct {
		input  int
		total  int
		expect Response
		err    error
	}{
		// All transmitted
		{
			input:  3,
			total:  3,
			expect: Response{Message: "sent records: 3", Ok: true},
			err:    nil,
		},
		{
			input:  1,
			total:  4,
			expect: Response{Message: "sent records: 1", Ok: false},
			err:    errors.New(""),
		},
	}

	for _, test := range tests {
		result, err := NewResponse(test.input, test.total)

		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, result)
	}
}
