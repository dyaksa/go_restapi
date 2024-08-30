package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang_restapi/internal/delivery/http/web"
	"golang_restapi/internal/entity"
	"golang_restapi/utils"
)

var _ ProfileNonPII = &ProfileNonPIIRepository{}

type ProfileNonPII interface {
	Create(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) error
	FetchProfile(ctx context.Context, db *sql.DB, id string) (entity.ProfileNonPII, error)
	FindAll(ctx context.Context, db *sql.DB, pagination utils.Pagination, params web.ProfileQueryParam) (p []entity.ProfileNonPII, err error)
	Update(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) error
}

type ProfileNonPIIRepository struct{}

func NewProfileNonPIIRepository() *ProfileNonPIIRepository {
	return &ProfileNonPIIRepository{}
}

func (repository *ProfileNonPIIRepository) Create(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) (err error) {
	query := "INSERT INTO profiles_not_pii (id, nik, name, phone, email) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.ExecContext(ctx, query,
		profile.ID,
		profile.Nik,
		profile.Name,
		profile.Phone,
		profile.Email)
	if err != nil {
		err = fmt.Errorf("error when inserting profile: %w", err)
	}
	return
}

func (repository *ProfileNonPIIRepository) FetchProfile(ctx context.Context, db *sql.DB, id string) (entity.ProfileNonPII, error) {
	query := "SELECT id, nik, name, phone, email FROM profiles_not_pii WHERE id = $1"
	row := db.QueryRowContext(ctx, query, id)
	var i entity.ProfileNonPII
	err := row.Scan(&i.ID, &i.Nik, &i.Name, &i.Phone, &i.Email)
	return i, err
}

func (repository *ProfileNonPIIRepository) FindAll(ctx context.Context, db *sql.DB, pagination utils.Pagination, params web.ProfileQueryParam) (p []entity.ProfileNonPII, err error) {
	query := "SELECT id, nik, name, phone, email FROM profiles_not_pii"
	args := []interface{}{}
	if params.Key != "" && params.Value != "" {
		query += " WHERE " + params.Key + " LIKE  $1"
		args = append(args, "%"+params.Value+"%")
	}

	query += pagination.PaginateQuery()
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	for rows.Next() {
		var i entity.ProfileNonPII
		err = rows.Scan(&i.ID, &i.Nik, &i.Name, &i.Phone, &i.Email)
		if err != nil {
			return
		}
		p = append(p, i)
	}

	return
}

func (repository *ProfileNonPIIRepository) Update(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) (err error) {
	query := "UPDATE profiles_not_pii SET nik = $1, name = $2, phone = $3, email = $4 WHERE id = $5"
	_, err = tx.ExecContext(ctx, query,
		profile.Nik,
		profile.Name,
		profile.Phone,
		profile.Email,
		profile.ID)
	if err != nil {
		err = fmt.Errorf("error when updating profile: %w", err)
	}
	return
}
