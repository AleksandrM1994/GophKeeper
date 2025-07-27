package bbolt

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type UserData struct {
	Login string `json:"login"`
	JWT   string `json:"jwt"`
	Key   []byte `json:"key"`
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
		b := tx.Bucket([]byte(bucketName))
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
		b := tx.Bucket([]byte(bucketName))
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
