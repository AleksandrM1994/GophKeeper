package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestConnectSQLite_Success(t *testing.T) {
	// 1) Переключаемся в изолированный временный каталог
	tmp := t.TempDir()
	require.NoError(t, os.Chdir(tmp))

	// 2) Создаём структуру папок internal/storage/sqlite/db
	dbDir := filepath.Join("internal", "storage", "sqlite", "db")
	require.NoError(t, os.MkdirAll(dbDir, 0o755))

	// 3) Создаём пустой файл sqlite.db
	dbFile := filepath.Join(dbDir, "sqlite.db")
	f, err := os.Create(dbFile)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// 4) Вызываем функцию подключения
	db, err := ConnectSQLite()
	require.NoError(t, err)
	require.NotNil(t, db)

	// 5) Проверяем, что нужные таблицы появились
	mustHaveTable(t, db, "user_data")
	mustHaveTable(t, db, "private_data")
}

func mustHaveTable(t *testing.T, db *gorm.DB, name string) {
	has := db.Migrator().HasTable(name)
	require.Truef(t, has, "ожидали таблицу %q, но её нет", name)
}
