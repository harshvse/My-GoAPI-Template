include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: dev up down
run: up

up:
	docker compose up -d db
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker compose exec db pg_isready -h localhost -p 5432 -U user; do \
		echo "PostgreSQL is unavailable - sleeping"; \
		sleep 1; \
	done
	@echo "PostgreSQL is ready!"
	air

down:
	podman-compose down

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_URL) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_URL) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: install-deps
install-deps:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	&& go install github.com/air-verse/air@latest \
	&& go install github.com/swaggo/swag/cmd/swag@latest