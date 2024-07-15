package service

import (
	"context"
	"database/sql"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"
	"net/url"

	crypt "github.com/dyaksa/encryption-pii/go-encrypt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	typeHeapEmail       = "email_text_heap"
	typeHeapProfileName = "profile_text_heap"
)

type ProfileService interface {
	Create(ctx context.Context, tenantID uuid.UUID, request web.ProfileRequest) web.ProfileRequest
	FetchProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) *web.ProfileResponse
	FindByTextHeapContent(ctx context.Context, tenantID uuid.UUID, values url.Values) (prs []*web.ProfileResponse, err error)
	Update(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, request web.ProfileRequest) web.ProfileRequest
}

type ProfileServiceImpl struct {
	db         *sql.DB
	repository repository.Profile
	textHeap   repository.TextHeap
	validator  *validator.Validate
	crypt      *crypt.Lib
}

func NewProfileService(db *sql.DB, repository repository.Profile, textHeap repository.TextHeap, validator *validator.Validate, crypt *crypt.Lib) *ProfileServiceImpl {
	return &ProfileServiceImpl{
		db:         db,
		repository: repository,
		textHeap:   textHeap,
		validator:  validator,
		crypt:      crypt,
	}
}

func (service *ProfileServiceImpl) Create(ctx context.Context, tenantID uuid.UUID, request web.ProfileRequest) web.ProfileRequest {
	err := service.validator.Struct(request)
	helper.PanicIf(err)

	tx, err := service.db.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	email, emailDict := service.crypt.BuildHeap(request.Email, typeHeapEmail)
	name, nameDict := service.crypt.BuildHeap(request.Name, typeHeapProfileName)

	profile := entity.Profile{
		ID:        uuid.New(),
		Nik:       service.crypt.AEADString(request.Nik),
		NikBidx:   service.crypt.BIDXString(request.Nik),
		Name:      service.crypt.AEADString(request.Name),
		NameBidx:  name,
		Phone:     service.crypt.AEADString(request.Phone),
		PhoneBidx: service.crypt.BIDXString(request.Phone),
		Email:     service.crypt.AEADString(request.Email),
		EmailBidx: email,
		DOB:       service.crypt.AEADString(request.DOB),
	}

	err = service.repository.Save(ctx, tx, profile)
	helper.PanicIf(err)

	err = service.crypt.SaveToHeap(ctx, tx, emailDict)
	helper.PanicIf(err)

	err = service.crypt.SaveToHeap(ctx, tx, nameDict)
	helper.PanicIf(err)

	return request
}

func (service *ProfileServiceImpl) FetchProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) *web.ProfileResponse {
	tx, err := service.db.Begin()
	helper.PanicIf(err)

	spr, err := service.repository.FetchProfile(ctx, entity.FetchProfileParams{ID: id}, tx, func(fpr *entity.FetchProfileRow) {
		fpr.Nik = service.crypt.BToString()
		fpr.Name = service.crypt.BToString()
		fpr.Phone = service.crypt.BToString()
		fpr.Email = service.crypt.BToString()
		fpr.DOB = service.crypt.BToString()
	})
	helper.PanicIf(err)

	pr := &web.ProfileResponse{
		ID:    id,
		Nik:   spr.Nik.To(),
		Name:  spr.Name.To(),
		Phone: spr.Phone.To(),
		Email: spr.Email.To(),
		DOB:   spr.DOB.To(),
	}

	return pr
}

func (service *ProfileServiceImpl) FindByTextHeapContent(ctx context.Context, tenantID uuid.UUID, values url.Values) (prs []*web.ProfileResponse, err error) {
	var isType string
	var column string
	var params = []entity.FindProfileByParams{
		{Type: typeHeapEmail, Content: values.Get("email")},
		{Type: typeHeapProfileName, Content: values.Get("name")},
	}

	for _, param := range params {
		if param.Content != "" {
			isType = param.Type
		}
	}

	tx, err := service.db.Begin()
	helper.PanicIf(err)

	heaps := []string{}
	for _, param := range params {
		heaps, err = service.crypt.SearchContents(
			ctx, tx, param.Type,
			crypt.FindTextHeapByContentParams{Content: param.Content})
		helper.PanicIf(err)
	}

	switch isType {
	case typeHeapEmail:
		column = "email_bidx"
	case typeHeapProfileName:
		column = "name_bidx"
	}

	spr, err := service.repository.Find(ctx, entity.FindProfileByBIDXParams{ColumnHeap: column, Hash: heaps}, tx, service.crypt)

	for _, sp := range spr {
		pr := &web.ProfileResponse{
			ID:    sp.ID,
			Nik:   sp.Nik.To(),
			Name:  sp.Name.To(),
			Phone: sp.Phone.To(),
			Email: sp.Email.To(),
			DOB:   sp.Dob.To(),
		}
		prs = append(prs, pr)
	}

	return
}

func (service *ProfileServiceImpl) Update(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, request web.ProfileRequest) web.ProfileRequest {
	err := service.validator.Struct(request)
	helper.PanicIf(err)

	tx, err := service.db.Begin()
	helper.PanicIf(err)
	defer helper.CommitAndRollbackError(tx)

	email, emailDict := service.crypt.BuildHeap(request.Email, typeHeapEmail)
	name, nameDict := service.crypt.BuildHeap(request.Name, typeHeapProfileName)

	profile := entity.Profile{
		ID:        id,
		Nik:       service.crypt.AEADString(request.Nik),
		NikBidx:   service.crypt.BIDXString(request.Nik),
		Name:      service.crypt.AEADString(request.Name),
		NameBidx:  name,
		Phone:     service.crypt.AEADString(request.Phone),
		PhoneBidx: service.crypt.BIDXString(request.Phone),
		Email:     service.crypt.AEADString(request.Email),
		EmailBidx: email,
		DOB:       service.crypt.AEADString(request.DOB),
	}

	err = service.crypt.SaveToHeap(ctx, tx, emailDict)
	helper.PanicIf(err)

	err = service.crypt.SaveToHeap(ctx, tx, nameDict)
	helper.PanicIf(err)

	err = service.repository.Update(ctx, tx, profile)
	helper.PanicIf(err)

	return request
}

// func pqByteArray(arr [][]byte) driver.Valuer {
// 	return pq.ByteaArray(arr)
// }
