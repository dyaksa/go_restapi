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

type ProfileController struct {
	ProfileUsecase *usecase.ProfileUseCase
}

func NewProfileController(profileUsecase *usecase.ProfileUseCase) *ProfileController {
	return &ProfileController{ProfileUsecase: profileUsecase}
}

func (c *ProfileController) Create(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	var profileRequest web.ProfileRequest

	if err := ctx.Bind(&profileRequest); err != nil {
		helper.ValidationErrorResponse(ctx, err)
		return nil
	}

	res, err := c.ProfileUsecase.Create(ctx.Request().Context(), &profileRequest)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusCreated, res)
	return nil
}

func (c *ProfileController) FetchProfile(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	userId := ctx.Param("id")

	res, err := c.ProfileUsecase.FetchProfile(ctx.Request().Context(), userId)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusOK, res)
	return nil
}

func (c *ProfileController) FindAll(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	var params web.ProfileQueryParam

	pagination := utils.Pagination{Page: params.Page, Limit: params.Limit}

	if err := ctx.Bind(&params); err != nil {
		helper.ValidationErrorResponse(ctx, err)
		return nil
	}

	res, err := c.ProfileUsecase.FindAll(ctx.Request().Context(), pagination, params)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusOK, res)
	return nil
}

func (c *ProfileController) Update(ctx echo.Context) error {
	var errOut []appError.ErrorDetail
	var profileRequest web.ProfileRequest
	var id = ctx.Param("id")

	if err := ctx.Bind(&profileRequest); err != nil {
		helper.ValidationErrorResponse(ctx, err)
		return nil
	}

	res, err := c.ProfileUsecase.Update(ctx.Request().Context(), id, &profileRequest)
	if err != nil {
		errOut = append(errOut, err.Details)
		helper.ErrorResponse(ctx, err.Code, errOut)
		return nil
	}

	helper.SuccessResponse(ctx, http.StatusOK, res)
	return nil
}
