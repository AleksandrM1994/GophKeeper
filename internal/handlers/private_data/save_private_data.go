package private_data

import (
	"net/http"

	"github.com/gin-gonic/gin"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service/private_data/dto"
)

type SavePrivateDataRequest struct {
	Data []byte                     `json:"data"`
	Type repository.PrivateDataType `json:"type"`
}

func (c *PrivateDataController) SavePrivateData(ctx *gin.Context) {
	var req *SavePrivateDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, custom_errs.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := c.service.SavePrivateData(ctx, &dto.SavePrivateDataRequest{
		Data: req.Data,
		Type: req.Type,
	})
	if err != nil {
		custom_errs.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
