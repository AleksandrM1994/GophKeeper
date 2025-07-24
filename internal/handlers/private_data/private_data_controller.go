package private_data

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/middlewares"
	"github.com/GophKeeper/internal/service/private_data"
	"github.com/GophKeeper/internal/service/user"
)

type PrivateDataController struct {
	cfg                config.Config
	lg                 *zap.SugaredLogger
	userService        *user.UserServiceImpl
	privateDataService *private_data.PrivateDataServiceImpl
}

func NewController(
	cfg config.Config,
	logger *zap.SugaredLogger,
	userService *user.UserServiceImpl,
	privateDataService *private_data.PrivateDataServiceImpl,
) *PrivateDataController {
	return &PrivateDataController{
		cfg:                cfg,
		lg:                 logger,
		userService:        userService,
		privateDataService: privateDataService,
	}
}

func (c *PrivateDataController) RegisterRoutes(r *gin.Engine) {
	privateDataGroup := r.Group("/api/private-data").Use(middlewares.Authorizer(c.lg, c.cfg, c.userService))
	privateDataGroup.POST("/save", c.SavePrivateData)
}
