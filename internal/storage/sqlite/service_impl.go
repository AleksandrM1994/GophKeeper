package sqlite

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ColumnInfo struct {
	CID        int
	Name       string
	Type       string
	NotNull    bool
	DefaultVal interface{}
	PK         bool
}

type ServiceImpl struct {
	lg *zap.SugaredLogger
	db *gorm.DB
}

func NewServiceImpl(lg *zap.SugaredLogger, db *gorm.DB) *ServiceImpl {
	return &ServiceImpl{lg: lg, db: db}
}

func (s *ServiceImpl) GetUserData(ctx context.Context, login string) (*UserData, error) {
	var ud UserData
	if err := s.db.WithContext(ctx).Where("login = ?", login).Take(&ud).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		s.lg.Errorf("failed to get user data for login %s: %v", login, err)
		return nil, err
	}
	return &ud, nil
}

// SaveUserData сохраняет или обновляет данные пользователя
func (s *ServiceImpl) SaveUserData(ctx context.Context, ud *UserData) error {
	if err := s.db.WithContext(ctx).Save(ud).Error; err != nil {
		s.lg.Errorf("failed to save user data for login %s: %v", ud.Login, err)
		return err
	}
	return nil
}

// GetPrivateData возвращает приватные данные по ID
func (s *ServiceImpl) GetPrivateData(ctx context.Context, login string) ([]*PrivateData, error) {
	var cols []ColumnInfo
	s.db.Raw("PRAGMA table_info('private_data')").Scan(&cols)
	s.lg.Infof("cols = %+v", cols)

	var raw []*PrivateData
	s.lg.Infof("GetPrivateData: login='%s'", login)
	if err := s.db.Raw("SELECT * FROM private_data WHERE user_login = ?", login).Scan(&raw).Error; err != nil {
		s.lg.Errorf("failed to get private data for login %s: %v", login, err)
		return nil, err
	}

	return raw, nil
}

// SavePrivateData сохраняет или обновляет приватные данные
func (s *ServiceImpl) SavePrivateData(ctx context.Context, pd *PrivateData) error {
	if pd.CreatedAt == nil {
		pd.CreatedAt = new(time.Time)
		*pd.CreatedAt = time.Now()
	}
	if pd.UpdatedAt == nil {
		pd.UpdatedAt = new(time.Time)
		*pd.UpdatedAt = time.Now()
	}

	if err := s.db.WithContext(ctx).Save(pd).Error; err != nil {
		s.lg.Errorf("failed to save private data with id %s: %v", pd.ID, err)
		return err
	}
	return nil
}
