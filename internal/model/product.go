package model

import "time"

type Product struct {
	ID          int64            `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	SKU         string           `db:"sku" json:"sku"`
	Status      int32            `db:"status" json:"status"`
	Description *string          `db:"description" json:"description"`
	Dimension   *string          `db:"dimension" json:"dimension"`
	Weight      *float64         `db:"weight" json:"weight"`
	Brand       *string          `db:"brand" json:"brand"`
	Material    *string          `db:"material" json:"material"`
	Origin      *string          `db:"origin" json:"origin"`
	ImgUrl      string           `db:"img_url" json:"img_url"`
	CategoryID  int64            `db:"category_id" json:"category_id"`
	CreatedAt   time.Time        `db:"create_at" json:"created_at"`
	UpdatedAt   time.Time        `db:"update_at" json:"updated_at"`
	Variant     []ProductVariant `db:"-" json:"variants,omitempty"`
}

type ProductVariant struct {
	ID           int64                 `db:"id" json:"id"`
	Name         string                `db:"name" json:"name"`
	DisplayName  string                `db:"display_name" json:"display_name"`
	DisplayOrder int64                 `db:"display_order" json:"display_order"`
	IsRequire    int16                 `db:"is_require" json:"is_require"`
	ProductID    int64                 `db:"product_id" json:"product_id"`
	Values       []ProductVariantValue `db:"-" json:"variants,omitempty"`
}

type ProductVariantValue struct {
	ID            int64     `db:"id" json:"id"`
	AttributeID   int64     `db:"attribute_id" json:"attribute_id"`
	Value         string    `db:"value" json:"value"`
	DisplayOrder  int       `db:"display_order" json:"display_order"`
	StockQuantity int       `db:"stock_quantity" json:"stock_quantity"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProductVariantValue struct {
	ID            int64  `json:"id" validate:"required"`
	AttributeID   int64  `json:"attribute_id" validate:"required"`
	Value         string `json:"value" validate:"required"`
	DisplayOrder  int    `json:"display_order"`
	StockQuantity int    `json:"stock_quantity,omitempty" validate:"required"`
}

type ProductVariantWithValues struct {
	VariantID    int64  `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	DisplayName  string `json:"display_name,omitempty"`
	DisplayOrder *int64 `json:"display_order,omitempty"`
	IsRequire    *int16 `json:"is_require"`
	ProductID    int64  `json:"product_id" validate:"required"`
	Values       []CreateProductVariantValue
}

type CreateProductRequest struct {
	Name        string                     `json:"name" validate:"required"`
	SKU         string                     `json:"sku" validate:"required"`
	CategoryID  int64                      `json:"category_id" validate:"required"`
	Status      int32                      `json:"status" validate:"required"`
	Description *string                    `json:"description,omitempty"`
	Dimension   *string                    `json:"dimension,omitempty"`
	Weight      *float64                   `json:"weight,omitempty"`
	Brand       *string                    `json:"brand,omitempty"`
	Material    *string                    `json:"material,omitempty"`
	Origin      *string                    `json:"origin,omitempty"`
	ImgUrl      string                     `json:"img_url" validate:"required"`
	Variants    []ProductVariantWithValues `json:"variants,omitempty"`
}
