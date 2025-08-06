package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

const (
	dbPath            = "internal/storage/sqlite/db/cli_data.db"
	bucketUserData    = "users_data"
	bucketPrivateData = "private_data"
)

func ConnectSQLite() (*sql.DB, error) {
	// Получаем абсолютный путь к файлу БД
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for DB: %w", err)
	}

	// Создаем папку, если её нет
	dbDir := filepath.Dir(absDBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create DB directory: %w", err)
	}

	db, err := sql.Open("sqlite3", absDBPath+"?_journal=WAL&_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite DB: %w", err)
	}

	// Включение WAL
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user_data ( " +
		"login TEXT NOT NULL UNIQUE, " +
		"jwt TEXT NOT NULL, " +
		"key BLOB NOT NULL)")

	if err != nil {
		return nil, fmt.Errorf("failed to create table 'user_data': %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS private_data ( " +
		"id TEXT PRIMARY KEY," +
		"type TEXT NOT NULL," +
		"data BLOB NOT NULL," +
		"created_at DATETIME," +
		"updated_at DATETIME)")

	if err != nil {
		return nil, fmt.Errorf("failed to create table 'private_data': %w", err)
	}

	return db, nil
}
