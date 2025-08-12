package middlewares

import (
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
	"github.com/GophKeeper/internal/mocks"
	"github.com/GophKeeper/internal/service/user/dto"
)

func performRequest(router *gin.Engine, method, path, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestAuthorizerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop().Sugar()
	cfg := config.Config{}

	tests := []struct {
		name         string
		authHeader   string
		svcResponse  *dto.CheckAuthResponse
		svcError     error
		expectedCode int
		expectUserID interface{}
	}{
		{
			name:         "no Authorization header",
			authHeader:   "",
			svcResponse:  nil,
			svcError:     nil,
			expectedCode: http.StatusUnauthorized,
			expectUserID: nil,
		},
		{
			name:         "invalid prefix",
			authHeader:   "Token abc",
			svcResponse:  nil,
			svcError:     nil,
			expectedCode: http.StatusUnauthorized,
			expectUserID: nil,
		},
		{
			name:         "service returns error",
			authHeader:   "Bearer bad.token",
			svcResponse:  nil,
			svcError:     errors.New("invalid token"),
			expectedCode: http.StatusUnauthorized,
			expectUserID: nil,
		},
		{
			name:         "successful authorization",
			authHeader:   "Bearer good.token",
			svcResponse:  &dto.CheckAuthResponse{UserID: "42"},
			svcError:     nil,
			expectedCode: http.StatusOK,
			expectUserID: "42",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockUserService(ctrl)

			if tc.svcResponse != nil || tc.svcError != nil {
				token := tc.authHeader[len("Bearer "):]
				svc.
					EXPECT().
					CheckAuthUser(gomock.Any(), &dto.CheckAuthRequest{JWT: token}).
					Return(tc.svcResponse, tc.svcError)
			}

			router := gin.New()
			router.Use(Authorizer(logger, cfg, svc))
			router.GET("/test", func(ctx *gin.Context) {
				uid, _ := ctx.Get("user_id")
				ctx.JSON(http.StatusOK, gin.H{"user_id": uid})
			})

			w := performRequest(router, "GET", "/test", tc.authHeader)
			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				var body map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectUserID, body["user_id"])
			} else {
				var resp custom_errs.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, resp.Code)
				assert.NotEmpty(t, resp.Error)
			}
		})
	}
}
