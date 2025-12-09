package model

import "time"

type Price struct {
	ID             int64      `db:"id" json:"id"`
	VariantID      int64      `db:"variant_id" json:"variant_id"`
	Price          float64    `db:"price" json:"price"`
	CompareAtPrice *float64   `db:"compare_at_price" json:"compare_at_price,omitempty"`
	CostPrice      *float64   `db:"cost_price" json:"cost_price,omitempty"`
	Status         int8       `db:"status" json:"status"`
	EffectiveFrom  time.Time  `db:"effective_from" json:"effective_from"`
	EffectiveTo    *time.Time `db:"effective_to" json:"effective_to,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

type CreatePriceRequest struct {
	VariantID      int64    `json:"variant_id" validate:"required"`
	Price          float64  `json:"price" validate:"required"`
	CompareAtPrice *float64 `json:"compare_at_price,omitempty"`
	CostPrice      *float64 `json:"cost_price,omitempty"`
	Status         int8     `json:"status" validate:"required"`
	EffectiveFrom  *string  `json:"effective_from,omitempty"`
	EffectiveTo    *string  `json:"effective_to,omitempty"`
}

type UpdatePriceRequest struct {
	Price          *float64 `json:"price,omitempty"`
	CompareAtPrice *float64 `json:"compare_at_price,omitempty"`
	CostPrice      *float64 `json:"cost_price,omitempty"`
	Status         *int8    `json:"status,omitempty"`
	EffectiveFrom  *string  `json:"effective_from,omitempty"`
	EffectiveTo    *string  `json:"effective_to,omitempty"`
}
