package service

import (
	"context"
	"database/sql"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
	"golang_restapi/model/web"
	"golang_restapi/repository"
	"net/url"

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/dyaksa/encryption-pii/crypto/aesx"
	"github.com/dyaksa/encryption-pii/crypto/query"
)

const (
	typeHeapName = "name_text_heap"
	typeHeapAddr = "address_text_heap"
)

type InvoiceService interface {
	FindAll(ctx context.Context) ([]entity.Invoice, error)
	Search(cctx context.Context, values url.Values) (irs []*web.InvoiceFormatterResponse, err error)
}

type InvoiceServiceImpl struct {
	db      *sql.DB
	invRepo repository.InvoiceRepository
	crypt   *crypto.Crypto
}

func NewInvoiceService(db *sql.DB) InvoiceService {
	invRepo := repository.NewInvoiceRepository()
	crypto, _ := crypto.New(crypto.Aes256KeySize)
	return &InvoiceServiceImpl{db: db, invRepo: invRepo, crypt: crypto}
}

func (s *InvoiceServiceImpl) Search(ctx context.Context, values url.Values) (irs []*web.InvoiceFormatterResponse, err error) {
	var isType string
	var column string
	var params = []entity.FindInvoiceParams{
		{Type: typeHeapName, Content: values.Get("name")},
		{Type: typeHeapAddr, Content: values.Get("address")},
	}

	for _, param := range params {
		if param.Content != "" {
			isType = param.Type
		}
	}

	tx, err := s.db.Begin()
	helper.PanicIf(err)

	heaps := []string{}
	for _, param := range params {
		if param.Content != "" {
			heaps, err = query.SearchContents(
				ctx, tx,
				param.Type,
				query.FindTextHeapByContentParams{Content: param.Content})
			helper.PanicIf(err)
		}
	}

	switch isType {
	case typeHeapName:
		column = "customer_name_bidx"
	case typeHeapAddr:
		column = "customer_address_bidx"
	}

	inv, err := s.invRepo.Find(s.crypt, ctx, entity.FindInvoiceByBIDXParams{ColumnHeap: column, Hash: heaps}, tx)
	helper.PanicIf(err)

	for _, i := range inv {
		ir := &web.InvoiceFormatterResponse{
			ID:            i.ID,
			InvoiceNumber: i.Number,
			CustomerName:  i.CustomerName.To(),
		}
		irs = append(irs, ir)
	}
	return
}

func (s *InvoiceServiceImpl) FindAll(ctx context.Context) ([]entity.Invoice, error) {
	tx, err := s.db.Begin()
	helper.PanicIf(err)

	invoices, err := s.invRepo.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	for _, invoice := range invoices {
		customerName := invoice.CustomerName
		customerAddr := invoice.CustomerAddr

		customerNameDec, err := aesx.Decrypt(aesx.AesCBC, []byte(*s.crypt.AESKey), []byte(customerName))
		if err != nil {
			return nil, err
		}

		customerAddrDec, err := aesx.Decrypt(aesx.AesCBC, []byte(*s.crypt.AESKey), []byte(customerAddr))
		if err != nil {
			return nil, err
		}

		invoicex := entity.Invoice{
			ID:            invoice.ID,
			CustomerNameB: s.crypt.Encrypt(string(customerNameDec), aesx.AesCBC),
			CustomerAddrB: s.crypt.Encrypt(string(customerAddrDec), aesx.AesCBC),
		}

		err = query.UpdateWithHeap(s.crypt, ctx, tx, "invoice", invoicex, invoicex.ID.String())
		helper.PanicIf(err)
	}

	defer helper.CommitAndRollbackError(tx)
	return invoices, nil
}
