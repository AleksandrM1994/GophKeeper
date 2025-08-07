package sqlite

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "internal/storage/sqlite/db/cli_data.db"
)

func ConnectSQLite() (*gorm.DB, error) {
	// Получаем абсолютный путь к файлу БД
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for DB: %w", err)
	}

	// Со здаем папку, если её нет
	dbDir := filepath.Dir(absDBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create DB directory: %w", err)
	}

	dsn := absDBPath + "?_journal=WAL&_foreign_keys=on&_busy_timeout=5000"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite DB: %w", err)
	}

	if errAutoMigrate := db.AutoMigrate(&UserData{}); errAutoMigrate != nil {
		return nil, fmt.Errorf("failed to auto migrate UserData: %w", errAutoMigrate)
	}

	if errAutoMigrate := db.AutoMigrate(&PrivateData{}); errAutoMigrate != nil {
		return nil, fmt.Errorf("failed to auto migrate PrivateData: %w", errAutoMigrate)
	}

	db = db.Debug()

	return db, nil
}
