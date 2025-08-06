package sqlite

import (
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/GophKeeper/internal/repository"
)

type UserData struct {
	Login string `json:"login"`
	JWT   string `json:"jwt"`
	Key   []byte `json:"key"`
}

type PrivateData struct {
	ID        string                     `json:"id"`
	Type      repository.PrivateDataType `json:"type"`
	Data      []byte                     `json:"data"`
	CreatedAt *time.Time                 `json:"created_at"`
	UpdatedAt *time.Time                 `json:"updated_at"`
}

type ServiceImpl struct {
	lg *zap.SugaredLogger
	db *sql.DB
}

func NewServiceImpl(lg *zap.SugaredLogger, db *sql.DB) *ServiceImpl {
	return &ServiceImpl{lg: lg, db: db}
}

func (s *ServiceImpl) GetUserData(login string) (*UserData, error) {
	var ud UserData
	err := s.db.QueryRow("SELECT login, jwt, key FROM user_data WHERE login = ?", login).
		Scan(&ud.Login, &ud.JWT, &ud.Key)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		s.lg.Errorf("failed to get user data for login %s: %v", login, err)
		return nil, err
	}
	return &ud, nil
}

func (s *ServiceImpl) SaveUserData(ud *UserData) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO user_data (login, jwt, key)
		VALUES (?, ?, ?)
	`, ud.Login, ud.JWT, ud.Key)

	if err != nil {
		s.lg.Errorf("failed to save user data for login %s: %v", ud.Login, err)
		return err
	}
	return nil
}

func (s *ServiceImpl) GetPrivateData(id string) (*PrivateData, error) {
	var pd PrivateData
	var createdAtStr, updatedAtStr string

	err := s.db.QueryRow("SELECT id, type, data, created_at, updated_at FROM private_data WHERE id = ?", id).
		Scan(&pd.ID, &pd.Type, &pd.Data, &createdAtStr, &updatedAtStr)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		s.lg.Errorf("failed to get private data with id %s: %v", id, err)
		return nil, err
	}

	// Парсинг времени
	if createdAtStr != "" {
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			s.lg.Errorf("failed to parse created_at for id %s: %v", id, err)
			return nil, err
		}
		pd.CreatedAt = &createdAt
	}

	if updatedAtStr != "" {
		updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			s.lg.Errorf("failed to parse updated_at for id %s: %v", id, err)
			return nil, err
		}
		pd.UpdatedAt = &updatedAt
	}

	return &pd, nil
}

func (s *ServiceImpl) SavePrivateData(pd *PrivateData) error {
	createdAtStr := ""
	updatedAtStr := ""

	if pd.CreatedAt != nil {
		createdAtStr = pd.CreatedAt.Format(time.RFC3339)
	}
	if pd.UpdatedAt != nil {
		updatedAtStr = pd.UpdatedAt.Format(time.RFC3339)
	}

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO private_data (id, type, data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, pd.ID, pd.Type, pd.Data, createdAtStr, updatedAtStr)

	if err != nil {
		s.lg.Errorf("failed to save private data with id %s: %v", pd.ID, err)
		return err
	}
	return nil
}
