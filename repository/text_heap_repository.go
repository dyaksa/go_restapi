package repository

import (
	"context"
	"database/sql"
	"golang_restapi/model/entity"
	"strings"
)

type TextHeap interface {
	Save(ctx context.Context, tx *sql.Tx, typeHeap string, textHeap entity.TextHeap) error
	FindByContent(ctx context.Context, tx *sql.Tx, typeHeap string, args entity.FindTextHeapByContentParams, iOptionalInitFunc func(*entity.FindTextHeapRow), iOptinalFilterFunc func(entity.FindTextHeapRow) (bool, error)) ([]entity.FindTextHeapRow, error)
	IsHashExist(ctx context.Context, tx *sql.Tx, typeHeap string, args entity.FindTextHeapByHashParams) (bool, error)
}

type TextHeapRepository struct{}

func NewTextHeapRepository() *TextHeapRepository {
	return &TextHeapRepository{}
}

func (repository *TextHeapRepository) IsHashExist(ctx context.Context, tx *sql.Tx, typeHeap string, args entity.FindTextHeapByHashParams) (bool, error) {
	var query = new(strings.Builder)
	query.WriteString("SELECT hash FROM ")
	query.WriteString(typeHeap)
	query.WriteString(" WHERE hash = $1")
	row := tx.QueryRowContext(ctx, query.String(), args.Hash)
	var i entity.FindTextHeapRow
	err := row.Scan(&i.Hash)
	if err != nil {
		return false, err
	}
	if i.Hash == args.Hash {
		return true, nil
	}
	return false, nil
}

func (repository *TextHeapRepository) Save(ctx context.Context, tx *sql.Tx, typeHeap string, textHeap entity.TextHeap) error {
	var query = new(strings.Builder)
	query.WriteString("INSERT INTO ")
	query.WriteString(typeHeap)
	query.WriteString(" (id, content, hash) VALUES ($1, $2, $3)")
	_, err := tx.ExecContext(ctx, query.String(), textHeap.ID, textHeap.Content, textHeap.Hash.HashString())
	return err
}

func (repository *TextHeapRepository) FindByContent(ctx context.Context, tx *sql.Tx, typeHeap string, args entity.FindTextHeapByContentParams, iOptionalInitFunc func(*entity.FindTextHeapRow), iOptinalFilterFunc func(entity.FindTextHeapRow) (bool, error)) ([]entity.FindTextHeapRow, error) {
	var query = new(strings.Builder)
	query.WriteString("SELECT id, content, hash FROM ")
	query.WriteString(typeHeap)
	query.WriteString(" WHERE content ILIKE $1")
	rows, err := tx.QueryContext(ctx, query.String(), "%"+args.Content+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []entity.FindTextHeapRow
	for rows.Next() {
		var i entity.FindTextHeapRow
		if iOptionalInitFunc != nil {
			iOptionalInitFunc(&i)
		}

		err = rows.Scan(&i.ID, &i.Content, &i.Hash)
		if err != nil {
			return nil, err
		}

		if iOptinalFilterFunc != nil {
			ok, err := iOptinalFilterFunc(i)
			if err != nil {
				return nil, err
			}

			if !ok {
				continue
			}
		}

		result = append(result, i)
	}

	return result, nil
}
