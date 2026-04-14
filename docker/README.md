# Docker Initialization Scripts

This directory contains initialization scripts for Docker services.

## Contents

### `minio-init.sh`
Initializes MinIO S3-compatible storage with:
- Bucket creation (if doesn't exist)
- Public read policy configuration
- MinIO client (mc) setup

**How it works:**
1. The `minio-init` service in docker-compose.yaml runs after MinIO is healthy
2. Executes the init script which configures the bucket automatically
3. No manual setup required after starting containers

**Usage:**
```bash
docker-compose up
```

The initialization happens automatically. You can then:
- Access MinIO console: http://localhost:9001
- Use MinIO CLI: `mc ls local`

**For development:**
- Bucket name: `trip-manager` (configurable via `S3_BUCKET` env var)
- Access key: `minioadmin` (configurable via `S3_ACCESS_KEY`)
- Secret key: `minioadmin` (configurable via `S3_SECRET_KEY`)

## Production Considerations

For production deployments:
1. Replace MinIO with AWS S3, Google Cloud Storage, or similar
2. Update `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY` environment variables
3. Set `S3_USE_SSL: true` for HTTPS connections
4. Use IAM roles instead of hardcoded credentials
5. The `minio-init` service can be removed or modified for cloud providers

