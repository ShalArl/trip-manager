.PHONY: help build test run run-dev migrate migrate-create clean migrate-down build-windows run-windows run-dev-win

# Variables
BINARY_NAME=api
BACKEND_DIR=backend
CMD_DIR=$(BACKEND_DIR)/cmd/api
BIN_DIR=$(BACKEND_DIR)/bin
BINARY_NAME_WIN=$(BINARY_NAME).exe

# Default target
help:
	@echo "Trip Manager - Makefile Commands"
	@echo "=================================="
	@echo ""
	@echo "Development Commands (Linux/macOS):"
	@echo "  make run-dev        Run the server with auto-reload (requires air)"
	@echo "  make run            Build and run the server"
	@echo "  make build          Build the backend binary"
	@echo ""
	@echo "Development Commands (Windows):"
	@echo "  make run-dev-win    Run the server for Windows with auto-reload (requires air)"
	@echo "  make run-windows    Build and run the server for Windows"
	@echo "  make build-windows  Build the backend binary for Windows"
	@echo ""
	@echo "Other Commands:"
	@echo "  make test           Run all tests"
	@echo "  make test-verbose   Run tests with verbose output"
	@echo ""
	@echo "Database Commands:"
	@echo "  make migrate        Run pending migrations (auto-runs on server start)"
	@echo "  make db-reset       Reset the database (CAUTION: deletes all data)"
	@echo "  make db-setup       Setup database with migrations"
	@echo ""
	@echo "Cleanup Commands:"
	@echo "  make clean          Remove built binaries"
	@echo "  make clean-all      Remove binaries and generated files"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make db-up          Start PostgreSQL with Docker"
	@echo "  make db-down        Stop PostgreSQL Docker container"
	@echo "  make storage-up     Start MinIO with Docker"
	@echo "  make storage-setup  Start MinIO and create bucket"
	@echo "  make storage-down   Stop MinIO Docker container"
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

# Build the backend binary for Windows
build-windows:
	@echo "Building backend binary for Windows..."
	@mkdir -p $(BIN_DIR)
	@cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME_WIN) ./cmd/api
	@echo "✓ Windows build complete: $(BIN_DIR)/$(BINARY_NAME_WIN)"

# Run the compiled Windows binary (requires WSL2 or Windows environment)
run-windows: build-windows
	@echo "Starting Windows server..."
	@./$(BIN_DIR)/$(BINARY_NAME_WIN)

# Run Windows build with auto-reload (requires 'air' - github.com/air-verse/air)
run-dev-win:
	@echo "Starting Windows server with auto-reload..."
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	@cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 air

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
	@bash -c 'read -r'
	@psql $(DATABASE_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" || true
	@echo "✓ Database reset complete"

# Start PostgreSQL with Docker
db-up:
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
db-down:
	@echo "Stopping PostgreSQL container..."
	@docker stop trip_manager_db 2>/dev/null || true
	@docker rm trip_manager_db 2>/dev/null || true
	@echo "✓ PostgreSQL stopped"

# Start MinIO with Docker
storage-up:
	@echo "Starting MinIO with Docker..."
	@docker run -d \
		--name trip_manager_minio \
		-e MINIO_ROOT_USER=minioadmin \
		-e MINIO_ROOT_PASSWORD=minioadmin \
		-p 9000:9000 \
		-p 9001:9001 \
		-v trip_manager_minio_data:/data \
		quay.io/minio/minio:latest server /data --console-address ":9001"
	@echo "✓ MinIO started"
	@echo "  S3 API:       http://localhost:9000"
	@echo "  Console:      http://localhost:9001"
	@echo "  Credentials:  minioadmin:minioadmin"
	@sleep 3
	@echo "⏳ Waiting for MinIO to be ready..."
	@until curl -s http://localhost:9000/minio/health/live > /dev/null 2>&1; do sleep 1; done
	@echo "✓ MinIO is ready!"

# Setup MinIO bucket (creates if not exists)
storage-setup: storage-up
	@echo ""
	@echo "Setting up MinIO bucket..."
	@command -v mc >/dev/null 2>&1 || { echo "Installing MinIO Client (mc)..."; curl -s https://dl.min.io/client/mc/release/linux-amd64/mc -o /tmp/mc && chmod +x /tmp/mc && sudo mv /tmp/mc /usr/local/bin/mc 2>/dev/null || echo "Please install mc manually: https://min.io/download#minio-client"; }
	@sleep 1
	@mc alias set minio http://localhost:9000 minioadmin minioadmin 2>/dev/null || true
	@mc mb minio/trip-manager 2>/dev/null || echo "✓ Bucket already exists"
	@echo "✓ MinIO bucket setup complete"
	@echo ""
	@echo "Ready to use MinIO:"
	@echo "  Access at:  http://localhost:9000"
	@echo "  Console at: http://localhost:9001"
	@echo "  Bucket:     trip-manager"

# Stop MinIO Docker container
storage-down:
	@echo "Stopping MinIO container..."
	@docker stop trip_manager_minio 2>/dev/null || true
	@docker rm trip_manager_minio 2>/dev/null || true
	@echo "✓ MinIO stopped"

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
dev-setup: db-up db-setup storage-setup build
	@echo ""
	@echo "✓ Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Set environment variables: export DATABASE_URL=postgres://postgres:postgres@localhost:5432/trip_manager"
	@echo "  2. Start the server: make run"
	@echo "     or with auto-reload: make run-dev"
	@echo ""

# Start everything (DB + Server)
dev: storage-down db-down storage-up db-up
	@echo "Database started. Starting server..."
	@$(MAKE) run

# Version info
version:
	@echo "Trip Manager Backend"
	@cd $(BACKEND_DIR) && go version

# All help and info
info: version help
	@echo ""
	@echo "Environment Variables (Backend):"
	@echo "  DATABASE_URL=postgres://postgres:postgres@localhost:5432/trip_manager (required)"
	@echo "  SERVER_PORT=8000 (default)"
	@echo "  JWT_SECRET=your-secret-key (default: your-secret-key-change-in-production)"
	@echo "  ENVIRONMENT=development (default)"
	@echo ""
	@echo "Storage Configuration:"
	@echo "  STORAGE_TYPE=local (default) | s3"
	@echo ""
	@echo "  For Local Storage:"
	@echo "    UPLOAD_DIR=./uploads (default)"
	@echo ""
	@echo "  For S3 Storage (MinIO/AWS):"
	@echo "    S3_ENDPOINT=http://minio:9000 (MinIO) or empty (AWS)"
	@echo "    S3_BUCKET=trip-manager (default)"
	@echo "    S3_REGION=us-east-1 (default)"
	@echo "    S3_ACCESS_KEY=minioadmin"
	@echo "    S3_SECRET_KEY=minioadmin"
	@echo "    S3_PUBLIC_URL=http://localhost:9000"
	@echo "    S3_USE_SSL=false (default)"
	@echo ""

