package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	custom_errs "github.com/GophKeeper/internal/errors"
)

type UserRepositoryImpl struct {
	*Repository
}

func NewUserRepository(repo *Repository) UserRepository {
	return &UserRepositoryImpl{repo}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return fmt.Errorf("error create user: %w", custom_errs.ErrDuplicateKey)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepositoryImpl) GetUserByLogPass(ctx context.Context, login, password string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Model(&User{}).Where("login = ? and password = ?", login, password).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Model(&user).Where("id = ?", userID).Take(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUserByID(ctx context.Context, userID string, updateFunc func(user *User) error) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, &User{ID: userID}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("tx.Take: %w", err)
		}
		err = updateFunc(&user)
		if err != nil {
			return fmt.Errorf("updateFunc: %w", err)
		}
		err = tx.Save(&user).Error
		if err != nil {
			return fmt.Errorf("tx.Save: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
