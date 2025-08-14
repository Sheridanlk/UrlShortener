package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomAlias(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 0",
			size: 0,
		},
		{
			name: "size = 20",
			size: 20,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str1 := NewRandomAlias(test.size)
			str2 := NewRandomAlias(test.size)

			assert.Len(t, str1, test.size)
			assert.Len(t, str2, test.size)
		})
	}
}
