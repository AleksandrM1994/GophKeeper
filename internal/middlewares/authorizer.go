package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/service/user"
)

func Authorizer(lg *zap.SugaredLogger, cfg config.Config, srv *user.UserServiceImpl) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lg.Info("START CHECK AUTHORIZATION")

		ctx.JSON(http.StatusUnauthorized, errors.ErrorResponse{
			Code: http.StatusUnauthorized,
			// Error: errCookie.Error(),
		})
	}
}
