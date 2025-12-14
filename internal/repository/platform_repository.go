package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/database"
	"simple-template/internal/model"
	"simple-template/internal/utils"
)

type PlatformRepository struct {
	db *database.DB
}

func NewPlatformRepository(db *database.DB) *PlatformRepository {
	return &PlatformRepository{
		db: db,
	}
}

func (r *PlatformRepository) GetAll(ctx context.Context) ([]*model.Platform, error) {
	query, args, err := r.db.Dialect.
		Select("id", "name", "api_endpoint", "feature_struct", "created_at", "updated_at").From("platform").ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build query get platform: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get platforms: %w", err)
	}
	defer rows.Close()

	var platforms []*model.Platform
	for rows.Next() {
		var platform model.Platform
		var ApiEndpoint sql.NullString
		var FeatureStruct sql.NullString

		err := rows.Scan(
			&platform.ID,
			&platform.Name,
			&ApiEndpoint,
			&FeatureStruct,
			&platform.CreatedAt,
			&platform.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan platform: %w", err)
		}
		platform.ApiEndpoint = utils.NullStringToString(ApiEndpoint)
		platform.FeatureStruct = utils.NullStringToString(FeatureStruct)

		platforms = append(platforms, &platform)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return platforms, nil
}
