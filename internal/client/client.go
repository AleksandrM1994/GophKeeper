package client

import (
	"context"

	api "github.com/GophKeeper/pkg/api"
)

type Client interface {
	AuthUser(ctx context.Context, request *api.AuthUserRequest) (*api.AuthUserResponse, error)
	SavePrivateData(ctx context.Context, request *api.SavePrivateDataRequest) error
}
