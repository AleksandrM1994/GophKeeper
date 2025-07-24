package bbolt

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const (
	bucketName = "users_data"
)

type UserData struct {
	Login string `json:"login"`
	JWT   string `json:"jwt"`
	Key   string `json:"key"`
}

type ServiceImpl struct {
	lg *zap.SugaredLogger
	db *bbolt.DB
}

func NewServiceImpl(lg *zap.SugaredLogger, db *bbolt.DB) *ServiceImpl {
	return &ServiceImpl{lg: lg, db: db}
}

func (s *ServiceImpl) SaveUserData(record UserData) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		buf, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}
		return b.Put([]byte(record.Login), buf)
	})
}

func (s *ServiceImpl) GetUserData(login string) error {
	var data UserData
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		buf := b.Get([]byte(login))
		if buf == nil {
			return fmt.Errorf("user not found")
		}
		return json.Unmarshal(buf, &data)
	})
}
