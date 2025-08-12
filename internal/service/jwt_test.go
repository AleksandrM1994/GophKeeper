package service

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	secret := "test-secret"
	userID := "12345"

	token, err := GenerateJWT(secret, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Разделим токен на части
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3, "токен должен состоять из трёх частей")

	// Декодируем payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	var payload map[string]interface{}
	err = json.Unmarshal(payloadBytes, &payload)
	assert.NoError(t, err)

	assert.Equal(t, "my-app", payload["iss"])
	assert.Equal(t, userID, payload["sub"])
	assert.Equal(t, userID, payload["user_id"])
	assert.NotNil(t, payload["exp"])
	assert.NotNil(t, payload["iat"])
	assert.NotNil(t, payload["nbf"])
}

func TestParseJWT_Success(t *testing.T) {
	secret := "test-secret"
	userID := "12345"

	// Сгенерируем валидный токен
	token, _ := GenerateJWT(secret, userID)

	// Парсим его
	claims, err := ParseJWT(secret, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "my-app", claims.Issuer)
	assert.Equal(t, userID, claims.Subject)
}

func TestParseJWT_InvalidSignature(t *testing.T) {
	secret := "test-secret"
	userID := "12345"

	// Сгенерируем токен
	token, _ := GenerateJWT(secret, userID)

	// Парсим с другим секретом — ожидаем ошибку
	_, err := ParseJWT("wrong-secret", token)
	assert.Error(t, err)
}

func TestParseJWT_InvalidAlgorithm(t *testing.T) {
	invalidToken := "invalid.jwt.token.with.rsa.signature"
	_, err := ParseJWT("test-secret", invalidToken)
	assert.Error(t, err)
}

func TestParseJWT_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	userID := "12345"

	// Генерация истёкшего токена
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-app",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // истёкший
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))

	// Парсим истёкший токен
	_, err := ParseJWT(secret, signedToken)
	assert.Error(t, err)
}
