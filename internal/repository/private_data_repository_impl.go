package repository

import (
	"context"
	"fmt"
)

type PrivateDataRepositoryImpl struct {
	*Repository
}

func NewPrivateDataRepository(repo *Repository) PrivateDataRepository {
	return &PrivateDataRepositoryImpl{repo}
}

func (r PrivateDataRepositoryImpl) CreatePrivateData(ctx context.Context, data *PrivateData) error {
	err := r.db.Create(data).Error
	if err != nil {
		return fmt.Errorf("failed to create private data: %w", err)
	}

	return nil
}
