package service

import (
	"context"
	"database/sql"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"
	"net/url"
	"regexp"
	"strings"

	"github.com/dyaksa/encryption-pii/crypt/sqlval"
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

	regex := regexp.MustCompile(`[a-zA-Z0-9]+`)

	partsMail := strings.Split(request.Email, "@")
	splitMails := []string{}
	for _, part := range partsMail {
		matches := regex.FindAllString(part, -1)
		splitMails = append(splitMails, matches...)
	}

	partsName := strings.Split(request.Name, " ")
	splitNames := []string{}
	for _, part := range partsName {
		matches := regex.FindAllString(part, -1)
		splitNames = append(splitNames, matches...)
	}

	var textHeaps = []entity.TextHeap{}
	var mailBuilder = new(strings.Builder)
	var nameBuilder = new(strings.Builder)

	for _, splitMail := range splitMails {
		textHeaps = append(textHeaps, entity.TextHeap{
			ID:      uuid.New(),
			Content: strings.ToLower(splitMail),
			Hash:    sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitMail),
			Type:    typeHeapEmail,
		})
		mailBuilder.WriteString(sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitMail).HashString())
	}

	for _, splitName := range splitNames {
		textHeaps = append(textHeaps, entity.TextHeap{
			ID:      uuid.New(),
			Content: strings.ToLower(splitName),
			Hash:    sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitName),
			Type:    typeHeapProfileName,
		})
		nameBuilder.WriteString(sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitName).HashString())
	}

	id := uuid.New()
	profile := entity.Profile{
		ID:        id,
		Nik:       sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Nik, id[:]),
		NikBidx:   sqlval.BIDXString(service.crypt.BIDXFunc(tenantID), request.Nik),
		Name:      sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Name, id[:]),
		NameBidx:  nameBuilder.String(),
		Phone:     sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Phone, id[:]),
		PhoneBidx: sqlval.BIDXString(service.crypt.BIDXFunc(tenantID), request.Phone),
		Email:     sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Email, id[:]),
		EmailBidx: mailBuilder.String(),
		DOB:       sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.DOB, id[:]),
	}

	for _, th := range textHeaps {
		if ok, _ := service.textHeap.IsHashExist(ctx, tx, th.Type, entity.FindTextHeapByHashParams{Hash: th.Hash.HashString()}); !ok {
			err = service.textHeap.Save(ctx, tx, th.Type, th)
			helper.PanicIf(err)
		}
	}

	err = service.repository.Save(ctx, tx, profile)
	helper.PanicIf(err)

	return request
}

func (service *ProfileServiceImpl) FetchProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) *web.ProfileResponse {
	tx, err := service.db.Begin()
	helper.PanicIf(err)

	spr, err := service.repository.FetchProfile(ctx, entity.FetchProfileParams{ID: id}, tx, func(fpr *entity.FetchProfileRow) {
		fpr.Nik = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", id[:])
		fpr.Name = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", id[:])
		fpr.Phone = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", id[:])
		fpr.Email = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", id[:])
		fpr.DOB = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", id[:])
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
		_, err := service.textHeap.FindByContent(ctx, tx,
			param.Type,
			entity.FindTextHeapByContentParams{Content: param.Content},
			nil,
			func(fthr entity.FindTextHeapRow) (bool, error) {
				heaps = append(heaps, fthr.Hash)
				return true, nil
			})
		helper.PanicIf(err)
	}

	switch isType {
	case typeHeapEmail:
		column = "email_bidx"
	case typeHeapProfileName:
		column = "name_bidx"
	}

	spr, err := service.repository.FindBy(ctx, column, entity.FindProfileByBIDXParams{Hash: heaps}, tx, func(fpbnr *entity.FindProfilesByNameRow) {
		fpbnr.Nik = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", fpbnr.ID[:])
		fpbnr.Name = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", fpbnr.ID[:])
		fpbnr.Phone = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", fpbnr.ID[:])
		fpbnr.Email = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", fpbnr.ID[:])
		fpbnr.Dob = sqlval.AEADString(service.crypt.AEADFunc(tenantID), "", fpbnr.ID[:])
	})
	helper.PanicIf(err)

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

	regex := regexp.MustCompile(`[a-zA-Z0-9]+`)

	partsMail := strings.Split(request.Email, "@")
	splitMails := []string{}
	for _, part := range partsMail {
		matches := regex.FindAllString(part, -1)
		splitMails = append(splitMails, matches...)
	}

	partsName := strings.Split(request.Name, " ")
	splitNames := []string{}
	for _, part := range partsName {
		matches := regex.FindAllString(part, -1)
		splitNames = append(splitNames, matches...)
	}

	var textHeaps = []entity.TextHeap{}
	var mailBuilder = new(strings.Builder)
	var nameBuilder = new(strings.Builder)

	for _, splitMail := range splitMails {
		textHeaps = append(textHeaps, entity.TextHeap{
			ID:      uuid.New(),
			Content: strings.ToLower(splitMail),
			Hash:    sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitMail),
			Type:    typeHeapEmail,
		})
		mailBuilder.WriteString(sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitMail).HashString())
	}

	for _, splitName := range splitNames {
		textHeaps = append(textHeaps, entity.TextHeap{
			ID:      uuid.New(),
			Content: strings.ToLower(splitName),
			Hash:    sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitName),
			Type:    typeHeapProfileName,
		})
		nameBuilder.WriteString(sqlval.HMACString(service.crypt.HMACFunc(tenantID), splitName).HashString())
	}

	profile := entity.Profile{
		ID:        id,
		Nik:       sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Nik, id[:]),
		NikBidx:   sqlval.BIDXString(service.crypt.BIDXFunc(tenantID), request.Nik),
		Name:      sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Name, id[:]),
		NameBidx:  nameBuilder.String(),
		Phone:     sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Phone, id[:]),
		PhoneBidx: sqlval.BIDXString(service.crypt.BIDXFunc(tenantID), request.Phone),
		Email:     sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.Email, id[:]),
		EmailBidx: mailBuilder.String(),
		DOB:       sqlval.AEADString(service.crypt.AEADFunc(tenantID), request.DOB, id[:]),
	}

	for _, th := range textHeaps {
		if ok, _ := service.textHeap.IsHashExist(ctx, tx, th.Type, entity.FindTextHeapByHashParams{Hash: th.Hash.HashString()}); !ok {
			err = service.textHeap.Save(ctx, tx, th.Type, th)
			helper.PanicIf(err)
		}
	}

	err = service.repository.Update(ctx, tx, profile)
	helper.PanicIf(err)

	return request
}

// func pqByteArray(arr [][]byte) driver.Valuer {
// 	return pq.ByteaArray(arr)
// }
