# Docker Setup Guide

## Quick Start

The project now uses Docker Compose for development and supports production deployment.

### Development Setup

**Using Docker Compose (Recommended):**
```bash
docker-compose up
```

This starts all services:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8000
- **MinIO Console**: http://localhost:9001 (user: minioadmin, pass: minioadmin)
- **PostgreSQL**: localhost:5432

All initialization (database migrations, MinIO bucket setup) happens automatically!

**Using Makefile:**
```bash
make docker-up      # Start all services
make docker-logs    # View all logs
make docker-down    # Stop all services
```

### MinIO Initialization

MinIO now initializes automatically via the `minio-init` service:

1. Creates the `trip-manager` bucket (configurable via `S3_BUCKET`)
2. Sets public read access policy
3. Uses environment variables for configuration

**No manual setup required!** Just run `docker-compose up`.

### Environment Variables

Create a `.env` file in the project root:

```dotenv
# Database
DB_USER=trip_nugget
DB_PASSWORD=trip_nugget_superpw
DB_NAME=trip_nugget

# JWT
JWT_SECRET=your-secret-key-here

# S3 Storage (MinIO in development)
S3_BUCKET=trip-manager
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_PUBLIC_URL=http://localhost:9000/trip-manager

# Optional: For production AWS S3
# S3_ENDPOINT=https://s3.amazonaws.com
# S3_BUCKET=your-bucket-name
# S3_REGION=us-east-1
# S3_USE_SSL=true
```

## Docker Compose Services

### Frontend
- **Service name**: `frontend`
- **Port**: 3000
- **Image**: ghcr.io/${GHCR_USERNAME}/trip-manager/frontend

### Backend
- **Service name**: `backend`
- **Port**: 8000
- **Image**: ghcr.io/${GHCR_USERNAME}/trip-manager/backend
- **Dependencies**: database, minio

### PostgreSQL Database
- **Service name**: `database`
- **Port**: 5432
- **Image**: postgres:16-alpine

### MinIO S3 Storage
- **Service name**: `minio`
- **Ports**: 9000 (API), 9001 (Console)
- **Image**: quay.io/minio/minio:latest

### MinIO Initialization
- **Service name**: `minio-init`
- **Runs after**: minio (when healthy)
- **Creates**: Bucket and policies

## Useful Commands

```bash
# View all services
docker-compose ps

# View logs
docker-compose logs -f              # All services
docker-compose logs -f backend      # Specific service
docker-compose logs -f minio-init   # See initialization

# Stop services
docker-compose down

# Remove volumes (careful!)
docker-compose down -v

# Rebuild images
docker-compose build

# View service details
docker-compose exec backend /bin/sh
docker-compose exec database psql -U trip_nugget -d trip_nugget
```

## Production Considerations

For production:

1. **Use AWS S3 or similar** - Update `S3_ENDPOINT` and credentials
2. **Use managed database** - Don't use PostgreSQL container
3. **Use Caddy or reverse proxy** - Uncomment Caddy service in docker-compose
4. **Environment variables** - Use secure secrets management (GitHub Secrets, etc.)
5. **Remove MinIO** - Use cloud provider instead
6. **SSL/TLS** - Set `S3_USE_SSL: true` for HTTPS

## Troubleshooting

### MinIO bucket not created
```bash
# Check minio-init logs
docker-compose logs minio-init

# Manually run init
docker-compose exec minio-init sh /minio-init.sh
```

### Database connection failed
```bash
# Check database is healthy
docker-compose ps

# View database logs
docker-compose logs database

# Connect to database
docker-compose exec database psql -U trip_nugget -d trip_nugget
```

### Backend can't connect to services
- Verify service names are correct in environment variables
- Use service name (e.g., `minio:9000`) not `localhost:9000`
- Check network: `docker network ls`

## Migration from Old Setup

The old `/deploy/minio-setup.sh` scripts are deprecated. Everything is now handled by Docker Compose.

See `deploy/MINIO_MIGRATION.md` for details.

