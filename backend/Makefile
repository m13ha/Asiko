.PHONY: dev build run docker-dev docker-build docker-up docker-down clean docs-gen docs-build client-gen

# Development with hot reloading
dev:
	@echo "Starting development server with hot reloading..."
	@./dev.sh

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
	@./generate-client.sh

# Install Air for hot reloading
install-air:
	@echo "Installing Air for hot reloading..."
	@go install github.com/air-verse/air@latest