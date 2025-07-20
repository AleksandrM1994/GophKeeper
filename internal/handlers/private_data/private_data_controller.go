package private_data

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"github.com/GophKeeper/config"
	"github.com/GophKeeper/internal/service/private_data"
)

type PrivateDataController struct {
	cfg     config.Config
	lg      *zap.SugaredLogger
	service *private_data.PrivateDataServiceImpl
}

func NewController(
	cfg config.Config,
	logger *zap.SugaredLogger,
	service *private_data.PrivateDataServiceImpl,
) *PrivateDataController {
	return &PrivateDataController{
		cfg:     cfg,
		lg:      logger,
		service: service,
	}
}

func (c *PrivateDataController) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/private-data/save", c.SavePrivateData)
}
