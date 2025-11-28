package model

import "time"

type Product struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Price      int64     `db:"price" json:"price"`
	CategoryID int64     `db:"category_id" json:"category_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name       string `json:"name" validate:"required"`
	Price      int64  `json:"price" validate:"required,gt=0"`
	CategoryID int64  `json:"category_id" validate:"required"`
}
