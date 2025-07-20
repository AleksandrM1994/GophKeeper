package repository

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/GophKeeper/config"
)

type Repository struct {
	cfg config.Config
	lg  *zap.SugaredLogger
	db  *gorm.DB
}

func NewRepository(cfg config.Config, lg *zap.SugaredLogger) (*Repository, error) {
	repo := &Repository{cfg: cfg, lg: lg}
	sqlDB, err := Connect(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = Migrate(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create db session with gorm: %w", err)
	}
	repo.db = gormDB
	return repo, nil
}

func Connect(dsn string) (*sql.DB, error) {
	db, errConnect := sql.Open("postgres", dsn)
	if errConnect != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", errConnect)
	}
	return db, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting SQL dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("error migration: %w", err)
	}
	return nil
}
