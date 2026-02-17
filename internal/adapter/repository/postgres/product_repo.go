package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (name, description, sale_price, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.SalePrice,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, description, sale_price, price, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var product domain.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SalePrice,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" { // specific pgx error check or import pgx
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, sale_price = $3, price = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query,
		product.Name,
		product.Description,
		product.SalePrice,
		product.Price,
		product.UpdatedAt,
		product.ID,
	)
	return err
}
