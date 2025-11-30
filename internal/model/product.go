package model

import "time"

type Product struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	SKU         string    `db:"sku" json:"sku"`
	Description *string   `db:"description" json:"description"`
	CategoryID  int64     `db:"category_id" json:"category_id"`
	CreatedAt   time.Time `db:"create_at" json:"created_at"`
	UpdatedAt   time.Time `db:"update_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	CategoryID  int64   `json:"category_id" validate:"required"`
	Description *string `json:"description,omitempty"`
}
