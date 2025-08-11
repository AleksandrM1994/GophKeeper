package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/mocks"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/user/dto"
)

func TestCheckAuthUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{JWTSecret: "jsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	// Создаём токен вручную
	userID := "uid123"
	jwtToken, _ := service.GenerateJWT(cfg.JWTSecret, userID)

	ctx := context.Background()
	req := &dto.CheckAuthRequest{JWT: jwtToken}

	// Мокируем GetUserByID, возвращаем пользователя с тем же JWT
	user := &repository.User{ID: userID, JWT: jwtToken}
	userRepo.
		EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil)

	resp, err := svc.CheckAuthUser(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, userID, resp.UserID)
}

func TestCheckAuthUser_InvalidJWT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{JWTSecret: "jsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.CheckAuthRequest{JWT: "invalid.token.value"}

	resp, err := svc.CheckAuthUser(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service.ParseJWT")
}

func TestCheckAuthUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{JWTSecret: "jsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	// Правильный токен
	userID := "u1"
	jwtToken, _ := service.GenerateJWT(cfg.JWTSecret, userID)

	ctx := context.Background()
	req := &dto.CheckAuthRequest{JWT: jwtToken}

	repoErr := errors.New("db fail")
	userRepo.
		EXPECT().
		GetUserByID(ctx, userID).
		Return(nil, repoErr)

	resp, err := svc.CheckAuthUser(ctx, req)

	assert.Nil(t, resp)
	assert.EqualError(t, err, "userRepo.GetUserByID:db fail")
}
