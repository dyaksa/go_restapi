package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang_restapi/helper"
	"golang_restapi/model/entity"
)

type Category interface {
	Save(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category
	Update(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category
	Delete(ctx context.Context, tx *sql.Tx, categoryId int)
	FindOneByID(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error)
	FindAll(ctx context.Context, tx *sql.Tx) []entity.Category
}

type CategoryRepository struct {
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (repository *CategoryRepository) Save(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category {
	query := "INSERT INTO category (name) VALUES (?)"
	result, err := tx.ExecContext(ctx, query, category.Name)
	helper.PanicIf(err)

	id, err := result.LastInsertId()
	helper.PanicIf(err)

	category = entity.Category{
		ID:   int(id),
		Name: category.Name,
	}

	return category
}

func (repository *CategoryRepository) Update(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category {
	query := "update category set name = ? where id = ?"
	_, err := tx.ExecContext(ctx, query, category.Name, category.ID)
	helper.PanicIf(err)

	return category
}

func (repository *CategoryRepository) Delete(ctx context.Context, tx *sql.Tx, categoryId int) {
	query := "DELETE FROM category WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, categoryId)
	helper.PanicIf(err)
}

func (repository *CategoryRepository) FindOneByID(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error) {
	query := "SELECT id, name FROM category WHERE id = ?"
	row, err := tx.QueryContext(ctx, query, categoryId)
	helper.PanicIf(err)
	defer row.Close()

	var category entity.Category
	if row.Next() {
		err := row.Scan(&category.ID, &category.Name)
		helper.PanicIf(err)
		return category, nil
	} else {
		return category, errors.New("category not found")
	}

}

func (repository *CategoryRepository) FindAll(ctx context.Context, tx *sql.Tx) []entity.Category {
	query := "SELECT id, name FROM category"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIf(err)

	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		err := rows.Scan(&category.ID, &category.Name)
		helper.PanicIf(err)
		categories = append(categories, category)
	}
	return categories
}
