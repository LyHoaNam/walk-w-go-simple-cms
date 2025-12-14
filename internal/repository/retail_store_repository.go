package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/database"
	"simple-template/internal/model"
	"simple-template/internal/utils"
)

type RetailStoreRepository struct {
	db *database.DB
}

func NewRetailStoreRepository(db *database.DB) *RetailStoreRepository {
	return &RetailStoreRepository{
		db: db,
	}
}

func (r *RetailStoreRepository) GetAll(ctx context.Context) ([]*model.RetailStore, error) {
	query, args, err := r.db.Dialect.
		Select("id", "name", "phone_number", "created_at", "updated_at").From("retail_stores").ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build query get retail store: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get retail store: %w", err)
	}
	defer rows.Close()

	var retailStores []*model.RetailStore
	for rows.Next() {
		var retailStore model.RetailStore
		var PhoneNumber sql.NullString

		err := rows.Scan(
			&retailStore.ID,
			&retailStore.Name,
			&PhoneNumber,
			&retailStore.CreatedAt,
			&retailStore.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan retail store: %w", err)
		}
		retailStore.PhoneNumber = utils.NullStringToString(PhoneNumber)

		retailStores = append(retailStores, &retailStore)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return retailStores, nil
}
