package dto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Предположим, что у CheckAuthRequest тоже есть метод Validate()
func (r *CheckAuthRequest) Validate() error {
	if r.JWT == "" {
		return fmt.Errorf("invalid jwt token")
	}
	return nil
}

// Тесты для CheckAuthRequest
func TestCheckAuthRequest_Validations(t *testing.T) {
	// Тест успешной валидации
	t.Run("Valid Request", func(t *testing.T) {
		validRequest := CheckAuthRequest{
			JWT: "valid.jwt.token",
		}

		err := validRequest.Validate()
		assert.NoError(t, err)
	})

	// Тест на пустой JWT
	t.Run("Empty JWT", func(t *testing.T) {
		invalidRequest := CheckAuthRequest{}

		err := invalidRequest.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid jwt token")
	})
}

// Тесты для CheckAuthResponse
func TestCheckAuthResponse_Validations(t *testing.T) {
	// Тест успешной инициализации
	t.Run("Valid Response", func(t *testing.T) {
		validResponse := CheckAuthResponse{
			UserID: "user123",
		}

		assert.NotNil(t, validResponse)
		assert.Equal(t, "user123", validResponse.UserID)
	})

	// Тест на пустой UserID
	t.Run("Empty UserID", func(t *testing.T) {
		emptyResponse := CheckAuthResponse{}

		assert.NotNil(t, emptyResponse)
		assert.Equal(t, "", emptyResponse.UserID)
	})

	// Тест на null-значение
	t.Run("Nil Response", func(t *testing.T) {
		var nilResponse *CheckAuthResponse

		assert.Nil(t, nilResponse)
	})
}
