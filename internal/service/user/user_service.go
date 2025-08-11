package user

import (
	"context"

	"github.com/GophKeeper/internal/service/user/dto"
)

type UserService interface {
	AuthUser(ctx context.Context, req *dto.AuthUserRequest) (*dto.AuthUserResponse, error)
	CheckAuthUser(ctx context.Context, req *dto.CheckAuthRequest) (*dto.CheckAuthResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) error
}
