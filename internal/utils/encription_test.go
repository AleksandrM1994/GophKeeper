package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validKey  = bytes.Repeat([]byte{0x01}, 32) // 256-bit key
	shortKey  = []byte("shortkey")             // invalid length
	plaintext = []byte("The quick brown fox jumps over the lazy dog")
)

func TestEncryptDecrypt_Success(t *testing.T) {
	// Encrypt
	ciphertext, nonce, err := Encrypt(plaintext, validKey)
	require.NoError(t, err, "expected no error from Encrypt")
	require.NotEmpty(t, ciphertext, "ciphertext should not be empty")

	// Decrypt
	decrypted, err := Decrypt(ciphertext, nonce, validKey)
	require.NoError(t, err, "expected no error from Decrypt")
	assert.Equal(t, plaintext, decrypted, "decrypted plaintext should match original")
}

func TestEncrypt_InvalidKey(t *testing.T) {
	_, _, err := Encrypt(plaintext, shortKey)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create cipher")
}

func TestDecrypt_InvalidKey(t *testing.T) {
	// First create valid ciphertext/nonce
	ciphertext, nonce, err := Encrypt(plaintext, validKey)
	require.NoError(t, err)

	// Decrypt with wrong key
	wrongKey := bytes.Repeat([]byte{0x02}, 32)
	_, err = Decrypt(ciphertext, nonce, wrongKey)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt data")
}

func TestDecrypt_CorruptedCiphertext(t *testing.T) {
	ciphertext, nonce, err := Encrypt(plaintext, validKey)
	require.NoError(t, err)

	// Corrupt the ciphertext
	if len(ciphertext) > 0 {
		ciphertext[0] ^= 0xFF
	}

	_, err = Decrypt(ciphertext, nonce, validKey)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt data")
}

func TestDecrypt_InvalidKeyLength(t *testing.T) {
	ciphertext, nonce, err := Encrypt(plaintext, validKey)
	require.NoError(t, err)

	// Use invalid key length
	_, err = Decrypt(ciphertext, nonce, shortKey)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create cipher")
}
