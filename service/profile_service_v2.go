package service

import (
	"context"
	"database/sql"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/dyaksa/encryption-pii/crypto/aesx"
	"github.com/dyaksa/encryption-pii/crypto/query"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProfileServiceV2 interface {
	Create(ctx context.Context, request web.ProfileRequest) web.ProfileRequest
	FindByID(ctx context.Context, id uuid.UUID) *web.ProfileResponse
	FindInvoice(ctx context.Context, id uuid.UUID) *web.InvoiceResponse
}

type ProfileServiceV2Impl struct {
	crypto       *crypto.Crypto
	validator    *validator.Validate
	db           *sql.DB
	repositoryv2 *repository.ProfileRepositoryV2
}

func NewProfileServiceV2(db *sql.DB, crypto *crypto.Crypto, validator *validator.Validate) *ProfileServiceV2Impl {
	profile := repository.NewProfileRepositoryV2()
	return &ProfileServiceV2Impl{
		db:           db,
		crypto:       crypto,
		validator:    validator,
		repositoryv2: profile,
	}
}

func (service *ProfileServiceV2Impl) FindByID(ctx context.Context, id uuid.UUID) *web.ProfileResponse {
	tx, err := service.db.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	args := entity.FetchProfileParams{
		ID: id,
	}

	spr, err := service.repositoryv2.FetchProfile(ctx, args, tx, func(fprv *entity.FetchProfileRowV2) {
		fprv.Nik = service.crypto.Encrypt("", aesx.AesCBC)
		fprv.Name = service.crypto.Encrypt("", aesx.AesCBC)
		fprv.Phone = service.crypto.Encrypt("", aesx.AesCBC)
		fprv.Email = service.crypto.Encrypt("", aesx.AesCBC)
		fprv.DOB = service.crypto.Encrypt("", aesx.AesCBC)
	})
	helper.PanicIf(err)

	return &web.ProfileResponse{
		ID:    id,
		Nik:   spr.Nik.To(),
		Name:  spr.Name.To(),
		Phone: spr.Phone.To(),
		Email: spr.Email.To(),
		DOB:   spr.DOB.To(),
	}
}

func (service *ProfileServiceV2Impl) FindInvoice(ctx context.Context, id uuid.UUID) *web.InvoiceResponse {
	tx, err := service.db.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	sir, err := service.repositoryv2.FetchInvoice(ctx, tx, id, nil, func(fir *entity.FetchInvoiceRow) {
		fir.CustomerNameCrypt = service.crypto.Decrypt(fir.CustomerName, aesx.AesCBC)
		fir.CustomerAddrCrypt = service.crypto.Decrypt(fir.CustomerAddress, aesx.AesCBC)
	})
	helper.PanicIf(err)

	return &web.InvoiceResponse{
		ID:           id,
		CustomerName: sir.CustomerName,
	}

}

func (service *ProfileServiceV2Impl) Create(ctx context.Context, request web.ProfileRequest) web.ProfileRequest {
	err := service.validator.Struct(request)
	helper.PanicIf(err)

	tx, err := service.db.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	profile := entity.ProfileV2{
		ID:    uuid.New(),
		Nik:   service.crypto.Encrypt(request.Nik, aesx.AesCBC),
		Name:  service.crypto.Encrypt(request.Name, aesx.AesCBC),
		Phone: service.crypto.Encrypt(request.Phone, aesx.AesCBC),
		Email: service.crypto.Encrypt(request.Email, aesx.AesCBC),
		DOB:   service.crypto.Encrypt(request.DOB, aesx.AesCBC),
	}

	err = query.InsertWithHeap(service.crypto, ctx, tx, "profile", profile)
	helper.PanicIf(err)

	return request
}
