package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/mocks"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service"
	"github.com/GophKeeper/internal/service/user/dto"
)

func TestAuthUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{HashSecret: "hsec", JWTSecret: "jsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.AuthUserRequest{Login: "user", Password: "pass"}

	// Хэши
	loginHash, _ := service.HashData(cfg.HashSecret, []byte(req.Login))
	passHash, _ := service.HashData(cfg.HashSecret, []byte(req.Password))

	// Мокируем GetUserByLogPass
	user := &repository.User{ID: "uid", JWT: ""}
	userRepo.
		EXPECT().
		GetUserByLogPass(ctx, loginHash, passHash).
		Return(user, nil)

	// Мокируем UpdateUserByID: применяем функцию
	userRepo.
		EXPECT().
		UpdateUserByID(ctx, user.ID, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, fn func(*repository.User) error) error {
			return fn(user)
		})

	resp, err := svc.AuthUser(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, resp.UserID)
	assert.NotEmpty(t, resp.JWT)
}

func TestAuthUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{HashSecret: "h", JWTSecret: "j"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.AuthUserRequest{Login: "u", Password: "p"}

	loginHash, _ := service.HashData(cfg.HashSecret, []byte(req.Login))
	passHash, _ := service.HashData(cfg.HashSecret, []byte(req.Password))

	userRepo.
		EXPECT().
		GetUserByLogPass(ctx, loginHash, passHash).
		Return(nil, nil)

	resp, err := svc.AuthUser(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), custom_errs.ErrUnauthorized.Error())
}

func TestAuthUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{HashSecret: "h", JWTSecret: "j"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.AuthUserRequest{Login: "u", Password: "p"}

	loginHash, _ := service.HashData(cfg.HashSecret, []byte(req.Login))
	passHash, _ := service.HashData(cfg.HashSecret, []byte(req.Password))

	repoErr := errors.New("db error")
	userRepo.
		EXPECT().
		GetUserByLogPass(ctx, loginHash, passHash).
		Return(nil, repoErr)

	resp, err := svc.AuthUser(ctx, req)

	assert.Nil(t, resp)
	assert.EqualError(t, err, "failed to get user: db error")
}
