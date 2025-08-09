package sqlite

import (
	"time"

	"github.com/GophKeeper/internal/repository"
)

type UserData struct {
	ID    uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Login string `gorm:"column:login;uniqueIndex"`
	JWT   string `gorm:"column:jwt"`
	Key   []byte `gorm:"column:key"`
}

func (UserData) TableName() string { return "user_data" }

type PrivateData struct {
	ID        string                     `gorm:"column:id;primaryKey"`
	Type      repository.PrivateDataType `gorm:"column:type"`
	Data      []byte                     `gorm:"column:data"`
	Nonce     []byte                     `gorm:"column:nonce"`
	CreatedAt *time.Time                 `gorm:"column:created_at"`
	UpdatedAt *time.Time                 `gorm:"column:updated_at"`
	UserLogin string                     `gorm:"column:user_login;index;not null"`
}

func (PrivateData) TableName() string { return "private_data" }
