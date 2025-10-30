include .env

.PHONY: dev build run docker-dev docker-build docker-up docker-down clean docs-gen docs-build client-gen install-air mocks migrate-create migrate-up migrate-down migrate-to

# Tool versions (override via environment if needed)
AIR_VERSION ?= latest

# Development with hot reloading
dev:
	@echo "Starting development server with hot reloading..."
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing Air ($(AIR_VERSION)) for hot reloading..."; \
		go install github.com/air-verse/air@$(AIR_VERSION); \
	fi
	@mkdir -p tmp
	@air -c .air.toml

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/app .

# Run the built application
run: build
	@echo "Running application..."
	@./bin/app

# Docker development with hot reloading
docker-dev:
	@echo "Starting Docker development environment with hot reloading..."
	@docker compose up --build

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t appointment-master-api .

# Start Docker services
docker-up:
	@echo "Starting Docker services..."
	@docker compose up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	@docker compose down

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ tmp/ build-errors.log
	@docker compose down --volumes --remove-orphans 2>/dev/null || true

# --- Code Generation ---

# Generate OpenAPI spec from code annotations
docs-gen:
	@echo "Generating OpenAPI spec..."
	@swag init -g main.go -o ../docs

# Build static HTML documentation
docs-build: docs-gen
	@echo "Building static documentation..."
	@cd ../docs && npm install && npm run build

# Generate TypeScript API client
client-gen: docs-gen
	@echo "Generating TypeScript API client..."
	@INPUT_SPEC=./docs/swagger.json OUTPUT_DIR=../api-client/src ; \
		rm -rf $$OUTPUT_DIR ; \
		openapi-generator-cli generate \
			-i $$INPUT_SPEC \
			-g typescript-fetch \
			-o $$OUTPUT_DIR \
			--additional-properties=typescriptThreePlus=true,usePromises=true ; \
		echo "✅ API client generated successfully in $$OUTPUT_DIR"

# Regenerate repository mocks
mocks:
	@echo "Deleting old mocks..."
	@rm -rf ./repository/mocks ./services/mocks
	@echo "Removing deprecated repository mocks directory if it exists..."
	@rm -rf ./mocks
	@echo "Regenerating repository mocks..."
	@mockery --all --dir=./repository --output=./repository/mocks
	@echo "Regenerating service mocks..."
	@mockery --all --dir=./services --output=./services/mocks
	@echo "Regenerating notification mocks..."
	@mockery --all --dir=./notifications --output=./notifications/mocks

# Install Air for hot reloading
install-air:
	@echo "Installing Air for hot reloading..."
	@go install github.com/air-verse/air@latest

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# --- Database Migrations ---

DB_URL := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<migration_name>"; exit 1; fi
	@echo "Creating migration file: $(name)"
	@migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	@echo "Applying all up migrations..."
	@migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	@echo "Applying all down migrations..."
	@migrate -path db/migrations -database "$(DB_URL)" -verbose down

migrate-to:
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-to version=<version_number>"; exit 1; fi
	@echo "Migrating to version $(version)..."
	@migrate -path db/migrations -database "$(DB_URL)" -verbose goto $(version)
