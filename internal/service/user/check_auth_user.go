package user

import (
	"context"
	"fmt"

	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/user/dto"
)

func (s *UserServiceImpl) CheckAuthUser(_ context.Context, req *dto.CheckAuthRequest) (*dto.CheckAuthResponse, error) {
	claims, errParseJWT := service.ParseJWT(s.cfg.JWTSecret, req.JWT)
	if errParseJWT != nil {
		return nil, fmt.Errorf("userRepository.UpdateUserByID:%w", errParseJWT)
	}

	return &dto.CheckAuthResponse{UserID: claims.UserID}, nil
}
