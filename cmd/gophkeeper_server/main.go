package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	privateDataHandlers "github.com/GophKeeper/internal/handlers/private_data"
	userHandlers "github.com/GophKeeper/internal/handlers/user"
	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/repository"
	privateDataService "github.com/GophKeeper/internal/service/private_data"
	userService "github.com/GophKeeper/internal/service/user"
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

	g := gin.Default()

	repo, err := repository.NewRepository(cfg, &lg)
	if err != nil {
		lg.Fatalf("repository.NewRepository, %w", err)
	}

	userRepo := repository.NewUserRepository(repo)
	privateDataRepo := repository.NewPrivateDataRepository(repo)

	userServiceImpl := userService.NewService(&lg, cfg, userRepo)
	privateDataServiceImpl := privateDataService.NewService(&lg, cfg, privateDataRepo)

	userController := userHandlers.NewController(cfg, &lg, userServiceImpl)
	userController.RegisterRoutes(g)
	privateDataController := privateDataHandlers.NewController(cfg, &lg, userServiceImpl, privateDataServiceImpl)
	privateDataController.RegisterRoutes(g)

	server := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      g,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	kafkaController := kafka.NewController(&lg, cfg.KafkaHost, nil)
	errInitKafkaTopics := kafkaController.InitKafkaTopics()
	if errInitKafkaTopics != nil {
		lg.Fatalf("kafkaController.InitKafkaTopics, %w", errInitKafkaTopics)
	}

	err = server.ListenAndServe()
	if err != nil {
		lg.Fatalf("g.Run, %w", err)
	}
}
