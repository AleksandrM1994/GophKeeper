package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/storage/sqlite"
)

func main() {
	logger, loggerErr := zap.NewDevelopment()
	if loggerErr != nil {
		panic("cannot initialize zap")
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	lg := *logger.Sugar()

	cfg, errNewConfig := config.NewConfig()
	if errNewConfig != nil {
		panic(errNewConfig)
	}

	db, err := sqlite.ConnectSQLite()
	if err != nil {
		lg.Fatal(err)
	}

	sqliteService := sqlite.NewServiceImpl(&lg, db)

	kafkaService := kafka.NewKafkaService(&lg, cfg.KafkaHost, sqliteService)
	errInitKafkaTopics := kafka.InitKafkaTopics(cfg.KafkaHost)
	if errInitKafkaTopics != nil {
		lg.Fatalf("kafkaController.InitKafkaTopics, %w", errInitKafkaTopics)
	}

	if err := kafkaService.ReadMessage(context.Background(), kafka.GophKeeperPrivateDataSavedTopic); err != nil {
		lg.Fatalf("ошибка при чтении сообщений из Kafka: %w", err)
	}
}
