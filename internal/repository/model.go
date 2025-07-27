package repository

import "time"

type User struct {
	ID       string `db:"primary_key"`
	Login    string `db:"not null;unique"`
	Password string `db:"not null"`
	JWT      string `db:"unique"`
}

func (User) TableName() string {
	return "users"
}

type PrivateData struct {
	ID        string          `db:"primary_key"`
	Data      []byte          `db:"not null"`
	Type      PrivateDataType `gorm:"type:enum('UNKNOWN', 'TEXT', 'FILE', 'AUTH', 'BANK')"`
	CreatedAt *time.Time      `db:"not null"`
	UpdatedAt *time.Time      `db:"not null"`
	UserID    string          `db:"not null"`
	Nonce     []byte          `db:"not null"`
}

func (PrivateData) TableName() string {
	return "private_data"
}
