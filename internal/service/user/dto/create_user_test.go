package dto

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	custom_errs "github.com/GophKeeper/internal/errors"
)

// Тесты для CreateUserRequest
func TestCreateUserRequest_Validations(t *testing.T) {
	// Тест успешной валидации
	t.Run("Valid Request", func(t *testing.T) {
		validRequest := CreateUserRequest{
			Login:    "test_user",
			Password: "secure_password",
		}

		err := validRequest.Validate()
		assert.NoError(t, err)
	})

	// Тест на пустой логин
	t.Run("Empty Login", func(t *testing.T) {
		invalidRequest := CreateUserRequest{
			Password: "secure_password",
		}

		err := invalidRequest.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid login")
		assert.True(t, errors.Is(err, custom_errs.ErrValidate))
	})

	// Тест на пустой пароль
	t.Run("Empty Password", func(t *testing.T) {
		invalidRequest := CreateUserRequest{
			Login: "test_user",
		}

		err := invalidRequest.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid password")
		assert.True(t, errors.Is(err, custom_errs.ErrValidate))
	})

	// Тест на полностью пустой запрос
	t.Run("Empty Request", func(t *testing.T) {
		emptyRequest := CreateUserRequest{}

		err := emptyRequest.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid login")
		assert.True(t, errors.Is(err, custom_errs.ErrValidate))
	})
}

// Тесты для CreateUserResponse
func TestCreateUserResponse_Validations(t *testing.T) {
	// Тест успешной инициализации
	t.Run("Valid Response", func(t *testing.T) {
		validResponse := CreateUserResponse{
			JWT:    "valid.jwt.token",
			UserID: "user123",
		}

		assert.NotNil(t, validResponse)
		assert.Equal(t, "valid.jwt.token", validResponse.JWT)
		assert.Equal(t, "user123", validResponse.UserID)
	})

	// Тест на пустой JWT
	t.Run("Empty JWT", func(t *testing.T) {
		emptyJWTResponse := CreateUserResponse{
			UserID: "user123",
		}

		assert.NotNil(t, emptyJWTResponse)
		assert.Equal(t, "", emptyJWTResponse.JWT)
		assert.Equal(t, "user123", emptyJWTResponse.UserID)
	})

	// Тест на пустой UserID
	t.Run("Empty UserID", func(t *testing.T) {
		emptyUserIDResponse := CreateUserResponse{
			JWT: "valid.jwt.token",
		}

		assert.NotNil(t, emptyUserIDResponse)
		assert.Equal(t, "valid.jwt.token", emptyUserIDResponse.JWT)
		assert.Equal(t, "", emptyUserIDResponse.UserID)
	})

	// Тест на полностью пустой ответ
	t.Run("Empty Response", func(t *testing.T) {
		emptyResponse := CreateUserResponse{}

		assert.NotNil(t, emptyResponse)
		assert.Equal(t, "", emptyResponse.JWT)
		assert.Equal(t, "", emptyResponse.UserID)
	})
}
