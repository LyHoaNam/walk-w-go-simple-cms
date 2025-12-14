package model

import "time"

type Platform struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	ApiEndpoint   string    `db:"api_endpoint" json:"api_endpoint"`
	FeatureStruct string    `db:"feature_structure" json:"feature_structure"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
