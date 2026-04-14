# Storage Configuration Guide

## Overview

Trip Manager supports S3-compatible storage for file uploads (avatars, etc.):

- **Development**: MinIO (local S3-compatible storage in Docker)
- **Production**: AWS S3 or any S3-compatible service

## Quick Start with Docker Compose

MinIO is automatically initialized when you run:

```bash
docker-compose up
```

This automatically:
1. Starts MinIO service
2. Creates the `trip-manager` bucket
3. Sets public read permissions
4. Configures all settings

**That's it!** No additional setup needed.

---

## Development Setup with MinIO

MinIO is an S3-compatible object storage that runs in Docker:
- ✅ Same API as AWS S3
- ✅ Runs locally in a container
- ✅ Easy to switch to production S3 later
- ✅ No AWS account needed for development

### MinIO Credentials (Development)
```
Endpoint:     http://localhost:9000 (from your machine)
              http://minio:9000 (from containers)
Console:      http://localhost:9001
Access Key:   minioadmin
Secret Key:   minioadmin
Bucket:       trip-manager
Public URL:   http://localhost:9000/trip-manager
```

### Access MinIO Console

1. Open http://localhost:9001
2. Login with `minioadmin` / `minioadmin`
3. Browse buckets and files

---

## Production Setup with AWS S3

To use AWS S3 in production, update your `.env` file:

```dotenv
S3_ENDPOINT=https://s3.amazonaws.com
S3_BUCKET=your-production-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=your-aws-access-key
S3_SECRET_KEY=your-aws-secret-key
S3_PUBLIC_URL=https://your-bucket.s3.amazonaws.com
S3_USE_SSL=true
```

Then restart services:
```bash
docker-compose up
```

---

## Environment Variables

All storage configuration is done via environment variables:

### Required
```dotenv
S3_BUCKET=trip-manager                      # Bucket name
S3_ACCESS_KEY=minioadmin                    # Access key
S3_SECRET_KEY=minioadmin                    # Secret key
```

### Optional (with defaults)
```dotenv
S3_ENDPOINT=http://minio:9000              # S3 endpoint
S3_REGION=us-east-1                        # AWS region
S3_PUBLIC_URL=http://localhost:9000/trip-manager  # Public URL
S3_USE_SSL=false                            # Use HTTPS
STORAGE_TYPE=s3                             # Storage backend
```

---

## File Upload Flow

1. **Frontend**: User selects avatar → uploads to `/api/users/me` with multipart form data
2. **Backend**: Receives file → uploads to MinIO/S3
3. **Storage**: File stored in `avatars/` prefix
4. **URL Generation**: Backend returns public URL to frontend
5. **Frontend**: Displays avatar from public URL

Example URL generated:
```
http://localhost:9000/trip-manager/avatars/user-id.jpg
```

---

## Troubleshooting

### MinIO bucket not created
```bash
# Check logs
docker-compose logs minio-init

# Manually initialize
docker-compose exec minio-init sh /minio-init.sh
```

### Can't access files from MinIO
```bash
# Check if bucket is public
docker-compose exec minio mc anonymous list local/trip-manager

# Make it public
docker-compose exec minio mc anonymous set public local/trip-manager
```

### Connection refused to S3 endpoint
- Verify S3_ENDPOINT is correct
- For containers use service name: `http://minio:9000`
- For local access use: `http://localhost:9000`

### Files uploaded but not visible in console
- Check bucket permissions: should be "public" for read access
- Verify avatar URL is correct: `S3_PUBLIC_URL/avatars/{filename}`

---

## For More Information

See `DOCKER_SETUP.md` for comprehensive Docker configuration guide.### Enable S3 Storage in Development

Set environment variables in `.env`:

```env
STORAGE_TYPE=s3
S3_ENDPOINT=http://minio:9000
S3_BUCKET=trip-manager
S3_REGION=us-east-1
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_PUBLIC_URL=http://localhost:9000
S3_USE_SSL=false
```

Or in `docker-compose.yaml`:

```yaml
services:
  backend:
    environment:
      STORAGE_TYPE: s3
      S3_ENDPOINT: http://minio:9000
      S3_BUCKET: trip-manager
      S3_REGION: us-east-1
      S3_ACCESS_KEY: minioadmin
      S3_SECRET_KEY: minioadmin
      S3_PUBLIC_URL: http://localhost:9000
      S3_USE_SSL: false
```

## Production Setup with AWS S3

To switch to production AWS S3:

### 1. Add AWS SDK Dependency

```bash
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/credentials
go get github.com/aws/aws-sdk-go-v2/feature/s3/manager
go get github.com/aws/aws-sdk-go-v2/service/s3
```

### 2. Set Environment Variables

```env
STORAGE_TYPE=s3
S3_BUCKET=my-trip-manager-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=<your-aws-access-key>
S3_SECRET_KEY=<your-aws-secret-key>
S3_PUBLIC_URL=https://s3.us-east-1.amazonaws.com/my-trip-manager-bucket
S3_USE_SSL=true
```

### 3. Create S3 Bucket

```bash
aws s3api create-bucket \
  --bucket my-trip-manager-bucket \
  --region us-east-1 \
  --acl private
```

### 4. Set Bucket CORS (for public uploads)

```bash
aws s3api put-bucket-cors \
  --bucket my-trip-manager-bucket \
  --cors-configuration file://cors.json
```

**cors.json:**
```json
{
  "CORSRules": [
    {
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "PUT", "POST"],
      "AllowedHeaders": ["*"],
      "MaxAgeSeconds": 3000
    }
  ]
}
```

## Switching Between Local and S3

### Development (Local Storage)

```env
STORAGE_TYPE=local
UPLOAD_DIR=./uploads
```

```bash
docker-compose up -d postgres
docker-compose up backend frontend
```

### Development (With MinIO)

```env
STORAGE_TYPE=s3
S3_ENDPOINT=http://minio:9000
S3_BUCKET=trip-manager
...
```

```bash
docker-compose up -d
```

### Production (AWS S3)

```env
STORAGE_TYPE=s3
S3_ENDPOINT=  # Leave empty for AWS
S3_BUCKET=my-trip-manager-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=<aws-key>
S3_SECRET_KEY=<aws-secret>
S3_USE_SSL=true
```

## Architecture

```
┌─────────────────────────────────────┐
│  HTTP Handler (UpdateMeHandler)     │
│  - Parse multipart/form-data        │
│  - Validate file                    │
└────────────┬────────────────────────┘
             │
┌────────────▼────────────────────────┐
│  UserService                        │
│  - Business logic                   │
│  - Calls MediaService               │
└────────────┬────────────────────────┘
             │
┌────────────▼────────────────────────┐
│  MediaService                       │
│  - Upload logic                     │
│  - Uses Storage interface           │
└────────────┬────────────────────────┘
             │
    ┌────────┴─────────┐
    │                  │
┌───▼──────┐    ┌──────▼──────┐
│ Local    │    │    S3       │
│ Storage  │    │  - MinIO    │
│          │    │  - AWS S3   │
└──────────┘    │  - GCS      │
                └─────────────┘
```

## Benefits

✅ **No Vendor Lock-in** - Same code works with MinIO, AWS S3, GCS, etc.
✅ **Easy Testing** - MinIO in Docker for local testing
✅ **Simple Migration** - Just change credentials to switch providers
✅ **Production Ready** - Same interface as AWS S3
✅ **Cost Efficient** - MinIO is open source, AWS S3 is pay-as-you-go

