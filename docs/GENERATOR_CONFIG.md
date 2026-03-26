# Code Generator Configuration Guide

## Overview

This project uses **two specialized code generators** to automatically generate API types from the OpenAPI specification:

- **TypeScript Frontend:** `openapi-typescript` (types-only)
- **Go Backend:** `oapi-codegen` (types-only generation)

Generation is orchestrated via **Turbo** for caching and parallel execution across workspaces.

## Quick Start

```bash
# Generate all code (TS + Go)
pnpm gen

# Or individually
pnpm gen:ts      # TypeScript client types only
pnpm gen:go      # Go models only

# Clean all generated code
pnpm clean
```

---

## Configuration Files

### 1. OpenAPI Specification (Source of Truth)

**File:** `api-spec/openapi.yaml`

This is the **single source of truth** for the API contract. All generated code derives from this.

Current resources:
- **Trips:** CRUD operations, list with pagination
- **Locations:** List, add, delete with coordinates and sequence
- **Activities:** CRUD operations with categories and timing

### 2. Root Scripts: `package.json`

```json
{
  "scripts": {
    "gen": "turbo run gen",
    "gen:ts": "mkdir -p ./frontend/generated && openapi-typescript ./api-spec/openapi.yaml -o ./frontend/generated/types.ts",
    "gen:go": "mkdir -p ./backend/internal/generated && rm -rf ./backend/internal/generated/*.go && oapi-codegen -generate types -package api ./api-spec/openapi.yaml > ./backend/internal/generated/models.go",
    "clean": "rm -rf ./frontend/generated/types.ts ./backend/internal/generated/models.go"
  }
}
```

---

## Generator Tools

### TypeScript: openapi-typescript

**Tool:** `openapi-typescript` v6.7.5
**Mode:** Types-only generation (no HTTP client)
**Output:** `frontend/generated/types.ts`

**Generates:**
- TypeScript interfaces for all data models
- Type definitions for all API operations

**Configuration:**
- Input: `./api-spec/openapi.yaml`

### Go: oapi-codegen

**Tool:** `github.com/oapi-codegen/oapi-codegen/v2`
**Mode:** `-generate types` (models-only)
**Output:** `backend/internal/generated/models.go`

**Generates:**
- Go structs for all API models (Trip, Activity, Location, etc.)
- Request/Response types (CreateTripRequest, UpdateActivityRequest, etc.)
- Enum types with validation methods
- Parameter types (ListTripsParams, etc.)

**Why models-only?**
- Clean separation: Generated code = data types only
- No HTTP scaffolding or client code
- You control routing via Chi router in `backend/cmd/api/main.go`
- Flexible: Handlers can be implemented independently

---

## Generated Code Structure

### Frontend (TypeScript)

```
frontend/generated/
└── types.ts                       # Ignored by git
    ├── Trip interface
    ├── Activity interface
    ├── Location interface
    ├── CreateTripRequest interface
    ├── CreateActivityRequest interface
    ├── ActivityCategory type
    ├── TripStatus type
    └── ... (all types)
```

### Backend (Go)

```
backend/internal/generated/
└── models.go                       # Ignored by git
    ├── Activity struct
    ├── Trip struct
    ├── Location struct
    ├── CreateActivityRequest struct
    ├── ActivityCategory enum
    ├── TripStatus enum
    ├── ListTripsParams type
    └── ... (all models with JSON tags)
```

**Example model:**
```go
// Using generated models in your handlers
type Activity struct {
    Id          *openapi_types.UUID `json:"id,omitempty"`
    Name        *string             `json:"name,omitempty"`
    Category    *ActivityCategory   `json:"category,omitempty"`
    Date        *openapi_types.Date `json:"date,omitempty"`
    Cost        *float32            `json:"cost,omitempty"`
    Currency    *string             `json:"currency,omitempty"`
}
```

---

## Custom Implementation

### Do NOT Edit Generated Files

Generated files have headers indicating automatic generation. Changes are **overwritten on next generation**.

### Keep Custom Code Separate

```
backend/
├── cmd/api/
│   └── main.go                 # ← Entry point (yours)
├── internal/
│   ├── generated/
│   │   └── models.go           # ← Generated (ignore)
│   ├── handlers/               # ← Your implementations
│   │   ├── trips.go
│   │   ├── activities.go
│   │   └── locations.go
│   ├── services/               # ← Business logic
│   │   └── trip_service.go
│   └── middleware/             # ← Custom middleware
│       └── cors.go
└── go.mod
```

**Example handler implementation:**
```go
// backend/internal/handlers/trips.go
package handlers

import (
    "github.com/ShalArl/trip-manager/internal/generated"
)

func (h *Handler) ListTrips(ctx context.Context, params generated.ListTripsParams) (*generated.ListTrips200Response, error) {
    // Your business logic here
    trips := // ... fetch from database
    return &generated.ListTrips200Response{
        Items: trips,
    }, nil
}
```

**Frontend service layer:**
```typescript
// frontend/api/trips.ts
import type { Trip, CreateTripRequest, ListTripsParams } from '@/generated/types'

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export async function listTrips(params?: ListTripsParams) {
  const query = new URLSearchParams()
  if (params?.limit) query.append('limit', String(params.limit))
  if (params?.offset) query.append('offset', String(params.offset))

  const response = await fetch(`${API_BASE}/trips?${query}`)
  if (!response.ok) throw new Error(`Failed to fetch trips: ${response.statusText}`)
  return response.json() as Promise<Trip[]>
}

export async function createTrip(req: CreateTripRequest) {
  const response = await fetch(`${API_BASE}/trips`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!response.ok) throw new Error(`Failed to create trip: ${response.statusText}`)
  return response.json() as Promise<Trip>
}
```

