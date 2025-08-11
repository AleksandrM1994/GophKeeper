package private_data

import (
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/kafka"
	"github.com/GophKeeper/internal/repository"
)

type PrivateDataServiceImpl struct {
	lg              *zap.SugaredLogger
	cfg             config.Config
	privateDataRepo repository.PrivateDataRepository
	cookie          *securecookie.SecureCookie
	kafkaService    kafka.KafkaService
}

func NewService(
	lg *zap.SugaredLogger,
	cfg config.Config,
	privateDataRepo repository.PrivateDataRepository,
	kafkaService kafka.KafkaService) PrivateDataService {
	return &PrivateDataServiceImpl{lg: lg, cfg: cfg, privateDataRepo: privateDataRepo, kafkaService: kafkaService}
}
