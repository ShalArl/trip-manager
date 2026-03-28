# Trip Manager

Developed for Cloud Application Development @HTWG Konstanz in SS26.

**Goal:**
Develop an application for managing trips for leisure travellers.
The trips should be managed as a list, which can be accessed and modified

---

## Quick Start

```bash
# Install & generate code
pnpm install
pnpm gen

# Start development
pnpm dev
```

👉 **[Full Setup Guide →](./SETUP.md)**

---

## Technologies

### Frontend
- **Next.js** (React + TypeScript)
- Auto-generated types from OpenAPI spec

### Backend
- **Go** (1.24+)
- Auto-generated server code from OpenAPI spec

### API Specification
- **OpenAPI 3.0.0** (`api-spec/openapi.yaml`)
- Shared contract between frontend & backend

---

## Project Structure

```
trip-manager
├─ api-spec/           # OpenAPI specification (source of truth)
├─ frontend/           # Next.js app (auto-generated types in src/generated/)
├─ backend/            # Go API server (auto-generated code in internal/api/)
├─ docs/               # Documentation
├─ package.json        # pnpm workspace config
├─ pnpm-workspace.yaml # Monorepo definition
├─ turbo.json          # Build tasks
└─ SETUP.md            # Detailed setup guide
```

---

## Code Generation

When `api-spec/openapi.yaml` changes, regenerate:

```bash
pnpm gen
```

This updates:
- `frontend/generated/types.ts` (TypeScript)
- `backend/internal/generated/models.go` (Go)

Turbo automatically invalidates cache on spec changes.

---

Developed by

Arlind Shala, André Königer & Reyhan Karamahmut
