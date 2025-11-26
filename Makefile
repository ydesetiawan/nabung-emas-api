.PHONY: help install run build test clean migrate db-create db-drop db-reset

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

run: ## Run the application
	@echo "🚀 Starting server..."
	go run cmd/server/main.go

build: ## Build the application
	@echo "🔨 Building application..."
	go build -o bin/nabung-emas-api cmd/server/main.go
	@echo "✅ Build complete! Binary: bin/nabung-emas-api"

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -rf uploads/
	@echo "✅ Clean complete!"

db-create: ## Create database
	@echo "📊 Creating database..."
	createdb nabung_emas
	@echo "✅ Database created!"

db-drop: ## Drop database
	@echo "⚠️  Dropping database..."
	dropdb nabung_emas
	@echo "✅ Database dropped!"

db-reset: db-drop db-create migrate ## Reset database (drop, create, migrate)
	@echo "✅ Database reset complete!"

migrate: ## Run database migrations
	@echo "🔄 Running migrations..."
	@for file in migrations/*.sql; do \
		echo "Running $$file..."; \
		psql -d nabung_emas -f $$file; \
	done
	@echo "✅ Migrations complete!"

dev: ## Run with hot reload (requires air)
	@echo "🔥 Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t nabung-emas-api:latest .
	@echo "✅ Docker image built!"

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env nabung-emas-api:latest

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not installed. Install from: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

setup: install db-create migrate ## Complete setup (install, create db, migrate)
	@echo "✅ Setup complete! Run 'make run' to start the server."
