package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/GophKeeper/internal/repository"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Открываем in-memory SQLite
	gormDB, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{},
	)
	require.NoError(t, err)

	// Автомиграция
	err = gormDB.AutoMigrate(&UserData{}, &PrivateData{})
	require.NoError(t, err)

	return gormDB, func() {
		db, err := gormDB.DB()
		if err == nil {
			db.Close()
		}
	}
}

func TestServiceImpl_GetUserData_Success(t *testing.T) {
	gormDB, teardown := setupTestDB(t)
	defer teardown()

	userData := &UserData{
		Login: "test",
		Key:   []byte("key"),
		JWT:   "jwt",
	}
	require.NoError(t, gormDB.Create(userData).Error)

	service := NewServiceImpl(zap.NewNop().Sugar(), gormDB)

	ud, err := service.GetUserData(context.Background(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, ud)
	assert.Equal(t, "test", ud.Login)
	assert.Equal(t, "jwt", ud.JWT)
}

func TestServiceImpl_GetUserData_NotFound(t *testing.T) {
	gormDB, teardown := setupTestDB(t)
	defer teardown()

	service := NewServiceImpl(zap.NewNop().Sugar(), gormDB)

	ud, err := service.GetUserData(context.Background(), "nonexistent")
	assert.Nil(t, ud)
	assert.NoError(t, err)
}

func TestServiceImpl_SaveUserData_Success(t *testing.T) {
	gormDB, teardown := setupTestDB(t)
	defer teardown()

	userData := &UserData{
		Login: "test",
		Key:   []byte("key"),
		JWT:   "jwt",
	}
	service := NewServiceImpl(zap.NewNop().Sugar(), gormDB)

	err := service.SaveUserData(context.Background(), userData)
	assert.NoError(t, err)

	var saved UserData
	err = gormDB.Where("login = ?", "test").First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "test", saved.Login)
	assert.Equal(t, "jwt", saved.JWT)
}

func TestServiceImpl_GetPrivateData_Success(t *testing.T) {
	gormDB, teardown := setupTestDB(t)
	defer teardown()

	privateData := &PrivateData{
		ID:        "1",
		Type:      repository.PrivateDataTypeText,
		UserLogin: "test",
		Data:      []byte("data"),
		Nonce:     []byte("nonce"),
	}
	require.NoError(t, gormDB.Create(privateData).Error)

	service := NewServiceImpl(zap.NewNop().Sugar(), gormDB)

	raw, err := service.GetPrivateData(context.Background(), "test")
	assert.NoError(t, err)
	assert.Len(t, raw, 1)
	assert.Equal(t, repository.PrivateDataTypeText, raw[0].Type)
	assert.Equal(t, "1", raw[0].ID)
}

func TestServiceImpl_GetPrivateData_Empty(t *testing.T) {
	gormDB, teardown := setupTestDB(t)
	defer teardown()

	service := NewServiceImpl(zap.NewNop().Sugar(), gormDB)

	raw, err := service.GetPrivateData(context.Background(), "test")
	assert.NoError(t, err)
	assert.Empty(t, raw)
}

func TestServiceImpl_SavePrivateData_Success(t *testing.T) {
	gormDB, teardown := setupTestDB(t)
	defer teardown()

	service := NewServiceImpl(zap.NewNop().Sugar(), gormDB)

	now := time.Now()
	pd := &PrivateData{
		ID:        "1",
		Type:      repository.PrivateDataTypeText,
		UserLogin: "test",
		Data:      []byte("data"),
		Nonce:     []byte("nonce"),
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	err := service.SavePrivateData(context.Background(), pd)
	assert.NoError(t, err)

	var saved PrivateData
	err = gormDB.Where("id = ?", "1").First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, repository.PrivateDataTypeText, saved.Type)
	assert.Equal(t, "test", saved.UserLogin)
}
