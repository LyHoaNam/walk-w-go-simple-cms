package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/database"
	"simple-template/internal/model"

	"github.com/doug-martin/goqu/v9"
)

type CustomerRepository struct {
	db *database.DB
}

func NewCustomerRepository(db *database.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	query, args, err := r.db.Dialect.
		Insert("customer").Rows(
		goqu.Record{
			"first_name":   customer.FirstName,
			"last_name":    customer.LastName,
			"address":      customer.Address,
			"email":        customer.Email,
			"phone_number": customer.PhoneNumber,
		}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("fail to build insert query to create customer %w", err)
	}

	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert: %w", err)
	}
	customer.ID = id
	return customer, nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	query, args, err := r.db.Dialect.
		Select("id",
			"first_name",
			"last_name",
			"address", "email",
			"phone_number",
			"created_at",
			"updated_at").From("customer").Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query %w", err)
	}

	var customer model.Customer
	err = r.db.SQL.QueryRowContext(ctx, query, args...).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.LastName,
		&customer.Address,
		&customer.Email,
		&customer.PhoneNumber,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer")
	}
	return &customer, nil
}

func (r *CustomerRepository) GetAll(ctx context.Context) ([]*model.Customer, error) {
	query, args, err := r.db.Dialect.
		Select("id",
			"first_name",
			"last_name",
			"address", "email",
			"phone_number",
			"created_at",
			"updated_at").
		From("customer").
		Order(goqu.I("created_at").Desc()).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer %w", err)
	}
	defer rows.Close()

	var customers []*model.Customer

	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.Address,
			&customer.Email,
			&customer.PhoneNumber,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer %w", err)
		}
		customers = append(customers, &customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return customers, nil
}

func (r *CustomerRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	query, args, err := r.db.Dialect.
		Update("customer").Set(updates).Where(goqu.Ex{"id": id}).ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get row effected: %w", err)
	}
	if rowsEffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	query, args, err := r.db.Dialect.Delete("customer").
		Where(goqu.Ex{"id": id}).ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows effected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}
