# Setup & Development Guide

## Prerequisites

- **Node.js** 18+ & pnpm 9.0.0+
- **Go** 1.21+ (für Backend)
- **Docker** (optional, für Containerisierung)

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
- **TypeScript types** → `frontend/src/generated/types.ts`
- **Go server code** → `backend/internal/api/generated.go`

### 3. Start Development

```bash
pnpm dev
```

Starts all workspaces in parallel:
- Frontend: Next.js dev server (localhost:3000)
- Backend: Go server (localhost:8080)

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

**Generated API Client** (`frontend/src/generated/types.ts`):
- Auto-generated TypeScript types from OpenAPI spec
- Use directly in API service layer

### Backend Development

```bash
cd backend
go run ./cmd/api   # Start API server
go test ./...      # Run tests
go mod tidy        # Manage dependencies
```

**Generated Server Code** (`backend/internal/api/generated.go`):
- Auto-generated Go types & HTTP client
- Implement interfaces as needed

---

## Local Setup for Go Code Generation

### Option 1: Using oapi-codegen (Recommended)

```bash
# Install oapi-codegen globally
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Generate Go code
pnpm gen:go
```

### Option 2: Using openapi-generator-cli

Already included in `package.json`. No additional setup needed.

---

## Project Structure

```
trip-manager/
├── api-spec/             # OpenAPI specification
│   ├── openapi.yaml      # API contract (source of truth)
│   └── README.md         # API documentation
│
├── frontend/             # Next.js TypeScript app
│   ├── src/
│   │   ├── generated/    # Auto-generated from openapi.yaml
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
# Clear Turbo cache and regenerate
rm -rf .turbo
pnpm gen
```

### Go code generation fails

Ensure `oapi-codegen` is installed:
```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
which oapi-codegen  # Verify installation
```

### Port conflicts

- Frontend default: `3000` (set `PORT=3001 pnpm dev` to override)
- Backend default: `8080` (check `backend/cmd/api/main.go`)

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
