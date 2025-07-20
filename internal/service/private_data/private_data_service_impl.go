package private_data

import (
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/repository"
)

type PrivateDataServiceImpl struct {
	lg              *zap.SugaredLogger
	cfg             config.Config
	privateDataRepo repository.PrivateDataRepository
	cookie          *securecookie.SecureCookie
}

func NewService(lg *zap.SugaredLogger, cfg config.Config, privateDataRepo repository.PrivateDataRepository) *PrivateDataServiceImpl {
	srv := &PrivateDataServiceImpl{lg: lg, cfg: cfg, privateDataRepo: privateDataRepo}
	return srv
}
