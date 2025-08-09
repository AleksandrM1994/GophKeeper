package sqlite

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "internal/storage/sqlite/db/sqlite.db"
)

func ConnectSQLite() (*gorm.DB, error) {
	// Получаем абсолютный путь к файлу БД
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for DB: %w", err)
	}

	fmt.Printf("Using SQLite DB at: %s", absDBPath)

	dsn := absDBPath
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
