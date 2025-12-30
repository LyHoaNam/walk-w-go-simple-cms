package model

import "time"

type Orders struct {
	ID            int64         `db:"id" json:"id"`
	PaymentStatus int8          `db:"payment_status" json:"payment_status"`
	CustomerID    int64         `db:"customer_id" json:"customer_id"`
	PlatformID    int64         `db:"platform_id" json:"platform_id"`
	RetailStoreID int64         `db:"retail_store_id" json:"retail_store_id"`
	PaymentID     int64         `db:"payment_id" json:"payment_id"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
	Items         []*OrderItems `json:"items,omitempty"`
}

type CreateOrders struct {
	CustomerID    int64              `json:"customer_id" validate:"required"`
	PlatformID    int64              `json:"platform_id" validate:"required"`
	RetailStoreID int64              `json:"retail_store_id" validate:"required"`
	PaymentID     int64              `json:"payment_id" validate:"required"`
	Items         []CreateOrderItems `json:"items,omitempty" validate:"dive"`
}

type OrderItems struct {
	ID             int64     `db:"id" json:"id"`
	OrderID        int64     `db:"order_id" json:"order_id"`
	PriceID        int64     `db:"price_id" json:"price_id"`
	VariantValueID int64     `db:"variant_value_id" json:"variant_value_id"`
	Quantity       int       `db:"quantity" json:"quantity"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type CreateOrderItems struct {
	Quantity         int64 `json:"quantity" validate:"required"`
	ProductVariantID int64 `json:"product_variant_id" validate:"required"`
	PriceID          int64 `json:"price_id" validate:"required"`
}

type OrdersPage struct {
	ID            int64
	PaymentStatus int8
	TotalAmount   float64
	CreatedAt     time.Time
	FirstName     string
	LastName      string
	Platform      string
	PaymentMethod string
	OrderStatus   int8
}
