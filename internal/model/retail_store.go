package model

import "time"

type RetailStore struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	PhoneNumber string    `db:"phone_number" json:"phone_number"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
