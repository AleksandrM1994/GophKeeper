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

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{HashSecret: "hsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.CreateUserRequest{Login: "test", Password: "pwd123"}

	userRepo.
		EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(nil)

	err := svc.CreateUser(ctx, req)

	assert.NoError(t, err)
}

func TestCreateUser_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{HashSecret: "hsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.CreateUserRequest{Login: "", Password: ""}

	err := svc.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validate")
}

func TestCreateUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := zap.NewExample().Sugar()
	cfg := config.Config{HashSecret: "hsec"}
	userRepo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(lg, cfg, userRepo)

	ctx := context.Background()
	req := &dto.CreateUserRequest{Login: "test", Password: "pass"}

	// Хэши
	loginHash, _ := service.HashData(cfg.HashSecret, []byte(req.Login))
	passHash, _ := service.HashData(cfg.HashSecret, []byte(req.Password))

	repoErr := errors.New("db write failed")
	userRepo.
		EXPECT().
		CreateUser(ctx, gomock.AssignableToTypeOf(&repository.User{})).
		DoAndReturn(func(_ context.Context, u *repository.User) error {
			// Проверяем, что ID, Login и Password корректно заданы
			assert.NotEmpty(t, u.ID)
			assert.Equal(t, loginHash, u.Login)
			assert.Equal(t, passHash, u.Password)
			return repoErr
		})

	err := svc.CreateUser(ctx, req)

	assert.EqualError(t, err, "userRepo.CreateUser: db write failed")
}
