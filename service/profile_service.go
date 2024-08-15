package service

import (
	"context"
	"database/sql"
	"fmt"
	"golang_restapi/dto"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"
	"golang_restapi/utils"

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/dyaksa/encryption-pii/crypto/aesx"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProfileService interface {
	Create(ctx context.Context, request web.ProfileRequest) (web.ProfileRequest, error)
	FetchProfile(ctx context.Context, id uuid.UUID) (*dto.FetchProfileDto, error)
	FindAll(ctx context.Context, pagination utils.Pagination, params dto.ParamsListProfile) ([]*dto.FetchProfileDto, error)
	Update(ctx context.Context, id uuid.UUID, request web.ProfileRequest) (web.ProfileRequest, error)
}

type ProfileServiceImpl struct {
	db         *sql.DB
	repository repository.Profile
	validator  *validator.Validate
	crypto     *crypto.Crypto
}

func NewProfileService(db *sql.DB, repository repository.Profile, validator *validator.Validate, crypto *crypto.Crypto) *ProfileServiceImpl {
	return &ProfileServiceImpl{
		db:         db,
		repository: repository,
		validator:  validator,
		crypto:     crypto,
	}
}

func (service *ProfileServiceImpl) Create(ctx context.Context, request web.ProfileRequest) (p web.ProfileRequest, err error) {
	err = service.validator.Struct(request)
	if err != nil {
		err = fmt.Errorf("error validation: %w", err)
		return
	}

	tx, err := service.db.Begin()
	if err != nil {
		return
	}

	defer helper.CommitAndRollbackError(tx)

	profile := entity.Profile{
		ID:    uuid.New(),
		Nik:   service.crypto.Encrypt(request.Nik, aesx.AesCBC),
		Name:  service.crypto.Encrypt(request.Name, aesx.AesCBC),
		Phone: service.crypto.Encrypt(request.Phone, aesx.AesCBC),
		Email: service.crypto.Encrypt(request.Email, aesx.AesCBC),
	}

	if err = service.crypto.BindHeap(&profile); err != nil {
		err = fmt.Errorf("error binding heap: %w", err)
		return
	}

	if err = service.repository.Create(ctx, tx, profile); err != nil {
		err = fmt.Errorf("error when inserting profile: %w", err)
		return
	}

	request.ID = profile.ID
	return request, nil
}

func (service *ProfileServiceImpl) FetchProfile(ctx context.Context, id uuid.UUID) (pr *dto.FetchProfileDto, err error) {
	tx, err := service.db.Begin()
	if err != nil {
		return
	}

	spr, err := service.repository.FetchProfile(ctx, id, tx, func(fpr *entity.Profile) {
		fpr.Nik = service.crypto.Decrypt(aesx.AesCBC)
		fpr.Name = service.crypto.Decrypt(aesx.AesCBC)
		fpr.Phone = service.crypto.Decrypt(aesx.AesCBC)
		fpr.Email = service.crypto.Decrypt(aesx.AesCBC)
	})
	helper.PanicIf(err)

	pr = &dto.FetchProfileDto{
		ID:    id,
		Nik:   spr.Nik.To(),
		Name:  spr.Name.To(),
		Phone: spr.Phone.To(),
		Email: spr.Email.To(),
	}

	return
}

func (service *ProfileServiceImpl) FindAll(ctx context.Context, pagination utils.Pagination, params dto.ParamsListProfile) (prs []*dto.FetchProfileDto, err error) {
	tx, err := service.db.Begin()
	if err != nil {
		return prs, err
	}

	_, err = service.repository.FindAll(ctx, pagination, params, tx, service.crypto, func(fpr *entity.Profile) {
		fpr.Nik = service.crypto.Decrypt(aesx.AesCBC)
		fpr.Name = service.crypto.Decrypt(aesx.AesCBC)
		fpr.Phone = service.crypto.Decrypt(aesx.AesCBC)
		fpr.Email = service.crypto.Decrypt(aesx.AesCBC)
	}, func(p entity.Profile) {
		pr := &dto.FetchProfileDto{
			ID:    p.ID,
			Nik:   p.Nik.To(),
			Name:  p.Name.To(),
			Phone: p.Phone.To(),
			Email: p.Email.To(),
		}
		prs = append(prs, pr)
	})

	if err != nil {
		return prs, err
	}
	return prs, nil
}

func (service *ProfileServiceImpl) Update(ctx context.Context, id uuid.UUID, request web.ProfileRequest) (web.ProfileRequest, error) {
	err := service.validator.Struct(request)
	helper.PanicIf(err)

	tx, err := service.db.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	profile := entity.Profile{
		ID:    id,
		Nik:   service.crypto.Encrypt(request.Nik, aesx.AesCBC),
		Name:  service.crypto.Encrypt(request.Name, aesx.AesCBC),
		Phone: service.crypto.Encrypt(request.Phone, aesx.AesCBC),
		Email: service.crypto.Encrypt(request.Email, aesx.AesCBC),
	}

	err = service.crypto.BindHeap(&profile)
	if err != nil {
		return request, err
	}

	err = service.repository.Update(ctx, tx, profile)
	if err != nil {
		return request, err
	}

	return request, nil
}
