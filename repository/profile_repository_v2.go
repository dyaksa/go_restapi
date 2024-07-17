package repository

import (
	"context"
	"database/sql"
	"golang_restapi/model/entity"

	"github.com/google/uuid"
)

type ProfileV2 interface {
	FetchProfile(ctx context.Context, args entity.FetchProfileParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FetchProfileRowV2)) (entity.FetchProfileRow, error)
	FetchInvoice(ctx context.Context, tx *sql.Tx, id uuid.UUID, iOptionalInitFunc func(*entity.FetchInvoiceRow)) (entity.FetchInvoiceRow, error)
}

type ProfileRepositoryV2 struct{}

func NewProfileRepositoryV2() *ProfileRepositoryV2 {
	return &ProfileRepositoryV2{}
}

func (repository *ProfileRepositoryV2) FetchInvoice(ctx context.Context, tx *sql.Tx, id uuid.UUID, iOptionalInitFunc func(*entity.FetchInvoiceRow), iOptsChange func(*entity.FetchInvoiceRow)) (entity.FetchInvoiceRow, error) {
	query := "SELECT id, customer_name, customer_address FROM invoice WHERE id = $1"
	row := tx.QueryRowContext(ctx, query, id)
	var i entity.FetchInvoiceRow
	if iOptionalInitFunc != nil {
		iOptionalInitFunc(&i)
	}

	err := row.Scan(&i.ID, &i.CustomerName, &i.CustomerAddress)

	if iOptsChange != nil {
		var iCrypt entity.FetchInvoiceRow
		iOptsChange(&iCrypt)
	}
	return i, err
}

func (repository *ProfileRepositoryV2) FetchProfile(ctx context.Context, args entity.FetchProfileParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FetchProfileRowV2)) (entity.FetchProfileRowV2, error) {
	query := "SELECT nik, name, phone, email, dob FROM profile WHERE id = $1"
	row := tx.QueryRowContext(ctx, query, args.ID)
	var i entity.FetchProfileRowV2
	if iOptionalInitFunc != nil {
		iOptionalInitFunc(&i)
	}

	err := row.Scan(&i.Nik, &i.Name, &i.Phone, &i.Email, &i.DOB)
	return i, err
}
