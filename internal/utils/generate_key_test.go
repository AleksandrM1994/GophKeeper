package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		expected bool
	}{
		{"Valid size", 16, true},
		{"Large size", 256, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.size)
			if tt.expected {
				assert.NoError(t, err)
				assert.Len(t, key, tt.size)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
