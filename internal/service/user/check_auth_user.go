package user

import (
	"context"
	"fmt"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/user/dto"
)

func (s *UserServiceImpl) CheckAuthUser(ctx context.Context, req *dto.CheckAuthRequest) (*dto.CheckAuthResponse, error) {
	claims, errParseJWT := service.ParseJWT(s.cfg.JWTSecret, req.JWT)
	if errParseJWT != nil {
		return nil, fmt.Errorf("service.ParseJWT:%w", errParseJWT)
	}

	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetUserByID:%w", err)
	}

	if user.JWT != req.JWT {
		return nil, fmt.Errorf("wrong JWT:%w", custom_errs.ErrUnauthorized)
	}

	return &dto.CheckAuthResponse{UserID: claims.UserID}, nil
}
