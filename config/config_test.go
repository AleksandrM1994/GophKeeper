package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тесты для Config
func TestConfig_Defaults(t *testing.T) {
	// Очищаем переменные окружения
	defer os.Clearenv()

	// Создаем конфигурацию без переменных окружения
	cfg, err := NewConfig()
	require.NoError(t, err)

	// Проверяем значения по умолчанию
	assert.Equal(t, ":8080", cfg.HTTPAddress)
	assert.Equal(t, "user=postgres password=postgres dbname=praktikum host=localhost port=5432 sslmode=disable", cfg.DSN)
	assert.Equal(t, "my_hash_secret", cfg.HashSecret)
	assert.Equal(t, "my_jwt_secret", cfg.JWTSecret)
	assert.Equal(t, "localhost:9092", cfg.KafkaHost)
}

func TestConfig_EnvOverrides(t *testing.T) {
	// Устанавливаем переменные окружения
	defer os.Clearenv()
	os.Setenv("RUN_ADDRESS", ":9000")
	os.Setenv("DATABASE_URI", "custom_dsn")
	os.Setenv("HASH_SECRET", "custom_hash")
	os.Setenv("JWT_SECRET", "custom_jwt")
	os.Setenv("KAFKA_HOST", "custom_kafka:9092")

	// Создаем конфигурацию
	cfg, err := NewConfig()
	require.NoError(t, err)

	// Проверяем переопределенные значения
	assert.Equal(t, ":9000", cfg.HTTPAddress)
	assert.Equal(t, "custom_dsn", cfg.DSN)
	assert.Equal(t, "custom_hash", cfg.HashSecret)
	assert.Equal(t, "custom_jwt", cfg.JWTSecret)
	assert.Equal(t, "custom_kafka:9092", cfg.KafkaHost)
}
