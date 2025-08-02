package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/viper"
)

func NewConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

type Config struct {
	HTTPAddress string `env:"RUN_ADDRESS" envDefault:":8080"`
	DSN         string `env:"DATABASE_URI" envDefault:"user=postgres password=postgres dbname=praktikum host=localhost port=5432 sslmode=disable"`
	HashSecret  string `env:"HASH_SECRET" envDefault:"my_hash_secret"`
	JWTSecret   string `env:"JWT_SECRET" envDefault:"my_jwt_secret"`
	KafkaHost   string `env:"KAFKA_HOST" envDefault:"localhost:9092"`
}

func NewCliConfig() (*viper.Viper, error) {
	cfg := viper.New()

	// Получаем путь к исполняемому файлу
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Не удалось получить путь к исполняемому файлу: %v", err)
	}
	exPath := filepath.Dir(ex)

	// Добавляем путь к директории с бинарником + поддиректория config
	cfg.AddConfigPath(filepath.Join(exPath, "config"))
	cfg.SetConfigName("config_cli")
	cfg.SetConfigType("yaml")

	cfg.AutomaticEnv()

	if err := cfg.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения конфига: %v", err)
	}
	return cfg, nil
}
