package bbolt

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

const (
	dbPath     = "internal/storage/bbolt/db/cli_data.db"
	bucketName = "users_data"
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

	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("could not create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db.Update: %v", err)
	}

	return db, nil
}
