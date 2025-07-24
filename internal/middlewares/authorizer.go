package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/service/user"
	"github.com/GophKeeper/internal/service/user/dto"
)

func Authorizer(lg *zap.SugaredLogger, cfg config.Config, srv *user.UserServiceImpl) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lg.Info("START CHECK AUTHORIZATION")

		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, custom_errs.ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: errors.New("authorization header is required").Error(),
			})
		}

		authToken, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, custom_errs.ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: errors.New("authorization header is invalid").Error(),
			})
		}

		res, err := srv.CheckAuthUser(ctx, &dto.CheckAuthRequest{
			JWT: authToken,
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, custom_errs.ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			})
		}
		ctx.Set("user_id", res.UserID)
	}
}
