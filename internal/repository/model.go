package repository

import "time"

type User struct {
	ID       string `db:"id;primary_key"`
	Login    string `db:"login;not null;unique"`
	Password string `db:"password;not null"`
	JWT      string `db:"jwt;unique"`
}

func (User) TableName() string {
	return "users"
}

type PrivateData struct {
	ID        string          `db:"id;primary_key"`
	Data      string          `db:"data;not null"`
	Type      PrivateDataType `gorm:"type:enum('UNKNOWN', 'TEXT', 'FILE', 'AUTH', 'BANK')"`
	CreatedAt *time.Time      `db:"created_at;not null"`
	UpdatedAt *time.Time      `db:"updated_at;not null"`
	UserID    string          `db:"user_id;not null"`
}

func (PrivateData) TableName() string {
	return "private_data"
}
