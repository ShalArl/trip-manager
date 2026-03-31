# CI/CD Pipeline

Diese GitHub Actions Pipeline automatisiert das Bauen und Pushen von Docker Images.

## Features

### ✅ Change Detection
- Nur Services mit Änderungen werden gebaut
- Spart CI/CD Zeit und Ressourcen
- Überwacht:
  - `backend/**` → Backend build
  - `frontend/**` → Frontend build
  - `api-spec/**` → Beide Services
  - `.github/workflows/deploy.yaml` → Beide Services

### ✅ Multi-Service Builds
- **Backend**: Go API
- **Frontend**: Next.js mit NEXT_PUBLIC_API_URL Injection

### ✅ Docker Registry
- **Registry**: GitHub Container Registry (ghcr.io)
- **Automatisches Tagging**:
  - `latest` (auf main branch)
  - Branch name (e.g., `dev`, `main`)
  - Git SHA (eindeutig pro commit)
  - Semver tags (bei Release Tags)

### ✅ Build Optimization
- Layer caching mit GitHub Actions Cache
- `.dockerignore` Dateien für kleinere Image Größen
- Parallel Builds für Backend und Frontend

## Environment Variables

### Secrets (GitHub Repository Settings erforderlich)

```
NEXT_PUBLIC_API_URL  # Optional, Standard: http://localhost:8000
```

Gehe zu: **Settings → Secrets and variables → Actions**

## Tagging Schema

### Bei Push auf `main`:
```
ghcr.io/your-org/trip-manager/backend:latest
ghcr.io/your-org/trip-manager/backend:main
ghcr.io/your-org/trip-manager/backend:sha-abc123
```

### Bei Push auf `dev`:
```
ghcr.io/your-org/trip-manager/backend:dev
ghcr.io/your-org/trip-manager/backend:sha-abc123
```

### Bei Git Tag `v1.0.0`:
```
ghcr.io/your-org/trip-manager/backend:1.0.0
ghcr.io/your-org/trip-manager/backend:1.0
ghcr.io/your-org/trip-manager/backend:sha-abc123
```

## Zugriff auf Images

### Frontend
```bash
docker pull ghcr.io/your-org/trip-manager/frontend:latest
docker run -e NEXT_PUBLIC_API_URL=https://api.example.com ghcr.io/your-org/trip-manager/frontend:latest
```

### Backend
```bash
docker pull ghcr.io/your-org/trip-manager/backend:latest
docker run -e DATABASE_URL="postgresql://..." ghcr.io/your-org/trip-manager/backend:latest
```

## Workflow

1. **Push zu main/dev**
2. **Change Detection**: Welche Services geändert?
3. **Build Backend** (wenn Änderungen in `backend/`)
4. **Build Frontend** (wenn Änderungen in `frontend/`)
5. **Push zu ghcr.io** (automatisch)
6. **Summary**: Status report

## Lokales Testen

Um Dockerfiles lokal zu bauen:

```bash
# Backend
cd backend && docker build -t trip-manager-backend:dev .

# Frontend (mit API URL)
cd frontend && docker build \
  --build-arg NEXT_PUBLIC_API_URL=http://localhost:8000 \
  -t trip-manager-frontend:dev .
```

## Troubleshooting

### ❌ "permission denied" beim Push
- Prüfe: **Settings → Actions → General → Workflow permissions**
- Setze auf: "Read and write permissions"

### ❌ "NEXT_PUBLIC_API_URL" ist undefined im Frontend
- Secret in GitHub hinterlegt? **Settings → Secrets**
- Build Argument wird übergeben? Siehe workflow YAML

### ❌ "Broken Dependencies" beim Build
- Cache invalidieren: GitHub Actions → Caches → Delete all
- oder: Manuell auf "Re-run all jobs" klicken

## Nächste Schritte

- Deployment zu Kubernetes/Docker Swarm
- Image Scanning (Trivy, Snyk)
- Notifications bei Build Failure (Slack, Email)

