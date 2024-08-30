package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang_restapi/internal/delivery/http/web"
	"golang_restapi/internal/entity"
	"golang_restapi/utils"
	"strings"

	"github.com/dyaksa/encryption-pii/crypto"
)

var _ Profile = &ProfileRepository{}

type Profile interface {
	Create(ctx context.Context, tx *sql.Tx, profile entity.Profile) error
	FetchProfile(ctx context.Context, db *sql.DB, id string, initProfile func(*entity.Profile)) (entity.Profile, error)
	FindAll(ctx context.Context, db *sql.DB, pagination utils.Pagination, params web.ProfileQueryParam, initProfile func(*entity.Profile), buildDataFunc func(entity.Profile)) (p []entity.Profile, err error)
	Update(ctx context.Context, tx *sql.Tx, id string, profile entity.Profile) error
}

type ProfileRepository struct {
	crypto *crypto.Crypto
}

func NewProfileRepository(c *crypto.Crypto) *ProfileRepository {
	return &ProfileRepository{crypto: c}
}

func (repository *ProfileRepository) Create(ctx context.Context, tx *sql.Tx, profile entity.Profile) (err error) {
	query := "INSERT INTO profiles (id, nik, nik_bidx, name, name_bidx, phone, phone_bidx, email, email_bidx) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err = tx.ExecContext(ctx, query,
		profile.ID,
		profile.Nik,
		profile.NikBidx,
		profile.Name,
		profile.NameBidx,
		profile.Phone,
		profile.PhoneBidx,
		profile.Email,
		profile.EmailBidx)
	if err != nil {
		err = fmt.Errorf("error when inserting profile: %w", err)
	}
	return
}

func (repository *ProfileRepository) FetchProfile(ctx context.Context, db *sql.DB, id string, initProfile func(*entity.Profile)) (entity.Profile, error) {
	query := "SELECT id, nik, name, phone, email FROM profiles WHERE id = $1"
	row := db.QueryRowContext(ctx, query, id)
	var i entity.Profile
	if initProfile != nil {
		initProfile(&i)
	}

	err := row.Scan(&i.ID, &i.Nik, &i.Name, &i.Phone, &i.Email)
	return i, err
}

func (repository *ProfileRepository) FindAll(ctx context.Context, db *sql.DB, pagination utils.Pagination, params web.ProfileQueryParam, initProfile func(*entity.Profile), buildDataFunc func(entity.Profile)) (p []entity.Profile, err error) {
	query := "SELECT id, nik, name, phone, email FROM profiles"
	var queryParams []interface{}
	if params.Key != "" && params.Value != "" {
		if params.Key == "name" {
			heaps, _ := repository.crypto.SearchContents(ctx, "name_text_heap", func(fthbcp *crypto.FindTextHeapByContentParams) {
				fthbcp.Content = params.Value
			})
			switch {
			case len(heaps) != 0:
				like := []string{}
				for _, heap := range heaps {
					like = append(like, fmt.Sprintf("name_bidx LIKE $%d", len(queryParams)+1))
					queryParams = append(queryParams, "%"+heap+"%")
				}
				query += " WHERE " + strings.Join(like, " OR ")
			}
		}
	}
	query += pagination.PaginateQuery()
	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.Profile
		if initProfile != nil {
			initProfile(&i)
		}
		err = rows.Scan(&i.ID, &i.Nik, &i.Name, &i.Phone, &i.Email)
		if err != nil {
			return nil, err
		}

		if buildDataFunc != nil {
			buildDataFunc(i)
		}

		p = append(p, i)
	}
	return
}

func (repository *ProfileRepository) Update(ctx context.Context, tx *sql.Tx, id string, profile entity.Profile) error {
	query := "UPDATE profiles SET nik = $1, nik_bidx = $2, name = $3, name_bidx = $4, phone = $5, phone_bidx = $6, email = $7, email_bidx = $8 WHERE id = $9"
	_, err := tx.ExecContext(ctx, query, profile.Nik, profile.NikBidx, profile.Name, profile.NameBidx, profile.Phone, profile.PhoneBidx, profile.Email, profile.EmailBidx, id)
	return err
}
