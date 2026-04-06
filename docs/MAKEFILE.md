# Makefile Quick Reference

A comprehensive Makefile is included in the project root to make development easier.

## Start Here

```bash
# One-command setup (installs everything)
make dev-setup

# Then start the server
make run
```

## Essential Commands

### Building & Running

```bash
make build              # Build the backend binary
make run                # Build and run the server
make run-dev            # Run with auto-reload (hot reload on code changes)
make clean              # Remove built binaries
```

### Testing

```bash
make test               # Run all unit tests
make test-verbose       # Run tests with detailed output
make test-coverage      # Generate code coverage report
```

### Database

```bash
make docker-up          # Start PostgreSQL in Docker
make docker-down        # Stop PostgreSQL Docker
make docker-logs        # Watch database logs
make db-setup           # Create database (if needed)
make db-reset           # ⚠️ Reset database (deletes all data!)
```

### Code Quality

```bash
make fmt                # Format code with 'go fmt'
make lint               # Run linter (golangci-lint)
make deps-check         # Verify dependencies
```

### Help & Info

```bash
make help               # Show all available commands
make info               # Show version + help + environment info
make version            # Show Go version
```

## Automatic Migrations

**Zero-config database setup!** Migrations run automatically when you start the server:

```bash
make run
# Server starts and automatically runs all migrations from backend/migrations/
```

## Development Environment Setup

```bash
# One-time setup
make dev-setup

# This:
# 1. Starts PostgreSQL in Docker
# 2. Creates the database
# 3. Builds the binary
# 4. Prints next steps

# Then:
make run        # Start the server
make run-dev    # Or with auto-reload for development
```

## Typical Development Workflow

```bash
# Terminal 1: Watch database logs (optional)
make docker-logs

# Terminal 2: Run with auto-reload
make run-dev

# Terminal 3: Run tests
make test

# Terminal 4: Format + lint code
make fmt && make lint
```

## Environment Variables

The Makefile uses these environment variables:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5432/trip_manager"
JWT_SECRET="your-secret-key-change-in-production"
SERVER_PORT="8000"
ENVIRONMENT="development"
```

For `make docker-up` to work with the default DATABASE_URL, no configuration needed!

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `command not found: make` | Install: `apt install make` (Linux) or `brew install make` (Mac) |
| Database connection failed | Run `make docker-up` to start PostgreSQL |
| Build fails | Run `make clean-all && make build` |
| Port already in use | Change: `export SERVER_PORT=8081 && make run` |

## See Also

- [SETUP.md](./SETUP.md) - Full setup guide
- [Makefile](./Makefile) - Complete source code with all 20+ commands
- [backend/AUTH.md](./backend/AUTH.md) - Authentication documentation

---

👈 **[Back to README](../README.md)**

