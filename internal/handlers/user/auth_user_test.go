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
	"github.com/GophKeeper/internal/service/user/dto"
	api "github.com/GophKeeper/pkg/api"
)

func TestAuthUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  *api.AuthUserRequest
		expectedCode int
		expectedJWT  string
		serviceError error
	}{
		{
			name:         "successful_auth",
			requestBody:  &api.AuthUserRequest{Login: "test", Password: "password"},
			expectedCode: http.StatusOK,
			expectedJWT:  "valid-jwt-token",
			serviceError: nil,
		},
		{
			name:         "service_error",
			requestBody:  &api.AuthUserRequest{Login: "test", Password: "password"},
			expectedCode: http.StatusInternalServerError,
			expectedJWT:  "",
			serviceError: errors.New("service error"),
		},
		{
			name:         "empty_request",
			requestBody:  &api.AuthUserRequest{},
			expectedCode: http.StatusBadRequest,
			expectedJWT:  "",
			serviceError: custom_errs.ErrValidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)
			logger, _ := zap.NewDevelopment()
			lg := *logger.Sugar()

			if tt.serviceError == nil {
				service.EXPECT().AuthUser(gomock.Any(), gomock.Any()).Return(&dto.AuthUserResponse{
					JWT: tt.expectedJWT,
				}, nil).Times(1)
			} else {
				service.EXPECT().AuthUser(gomock.Any(), gomock.Any()).Return(nil, tt.serviceError).Times(1)
			}

			controller := &UserController{
				service: service,
				lg:      &lg,
			}

			r := gin.Default()
			r.POST("/auth", controller.AuthUserHandler)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var resp api.AuthUserResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedJWT, resp.Jwt)
			} else {
				var errResp custom_errs.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Code)
			}
		})
	}
}
