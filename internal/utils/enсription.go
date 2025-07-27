package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func Encrypt(plaintext, encryptionKey []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce, err := GenerateKey(aesgcm.NonceSize())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func Decrypt(ciphertextStr, nonceStr, encryptionKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	plaintext, err := aesgcm.Open(nil, nonceStr, ciphertextStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plaintext, nil
}
