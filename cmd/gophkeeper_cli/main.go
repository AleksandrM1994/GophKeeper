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

	gophKeeperClient := client.NewClient(&lg, cfg)

	rootCmd := cobra_cli.NewRootCmd()

	authCmd := cobra_cli.NewAuthCmd(gophKeeperClient)
	rootCmd.AddCommand(authCmd)

	saveCmd := cobra_cli.NewSaveCmd(gophKeeperClient)
	authCmd.AddCommand(saveCmd)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
