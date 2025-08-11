package sqlite

import "context"

type SqliteService interface {
	GetUserData(ctx context.Context, login string) (*UserData, error)
	SaveUserData(ctx context.Context, ud *UserData) error
	GetPrivateData(ctx context.Context, login string) ([]*PrivateData, error)
	SavePrivateData(ctx context.Context, pd *PrivateData) error
}
