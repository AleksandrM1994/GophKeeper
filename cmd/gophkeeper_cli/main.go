package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"

	"github.com/GophKeeper/cmd/gophkeeper_cli/cobra_cli"
	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger, loggerErr := zap.NewDevelopment()
	if loggerErr != nil {
		panic("cannot initialize zap")
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			// преобразуем заведомо «safe» ошибки в no-op
			if !strings.Contains(err.Error(), "inappropriate ioctl") {
				fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
			}
		}
	}()

	lg := *logger.Sugar()

	cfg, err := config.NewCliConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlite.ConnectSQLite()
	if err != nil {
		lg.Fatal(err)
	}

	sqliteService := sqlite.NewServiceImpl(&lg, db)

	gophKeeperClient := client.NewClient(&lg, cfg)

	rootCmd := cobra_cli.NewRootCmd()

	authCmd := cobra_cli.NewAuthCmd(gophKeeperClient, sqliteService)
	rootCmd.AddCommand(authCmd)

	saveCmd := cobra_cli.NewSaveCmd(gophKeeperClient, sqliteService)
	rootCmd.AddCommand(saveCmd)

	kafkaController := kafka.NewController(&lg, cfg.GetString("kafka.host"), sqliteService)
	errInitKafkaTopics := kafkaController.InitKafkaTopics()
	if errInitKafkaTopics != nil {
		lg.Fatalf("kafkaController.InitKafkaTopics, %w", errInitKafkaTopics)
	}

	lg.Info("запуск Kafka-потребителя для топика gophkeeper.privateDataSaved")
	go func() {
		if err := kafkaController.ReadMessage(context.Background(), kafka.GophKeeperPrivateDataSavedTopic); err != nil {
			lg.Fatalf("ошибка при чтении сообщений из Kafka: %w", err)
		}
	}()

	// Обработка сигналов завершения работы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск CLI
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	// Ожидание сигнала остановки
	<-sigChan
	lg.Info("Получен сигнал завершения. Выход...")
}
