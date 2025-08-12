package repository

import "context"

type PrivateDataRepository interface {
	CreatePrivateData(ctx context.Context, data *PrivateData) error
}
