
# LOCAL DEV Variables
BINARY_NAME=api
BACKEND_DIR=backend
CMD_DIR=$(BACKEND_DIR)/cmd/api
BIN_DIR=$(BACKEND_DIR)/bin
BINARY_NAME_WIN=$(BINARY_NAME).exe
STORAGE_TYPE?=s3

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
	@echo "Docker Commands (Recommended):"
	@echo "  make docker-up      Start all services with docker-compose"
	@echo "  make docker-down    Stop all services"
	@echo "  make docker-logs    View logs from all services"
	@echo "  make docker-logs-SERVICE  View logs for a specific service"
	@echo "Legacy Commands (Use docker-compose instead):"
	@echo "  make db-up          Start PostgreSQL (use docker-compose up database)"
	@echo "  make db-down        Stop PostgreSQL (use docker-compose down)"
	@echo "  make storage-up     Start MinIO (use docker-compose up minio minio-init)"
	@echo "  make storage-down   Stop MinIO"
	@echo ""
	@echo "Firebase Commands:"
	@echo "  make firebase-up    Start Firebase Emulators with Docker"
	@echo "  make firebase-down  Stop Firebase Emulators"
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
	@echo "Seeding Commands:"
	@echo "  make provision-only Provision users from CSV (requires running backend & Firebase)"
	@echo ""
	@echo "Cleanup Commands:"
	@echo "  make clean          Remove built binaries"
	@echo "  make clean-all      Remove binaries and generated files"
	@echo ""

# Setup MinIO bucket for local development
minio-setup:
	@echo "Setting up MinIO bucket..."
	@docker exec trip_manager_minio mc alias set local http://localhost:9000 minioadmin minioadmin --api S3v4
	@docker exec trip_manager_minio mc mb local/trip-manager --ignore-existing
	@docker exec trip_manager_minio mc anonymous set public local/trip-manager
	@echo "✓ MinIO bucket ready"

# Build the backend binary
build:
	@echo "Building backend binary..."
	@mkdir -p $(BIN_DIR)
	@cd $(BACKEND_DIR) && go build -o bin/$(BINARY_NAME) ./cmd/api
	@echo "✓ Build complete: $(BIN_DIR)/$(BINARY_NAME)"

# Run the compiled binary
run: build
	@echo "Starting server..."
	@STORAGE_TYPE=$(STORAGE_TYPE) FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 FIREBASE_PROJECT_ID=trip-manager-local ./$(BIN_DIR)/$(BINARY_NAME)

# Run with auto-reload (requires 'air' - github.com/air-verse/air)
run-dev:
	@echo "Starting server with auto-reload..."
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	@cd $(BACKEND_DIR) && STORAGE_TYPE=$(STORAGE_TYPE) FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 FIREBASE_PROJECT_ID=trip-manager-local air

# Build the backend binary for Windows
build-windows:
	@echo "Building backend binary for Windows..."
	@mkdir -p $(BIN_DIR)
	@cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME_WIN) ./cmd/api
	@echo "✓ Windows build complete: $(BIN_DIR)/$(BINARY_NAME_WIN)"

# Run the compiled Windows binary (requires WSL2 or Windows environment)
run-windows: build-windows
	@echo "Starting Windows server..."
	@STORAGE_TYPE=$(STORAGE_TYPE) FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 FIREBASE_PROJECT_ID=trip-manager-local ./$(BIN_DIR)/$(BINARY_NAME_WIN)

# Run Windows build with auto-reload (requires 'air' - github.com/air-verse/air)
run-dev-win:
	@echo "Starting Windows server with auto-reload..."
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	@cd $(BACKEND_DIR) && STORAGE_TYPE=$(STORAGE_TYPE) FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 FIREBASE_PROJECT_ID=trip-manager-local GOOS=windows GOARCH=amd64 air

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
	@echo ""
	@echo "To set up bucket automatically, use docker-compose:"
	@echo "  docker-compose up minio minio-init"

# Stop MinIO Docker container
storage-down:
	@echo "Stopping MinIO container..."
	@docker stop trip_manager_minio 2>/dev/null || true
	@docker rm trip_manager_minio 2>/dev/null || true
	@echo "✓ MinIO stopped"

# Start Firebase Emulators locally
firebase-up:
	@echo "Building Firebase Emulator Docker image..."
	@docker build -f firebase/Dockerfile -t trip-manager-firebase:latest .
	@echo "✓ Firebase image built"
	@echo ""
	@echo "Starting Firebase Emulators..."
	@docker run -d \
		--name trip_manager_firebase \
		-p 8080:8080 \
		-p 4000:4000 \
		-p 9099:9099 \
		-v $(PWD)/firebase:/firebase \
		trip-manager-firebase:latest
	@echo "✓ Firebase Emulators started"
	@echo "  Emulator UI:  http://localhost:4000"
	@echo "  Firestore:    http://localhost:8080"
	@echo "  Auth:         http://localhost:9099"
	@echo "⏳ Waiting for Firebase Emulators to be ready (this may take 30-60 seconds)..."
	@counter=0; \
	until [ $$counter -gt 120 ] || curl -s http://localhost:8080 > /dev/null 2>&1; do \
		counter=$$((counter+1)); \
		sleep 1; \
	done; \
	if [ $$counter -le 120 ]; then \
		echo "✓ Firebase Emulators are ready!"; \
	else \
		echo "⚠ Firebase Emulators may still be starting. Check logs with: docker logs trip_manager_firebase"; \
	fi

# Stop Firebase Emulators
firebase-down:
	@echo "Stopping Firebase Emulators container..."
	@docker stop --timeout=30 trip_manager_firebase 2>/dev/null || true
	@docker rm trip_manager_firebase 2>/dev/null || true
	@echo "✓ Firebase Emulators stopped"

# Start all services with docker-compose
docker-up:
	@echo "Starting all services with docker-compose..."
	@docker-compose up -d
	@echo "✓ All services started"
	@echo ""
	@echo "Services:"
	@echo "  Frontend:  http://localhost:3000"
	@echo "  Backend:   http://localhost:8000"
	@echo "  MinIO:     http://localhost:9000 (API) & http://localhost:9001 (Console)"
	@echo "  Firebase:  http://localhost:4000 (UI) & http://localhost:9099 (Auth) & http://localhost:8080 (Firestore)"
	@echo "  Database:  localhost:5432"

# Stop all services with docker-compose
docker-down:
	@echo "Stopping all services..."
	@docker-compose down
	@echo "✓ All services stopped"

# View docker-compose logs
docker-logs:
	@docker-compose logs -f

# View logs for a specific service
docker-logs-%:
	@docker-compose logs -f $*

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
dev-setup: db-up db-setup build
	@echo ""
	@echo "✓ Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Set environment variables: export DATABASE_URL=postgres://postgres:postgres@localhost:5432/trip_manager"
	@echo "  2. Start the server: make run"
	@echo "     or with auto-reload: make run-dev"
	@echo ""

# Start everything locally (DB + Storage + Firebase + Server)
dev: storage-down db-down firebase-down storage-up db-up firebase-up
	@echo "All services started. Starting server..."
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
	@echo "  STORAGE_TYPE=s3 (default: s3, options: s3, local)"
	@echo "  SERVER_PORT=8000 (default)"
	@echo "  JWT_SECRET=your-secret-key (default: your-secret-key-change-in-production)"
	@echo "  ENVIRONMENT=development (default)"
	@echo ""

# Provision users from CSV (requires running backend & Firebase)
provision-only:
	@echo "Starting user provisioning..."
	@cd tests/loadtests && $(MAKE) provision-only
	@echo "✓ Provisioning complete"

