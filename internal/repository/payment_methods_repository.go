package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/database"
	"simple-template/internal/model"
	"simple-template/internal/utils"
)

type PaymentMethodsRepository struct {
	db *database.DB
}

func NewPaymentMethodsRepository(db *database.DB) *PaymentMethodsRepository {
	return &PaymentMethodsRepository{
		db: db,
	}
}

func (r *PaymentMethodsRepository) GetAll(ctx context.Context) ([]*model.PaymentMethods, error) {
	query, args, err := r.db.Dialect.
		Select("id", "name", "code", "description", "is_active", "created_at", "updated_at").From("payment_methods").ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build query get payment methods: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer rows.Close()

	var PaymentMethods []*model.PaymentMethods
	for rows.Next() {
		var PaymentMethod model.PaymentMethods
		var Description sql.NullString

		err := rows.Scan(
			&PaymentMethod.ID,
			&PaymentMethod.Name,
			&PaymentMethod.Code,
			&Description,
			&PaymentMethod.IsActive,
			&PaymentMethod.CreatedAt,
			&PaymentMethod.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan retail store: %w", err)
		}
		PaymentMethod.Description = utils.NullStringToString(Description)

		PaymentMethods = append(PaymentMethods, &PaymentMethod)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return PaymentMethods, nil
}
