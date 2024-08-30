package usecase

import (
	"context"
	"database/sql"
	"golang_restapi/internal/delivery/http/web"
	"golang_restapi/internal/entity"
	"golang_restapi/internal/repository"
	appError "golang_restapi/pkg/errors"
	"golang_restapi/utils"
	"net/http"

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/dyaksa/encryption-pii/crypto/aesx"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var _ Profile = &ProfileUseCase{}

type Profile interface {
	Create(ctx context.Context, req *web.ProfileRequest) (*web.ProfileResponse, *appError.AppError)
	FetchProfile(ctx context.Context, id string) (*web.ProfileResponse, *appError.AppError)
	FindAll(ctx context.Context, pagination utils.Pagination, params web.ProfileQueryParam) ([]*web.ProfileResponse, *appError.AppError)
	Update(ctx context.Context, id string, req *web.ProfileRequest) (*web.ProfileResponse, *appError.AppError)
}

type ProfileUseCase struct {
	db                *sql.DB
	crypto            *crypto.Crypto
	profileRepository *repository.ProfileRepository
	validation        *validator.Validate
}

func NewProfileUseCase(db *sql.DB, crypto *crypto.Crypto, validation *validator.Validate, profileRepository *repository.ProfileRepository) *ProfileUseCase {
	return &ProfileUseCase{db: db, crypto: crypto, validation: validation, profileRepository: profileRepository}
}

func (uc *ProfileUseCase) Create(ctx context.Context, req *web.ProfileRequest) (*web.ProfileResponse, *appError.AppError) {
	if err := uc.validation.Struct(req); err != nil {
		return nil, appError.NewAppError(http.StatusBadRequest, appError.ErrorDetail{
			Field:          "validation",
			ValidationCode: string(appError.ValidationError),
			Param:          "",
			Message:        err.Error(),
		})
	}

	tx, err := uc.db.Begin()
	defer tx.Rollback()

	var profile entity.Profile

	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "db",
			ValidationCode: string(appError.ServerError),
			Param:          "",
			Message:        "error begin transaction",
		})
	}

	profile.ID = uuid.New()
	profile.Nik = uc.crypto.Encrypt(req.Nik, aesx.AesCBC)
	profile.Name = uc.crypto.Encrypt(req.Name, aesx.AesCBC)
	profile.Phone = uc.crypto.Encrypt(req.Phone, aesx.AesCBC)
	profile.Email = uc.crypto.Encrypt(req.Email, aesx.AesCBC)

	err = uc.crypto.BindHeap(&profile)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "crypto",
			ValidationCode: string(appError.BindHeap),
			Param:          "",
			Message:        "error bind heap",
		})
	}

	err = uc.profileRepository.Create(ctx, tx, profile)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "profileRepository",
			ValidationCode: string(appError.InsertData),
			Param:          "",
			Message:        "error insert data",
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

	return &web.ProfileResponse{
		ID:    profile.ID,
		Nik:   req.Nik,
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}, nil
}

func (uc *ProfileUseCase) FetchProfile(ctx context.Context, id string) (*web.ProfileResponse, *appError.AppError) {
	profile, err := uc.profileRepository.FetchProfile(ctx, uc.db, id, func(p *entity.Profile) {
		p.Nik = uc.crypto.Decrypt(aesx.AesCBC)
		p.Name = uc.crypto.Decrypt(aesx.AesCBC)
		p.Phone = uc.crypto.Decrypt(aesx.AesCBC)
		p.Email = uc.crypto.Decrypt(aesx.AesCBC)
	})

	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "profileRepository",
			ValidationCode: string(appError.FetchData),
			Param:          "",
			Message:        "error fetch data",
		})
	}

	return &web.ProfileResponse{
		ID:    profile.ID,
		Nik:   profile.Nik.To(),
		Name:  profile.Name.To(),
		Phone: profile.Phone.To(),
		Email: profile.Email.To(),
	}, nil
}

func (us *ProfileUseCase) FindAll(ctx context.Context, pagination utils.Pagination, params web.ProfileQueryParam) ([]*web.ProfileResponse, *appError.AppError) {
	var prs []*web.ProfileResponse
	_, err := us.profileRepository.FindAll(ctx, us.db, pagination, params, func(p *entity.Profile) {
		p.Nik = us.crypto.Decrypt(aesx.AesCBC)
		p.Name = us.crypto.Decrypt(aesx.AesCBC)
		p.Phone = us.crypto.Decrypt(aesx.AesCBC)
		p.Email = us.crypto.Decrypt(aesx.AesCBC)
	}, func(p entity.Profile) {
		prs = append(prs, &web.ProfileResponse{
			ID:    p.ID,
			Nik:   p.Nik.To(),
			Name:  p.Name.To(),
			Phone: p.Phone.To(),
			Email: p.Email.To(),
		})
	})

	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "profileRepository",
			ValidationCode: string(appError.FetchData),
			Param:          "",
			Message:        "error fetch data",
		})
	}

	return prs, nil
}

func (uc *ProfileUseCase) Update(ctx context.Context, id string, req *web.ProfileRequest) (*web.ProfileResponse, *appError.AppError) {
	if err := uc.validation.Struct(req); err != nil {
		return nil, appError.NewAppError(http.StatusBadRequest, appError.ErrorDetail{
			Field:          "validation",
			ValidationCode: string(appError.ValidationError),
			Param:          "",
			Message:        err.Error(),
		})
	}

	tx, err := uc.db.Begin()
	defer tx.Rollback()

	var profile entity.Profile

	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "db",
			ValidationCode: string(appError.ServerError),
			Param:          "",
			Message:        "error begin transaction",
		})
	}

	profile.Nik = uc.crypto.Encrypt(req.Nik, aesx.AesCBC)
	profile.Name = uc.crypto.Encrypt(req.Name, aesx.AesCBC)
	profile.Phone = uc.crypto.Encrypt(req.Phone, aesx.AesCBC)
	profile.Email = uc.crypto.Encrypt(req.Email, aesx.AesCBC)

	err = uc.crypto.BindHeap(&profile)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "crypto",
			ValidationCode: string(appError.BindHeap),
			Param:          "",
			Message:        "error bind heap",
		})
	}

	err = uc.profileRepository.Update(ctx, tx, id, profile)
	if err != nil {
		return nil, appError.NewAppError(http.StatusInternalServerError, appError.ErrorDetail{
			Field:          "profileRepository",
			ValidationCode: string(appError.UpdateData),
			Param:          "",
			Message:        "error update data",
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

	return &web.ProfileResponse{
		ID:    uuid.MustParse(id),
		Nik:   req.Nik,
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}, nil
}
