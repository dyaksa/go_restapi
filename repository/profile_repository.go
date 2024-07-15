package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang_restapi/model/entity"
	"strings"

	crypt "github.com/dyaksa/encryption-pii/go-encrypt"
)

type Entity interface{}

type Profile interface {
	Save(ctx context.Context, tx *sql.Tx, profile entity.Profile) error
	FetchProfile(ctx context.Context, args entity.FetchProfileParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FetchProfileRow)) (entity.FetchProfileRow, error)
	Find(ctx context.Context, args entity.FindProfileByBIDXParams, tx *sql.Tx, c *crypt.Lib) ([]entity.FindProfilesByNameRow, error)
	FindBy(ctx context.Context, column string, args entity.FindProfileByBIDXParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FindProfilesByNameRow)) ([]entity.FindProfilesByNameRow, error)
	Update(ctx context.Context, tx *sql.Tx, profile entity.Profile) error
}

type ProfileRepository struct{}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{}
}

func (repository *ProfileRepository) Save(ctx context.Context, tx *sql.Tx, profile entity.Profile) error {
	query := "INSERT INTO profile (id, nik, nik_bidx, name, name_bidx, phone, phone_bidx, email, email_bidx, dob) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	_, err := tx.ExecContext(ctx, query, profile.ID, profile.Nik, profile.NikBidx, profile.Name, profile.NameBidx, profile.Phone, profile.PhoneBidx, profile.Email, profile.EmailBidx, profile.DOB)
	return err
}

func (repository *ProfileRepository) FetchProfile(ctx context.Context, args entity.FetchProfileParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FetchProfileRow)) (entity.FetchProfileRow, error) {
	query := "SELECT nik, name, phone, email, dob FROM profile WHERE id = $1"
	row := tx.QueryRowContext(ctx, query, args.ID)
	var i entity.FetchProfileRow
	if iOptionalInitFunc != nil {
		iOptionalInitFunc(&i)
	}

	err := row.Scan(&i.Nik, &i.Name, &i.Phone, &i.Email, &i.DOB)
	return i, err
}

func (repository *ProfileRepository) FindBy(ctx context.Context, column string, args entity.FindProfileByBIDXParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FindProfilesByNameRow)) (pbnr []entity.FindProfilesByNameRow, err error) {
	baseQuery := "SELECT id, nik, name, email, phone, dob FROM profile"
	query, arg := buildLikeQuery(column, baseQuery, args.Hash)
	rows, err := tx.QueryContext(ctx, query, arg...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var i entity.FindProfilesByNameRow
		if iOptionalInitFunc != nil {
			iOptionalInitFunc(&i)
		}

		err = rows.Scan(&i.ID, &i.Nik, &i.Name, &i.Email, &i.Phone, &i.Dob)
		if err != nil {
			return
		}

		pbnr = append(pbnr, i)
	}
	return
}

func (repository *ProfileRepository) Find(ctx context.Context, args entity.FindProfileByBIDXParams, tx *sql.Tx, c *crypt.Lib) ([]entity.FindProfilesByNameRow, error) {
	baseQuery := "SELECT id, nik, name, email, phone, dob FROM profile"
	return crypt.QueryLike(ctx, baseQuery, tx, func(ip *crypt.ILikeParams) {
		ip.ColumnHeap = args.ColumnHeap
		ip.Hash = args.Hash
	}, func(t *entity.FindProfilesByNameRow) {
		t.Nik = c.BToString()
		t.Name = c.BToString()
		t.Phone = c.BToString()
		t.Email = c.BToString()
		t.Dob = c.BToString()
	})
}

func (repository *ProfileRepository) Update(ctx context.Context, tx *sql.Tx, profile entity.Profile) error {
	fmt.Println(profile.Email)
	fmt.Println(profile.ID)
	query := "UPDATE profile SET nik = $1, nik_bidx = $2, name = $3, name_bidx = $4, phone = $5, phone_bidx = $6, email = $7, email_bidx = $8, dob = $9 WHERE id = $10"
	_, err := tx.ExecContext(ctx, query, profile.Nik, profile.NikBidx, profile.Name, profile.NameBidx, profile.Phone, profile.PhoneBidx, profile.Email, profile.EmailBidx, profile.DOB, profile.ID)
	return err
}

func buildLikeQuery(column, baseQuery string, terms []string) (string, []interface{}) {
	var likeClauses []string
	var args []interface{}

	for _, term := range terms {
		likeClauses = append(likeClauses, column+" LIKE $"+fmt.Sprint(len(args)+1))
		args = append(args, "%"+term+"%")
	}

	fullQuery := fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(likeClauses, " OR "))

	return fullQuery, args
}

// func (repository *ProfileRepository) FindProfileByID(ctx context.Context, args entity.FindProfileByIDParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FetchProfileRow)) ([]entity.FetchProfileRow, error) {
// 	query := "SELECT id, nik, name, phone, email, dob FROM profile WHERE id = ANY($1)"
// 	rows, err := tx.QueryContext(ctx, query, pq.Array(args.ID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []entity.FetchProfileRow
// 	for rows.Next() {
// 		var i entity.FetchProfileRow
// 		if iOptionalInitFunc != nil {
// 			iOptionalInitFunc(&i)
// 		}

// 		err = rows.Scan(&i.ID, &i.Nik, &i.Name, &i.Phone, &i.Email, &i.DOB)
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, i)
// 	}
// 	return result, nil
// }

// func (repository *ProfileRepository) FindProfileByName(ctx context.Context, args entity.FindProfileByNameParams, tx *sql.Tx, iOptionalInitFunc func(*entity.FindProfilesByNameRow), iOptionalFilterFunc func(entity.FindProfilesByNameRow) (bool, error)) ([]entity.FindProfilesByNameRow, error) {
// 	query := "SELECT id, nik, name, phone, email, dob FROM profile WHERE name_bidx = ANY($1)"
// 	rows, err := tx.QueryContext(ctx, query, args.NameBidx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []entity.FindProfilesByNameRow
// 	for rows.Next() {
// 		var i entity.FindProfilesByNameRow
// 		if iOptionalInitFunc != nil {
// 			iOptionalInitFunc(&i)
// 		}

// 		err = rows.Scan(&i.ID, &i.Nin, &i.Name, &i.Phone, &i.Email, &i.Dob)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if iOptionalFilterFunc != nil {
// 			ok, err := iOptionalFilterFunc(i)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if !ok {
// 				continue
// 			}
// 		}

// 		result = append(result, i)
// 	}
// 	return result, nil
// }
