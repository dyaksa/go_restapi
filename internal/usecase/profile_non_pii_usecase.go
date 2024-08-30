package usecase

import (
	"context"
	"database/sql"
	"golang_restapi/internal/delivery/http/web"
	"golang_restapi/internal/dto"
	"golang_restapi/internal/entity"
	"golang_restapi/internal/repository"
	appError "golang_restapi/pkg/errors"
	"golang_restapi/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var _ ProfileNonPII = &ProfileNonPIIUseCase{}

type ProfileNonPII interface {
	Create(ctx context.Context, request web.ProfileRequest) (*web.ProfileRequest, *appError.AppError)
	FetchProfile(ctx context.Context, id string) (*dto.FetchProfileDto, *appError.AppError)
	FindAll(ctx context.Context, pagination utils.Pagination, params web.ProfileQueryParam) ([]*dto.FetchProfileDto, *appError.AppError)
	Update(ctx context.Context, id string, request web.ProfileRequest) (*web.ProfileRequest, *appError.AppError)
}

type ProfileNonPIIUseCase struct {
	db         *sql.DB
	validator  *validator.Validate
	repository *repository.ProfileNonPIIRepository
}

func NewProfileNonPIIUseCase(db *sql.DB, repository *repository.ProfileNonPIIRepository, validator *validator.Validate) *ProfileNonPIIUseCase {
	return &ProfileNonPIIUseCase{
		db:         db,
		validator:  validator,
		repository: repository,
	}
}

func (uc *ProfileNonPIIUseCase) Create(ctx context.Context, request web.ProfileRequest) (*web.ProfileRequest, *appError.AppError) {
	var entity entity.ProfileNonPII
	err := uc.validator.Struct(request)
	if err != nil {
		return nil, appError.NewAppError(http.StatusBadRequest, appError.ErrorDetail{
			Field:          "validation",
			ValidationCode: string(appError.ValidationError),
			Param:          "",
			Message:        err.Error(),
		})
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "db",
			ValidationCode: string(appError.ServerError),
			Param:          "",
			Message:        "error begin transaction",
		})
	}

	defer tx.Rollback()

	entity.ID = uuid.New()
	entity.Nik = request.Nik
	entity.Name = request.Name
	entity.Phone = request.Phone
	entity.Email = request.Email

	err = uc.repository.Create(ctx, tx, entity)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "repository",
			ValidationCode: string(appError.CommitError),
			Param:          "",
			Message:        err.Error(),
		})
	}

	if err = tx.Commit(); err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "db",
			ValidationCode: string(appError.CommitError),
			Param:          "",
			Message:        "error commit transaction",
		})
	}

	return &web.ProfileRequest{
		ID:    entity.ID,
		Nik:   entity.Nik,
		Name:  entity.Name,
		Phone: entity.Phone,
		Email: entity.Email,
	}, nil
}

func (uc *ProfileNonPIIUseCase) FetchProfile(ctx context.Context, id string) (*dto.FetchProfileDto, *appError.AppError) {
	entity, err := uc.repository.FetchProfile(ctx, uc.db, id)
	if err != nil {
		return nil, appError.NewAppError(http.StatusNotFound, appError.ErrorDetail{
			Field:          "repository",
			ValidationCode: string(appError.NotFound),
			Param:          "",
			Message:        "error fetch data",
		})
	}

	return &dto.FetchProfileDto{
		ID:    entity.ID,
		Nik:   entity.Nik,
		Name:  entity.Name,
		Phone: entity.Phone,
		Email: entity.Email,
	}, nil
}

func (uc *ProfileNonPIIUseCase) FindAll(ctx context.Context, pagination utils.Pagination, params web.ProfileQueryParam) ([]*dto.FetchProfileDto, *appError.AppError) {
	entities, err := uc.repository.FindAll(ctx, uc.db, pagination, params)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "repository",
			ValidationCode: string(appError.FetchData),
			Param:          "",
			Message:        "error fetch data",
		})
	}

	var result []*dto.FetchProfileDto
	for _, entity := range entities {
		result = append(result, &dto.FetchProfileDto{
			ID:    entity.ID,
			Nik:   entity.Nik,
			Name:  entity.Name,
			Phone: entity.Phone,
			Email: entity.Email,
		})
	}

	return result, nil
}

func (uc *ProfileNonPIIUseCase) Update(ctx context.Context, id string, request web.ProfileRequest) (*web.ProfileRequest, *appError.AppError) {
	var entity entity.ProfileNonPII
	err := uc.validator.Struct(request)
	if err != nil {
		return nil, appError.NewAppError(http.StatusBadRequest, appError.ErrorDetail{
			Field:          "validation",
			ValidationCode: string(appError.ValidationError),
			Param:          "",
			Message:        err.Error(),
		})
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "db",
			ValidationCode: string(appError.ServerError),
			Param:          "",
			Message:        "error begin transaction",
		})
	}

	defer tx.Rollback()

	entity.ID = uuid.MustParse(id)
	entity.Nik = request.Nik
	entity.Name = request.Name
	entity.Phone = request.Phone
	entity.Email = request.Email

	err = uc.repository.Update(ctx, tx, entity)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "repository",
			ValidationCode: string(appError.CommitError),
			Param:          "",
			Message:        "error when updating profile",
		})
	}

	if err = tx.Commit(); err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "db",
			ValidationCode: string(appError.CommitError),
			Param:          "",
			Message:        "error commit transaction",
		})
	}

	return &web.ProfileRequest{
		ID:    entity.ID,
		Nik:   entity.Nik,
		Name:  entity.Name,
		Phone: entity.Phone,
		Email: entity.Email,
	}, nil
}
