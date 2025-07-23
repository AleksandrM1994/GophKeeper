package client

import (
	"context"

	"github.com/GophKeeper/internal/handlers/private_data"
	"github.com/GophKeeper/internal/handlers/user"
)

type Client interface {
	AuthUser(ctx context.Context, request *user.AuthUserRequest) (*user.AuthUserResponse, error)
	SavePrivateData(ctx context.Context, request *private_data.SavePrivateDataRequest) error
}
