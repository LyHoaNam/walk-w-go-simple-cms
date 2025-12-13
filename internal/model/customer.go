package model

import "time"

type Customer struct {
	ID          int64     `db:"id" json:"id"`
	FirstName   string    `db:"first_name" json:"first_name"`
	LastName    string    `db:"last_name" json:"last_name"`
	Address     string    `db:"address" json:"address"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber string    `db:"phone_number" json:"phone_number"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCustomerRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name,omitempty"`
	Address     string `json:"address,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateCustomerRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Address     *string `json:"address,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}
