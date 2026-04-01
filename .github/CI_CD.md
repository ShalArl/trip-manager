# CI/CD Pipeline

This GitHub Actions pipeline automates building, pushing, and deploying this project.

It performs the following steps:
- Detection of changes in backend, frontend, or API specs
- Parallel building of backend and frontend Docker images based on detected changes
- Pushing built images to GitHub Container Registry (ghcr.io) with automatic tagging based on branch, Git SHA, and semver tags
- Deployment of the application to a Hetzner Cloud Server via `docker-compose` with `Caddy` (optional, only on main/dev)


## Features

### ✅ Change Detection
- Only services with changes are built
- Saves CI/CD time and resources
- Monitors:
  - `backend/**` → Backend build
  - `frontend/**` → Frontend build
  - `api-spec/**` → Both services
  - `.github/workflows/deploy-iaas.yaml` → Both services
- Uses: 
  - [actions/checkout@v6](https://github.com/actions/checkout)
  - [dorny/paths-filter@v4](https://github.com/dorny/paths-filter)


### ✅ Optimized Docker Builds
- **Registry**: GitHub Container Registry (ghcr.io)
- **Automatic Tagging**:
  - `latest` for deployments and `sha-abc123` for unique identification
  - Branch name (e.g., `dev`, `main`)
  - Git SHA (unique per commit)
  - Semver tags (on release tags)
  - Layer caching for faster builds
  - Parallel builds for backend and frontend (based on detected changes)
  - `.dockerignore` files for smaller image sizes
  - Uses: 
    - [docker/setup-buildx-action@v4](https://github.com/docker/setup-buildx-action)
    - [docker/login-action@v4](https://github.com/docker/login-action)
    - [docker/metadata-action@v6](https://github.com/docker/metadata-action)
    - [docker/build-push-action@v7](https://github.com/docker/build-push-action)

### SSH Deployment
- Automatic SSH access to server
- Deployment and startup of `docker-compose` with `Caddy` as reverse proxy via `scp` and `ssh` 
- Secure handling of SSH key via GitHub Secrets
- Uses: [webfactory/ssh-agent@v0.9.0](https://github.com/webfactory/ssh-agent)

## Environment Variables

### Secrets (GitHub Repository Settings required)

```
NEXT_PUBLIC_API_URL    # Optional, default: http://localhost:8000
GITHUB_TOKEN           # Automatically provided by GitHub for API access in Actions
DB_NAME                # Database name for backend
DB_PASSWORD            # Database password for backend
DB_USER                # Database user for backend
DOMAIN                 # Domain for Caddy reverse proxy (e.g., example.com)
GHCR_PAT               # Personal Access Token with read:packages permission for ghcr.io
GHCR_USERNAME          # GitHub username for ghcr.io authentication (lowercase!!!)
SERVER_DEPLOY_KEY      # SSH private key for server access
SERVER_HOST            # Server IP or hostname
SERVER_PORT            # SSH port
SERVER_USER            # SSH username for server
JWT_SECRET             # JWT secret key for backend
```

Go to: **Settings → Secrets and variables → Actions**

## Tagging Schema

### On push to `main`:
```
ghcr.io/your-org/trip-manager/backend:latest
ghcr.io/your-org/trip-manager/backend:main
ghcr.io/your-org/trip-manager/backend:sha-abc123
```

### On push to `dev`:
```
ghcr.io/your-org/trip-manager/backend:latest
ghcr.io/your-org/trip-manager/backend:dev
ghcr.io/your-org/trip-manager/backend:sha-abc123
```

### On Git tag `v1.0.0`:
```
ghcr.io/your-org/trip-manager/backend:latest
ghcr.io/your-org/trip-manager/backend:1.0.0
ghcr.io/your-org/trip-manager/backend:1.0
ghcr.io/your-org/trip-manager/backend:sha-abc123
```

## Accessing Images

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

1. **Push to feature-branch**
2. **Pull request on main/dev**
3. **Change Detection**: Which services changed?
4. **Build Backend** (if changes in `backend/`)
5. **Build Frontend** (if changes in `frontend/`)
6. **Push to ghcr.io** (automatic)
7. **Summary**: Status report
8. **Deployment**: Automatic deployment to server (on push to `main` or `dev`)


## Troubleshooting

### docker-compose error pulling images `frontend` & `backend`
- Ensure that the secret `GHCR_PAT` (Personal Access Token with `read:packages` permission) is correctly set.
- Ensure that the username in `GHCR_USERNAME` is lowercase, regardless of capitalization on GitHub.
  - Example: `GHCR_USERNAME=your-github-username` (not `Your-GitHub-Username`)
