package user

import (
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/repository"
)

type UserServiceImpl struct {
	lg       *zap.SugaredLogger
	cfg      config.Config
	userRepo repository.UserRepository
}

func NewService(lg *zap.SugaredLogger, cfg config.Config, userRepo repository.UserRepository) *UserServiceImpl {
	srv := &UserServiceImpl{lg: lg, cfg: cfg, userRepo: userRepo}
	return srv
}
