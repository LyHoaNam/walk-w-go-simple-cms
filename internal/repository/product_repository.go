package repository

import (
	"context"
	"database/sql"
	"fmt"

	"simple-template/internal/database"
	"simple-template/internal/model"

	"github.com/doug-martin/goqu/v9"
)

type ProductRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {

	query, args, err := r.db.Dialect.
		Insert("products").
		Rows(goqu.Record{
			"name":        product.Name,
			"price":       product.Price,
			"category_id": product.CategoryID,
		}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	product.ID = id
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query, args, err := r.db.Dialect.
		Select("id", "name", "price", "category_id", "created_at", "updated_at").
		From("products").
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}
	var product model.Product
	err = r.db.SQL.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return &product, nil
}
