package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	internalHttp "github.com/nutchapon-m/go-pre-test-11/internal/adapter/handler/http"
	"github.com/nutchapon-m/go-pre-test-11/internal/adapter/repository/postgres"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/port"
	"github.com/nutchapon-m/go-pre-test-11/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgresC "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupIntegrationTest(t *testing.T) (*pgxpool.Pool, *echo.Echo, string) {
	ctx := context.Background()
	// Start DB
	pgContainer, err := postgresC.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgresC.WithDatabase("testdb"),
		postgresC.WithUsername("user"),
		postgresC.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	// Cleanup container after test
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbPool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() { dbPool.Close() })

	// Migration
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

	// Wiring
	repo := postgres.NewProductRepository(dbPool)
	svc := service.NewProductService(repo)
	handler := internalHttp.NewProductHandler(svc)

	e := echo.New()
	e.POST("/product", handler.CreateProduct)
	e.PATCH("/product/:id", handler.PatchProduct)

	// Start server on random port
	// In a real scenario we'd bind to :0 and get the listener port,
	// but here we can just use httptest or standard http calls if we run the handler directly?
	// Actually bootstrapping echo to listen on a port is better for E2E.
	// Let's rely on httptest.Server or just run it. Using httptest with Echo is easier for handler testing,
	// but full E2E usually implies network.
	// Let's skip valid network binding and use httptest which is cleaner for Go logic verification.
	// BUT the user asked for "Component test (E2E within service)".
	// I will return the Echo instance and the repository to assert state if needed, but primarily test via http requests.
	// Actually, let's just use httptest.NewServer(e)

	return dbPool, e, "" // string for url if we used real server
}

func TestE2E_CreateAndPatchProduct(t *testing.T) {
	dbPool, e, _ := setupIntegrationTest(t)

	// We can use net/http/httptest
	// server := httptest.NewServer(e) // Echo implements http.Handler
	// defer server.Close()
	// url := server.URL

	// Wait, Echo generic handler is `e`.

	// 1. Create Product
	createReq := port.CreateProductRequest{
		Name:  "E2E Product",
		Price: 50.0,
	}
	body, _ := json.Marshal(createReq)

	req := httptestRequest(http.MethodPost, "/product", bytes.NewReader(body))
	rec := httptestRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var createResp internalHttp.Response
	json.Unmarshal(rec.Body.Bytes(), &createResp)
	assert.True(t, createResp.Successful)

	// Parse response data to get ID (assuming `data` is map or struct)
	// The response data is `*domain.Product`.
	// We need to marshal/unmarshal mapstructure or just simple generic map or struct.
	dataBytes, _ := json.Marshal(createResp.Data)
	var product domain.Product
	json.Unmarshal(dataBytes, &product)

	assert.NotZero(t, product.ID)
	assert.Equal(t, "E2E Product", product.Name)

	// 2. Patch Product
	newName := "E2E Product Patched"
	patchReq := port.PatchProductRequest{
		Name: &newName,
	}
	body, _ = json.Marshal(patchReq)

	req = httptestRequest(http.MethodPatch, fmt.Sprintf("/product/%d", product.ID), bytes.NewReader(body))
	rec = httptestRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify in DB
	var storedName string
	err := dbPool.QueryRow(context.Background(), "SELECT name FROM products WHERE id=$1", product.ID).Scan(&storedName)
	assert.NoError(t, err)
	assert.Equal(t, newName, storedName)
}

func httptestRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

func httptestRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
