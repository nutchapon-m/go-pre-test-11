DATABASE_URL ?= postgres://user:password@localhost:5432/product_db?sslmode=disable
MIGRATE := docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "$(DATABASE_URL)"

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: migrate-up
migrate-up:
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE) down

.PHONY: test
test:
	go test -v ./...

.PHONY: tidy
tidy:
	go mod tidy
