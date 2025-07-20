package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
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

func LunaCheck(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	if len(number) == 0 {
		return false
	}

	sum := 0
	shouldDouble := false

	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}
