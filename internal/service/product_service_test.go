package service

import (
	"context"
	"testing"
	"time"

	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	price := 100.0
	req := &port.CreateProductRequest{
		Name:  "Test Product",
		Price: price,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
		return p.Name == "Test Product" && p.Price == 100.0
	})).Return(nil)

	product, err := service.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestPatchProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	existingProduct := &domain.Product{
		ID:        1,
		Name:      "Old Name",
		Price:     100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newName := "New Name"
	req := &port.PatchProductRequest{
		Name: &newName,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingProduct, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
		return p.Name == "New Name" && p.Price == 100.0 // Price unchanged
	})).Return(nil)

	err := service.PatchProduct(context.Background(), 1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
