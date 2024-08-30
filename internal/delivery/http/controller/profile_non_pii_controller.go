package controller

import (
	"golang_restapi/internal/delivery/http/web"
	"golang_restapi/internal/helper"
	"golang_restapi/internal/usecase"
	appError "golang_restapi/pkg/errors"
	"golang_restapi/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProfileNonPIIController struct {
	ProfileNonPIIUsecase *usecase.ProfileNonPIIUseCase
}

func NewProfileNonPIIController(profileNonPIIUsecase *usecase.ProfileNonPIIUseCase) *ProfileNonPIIController {
	return &ProfileNonPIIController{ProfileNonPIIUsecase: profileNonPIIUsecase}
}

func (c *ProfileNonPIIController) Create(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	var profileRequest web.ProfileRequest

	if err := ctx.Bind(&profileRequest); err != nil {
		helper.ValidationErrorResponse(ctx, err)
		return nil
	}

	res, err := c.ProfileNonPIIUsecase.Create(ctx.Request().Context(), profileRequest)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusCreated, res)
	return nil
}

func (c *ProfileNonPIIController) FetchProfile(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	userId := ctx.Param("id")

	res, err := c.ProfileNonPIIUsecase.FetchProfile(ctx.Request().Context(), userId)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusOK, res)
	return nil
}

func (c *ProfileNonPIIController) FindAll(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	var params web.ProfileQueryParam

	pagination := utils.Pagination{Page: params.Page, Limit: params.Limit}

	if err := ctx.Bind(&params); err != nil {
		helper.ValidationErrorResponse(ctx, err)
		return nil
	}

	res, err := c.ProfileNonPIIUsecase.FindAll(ctx.Request().Context(), pagination, params)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusOK, res)
	return nil
}

func (c *ProfileNonPIIController) Update(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	var profileRequest web.ProfileRequest
	var id = ctx.Param("id")

	if err := ctx.Bind(&profileRequest); err != nil {
		helper.ValidationErrorResponse(ctx, err)
		return nil
	}

	res, err := c.ProfileNonPIIUsecase.Update(ctx.Request().Context(), id, profileRequest)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusOK, res)
	return nil
}
