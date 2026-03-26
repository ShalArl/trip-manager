.PHONY: help build test run run-dev migrate migrate-create clean migrate-down

# Variables
BINARY_NAME=api
BACKEND_DIR=backend
CMD_DIR=$(BACKEND_DIR)/cmd/api
BIN_DIR=$(BACKEND_DIR)/bin

# Default target
help:
	@echo "Trip Manager - Makefile Commands"
	@echo "=================================="
	@echo ""
	@echo "Development Commands:"
	@echo "  make run-dev       Run the server with auto-reload (requires air)"
	@echo "  make run           Build and run the server"
	@echo "  make build         Build the backend binary"
	@echo "  make test          Run all tests"
	@echo "  make test-verbose  Run tests with verbose output"
	@echo ""
	@echo "Database Commands:"
	@echo "  make migrate       Run pending migrations (auto-runs on server start)"
	@echo "  make db-reset      Reset the database (CAUTION: deletes all data)"
	@echo "  make db-setup      Setup database with migrations"
	@echo ""
	@echo "Cleanup Commands:"
	@echo "  make clean         Remove built binaries"
	@echo "  make clean-all     Remove binaries and generated files"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-up     Start PostgreSQL with Docker"
	@echo "  make docker-down   Stop PostgreSQL Docker container"
	@echo ""

# Build the backend binary
build:
	@echo "Building backend binary..."
	@mkdir -p $(BIN_DIR)
	@cd $(BACKEND_DIR) && go build -o bin/$(BINARY_NAME) ./cmd/api
	@echo "✓ Build complete: $(BIN_DIR)/$(BINARY_NAME)"

# Run the compiled binary
run: build
	@echo "Starting server..."
	@./$(BIN_DIR)/$(BINARY_NAME)

# Run with auto-reload (requires 'air' - github.com/air-verse/air)
run-dev:
	@echo "Starting server with auto-reload..."
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	@cd $(BACKEND_DIR) && air

# Run all tests
test:
	@echo "Running tests..."
	@cd $(BACKEND_DIR) && go test ./...
	@echo "✓ Tests complete"

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@cd $(BACKEND_DIR) && go test -v ./...
	@echo "✓ Tests complete"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@cd $(BACKEND_DIR) && go test -cover ./...
	@echo "✓ Coverage report generated"

# Run migrations (executed automatically on server start)
migrate:
	@echo "Migrations run automatically when the server starts."
	@echo "To manually run migrations, start the server:"
	@echo "  make run"
	@echo ""
	@echo "Database URL: $(DATABASE_URL)"

# Setup database (create DB if not exists)
db-setup:
	@echo "Setting up database..."
	@psql $(DATABASE_URL) -c "SELECT 1" > /dev/null 2>&1 || { \
		echo "Creating database..."; \
		psql -c "CREATE DATABASE trip_manager;" 2>/dev/null || true; \
	}
	@echo "✓ Database ready. Migrations will run on server start."

# Reset the database (WARNING: deletes all data)
db-reset:
	@echo "⚠️  WARNING: This will DELETE ALL DATA from the database!"
	@echo "Press Ctrl+C to cancel, or press Enter to continue..."
	@read -r
	@psql $(DATABASE_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" || true
	@echo "✓ Database reset complete"

# Start PostgreSQL with Docker
docker-up:
	@echo "Starting PostgreSQL with Docker..."
	@docker run -d \
		--name trip_manager_db \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=trip_manager \
		-p 5432:5432 \
		-v trip_manager_data:/var/lib/postgresql/data \
		postgres:16-alpine
	@echo "✓ PostgreSQL started"
	@echo "  Connection: postgres://postgres:postgres@localhost:5432/trip_manager"
	@sleep 2
	@echo "⏳ Waiting for database to be ready..."
	@until pg_isready -h localhost -p 5432 > /dev/null 2>&1; do sleep 1; done
	@echo "✓ Database is ready!"

# Stop PostgreSQL Docker container
docker-down:
	@echo "Stopping PostgreSQL container..."
	@docker stop trip_manager_db 2>/dev/null || true
	@docker rm trip_manager_db 2>/dev/null || true
	@echo "✓ PostgreSQL stopped"

# View Docker logs
docker-logs:
	@docker logs -f trip_manager_db

# Clean up built binaries
clean:
	@echo "Cleaning up binaries..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "✓ Clean complete"

# Clean up everything
clean-all: clean
	@echo "Cleaning up all generated files..."
	@cd $(BACKEND_DIR) && go clean -cache -testcache
	@echo "✓ Full clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@cd $(BACKEND_DIR) && go fmt ./...
	@echo "✓ Format complete"

# Lint code (requires golangci-lint)
lint:
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	@cd $(BACKEND_DIR) && golangci-lint run ./...
	@echo "✓ Lint complete"

# Check dependencies
deps-check:
	@echo "Checking dependencies..."
	@cd $(BACKEND_DIR) && go mod tidy
	@cd $(BACKEND_DIR) && go mod verify
	@echo "✓ Dependencies verified"

# Setup development environment
dev-setup: docker-up db-setup build
	@echo ""
	@echo "✓ Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Set environment variables: export DATABASE_URL=postgres://postgres:postgres@localhost:5432/trip_manager"
	@echo "  2. Start the server: make run"
	@echo "     or with auto-reload: make run-dev"
	@echo ""

# Start everything (DB + Server)
dev: docker-up
	@echo "Database started. Starting server..."
	@$(MAKE) run

# Version info
version:
	@echo "Trip Manager Backend"
	@cd $(BACKEND_DIR) && go version

# All help and info
info: version help
	@echo ""
	@echo "Environment Variables:"
	@echo "  DATABASE_URL=postgres://postgres:postgres@localhost:5432/trip_manager (required)"
	@echo "  SERVER_PORT=8080 (default)"
	@echo "  JWT_SECRET=your-secret-key (default: your-secret-key-change-in-production)"
	@echo "  ENVIRONMENT=development (default)"
	@echo ""

