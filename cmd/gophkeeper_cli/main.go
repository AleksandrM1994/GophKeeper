package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/GophKeeper/cmd/gophkeeper_cli/cobra_cli"
	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/storage/bbolt"
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

	db, err := bbolt.ConnectBbolt()
	if err != nil {
		lg.Fatal(err)
	}
	defer db.Close()

	bboltService := bbolt.NewServiceImpl(&lg, db)

	gophKeeperClient := client.NewClient(&lg, cfg)

	rootCmd := cobra_cli.NewRootCmd()

	authCmd := cobra_cli.NewAuthCmd(gophKeeperClient, bboltService)
	rootCmd.AddCommand(authCmd)

	saveCmd := cobra_cli.NewSaveCmd(gophKeeperClient, bboltService)
	rootCmd.AddCommand(saveCmd)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	kafkaController := kafka.NewController(&lg, cfg.GetString("kafka.host"), bboltService)
	errInitKafkaTopics := kafkaController.InitKafkaTopics()
	if errInitKafkaTopics != nil {
		lg.Fatalf("kafkaController.InitKafkaTopics, %w", errInitKafkaTopics)
	}
}
