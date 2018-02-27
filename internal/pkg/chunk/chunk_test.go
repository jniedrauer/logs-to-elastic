package chunk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		len    int
		sz     int
		expect [][]int
	}{
		{
			len:    4,
			sz:     2,
			expect: [][]int{{0, 2}, {2, 4}},
		},
		{
			len:    2,
			sz:     1,
			expect: [][]int{{0, 1}, {1, 2}},
		},
		{
			len:    5,
			sz:     2,
			expect: [][]int{{0, 2}, {2, 4}, {4, 5}},
		},
	}

	for _, test := range tests {
		var result [][]int

		Chunk(test.sz, test.len, func(idx int, end int) {
			result = append(result, []int{idx, end})
		})

		assert.Equal(t, test.expect, result)
	}
}
