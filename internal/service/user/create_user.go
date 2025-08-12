package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/user/dto"
)

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *dto.CreateUserRequest) error {
	errValidate := req.Validate()
	if errValidate != nil {
		return fmt.Errorf("validate: %w", errValidate)
	}

	loginHash, err := service.HashData(s.cfg.HashSecret, []byte(req.Login))
	if err != nil {
		return fmt.Errorf("failed to hash login: %w", err)
	}

	passHash, err := service.HashData(s.cfg.HashSecret, []byte(req.Password))
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	id := uuid.New().String()

	errCreateUser := s.userRepo.CreateUser(
		ctx,
		&repository.User{
			ID:       id,
			Login:    loginHash,
			Password: passHash,
		})
	if errCreateUser != nil {
		return fmt.Errorf("userRepo.CreateUser: %w", errCreateUser)
	}

	return nil
}
