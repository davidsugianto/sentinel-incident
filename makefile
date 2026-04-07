.PHONY: help dev build test test-coverage clean docker-dev docker-build docker-up docker-down docker-logs docker-clean db-migrate db-reset lint fmt imports tidy deps swagger mocks run hooks-install hooks-uninstall pre-commit

# Variables
APP_NAME := sentinel-incident
GO := go
GOFLAGS := -v
DOCKER := docker
DOCKER_COMPOSE := docker compose

# Database
DB_URL := postgres://user:pass@localhost:5432/sentinel_incident?sslmode=disable
MIGRATE_CMD := $(DOCKER) run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations -database "$(DB_URL)"

# Binary
BINARY := ./bin/$(APP_NAME)

# ============================================
# Default target
# ============================================

help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  dev           Run with Air hot-reload (local, requires air installed)"
	@echo "  build         Build the binary"
	@echo "  run           Run the binary directly"
	@echo ""
	@echo "Docker:"
	@echo "  docker-dev    Run with Docker Compose + Air hot-reload"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-up     Start all services in background"
	@echo "  docker-down   Stop all services"
	@echo "  docker-logs   View container logs (follow mode)"
	@echo "  docker-clean  Remove containers, volumes, and images"
	@echo ""
	@echo "Database:"
	@echo "  db-migrate    Run database migrations (up)"
	@echo "  db-rollback   Rollback last migration"
	@echo "  db-reset      Reset database (drop and recreate)"
	@echo "  db-drop       Drop all migrations"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  test          Run all tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  lint          Run golangci-lint"
	@echo "  fmt           Format code with go fmt"
	@echo "  imports       Organize imports with goimports"
	@echo "  tidy          Run go mod tidy"
	@echo ""
	@echo "Utilities:"
	@echo "  clean         Clean build artifacts"
	@echo "  deps          Download dependencies"
	@echo "  swagger       Generate Swagger documentation"
	@echo "  mocks         Generate mocks with mockery"
	@echo ""
	@echo "Git Hooks:"
	@echo "  hooks-install Install pre-commit hooks"
	@echo "  hooks-uninstall Remove pre-commit hooks"
	@echo "  pre-commit    Run pre-commit checks manually"

# ============================================
# Development
# ============================================

dev:
	@echo "🔄 Starting development server with Air..."
	air -c .air.toml

build:
	@echo "🚀 Building..."
	$(GO) build $(GOFLAGS) -o $(BINARY) ./cmd/http

run: build
	@echo "▶️  Running..."
	$(BINARY)

# ============================================
# Docker
# ============================================

docker-dev:
	@echo "🐳 Starting Docker development environment..."
	$(DOCKER_COMPOSE) up --build

docker-build:
	@echo "🏗️  Building Docker image..."
	$(DOCKER) build -t $(APP_NAME):latest .

docker-up:
	@echo "🐳 Starting services in background..."
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Services started. Use 'make docker-logs' to view logs."

docker-down:
	@echo "🛑 Stopping services..."
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-clean:
	@echo "🧹 Cleaning Docker artifacts..."
	$(DOCKER_COMPOSE) down -v --rmi local
	$(DOCKER) system prune -f

# ============================================
# Database
# ============================================

db-migrate:
	@echo "⬆️  Running migrations..."
	$(MIGRATE_CMD) up

db-rollback:
	@echo "⬇️  Rolling back last migration..."
	$(MIGRATE_CMD) down 1

db-reset:
	@echo "💣 Resetting database..."
	$(MIGRATE_CMD) drop -f || true
	$(MIGRATE_CMD) up

db-drop:
	@echo "💣 Dropping all migrations..."
	$(MIGRATE_CMD) drop -f

# ============================================
# Testing & Quality
# ============================================

test:
	@echo "🧪 Running tests..."
	$(GO) test ./... -v

test-coverage:
	@echo "📊 Running tests with coverage..."
	$(GO) test ./... -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "❌ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt:
	@echo "📝 Formatting code..."
	$(GO) fmt ./...

imports:
	@echo "📝 Organizing imports..."
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "❌ goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

tidy:
	@echo "📦 Tidying modules..."
	$(GO) mod tidy

# ============================================
# Utilities
# ============================================

clean:
	@echo "🧹 Cleaning..."
	rm -rf ./bin
	rm -rf ./tmp
	rm -f coverage.out coverage.html

deps:
	@echo "📦 Downloading dependencies..."
	$(GO) mod download

swagger:
	@echo "📚 Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/http/main.go -o ./docs; \
	else \
		echo "❌ swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

mocks:
	@echo "🎭 Generating mocks..."
	@if command -v mockery >/dev/null 2>&1; then \
		mockery; \
	else \
		echo "❌ mockery not installed. Install with: go install github.com/vektra/mockery/v2@latest"; \
	fi

# ============================================
# Git Hooks
# ============================================

hooks-install:
	@echo "🪝 Installing git hooks..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
		echo "✅ Pre-commit hooks installed (via pre-commit framework)"; \
	else \
		cp .githooks/pre-commit .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit; \
		echo "✅ Pre-commit hooks installed (via .git/hooks)"; \
	fi
	@echo "💡 Tip: You can also use 'git config core.hooksPath .githooks'"

hooks-uninstall:
	@echo "🪝 Uninstalling git hooks..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit uninstall; \
	fi
	rm -f .git/hooks/pre-commit
	@echo "✅ Pre-commit hooks uninstalled"

pre-commit:
	@echo "🔍 Running pre-commit checks..."
	@./.githooks/pre-commit
