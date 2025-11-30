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

	query, arg, err := r.db.Dialect.Insert("product").Rows(goqu.Record{
		"name":        product.Name,
		"sku":         product.SKU,
		"category_id": product.CategoryID,
		"description": product.Description,
	}).ToSQL()

	if err != nil {
		return fmt.Errorf("fail: %w", err)
	}

	result, err := r.db.SQL.ExecContext(ctx, query, arg...)
	if err != nil {
		return fmt.Errorf("fail: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("fail: %w", err)
	}
	product.ID = id
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query, args, err := r.db.Dialect.Select("id", "name", "sku", "description", "category_id", "created_at", "updated_at").From("product").Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var product model.Product
	err = r.db.SQL.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get the product with id %d", id)
	}
	return &product, nil

}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*model.Product, error) {
	query, args, err := r.db.Dialect.Select("id", "name", "sku", "description", "category_id", "created_at", "updated_at").From("product").Order(goqu.I("created_at").Desc()).ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.CategoryID,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) DeleteByID(ctx context.Context, id int64) error {

	query, args, err := r.db.Dialect.Delete("product").Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build a delete query: %w", err)
	}
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	// check the rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
