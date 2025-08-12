package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims расширяет jwt.RegisteredClaims своими полями
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Генерация токена
func GenerateJWT(secret string, userID string) (string, error) {
	// Формируем пользовательские claims
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-app",                                          // кто выпустил токен
			Subject:   fmt.Sprint(userID),                                // тема (обычно ID пользователя)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // время жизни
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // время выдачи
			NotBefore: jwt.NewNumericDate(time.Now()),                    // токен действует не раньше
		},
	}

	// Создаём объект токена и подписываем его
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return signedToken, nil
}

// Проверка и разбор токена
func ParseJWT(secret string, jwtString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Приводим данные к типу *MyClaims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
