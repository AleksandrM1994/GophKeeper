package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrUniqueViolation = errors.New("unique violation error")
var ErrNotFound = errors.New("not found")
var ErrNoContent = errors.New("not content")
var ErrValidate = errors.New("empty value")
var ErrUnauthorized = errors.New("unauthorized")
var ErrBadRequest = errors.New("bad request")
var ErrResourceGone = errors.New("resource gone")
var ErrDuplicateKey = errors.New("duplicate key")
var ErrWrongFormat = errors.New("wrong format")
var ErrNotFunds = errors.New("not enough funds")
var ErrManyRequests = errors.New("many requests")

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func RespondWithError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrValidate):
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
	case errors.Is(err, ErrUnauthorized):
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: err.Error(),
		})
	case errors.Is(err, ErrNotFound):
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Code:  http.StatusNotFound,
			Error: err.Error(),
		})
	case errors.Is(err, ErrDuplicateKey):
		ctx.JSON(http.StatusConflict, ErrorResponse{
			Code:  http.StatusConflict,
			Error: err.Error(),
		})
	case errors.Is(err, ErrWrongFormat):
		ctx.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Code:  http.StatusUnprocessableEntity,
			Error: err.Error(),
		})
	case errors.Is(err, ErrNoContent):
		ctx.JSON(http.StatusNoContent, ErrorResponse{
			Code:  http.StatusNoContent,
			Error: err.Error(),
		})
	case errors.Is(err, ErrNotFunds):
		ctx.JSON(http.StatusPaymentRequired, ErrorResponse{
			Code:  http.StatusPaymentRequired,
			Error: err.Error(),
		})
	case errors.Is(err, ErrManyRequests):
		ctx.JSON(http.StatusTooManyRequests, ErrorResponse{
			Code:  http.StatusTooManyRequests,
			Error: err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}
}
