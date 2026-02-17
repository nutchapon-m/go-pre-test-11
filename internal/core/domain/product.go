package domain

import (
	"errors"
	"time"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type Product struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"` // Nullable
	SalePrice   *float64  `json:"sale_price,omitempty" db:"sale_price"`   // Nullable
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
