package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Тестовые случаи
func TestRespondWithError(t *testing.T) {
	// Создаем тестовые случаи
	testCases := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "BadRequest",
			err:          ErrBadRequest,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ValidateError",
			err:          ErrValidate,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Unauthorized",
			err:          ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "NotFound",
			err:          ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "DuplicateKey",
			err:          ErrDuplicateKey,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "WrongFormat",
			err:          ErrWrongFormat,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "NotFunds",
			err:          ErrNotFunds,
			expectedCode: http.StatusPaymentRequired,
		},
		{
			name:         "ManyRequests",
			err:          ErrManyRequests,
			expectedCode: http.StatusTooManyRequests,
		},
		{
			name:         "UnknownError",
			err:          errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем тестовый контекст
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Вызываем тестируемую функцию
			RespondWithError(c, tc.err)

			// Проверяем статус ответа
			assert.Equal(t, tc.expectedCode, w.Code)

			// Проверяем тело ответа
			var response ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, response.Code)
			assert.Contains(t, response.Error, tc.err.Error())
		})
	}
}
