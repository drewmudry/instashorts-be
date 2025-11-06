# Monorepo Makefile for Instashorts Backend

include .env
export

.PHONY: help build build-api build-worker build-renderer up down restart logs clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: build-api build-worker build-renderer ## Build all services

build-api: ## Build API service
	@echo "Building API..."
	cd is-api && go build -o ../bin/api ./cmd/api

build-worker: ## Build Worker service
	@echo "Building Worker..."
	cd is-worker && go build -o ../bin/worker ./cmd/worker

build-renderer: ## Build Renderer service
	@echo "Building Renderer..."
	cd is-render && npm install && npm run build

# Docker Compose targets
up: ## Start all services with docker-compose
	docker compose up -d

up-build: ## Build and start all services
	docker compose up -d --build

down: ## Stop all services
	docker compose down

restart: ## Restart all services
	docker compose restart

logs: ## Show logs from all services
	docker compose logs -f

logs-api: ## Show API logs
	docker compose logs -f api

logs-worker: ## Show Worker logs
	docker compose logs -f worker

logs-renderer: ## Show Renderer logs
	docker compose logs -f renderer

# Development targets
dev-api: ## Run API in development mode
	cd is-api && go run ./cmd/api

dev-worker: ## Run Worker in development mode
	cd is-worker && go run ./cmd/worker

dev-renderer: ## Run Renderer in development mode
	cd is-render && npm run dev

# Database migrations
migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@migrate -path is-api/migrations -database "postgres://$${BLUEPRINT_DB_USERNAME}:$${BLUEPRINT_DB_PASSWORD}@$${BLUEPRINT_DB_HOST}:$${BLUEPRINT_DB_PORT}/$${BLUEPRINT_DB_DATABASE}?sslmode=disable&search_path=$${BLUEPRINT_DB_SCHEMA}" up

migrate-down: ## Rollback last migration
	@echo "Rolling back migration..."
	@migrate -path is-api/migrations -database "postgres://$${BLUEPRINT_DB_USERNAME}:$${BLUEPRINT_DB_PASSWORD}@$${BLUEPRINT_DB_HOST}:$${BLUEPRINT_DB_PORT}/$${BLUEPRINT_DB_DATABASE}?sslmode=disable&search_path=$${BLUEPRINT_DB_SCHEMA}" down 1

# Testing
test: ## Run all tests
	@echo "Running tests..."
	go work test ./...

test-api: ## Run API tests
	cd is-api && go test ./...

test-worker: ## Run Worker tests
	cd is-worker && go test ./...

# Go workspace management
work-sync: ## Sync Go workspace
	go work sync

# Clean targets
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf is-render/dist/
	rm -rf is-render/node_modules/
	rm -rf is-render/output/

clean-all: clean ## Clean everything including Docker volumes
	docker compose down -v

# Install dependencies
install: install-go install-node ## Install all dependencies

install-go: ## Install Go dependencies
	go work sync
	cd pkg && go mod download
	cd is-api && go mod download
	cd is-worker && go mod download

install-node: ## Install Node dependencies
	cd is-render && npm install

