package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/service/user/dto"
)

type AuthUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUserResponse struct {
	JWT string `json:"jwt"`
}

func (c *UserController) AuthUserHandler(ctx *gin.Context) {
	var req *AuthUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, custom_errs.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	res, err := c.service.AuthUser(ctx, &dto.AuthUserRequest{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		custom_errs.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &AuthUserResponse{JWT: res.JWT})
}
