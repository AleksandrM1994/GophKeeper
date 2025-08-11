package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func DatePtr(date time.Time) *time.Time {
	return &date
}

func HashData(secret string, data []byte) (string, error) {
	hash := hmac.New(sha256.New, []byte(secret))
	_, err := hash.Write(data)
	if err != nil {
		return "", fmt.Errorf("error hashing data: %v", err)
	}
	sum := hash.Sum(nil)
	hashCode := base64.StdEncoding.EncodeToString(sum)
	return hashCode, nil
}
