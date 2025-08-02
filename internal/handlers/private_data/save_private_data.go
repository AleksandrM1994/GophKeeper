package private_data

import (
	"net/http"

	"github.com/gin-gonic/gin"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/service/private_data/dto"
	api "github.com/GophKeeper/pkg/api"
)

func (c *PrivateDataController) SavePrivateData(ctx *gin.Context) {
	value, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, custom_errs.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: "empty user id",
		})
		return
	}

	userID := value.(string)

	var req api.SavePrivateDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, custom_errs.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	c.lg.Infow("server save private data request", "req", req)

	err := c.privateDataService.SavePrivateData(ctx, &dto.SavePrivateDataRequest{
		Data:   req.Data,
		Type:   FromProto(req.Type),
		UserID: userID,
		Nonce:  req.Nonce,
	})
	if err != nil {
		custom_errs.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func FromProto(in api.PrivateDataType) repository.PrivateDataType {
	switch in {
	case api.PrivateDataType_TEXT:
		return repository.PrivateDataTypeText
	case api.PrivateDataType_FILE:
		return repository.PrivateDataTypeFile
	case api.PrivateDataType_AUTH:
		return repository.PrivateDataTypeAuth
	case api.PrivateDataType_BANK:
		return repository.PrivateDataTypeBank
	default:
		return repository.PrivateDataTypeUnknown
	}
}
