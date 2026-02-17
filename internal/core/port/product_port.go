package port

import (
	"context"

	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	// Patch might not be needed in repo if we fetch -> update fields -> save, or we can have a specific patch method.
	// Ideally for patch, we fetch the existing entity, update the fields, and save it back.
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*domain.Product, error)
	PatchProduct(ctx context.Context, id int64, req *PatchProductRequest) error
}

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price" validate:"required,gt=0"`
}

type PatchProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       *float64 `json:"price"`
}
