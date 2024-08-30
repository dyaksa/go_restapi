package helper

import (
	"errors"
	"golang_restapi/internal/delivery/http/web"
	appError "golang_restapi/pkg/errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func SuccessResponse(ctx echo.Context, code int, data any) {
	ctx.JSON(code, web.ApiResponse{Code: code, Data: data})
}

func ErrorResponse(ctx echo.Context, code int, details any) {
	ctx.JSON(code, web.ApiResponse{Code: code, Details: details})
}

func ErrorBindingResponse(ctx echo.Context, validationErrors validator.ValidationErrors) {
	var errOut []appError.ErrorDetail
	for _, err := range validationErrors {
		errOut = append(errOut, appError.ErrorDetail{
			Field:          err.Field(),
			ValidationCode: err.Tag(),
			Param:          err.Param(),
			Message:        err.Error(),
		})
	}
	ctx.JSON(http.StatusUnprocessableEntity, web.ApiResponse{Details: errOut})
}

func ValidationErrorResponse(ctx echo.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		ErrorBindingResponse(ctx, validationErrors)
	} else {
		var outError []appError.ErrorDetail

		outError = append(outError, appError.ErrorDetail{
			Field:          "general",
			ValidationCode: string(appError.InvalidJson),
			Message:        err.Error(),
			Param:          "",
		})

		ErrorResponse(ctx, http.StatusBadRequest, outError)
	}
}
