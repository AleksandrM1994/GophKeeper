package bbolt

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

const (
	dbPath = "internal/storage/bbolt/db/cli_data.db"
)

func ConnectBbolt() (*bbolt.DB, error) {
	// Получаем путь к исполняемому файлу
	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not get executable path: %v", err)
	}

	exPath := filepath.Dir(ex)
	db, err := bbolt.Open(filepath.Join(exPath, dbPath), 0600, &bbolt.Options{
		Timeout: 15 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open bbolt db: %v", err)
	}

	return db, nil
}
