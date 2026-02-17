package service

import (
	"context"
	"time"

	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/port"
)

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *port.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) PatchProduct(ctx context.Context, id int64, req *port.PatchProductRequest) error {
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingProduct == nil {
		return domain.ErrProductNotFound
	}

	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = req.Description // Can be set to nil (or a new value)
	}
	if req.SalePrice != nil {
		existingProduct.SalePrice = req.SalePrice
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}

	existingProduct.UpdatedAt = time.Now()

	return s.repo.Update(ctx, existingProduct)
}
