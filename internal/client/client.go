package client

import (
	"context"

	"github.com/GophKeeper/internal/client/dto"
	api "github.com/GophKeeper/pkg/api"
)

type Client interface {
	AuthUser(ctx context.Context, request *api.AuthUserRequest) (*api.AuthUserResponse, error)
	SavePrivateData(ctx context.Context, request *dto.SavePrivateDataRequest) error
}
