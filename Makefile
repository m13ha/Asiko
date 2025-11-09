BACKEND_DIR := backend
WEB_DIR := web
DOCS_DIR := docs
ENV_FILE := $(BACKEND_DIR)/.env

-include $(ENV_FILE)

AIR_VERSION ?= latest
DB_USER ?= $(DB_USERNAME)
DB_NAME ?= $(DB_DATABASE)
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: dev build run docker-dev docker-build docker-up docker-down clean docs-gen docs-build client-gen install-air mocks migrate-create migrate-up migrate-down migrate-to test web-dev

dev:
	@echo "Starting development server with hot reloading..."
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing Air ($(AIR_VERSION)) for hot reloading..."; \
		go install github.com/air-verse/air@$(AIR_VERSION); \
	fi
	@mkdir -p $(BACKEND_DIR)/tmp
	@cd $(BACKEND_DIR) && air -c .air.toml

build:
	@echo "Building application..."
	@cd $(BACKEND_DIR) && go build -o bin/app .

run: build
	@echo "Running application..."
	@cd $(BACKEND_DIR) && ./bin/app

docker-dev:
	@echo "Starting Docker development environment with hot reloading..."
	@docker compose up --build

docker-build:
	@echo "Building backend Docker image..."
	@docker build -t appointment-master-api $(BACKEND_DIR)

docker-up:
	@echo "Starting Docker services..."
	@docker compose up -d

docker-down:
	@echo "Stopping Docker services..."
	@docker compose down

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BACKEND_DIR)/bin $(BACKEND_DIR)/tmp $(BACKEND_DIR)/build-errors.log
	@docker compose down --volumes --remove-orphans 2>/dev/null || true

docs-gen:
	@echo "Generating OpenAPI spec..."
	@cd $(BACKEND_DIR) && swag init -g main.go -o ../$(DOCS_DIR)

docs-build: docs-gen
	@echo "Building static documentation..."
	@cd $(DOCS_DIR) && npm install && npm run build

client-gen: docs-gen
	@echo "Generating TypeScript API client..."
	@INPUT_SPEC=$(DOCS_DIR)/swagger.json OUTPUT_DIR=$(WEB_DIR)/api-client/src TMP_DIR=$$(mktemp -d) ; \
		set -e ; \
		if openapi-generator-cli generate \
			-i $$INPUT_SPEC \
			-g typescript-fetch \
			-o $$TMP_DIR \
			--additional-properties=typescriptThreePlus=true,usePromises=true ; then \
			rm -rf $$OUTPUT_DIR ; \
			mkdir -p $$OUTPUT_DIR ; \
			rsync -a $$TMP_DIR/ $$OUTPUT_DIR/ ; \
			rm -rf $$TMP_DIR ; \
			echo "✅ API client generated successfully from $$INPUT_SPEC into $$OUTPUT_DIR" ; \
		else \
			echo "⚠️  Skipping update: openapi-generator-cli failed (offline?). Existing client preserved at $$OUTPUT_DIR" ; \
			rm -rf $$TMP_DIR ; \
			exit 1 ; \
		fi

mocks:
	@echo "Deleting old mocks..."
	@rm -rf $(BACKEND_DIR)/repository/mocks $(BACKEND_DIR)/services/mocks
	@echo "Removing deprecated repository mocks directory if it exists..."
	@rm -rf $(BACKEND_DIR)/mocks
	@echo "Regenerating repository mocks..."
	@mockery --all --dir=$(BACKEND_DIR)/repository --output=$(BACKEND_DIR)/repository/mocks
	@echo "Regenerating service mocks..."
	@mockery --all --dir=$(BACKEND_DIR)/services --output=$(BACKEND_DIR)/services/mocks
	@echo "Regenerating notification mocks..."
	@mockery --all --dir=$(BACKEND_DIR)/notifications --output=$(BACKEND_DIR)/notifications/mocks

install-air:
	@echo "Installing Air for hot reloading..."
	@go install github.com/air-verse/air@latest

test:
	@echo "Running tests..."
	@cd $(BACKEND_DIR) && go test ./... -v

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<migration_name>"; exit 1; fi
	@echo "Creating migration file: $(name)"
	@cd $(BACKEND_DIR) && migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	@echo "Applying all up migrations..."
	@cd $(BACKEND_DIR) && migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	@echo "Applying all down migrations..."
	@cd $(BACKEND_DIR) && migrate -path db/migrations -database "$(DB_URL)" -verbose down

migrate-to:
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-to version=<version_number>"; exit 1; fi
	@echo "Migrating to version $(version)..."
	@cd $(BACKEND_DIR) && migrate -path db/migrations -database "$(DB_URL)" -verbose goto $(version)

web-dev:
	@echo "Starting web dev server..."
	@cd $(WEB_DIR) && npm run dev
