package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestProductRepository(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbPool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer dbPool.Close()

	// Initial Migration
	_, err = dbPool.Exec(ctx, `
		CREATE TABLE products (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			sale_price NUMERIC(10, 2),
			price NUMERIC(10, 2) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	require.NoError(t, err)

	repo := NewProductRepository(dbPool)

	t.Run("CreateProduct", func(t *testing.T) {
		price := 100.0
		product := &domain.Product{
			Name:      "Test Product",
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, product)
		assert.NoError(t, err)
		assert.NotZero(t, product.ID)
	})

	t.Run("GetProduct", func(t *testing.T) {
		// Assuming ID 1 from previous test, but better to create new or clean up
		// Let's create a fresh one just to be safe or rely on state if sequential.
		// For simplicity, let's just insert another one.
		price := 200.0
		p := &domain.Product{Name: "Get Me", Price: price, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err := repo.Create(ctx, p)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, p.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, p.Name, fetched.Name)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		price := 300.0
		p := &domain.Product{Name: "Update Me", Price: price, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err := repo.Create(ctx, p)
		require.NoError(t, err)

		newName := "Updated Name"
		p.Name = newName
		p.UpdatedAt = time.Now()

		err = repo.Update(ctx, p)
		assert.NoError(t, err)

		fetched, err := repo.GetByID(ctx, p.ID)
		assert.NoError(t, err)
		assert.Equal(t, newName, fetched.Name)
	})
}
