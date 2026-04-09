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

👉 **[Full Setup Guide →](./docs/SETUP.md)**

---

## Technologies

### Frontend
- **Next.js** (React + TypeScript)
- Auto-generated types from OpenAPI spec

### Backend
- **Go** (1.24+)
- Auto-generated server code from OpenAPI spec
- PostgreSQL database
- JWT-based authentication

See the following guides for further information:
- **[Authentication Setup](./backend/AUTH.md)** - JWT configuration and usage
- **[Makefile](./docs/MAKEFILE.md)** - Development environment setup and common tasks
- 

### API Specification
- **OpenAPI 3.0.0** (`api-spec/openapi.yaml`)
- Shared contract between frontend & backend

Checkout: **[API Spec](./api-spec/README.md)** for details on endpoints, request/response schemas, and how to extend the API.
See **[Code Generation](#code-generation)** for details on how to keep frontend and backend in sync with the API spec.
---

## Project Structure

```
trip-manager
├─ api-spec/           # OpenAPI specification (source of truth)
├─ frontend/           # Next.js app (auto-generated types in src/generated/)
├─ backend/            # Go API server (auto-generated code in internal/generated/)
├─ docs/               # Documentation
├─ deploy/             # Deployment scripts and configurations
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

Checkout **[Setup Guide](./docs/SETUP.md)** for detailed instructions on code generation and troubleshooting. as well as **[Generator Config](./docs/GENERATOR_CONFIG.md)**
---

## Deployment

Comprehensive deployment guides:

- **[Runbook](./docs/RUNBOOK.md)** - Complete deployment guide (Automated & Manual)
- **[CI/CD Pipeline](./docs/CI_CD.md)** - GitHub Actions workflow details

**Quick Deploy:**
- Merge PR to `main` → Automatic GitHub Actions deployment
- Manual deployment: `cd deploy/hetzner && ./manual-deploy.sh <server-ip> <user> <port>`

---

Developed by

Arlind Shala, André Königer & Reyhan Karamahmut
