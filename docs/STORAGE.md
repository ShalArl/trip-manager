# Storage Configuration Guide

## Overview

The Trip Manager supports two storage backends:

1. **Local Storage** (Development/Testing) - Default
2. **S3-Compatible Storage** (Production) - MinIO, AWS S3, etc.

## Development Setup with MinIO

MinIO is an S3-compatible object storage that runs in Docker. It's perfect for development because:
- ✅ Same API as AWS S3
- ✅ Runs locally in a container
- ✅ Easy to switch to production S3 later
- ✅ No AWS account needed for development

### Starting MinIO with Docker Compose

```bash
docker-compose up -d minio
```

This starts MinIO with:
- **Access Key:** `minioadmin`
- **Secret Key:** `minioadmin`
- **Endpoint:** `http://minio:9000` (internal)
- **Public URL:** `http://localhost:9000` (from browser)
- **Console:** `http://localhost:9001` (manage buckets)

### Enable S3 Storage in Development

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

