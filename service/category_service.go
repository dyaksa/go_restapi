package service

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"golang_restapi/exception"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"
)

type CategoryService interface {
	Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse
	Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse
	Delete(ctx context.Context, categoryId int)
	FindByID(ctx context.Context, categoryId int) web.CategoryResponse
	FindAll(ctx context.Context) []web.CategoryResponse
}

type CategoryServiceImpl struct {
	tx         *sql.DB
	repository repository.Category
	validator  *validator.Validate
}

func NewCategoryService(tx *sql.DB, repository repository.Category, validator *validator.Validate) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		tx:         tx,
		repository: repository,
		validator:  validator,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	err := service.validator.Struct(request)
	helper.PanicIf(err)

	tx, err := service.tx.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	category := entity.Category{
		Name: request.Name,
	}

	category = service.repository.Save(ctx, tx, category)

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.validator.Struct(request)
	helper.PanicIf(err)

	tx, err := service.tx.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	category := entity.Category{
		ID:   request.ID,
		Name: request.Name,
	}

	_, err = service.repository.FindOneByID(ctx, tx, category.ID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	category = service.repository.Update(ctx, tx, category)
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.tx.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	_, err = service.repository.FindOneByID(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.repository.Delete(ctx, tx, categoryId)
}

func (service *CategoryServiceImpl) FindByID(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.tx.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	category, err := service.repository.FindOneByID(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.tx.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	categories := service.repository.FindAll(ctx, tx)
	return helper.ToCategorySliceResponse(categories)
}
