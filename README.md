# Go Pre-Test Project

This project implements a Product API using Go, Echo, PostgreSQL, and Clean Architecture.

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Make (optional)

## Setup

1. **Start Database**:
   ```bash
   make up
   # or
   docker compose up -d
   ```

2. **Run Migrations**:
   ```bash
   make migrate-up
   ```

3. **Run Application**:
   ```bash
   make run
   # or
   go run cmd/api/main.go
   ```

## API Documentation

Swagger documentation is available at:
http://localhost:8080/swagger/index.html

## Testing

Run unit and integration tests:
```bash
make test
```

## Structure

- `cmd/api`: Main entry point.
- `internal/core`: Domain entities and ports (interfaces).
- `internal/service`: Application business logic.
- `internal/adapter`: Implementations of ports (HTTP handlers, PostgreSQL repo).
- `migrations`: Database migration files.
- `test`: E2E/Component tests.

## API Endpoints

- `POST /product`: Create a new product.
- `PATCH /product/:id`: Patch an existing product.
