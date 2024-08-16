package service

import (
	"context"
	"database/sql"
	"golang_restapi/dto"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"
	"golang_restapi/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProfileNonPIIService interface {
	Create(ctx context.Context, request web.ProfileRequest) (web.ProfileRequest, error)
	FetchProfile(ctx context.Context, id uuid.UUID) (*dto.FetchProfileDto, error)
	FindAll(ctx context.Context, pagination utils.Pagination, params dto.ParamsListProfile) ([]*dto.FetchProfileDto, error)
	Update(ctx context.Context, id uuid.UUID, request web.ProfileRequest) (web.ProfileRequest, error)
}

type ProfileNonPIIServiceImpl struct {
	db         *sql.DB
	validator  *validator.Validate
	repository repository.ProfileNonPII
}

func NewProfileNonPIIService(db *sql.DB, repository repository.ProfileNonPII, validator *validator.Validate) *ProfileNonPIIServiceImpl {
	return &ProfileNonPIIServiceImpl{
		db:         db,
		validator:  validator,
		repository: repository,
	}
}

func (service *ProfileNonPIIServiceImpl) Create(ctx context.Context, request web.ProfileRequest) (p web.ProfileRequest, err error) {
	err = service.validator.Struct(request)
	if err != nil {
		return
	}

	tx, err := service.db.Begin()
	if err != nil {
		return
	}

	defer helper.CommitAndRollbackError(tx)

	profile := entity.ProfileNonPII{
		ID:    uuid.New(),
		Nik:   request.Nik,
		Name:  request.Name,
		Phone: request.Phone,
		Email: request.Email,
	}

	err = service.repository.Create(ctx, tx, profile)
	p = request
	return
}

func (service *ProfileNonPIIServiceImpl) FetchProfile(ctx context.Context, id uuid.UUID) (res *dto.FetchProfileDto, err error) {
	tx, err := service.db.Begin()
	if err != nil {
		return
	}

	defer helper.CommitAndRollbackError(tx)

	profile, err := service.repository.FetchProfile(ctx, id, tx)
	if err != nil {
		return
	}

	res = &dto.FetchProfileDto{
		ID:    profile.ID,
		Nik:   profile.Nik,
		Name:  profile.Name,
		Phone: profile.Phone,
		Email: profile.Email,
	}

	return
}

func (service *ProfileNonPIIServiceImpl) FindAll(ctx context.Context, pagination utils.Pagination, params dto.ParamsListProfile) (p []*dto.FetchProfileDto, err error) {
	tx, err := service.db.Begin()
	if err != nil {
		return
	}

	defer helper.CommitAndRollbackError(tx)

	profiles, err := service.repository.FindAll(ctx, pagination, params, tx)
	if err != nil {
		return
	}

	for _, profile := range profiles {
		p = append(p, &dto.FetchProfileDto{
			ID:    profile.ID,
			Nik:   profile.Nik,
			Name:  profile.Name,
			Phone: profile.Phone,
			Email: profile.Email,
		})
	}

	return
}

func (service *ProfileNonPIIServiceImpl) Update(ctx context.Context, id uuid.UUID, request web.ProfileRequest) (p web.ProfileRequest, err error) {
	err = service.validator.Struct(request)
	if err != nil {
		return
	}

	tx, err := service.db.Begin()
	if err != nil {
		return
	}

	defer helper.CommitAndRollbackError(tx)

	profile := entity.ProfileNonPII{
		ID:    id,
		Nik:   request.Nik,
		Name:  request.Name,
		Phone: request.Phone,
		Email: request.Email,
	}

	err = service.repository.Update(ctx, tx, profile)
	p = request
	return
}
