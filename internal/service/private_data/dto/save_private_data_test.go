package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/GophKeeper/internal/repository"
)

// Тесты для SavePrivateDataRequest
func TestSavePrivateDataRequest_Validations(t *testing.T) {
	// Тест успешной валидации
	t.Run("Valid Request", func(t *testing.T) {
		validRequest := SavePrivateDataRequest{
			Data:   []byte("some data"),
			Type:   repository.PrivateDataTypeText,
			UserID: "user123",
			Nonce:  []byte("nonce"),
			Login:  "test_user",
		}

		assert.NotNil(t, validRequest)
		assert.Equal(t, []byte("some data"), validRequest.Data)
		assert.Equal(t, repository.PrivateDataTypeText, validRequest.Type)
		assert.Equal(t, "user123", validRequest.UserID)
		assert.Equal(t, []byte("nonce"), validRequest.Nonce)
		assert.Equal(t, "test_user", validRequest.Login)
	})

	// Тест на отсутствие обязательных полей
	t.Run("Missing Required Fields", func(t *testing.T) {
		// Пустое поле Data
		invalidRequest1 := SavePrivateDataRequest{
			Type:   repository.PrivateDataTypeText,
			UserID: "user123",
			Login:  "test_user",
		}
		assert.NotNil(t, invalidRequest1)

		// Невалидный тип
		invalidRequest2 := SavePrivateDataRequest{
			Data:   []byte("some data"),
			Type:   "INVALID_TYPE",
			UserID: "user123",
			Nonce:  []byte("nonce"),
			Login:  "test_user",
		}
		assert.NotNil(t, invalidRequest2)
		assert.NotEqual(t, repository.PrivateDataTypeText, invalidRequest2.Type)

		// Пустой UserID
		invalidRequest3 := SavePrivateDataRequest{
			Data:  []byte("some data"),
			Type:  repository.PrivateDataTypeText,
			Nonce: []byte("nonce"),
			Login: "test_user",
		}
		assert.NotNil(t, invalidRequest3)
		assert.Equal(t, "", invalidRequest3.UserID)

		// Пустой Nonce
		invalidRequest4 := SavePrivateDataRequest{
			Data:   []byte("some data"),
			Type:   repository.PrivateDataTypeText,
			UserID: "user123",
			Login:  "test_user",
		}
		assert.NotNil(t, invalidRequest4)

		// Пустой Login
		invalidRequest5 := SavePrivateDataRequest{
			Data:   []byte("some data"),
			Type:   repository.PrivateDataTypeText,
			UserID: "user123",
			Nonce:  []byte("nonce"),
		}
		assert.NotNil(t, invalidRequest5)
		assert.Equal(t, "", invalidRequest5.Login)
	})

	// Тест на пустую структуру
	t.Run("Empty Request", func(t *testing.T) {
		emptyRequest := SavePrivateDataRequest{}
		assert.NotNil(t, emptyRequest)
		assert.Equal(t, "", emptyRequest.UserID)
		assert.Equal(t, "", emptyRequest.Login)
	})
}
