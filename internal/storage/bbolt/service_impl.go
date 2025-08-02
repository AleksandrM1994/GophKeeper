package bbolt

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
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
	db *bbolt.DB
}

func NewServiceImpl(lg *zap.SugaredLogger, db *bbolt.DB) *ServiceImpl {
	return &ServiceImpl{lg: lg, db: db}
}

func (s *ServiceImpl) SaveUserData(data *UserData) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketUserData))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		buf, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}
		return b.Put([]byte(data.Login), buf)
	})
}

func (s *ServiceImpl) GetUserData(login string) (*UserData, error) {
	data := &UserData{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketUserData))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		buf := b.Get([]byte(login))
		if buf != nil {
			err := json.Unmarshal(buf, data)
			if err != nil {
				return fmt.Errorf("json unmarshal: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		s.lg.Errorf("get userdata: %s", err)
	}
	return data, nil
}

func (s *ServiceImpl) SavePrivateData(login string, data *PrivateData) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketPrivateData))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		buf, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}
		return b.Put([]byte(login), buf)
	})
}

func (s *ServiceImpl) GetPrivateData(login string) (*PrivateData, error) {
	data := &PrivateData{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketPrivateData))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		buf := b.Get([]byte(login))
		if buf != nil {
			err := json.Unmarshal(buf, data)
			if err != nil {
				return fmt.Errorf("json unmarshal: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		s.lg.Errorf("get userdata: %s", err)
	}
	return data, nil
}
