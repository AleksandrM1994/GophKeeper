package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatePtr(t *testing.T) {
	date := time.Now()
	ptr := DatePtr(date)
	assert.Equal(t, date, *ptr)
}

func TestHashData(t *testing.T) {
	secret := "my-secret"
	data := []byte("test-data")

	hash, err := HashData(secret, data)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Проверка длины хэша после Base64
	assert.Len(t, hash, 44) // sha256 = 32 bytes, base64 = 44 chars
}
