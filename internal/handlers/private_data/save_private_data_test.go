package private_data

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

	"github.com/GophKeeper/config"
	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/middlewares"
	"github.com/GophKeeper/internal/mocks"
	"github.com/GophKeeper/internal/service/user/dto"
	api "github.com/GophKeeper/pkg/api"
)

func performRequest(router *gin.Engine, method, path string, headers map[string]string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestSavePrivateData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		userID       string
		requestBody  *api.PrivateDataSaved
		JWTToken     string
		expectedCode int
	}{
		{
			name:         "missing user_id in context",
			userID:       "",
			requestBody:  nil,
			JWTToken:     "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid json request body",
			userID:       "user-123",
			requestBody:  nil, // Пустое тело для проверки
			JWTToken:     "fake-token",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "successful save",
			userID: "user-123",
			requestBody: &api.PrivateDataSaved{
				Id:    "1",
				Type:  api.PrivateDataSaved_TEXT_TYPE,
				Data:  []byte("test"),
				Nonce: []byte("nonce"),
				Login: "login",
			},
			JWTToken:     "fake-token",
			expectedCode: http.StatusOK,
		},
		{
			name:   "service error",
			userID: "user-123",
			requestBody: &api.PrivateDataSaved{
				Id:    "1",
				Type:  api.PrivateDataSaved_TEXT_TYPE,
				Data:  []byte("test"),
				Nonce: []byte("nonce"),
				Login: "login",
			},
			JWTToken:     "fake-token",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svcPd := mocks.NewMockPrivateDataService(ctrl)
			svcUser := mocks.NewMockUserService(ctrl)

			logger, _ := zap.NewDevelopment()
			lg := *logger.Sugar()

			if tt.expectedCode == http.StatusOK {
				svcPd.EXPECT().SavePrivateData(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				svcUser.EXPECT().
					CheckAuthUser(gomock.Any(), &dto.CheckAuthRequest{JWT: tt.JWTToken}).
					Return(&dto.CheckAuthResponse{UserID: "user-123"}, nil).Times(1)
			} else if tt.expectedCode == http.StatusBadRequest {
				svcUser.EXPECT().
					CheckAuthUser(gomock.Any(), &dto.CheckAuthRequest{JWT: tt.JWTToken}).
					Return(&dto.CheckAuthResponse{UserID: "user-123"}, nil).Times(1)
			} else if tt.expectedCode == http.StatusInternalServerError {
				svcPd.EXPECT().SavePrivateData(gomock.Any(), gomock.Any()).Return(errors.New("service error")).Times(1)
				svcUser.EXPECT().
					CheckAuthUser(gomock.Any(), &dto.CheckAuthRequest{JWT: tt.JWTToken}).
					Return(&dto.CheckAuthResponse{UserID: "user-123"}, nil).Times(1)
			} else {
				svcUser.EXPECT().
					CheckAuthUser(gomock.Any(), &dto.CheckAuthRequest{JWT: tt.JWTToken}).
					Return(nil, custom_errs.ErrUnauthorized).Times(1)
			}

			cfg, err := config.NewConfig()
			if err != nil {
				t.Fatalf("Failed to create config: %v", err)
			}
			controller := NewController(cfg, &lg, svcUser, svcPd)

			r := gin.Default()
			privateDataGroup := r.Group("/api/private-data")
			privateDataGroup.Use(middlewares.Authorizer(&lg, cfg, svcUser))
			privateDataGroup.POST("/save", controller.SavePrivateData)

			headers := map[string]string{
				"Authorization": "Bearer " + tt.JWTToken,
				"Content-Type":  "application/json",
			}

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			w := performRequest(r, "POST", "/api/private-data/save", headers, body)
			assert.Equal(t, tt.expectedCode, w.Code)

			var resp custom_errs.ErrorResponse
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			if err == nil && tt.expectedCode != http.StatusOK {
				assert.Equal(t, tt.expectedCode, resp.Code)
				assert.NotEmpty(t, resp.Error)
			}
		})
	}
}
