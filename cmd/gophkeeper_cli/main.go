package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"

	"github.com/mattn/go-shellwords"

	"github.com/GophKeeper/cmd/gophkeeper_cli/cobra_cli"
	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/client"
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

	authCmd := cobra_cli.NewAuthCmd(&lg, gophKeeperClient, sqliteService)
	rootCmd.AddCommand(authCmd)

	saveCmd := cobra_cli.NewSaveCmd(&lg, gophKeeperClient, sqliteService)
	rootCmd.AddCommand(saveCmd)

	getCmd := cobra_cli.NewGetCmd(&lg, sqliteService)
	rootCmd.AddCommand(getCmd)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	reader := bufio.NewReader(os.Stdin)
	parser := shellwords.NewParser()

	fmt.Println("CLI GophKeeper запущен. Введите команду (например: auth, save --login \"test\" --text-data \"раз два три четыре\") или exit для выхода.")

	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			lg.Errorf("Ошибка ввода команды: %v", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			fmt.Println("Выход из CLI.")
			return
		}

		// Разбираем строку на аргументы, учитывая кавычки
		args, err := parser.Parse(line)
		if err != nil {
			lg.Errorf("Не удалось распарсить команду: %v", err)
			continue
		}

		// Передаём args в Cobra
		rootCmd.SetArgs(args)
		if err := rootCmd.ExecuteContext(context.Background()); err != nil {
			lg.Errorf("Ошибка выполнения команды: %v", err)
		}
	}
}
