package sqlite

import (
	"time"

	"github.com/GophKeeper/internal/repository"
)

type UserData struct {
	Login string `gorm:"primaryKey"`
	JWT   string
	Key   []byte
}

func (UserData) TableName() string { return "user_data" }

type PrivateData struct {
	ID        string `gorm:"primaryKey"`
	Type      repository.PrivateDataType
	Data      []byte
	CreatedAt *time.Time
	UpdatedAt *time.Time
	UserLogin string `gorm:"index;not null"`
}

func (PrivateData) TableName() string { return "private_data" }
