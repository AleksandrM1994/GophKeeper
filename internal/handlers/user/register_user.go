package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	custom_errs "github.com/GophKeeper/internal/errors"
	"github.com/GophKeeper/internal/service/user/dto"
)

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c *UserController) RegisterUserHandler(ctx *gin.Context) {
	var req *RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, custom_errs.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := c.service.CreateUser(ctx, &dto.CreateUserRequest{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		custom_errs.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
