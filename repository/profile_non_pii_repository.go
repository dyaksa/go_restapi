package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang_restapi/dto"
	"golang_restapi/model/entity"
	"golang_restapi/utils"

	"github.com/google/uuid"
)

type ProfileNonPII interface {
	Create(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) error
	FetchProfile(ctx context.Context, id uuid.UUID, tx *sql.Tx) (entity.ProfileNonPII, error)
	FindAll(ctx context.Context, pagination utils.Pagination, params dto.ParamsListProfile, tx *sql.Tx) (p []entity.ProfileNonPII, err error)
	Update(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) error
}

type ProfileNonPIIRepository struct{}

func NewProfileNonPIIRepository() *ProfileNonPIIRepository {
	return &ProfileNonPIIRepository{}
}

func (repository *ProfileNonPIIRepository) Create(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) (err error) {
	query := "INSERT INTO profiles_not_pii (id, nik, name, phone, email) VALUES ($1, $2, $3, $4, $5)"
	if _, err = tx.ExecContext(ctx, query,
		profile.ID,
		profile.Nik,
		profile.Name,
		profile.Phone,
		profile.Email); err != nil {
		err = fmt.Errorf("error when inserting profile: %w", err)
	}

	return
}

func (repository *ProfileNonPIIRepository) FetchProfile(ctx context.Context, id uuid.UUID, tx *sql.Tx) (p entity.ProfileNonPII, err error) {
	query := "SELECT id, nik, name, phone, email FROM profiles_not_pii WHERE id = $1"
	row := tx.QueryRowContext(ctx, query, id)
	if err = row.Scan(&p.ID, &p.Nik, &p.Name, &p.Phone, &p.Email); err != nil {
		err = fmt.Errorf("error when scanning profile: %w", err)
	}
	return
}

func (repository *ProfileNonPIIRepository) FindAll(ctx context.Context, pagination utils.Pagination, params dto.ParamsListProfile, tx *sql.Tx) (p []entity.ProfileNonPII, err error) {
	query := "SELECT id, nik, name, phone, email FROM profiles_not_pii"
	args := []interface{}{}
	if params.Key != "" && params.Value != "" {
		query += " WHERE " + params.Key + " LIKE  $1"
		args = append(args, "%"+params.Value+"%")
	}

	query += pagination.PaginateQuery()
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var profile entity.ProfileNonPII
		if err = rows.Scan(&profile.ID, &profile.Nik, &profile.Name, &profile.Phone, &profile.Email); err != nil {
			return
		}

		p = append(p, profile)
	}
	return
}

func (repository *ProfileNonPIIRepository) Update(ctx context.Context, tx *sql.Tx, profile entity.ProfileNonPII) (err error) {
	query := "UPDATE profiles_not_pii SET nik = $1, name = $2, phone = $3, email = $4 WHERE id = $5"
	if _, err = tx.ExecContext(ctx, query, profile.Nik, profile.Name, profile.Phone, profile.Email, profile.ID); err != nil {
		err = fmt.Errorf("error when updating profile: %w", err)
	}
	return
}
