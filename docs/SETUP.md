# Setup & Development Guide

## 🚀 Quick Start with Makefile

### Fastest Way to Get Started

```bash
cd trip-manager

# Automated setup: Docker + database + build
make dev-setup

# Start the server
make run
```

### Available Makefile Commands

```bash
make help              # Show all commands
make build             # Build the backend binary
make run               # Build and run
make run-dev           # Run with auto-reload
make test              # Run tests
make docker-up         # Start PostgreSQL
make docker-down       # Stop PostgreSQL
make clean             # Clean binaries
```

See `Makefile` for all 20+ commands!

---

## ✅ Automatic Database Migrations

**No manual migrations needed!** The server automatically runs all `.sql` files from `backend/internal/database/migrations/` on startup:

```bash
make run
# Output:
# Found 1 migration files
# Running migration: 001_init_schema.sql
# ✓ Migration 001_init_schema.sql completed
# All migrations completed successfully
```

### Adding New Migrations

```bash
# Create new migration file
touch backend/internal/database/migrations/002_my_migration.sql

# Add your SQL
# Write migration SQL (uses CREATE IF NOT EXISTS for safety)

# Next server start automatically runs it!
make run
```

---

## Prerequisites

- **Node.js** 18+ & pnpm 9.0.0+
- **Go** 1.24+ (für Backend)
- **Docker** (optional, für Containerisierung)

## Traditional Setup (without Makefile)

## Quick Start

### 1. Install Dependencies

```bash
cd trip-manager
pnpm install
```

This installs:
- Turbo for monorepo management
- OpenAPI code generators (TypeScript & Go)
- Dev tools (Prettier, ESLint)

### 2. Initialize Generated Files

```bash
pnpm gen
```

This auto-generates:
- **TypeScript types** → files in `frontend/src/generated/`
- **Go server code** → `backend/internal/generated/models.go`

### 3. Start Development

```bash
pnpm dev
```

Starts all workspaces in parallel:
- Frontend: Next.js dev server (localhost:3000)
- Backend: Go server (localhost:8000)

---

## Development Workflow

### Code Generation

When you modify `api-spec/openapi.yaml`:

```bash
# Generate from updated spec
pnpm gen

# Or individually
pnpm gen:ts   # TypeScript only
pnpm gen:go   # Go only
```

Turbo automatically invalidates cache when `api-spec/openapi.yaml` changes.

### Frontend Development

```bash
cd frontend
pnpm dev        # Start Next.js dev server
pnpm build      # Production build
pnpm lint       # Run ESLint
```

**Generated Types** (`frontend/src/generated/types.ts`):
- Auto-generated TypeScript interfaces from OpenAPI spec
- Use directly in your API service layer (build your own fetch wrapper)

### Backend Development

```bash
cd backend
go run ./cmd/api   # Start API server
go test ./...      # Run tests
go mod tidy        # Manage dependencies
```

**Generated Models** (`backend/internal/generated/models.go`):
- Auto-generated Go structs from OpenAPI spec
- Import and use directly in your handlers

---

## Local Setup for Code Generation

### Install Go Code Generator

```bash
# Install oapi-codegen globally (needed for Go type generation)
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Verify installation
which oapi-codegen
```

### Generate Code

```bash
# Generate both TypeScript and Go
pnpm gen

# Or individually
pnpm gen:ts   # TypeScript types only
pnpm gen:go   # Go models only
```

---

## Project Structure

```
trip-manager/
├── api-spec/             # OpenAPI specification
│   ├── openapi.yaml      # API contract (source of truth)
│   └── README.md         # API documentation
│
├── frontend/             # Next.js TypeScript app
│   ├── generated/        # Auto-generated from openapi.yaml
│   ├── src/
│   │   ├── app/          # Next.js app directory
│   │   └── components/   # React components
│   └── package.json
│
├── backend/              # Go API server
│   ├── internal/
│   │   └── generated/    # Generated Types
│   ├── cmd/api/          # Server entrypoint
│   ├── go.mod
│   └── go.sum
│
├── docs/                 # Additional documentation
├── package.json          # Root workspace config
├── pnpm-workspace.yaml   # pnpm monorepo config
└── turbo.json            # Turbo build tasks
```

---

## Build & Deployment

### Development Build

```bash
pnpm build
```

Runs build in all workspaces, respecting dependencies defined in turbo.json.

### Production Build

```bash
# Build and lint
pnpm build
pnpm lint

# Docker (optional)
docker build -f backend/Dockerfile -t trip-manager-api ./backend
docker build -f frontend/Dockerfile -t trip-manager-web ./frontend
```

---

## Troubleshooting

### "Cannot find openapi-typescript" or similar

```bash
pnpm install
```

Reinstall dependencies if code generators are missing.

### Generated files not updating

```bash
# Clean generated files and regenerate
pnpm clean
pnpm gen
```

### "oapi-codegen: not found"

Ensure oapi-codegen is installed globally:
```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
which oapi-codegen  # Verify installation
```

### Port conflicts

- Frontend default: `3000` (set `PORT=3001 pnpm dev` to override)
- Backend default: `8000` (check `backend/cmd/api/main.go`)

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v2
      - uses: actions/setup-node@v4
        with:
          node-version: 18
          cache: 'pnpm'
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: pnpm install
      - run: pnpm gen
      - run: pnpm build
      - run: pnpm lint
```

---

## Next Steps

1. **Start Frontend**: `cd frontend && pnpm dev`
2. **Start Backend**: `cd backend && go run ./cmd/api`
3. **Modify API Spec**: Edit `api-spec/openapi.yaml`
4. **Regenerate Code**: `pnpm gen`

Questions? Check the individual README files in `frontend/`, `backend/`, and `api-spec/`.

---

👈 **[Back to README](../README.md)**

