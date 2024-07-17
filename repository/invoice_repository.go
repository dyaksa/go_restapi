package repository

import (
	"context"
	"database/sql"
	"golang_restapi/model/entity"

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/dyaksa/encryption-pii/crypto/aesx"
	"github.com/dyaksa/encryption-pii/crypto/query"
)

type InvoiceRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Invoice, error)
	Find(c *crypto.Crypto, ctx context.Context, args entity.FindInvoiceByBIDXParams, tx *sql.Tx) ([]entity.FindInvoiceNameRows, error)
}

type InvoiceRepositoryImpl struct{}

func NewInvoiceRepository() InvoiceRepository {
	return &InvoiceRepositoryImpl{}
}

func (r *InvoiceRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Invoice, error) {
	var invoice entity.Invoice
	rows, err := tx.QueryContext(ctx, "SELECT id, customer_name, customer_name_b, customer_address, customer_address_b FROM invoice LIMIT 2")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var invoices []entity.Invoice
	for rows.Next() {
		err := rows.Scan(&invoice.ID, &invoice.CustomerName, &invoice.CustomerNameB, &invoice.CustomerAddr, &invoice.CustomerAddrB)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (service *InvoiceRepositoryImpl) Find(c *crypto.Crypto, ctx context.Context, args entity.FindInvoiceByBIDXParams, tx *sql.Tx) ([]entity.FindInvoiceNameRows, error) {
	q := `SELECT id, number, customer_name_b FROM invoice`
	return query.QueryLike(ctx, q, tx, func(ip *query.ILikeParams) {
		ip.ColumnHeap = args.ColumnHeap
		ip.Hash = args.Hash
	}, func(t *entity.FindInvoiceNameRows) {
		t.CustomerName = c.Decrypt("", aesx.AesCBC)
	})
}
