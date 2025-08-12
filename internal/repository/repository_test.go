package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	custom_errs "github.com/GophKeeper/internal/errors"
)

type StorageTestSuite struct {
	suite.Suite
	db                  *sql.DB
	gormDB              *gorm.DB
	testUserRepo        *UserRepositoryImpl
	testPrivateDataRepo *PrivateDataRepositoryImpl
	cleanup             func()
}

func (s *StorageTestSuite) SetupSuite() {
	// Build DSN for test Postgres
	dsn := fmt.Sprintf(
		"host=localhost port=5433 user=test password=test dbname=test sslmode=disable",
	)
	db, errConnect := sql.Open("postgres", dsn)
	if errConnect != nil {
		panic(fmt.Sprintf("failed to connect to test DB: %v", errConnect))
	}

	s.db = db

	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Errorf("error setting SQL dialect: %w", err))
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(fmt.Errorf("error migration: %w", err))
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error open gorm: %w", err))
	}
	s.gormDB = gormDB

	s.testPrivateDataRepo = &PrivateDataRepositoryImpl{
		Repository: &Repository{db: gormDB},
	}

	s.testUserRepo = &UserRepositoryImpl{
		Repository: &Repository{db: gormDB},
	}

	s.cleanup = func() {
		s.db.Close()
	}
}

func (s *StorageTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *StorageTestSuite) TearDownTest() {
	// Очищаем данные после каждого теста
	ctx := context.Background()
	_, err := s.db.ExecContext(ctx, "TRUNCATE TABLE private_data CASCADE")
	require.NoError(s.T(), err)
	_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE users CASCADE")
	require.NoError(s.T(), err)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestCreatePrivateData_Success() {
	ctx := context.Background()
	now := time.Now().UTC()

	pd := &PrivateData{
		ID:        "12345",
		Data:      []byte("secret"),
		Type:      PrivateDataTypeText,
		CreatedAt: &now,
		UpdatedAt: &now,
		UserID:    "user1",
		Nonce:     []byte("nonce"),
	}

	// Выполняем создание
	err := s.testPrivateDataRepo.CreatePrivateData(ctx, pd)
	s.Require().NoError(err, "expected no error on CreatePrivateData")

	// Проверяем, что запись действительно сохранилась
	var fetched PrivateData
	s.Require().NoError(
		s.gormDB.WithContext(ctx).
			First(&fetched, "id = ?", pd.ID).
			Error,
	)
	s.Equal(pd.Data, fetched.Data)
	s.Equal(pd.Type, fetched.Type)
	s.Equal(pd.UserID, fetched.UserID)
	s.Equal(pd.Nonce, fetched.Nonce)
}

func (s *StorageTestSuite) TestCreatePrivateData_DuplicateID() {
	ctx := context.Background()
	now := time.Now().UTC()

	pd1 := &PrivateData{
		ID:        "dup-" + now.Format("150405"),
		Data:      []byte("first"),
		Type:      PrivateDataTypeAuth,
		CreatedAt: &now,
		UpdatedAt: &now,
		UserID:    "userX",
		Nonce:     []byte("n1"),
	}
	s.Require().NoError(s.testPrivateDataRepo.CreatePrivateData(ctx, pd1))

	pd2 := &PrivateData{
		ID:        pd1.ID, // дублируем ID
		Data:      []byte("second"),
		Type:      PrivateDataTypeBank,
		CreatedAt: &now,
		UpdatedAt: &now,
		UserID:    "userY",
		Nonce:     []byte("n2"),
	}
	err := s.testPrivateDataRepo.CreatePrivateData(ctx, pd2)
	s.Require().Error(err, "expected error on duplicate ID")
	s.Contains(err.Error(), "failed to create private data")
}

func (s *StorageTestSuite) TestCreateUserAndGetByID() {
	ctx := context.Background()
	now := time.Now().UTC().Format(time.RFC3339Nano)

	user := &User{
		ID:       "u-" + now,
		Login:    "login-" + now,
		Password: "pass",
		JWT:      "jwt-" + now,
	}
	// CreateUser
	err := s.testUserRepo.CreateUser(ctx, user)
	s.Require().NoError(err)

	// GetUserByID
	fetched, err := s.testUserRepo.GetUserByID(ctx, user.ID)
	s.Require().NoError(err)
	s.Equal(user.Login, fetched.Login)
	s.Equal(user.Password, fetched.Password)
}

func (s *StorageTestSuite) TestCreateUser_DuplicateLogin() {
	ctx := context.Background()
	u1 := &User{ID: "123", Login: "test", Password: "p"}
	s.Require().NoError(s.testUserRepo.CreateUser(ctx, u1))

	u2 := &User{ID: "456", Login: "test", Password: "p2"}
	err := s.testUserRepo.CreateUser(ctx, u2)
	s.Require().Error(err)
	s.Contains(err.Error(), custom_errs.ErrDuplicateKey.Error())
}

func (s *StorageTestSuite) TestGetUserByLogPass_Success() {
	ctx := context.Background()
	user := &User{
		ID:       "12345",
		Login:    "login",
		Password: "pass",
	}
	s.Require().NoError(s.testUserRepo.CreateUser(ctx, user))

	fetched, err := s.testUserRepo.GetUserByLogPass(ctx, "login", "pass")
	s.Require().NoError(err)
	s.Equal(user.ID, fetched.ID)
	s.Equal(user.Login, fetched.Login)
}

func (s *StorageTestSuite) TestGetUserByLogPass_NotFound() {
	ctx := context.Background()
	_, err := s.testUserRepo.GetUserByLogPass(ctx, "", "")
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to get user")
}

func (s *StorageTestSuite) TestGetUserByID_Success() {
	ctx := context.Background()
	user := &User{
		ID:       "12345",
		Login:    "login",
		Password: "pass",
	}
	s.Require().NoError(s.testUserRepo.CreateUser(ctx, user))

	fetched, err := s.testUserRepo.GetUserByID(ctx, "12345")
	s.Require().NoError(err)
	s.Equal(user.Login, fetched.Login)
}

func (s *StorageTestSuite) TestGetUserByID_NotFound() {
	ctx := context.Background()
	_, err := s.testUserRepo.GetUserByID(ctx, "doesnt-exist")
	s.Require().Error(err)
	s.Contains(err.Error(), "failed to get user")
}

func (s *StorageTestSuite) TestUpdateUserByID_Success() {
	ctx := context.Background()
	user := &User{
		ID:       "12345",
		Login:    "login",
		Password: "pass",
		JWT:      "",
	}
	s.Require().NoError(s.testUserRepo.CreateUser(ctx, user))

	err := s.testUserRepo.UpdateUserByID(ctx, user.ID, func(u *User) error {
		u.Password = "pass-new"
		u.JWT = "jwt-new"
		return nil
	})
	s.Require().NoError(err)

	fetched, err := s.testUserRepo.GetUserByID(ctx, user.ID)
	s.Require().NoError(err)
	s.Equal("pass-new", fetched.Password)
	s.Equal("jwt-new", fetched.JWT)
}
