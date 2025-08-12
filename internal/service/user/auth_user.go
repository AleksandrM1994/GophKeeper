package user

import (
	"context"
	"fmt"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/user/dto"
)

func (s *UserServiceImpl) AuthUser(ctx context.Context, req *dto.AuthUserRequest) (*dto.AuthUserResponse, error) {
	errValidate := req.Validate()
	if errValidate != nil {
		return nil, fmt.Errorf("validate: %w", errValidate)
	}

	loginHash, err := service.HashData(s.cfg.HashSecret, []byte(req.Login))
	if err != nil {
		return nil, fmt.Errorf("failed to hash login: %w", err)
	}

	passHash, err := service.HashData(s.cfg.HashSecret, []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.GetUserByLogPass(ctx, loginHash, passHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found: %w", custom_errs.ErrUnauthorized)
	}

	jwt, errGenerateJWT := service.GenerateJWT(s.cfg.JWTSecret, user.ID)
	if errGenerateJWT != nil {
		return nil, fmt.Errorf("failed to generate jwt: %w", errGenerateJWT)
	}

	errUpdateUserByID := s.userRepo.UpdateUserByID(ctx, user.ID, func(currentUser *repository.User) error {
		currentUser.JWT = jwt
		return nil
	})
	if errUpdateUserByID != nil {
		return nil, fmt.Errorf("userRepository.UpdateUserByID:%w", errUpdateUserByID)
	}

	return &dto.AuthUserResponse{JWT: jwt, UserID: user.ID}, nil
}