---

## Git Strategy

### What's Committed

**Nothing from generated directories** - all generated code is in `.gitignore`:
```gitignore
/frontend/generated/types.ts
/backend/internal/generated/models.go
```

### What's Ignored

**Generated code files:**
```gitignore
# Ignore generated files in frontend
/frontend/generated/types.ts

# Ignore generated files in backend
/backend/internal/generated/models.go
```

**Benefits:**
- ✅ Generated code is never committed (clean repo)
- ✅ New clones regenerate: `pnpm install && pnpm gen && pnpm build`
- ✅ No merge conflicts on generation changes
- ✅ Smaller git history

---

## Workflow

### 1. Modify API Specification

Edit `api-spec/openapi.yaml`:
```yaml
/trips/{tripId}/comments:          # New endpoint
  post:
    operationId: addComment
    requestBody:
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/CreateCommentRequest'
    responses:
      '201':
        description: Comment added
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Comment'
```

### 2. Regenerate Code

```bash
# Turbo detects spec changes and regenerates automatically
pnpm gen
```

This generates:
- New TypeScript types in `frontend/generated/types.ts`
- New Go models in `backend/internal/generated/models.go`

### 3. Implement Backend Handlers

Add handler in `backend/internal/handlers/`:
```go
func (h *Handler) AddComment(ctx context.Context, tripId string, req generated.CreateCommentRequest) (*generated.Comment, error) {
    // Your business logic
    comment := &generated.Comment{
        Id:      uuid.New(),
        TripId:  tripId,
        Text:    req.Text,
    }
    return comment, nil
}
```

### 4. Mount Handlers in Router

Update `backend/cmd/api/main.go`:
```go
r := chi.NewRouter()

// Mount your handlers
h := handlers.NewHandler(db, logger)
r.Post("/trips/{tripId}/comments", h.AddComment)

http.ListenAndServe(":8080", r)
```

### 5. Use in Frontend

Create `frontend/api/trips.ts`:
```typescript
import type { Comment } from '@/generated/types'

export async function addComment(tripId: string, text: string) {
    const response = await fetch(`/api/trips/${tripId}/comments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text }),
    })
    if (!response.ok) throw new Error('Failed to add comment')
    return response.json() as Promise<Comment>
}
```

### 6. Test & Commit

```bash
# Verify clean generation
pnpm gen

# Build both packages
pnpm build

# Commit your custom code (generated files are ignored)
git add backend/internal/handlers/ frontend/api/
git commit -m "feat: add comment endpoint"
```

---

## Cleaning Generated Code

To remove all generated files:

```bash
pnpm clean
```

This removes:
- `backend/internal/generated/models.go`
- `frontend/generated/types.ts`

Run `pnpm gen` to regenerate them.

---

## Changing Generators or Configuration

### Upgrading openapi-typescript

```bash
pnpm add -D openapi-typescript@latest
```

### Installing New Tools

```bash
# Update pnpm/Node packages
pnpm install

# Install Go tools globally
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

---

## Troubleshooting

### Generation Fails: "Directory nonexistent"

The script now creates directories automatically with `mkdir -p`:
```json
"gen:ts": "mkdir -p ./frontend/generated && openapi-typescript ...",
"gen:go": "mkdir -p ./backend/internal/generated && rm -rf ... && oapi-codegen ..."
```

If issue persists:
```bash
mkdir -p backend/internal/generated frontend/generated
pnpm gen
```

### "Cannot find module openapi-typescript"

Ensure packages are installed:
```bash
pnpm install
```

### Generated Files Not Updated

Clear and regenerate:
```bash
pnpm clean
pnpm gen:ts
```

### OpenAPI YAML Syntax Errors

Validate your spec:
```bash
# Using openapi-typescript's built-in validation
pnpm gen:ts
```

If you see parsing errors, check:
- YAML indentation (must be 2 spaces)
- Schema definitions under `components.schemas`
- Path definitions have proper operationId

### Go Compilation Errors

Verify dependencies:
```bash
cd backend
go mod tidy
go build ./cmd/api
```

### Import Path Issues in Go

Verify `go.mod` has correct module name:
```bash
cat backend/go.mod
# Should be: module github.com/ShalArl/trip-manager
```

Generated code imports from:
```go
import "github.com/oapi-codegen/runtime/types"
```

### "oapi-codegen: not found"

Ensure oapi-codegen is installed globally:
```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
which oapi-codegen  # Verify installation
```

---

## Dependencies

### Root Level

```json
{
  "devDependencies": {
    "turbo": "^2.0.0",
    "openapi-typescript": "^6.7.5",
    "prettier": "latest"
  }
}
```

### Backend Go Module

```
github.com/go-chi/chi/v5 v5.0.11           # Router
github.com/oapi-codegen/runtime v1.3.0     # Runtime types for generated code
google.golang.org/uuid                      # UUID generation
```

Added automatically by `go mod tidy` after generation.

---

## References

- [openapi-typescript Documentation](https://openapi-ts.dev/)
- [oapi-codegen - Go Code Generator](https://github.com/oapi-codegen/oapi-codegen)
- [OpenAPI 3.0 Specification](https://spec.openapis.org/oas/v3.0.3)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [Turbo Caching](https://turbo.build/repo/docs/core-concepts/caching)
