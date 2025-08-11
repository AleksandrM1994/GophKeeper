package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/mocks"
	api "github.com/GophKeeper/pkg/api"
)

func TestRegisterUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  *api.RegisterUserRequest
		expectedCode int
		serviceError error
	}{
		{
			name:         "successful_registration",
			requestBody:  &api.RegisterUserRequest{Login: "test", Password: "password"},
			expectedCode: http.StatusOK,
			serviceError: nil,
		},
		{
			name:         "service_error",
			requestBody:  &api.RegisterUserRequest{Login: "test", Password: "password"},
			expectedCode: http.StatusInternalServerError,
			serviceError: errors.New("service error"),
		},
		{
			name:         "empty_request",
			requestBody:  &api.RegisterUserRequest{},
			expectedCode: http.StatusBadRequest,
			serviceError: custom_errs.ErrValidate,
		},
		{
			name:         "short_password",
			requestBody:  &api.RegisterUserRequest{Login: "test", Password: "123"},
			expectedCode: http.StatusBadRequest,
			serviceError: custom_errs.ErrBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)
			logger, _ := zap.NewDevelopment()
			lg := *logger.Sugar()

			// Настройка ожиданий для сервиса
			if tt.serviceError == nil {
				service.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			} else {
				service.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.serviceError).Times(1)
			}

			controller := &UserController{
				service: service,
				lg:      &lg,
			}

			r := gin.Default()
			r.POST("/register", controller.RegisterUserHandler)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusOK {
				var errResp custom_errs.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Code)
			}
		})
	}
}
