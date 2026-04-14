#!/bin/sh

# MinIO initialization script
# Sets up the bucket, policies, and CORS for presigned URL uploads

set -e

# Get environment variables (with defaults)
BUCKET_NAME="${S3_BUCKET:-trip-manager}"
MINIO_ENDPOINT="${S3_ENDPOINT:-http://minio:9000}"
MINIO_ROOT_USER="${MINIO_ROOT_USER:-minioadmin}"
MINIO_ROOT_PASSWORD="${MINIO_ROOT_PASSWORD:-minioadmin}"
DOMAIN="${DOMAIN:-localhost}"

echo "🔍 MinIO Initialization Debug:"
echo "  Endpoint: $MINIO_ENDPOINT"
echo "  Bucket: $BUCKET_NAME"
echo "  User: $MINIO_ROOT_USER"

echo "Waiting for MinIO to start..."
sleep 5

# Configure mc (MinIO client)
echo "Configuring MinIO client with credentials..."
/usr/bin/mc alias set local "$MINIO_ENDPOINT" "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD" --api S3v4

# Test connection
echo "Testing MinIO connection..."
if /usr/bin/mc ready local; then
    echo "✓ MinIO is ready!"
else
    echo "❌ MinIO is not ready, but continuing anyway..."
fi

# Create bucket if it doesn't exist
echo "Creating bucket '$BUCKET_NAME' if it doesn't exist..."
/usr/bin/mc mb "local/$BUCKET_NAME" --ignore-existing || echo "⚠️  Bucket creation may have failed, but continuing..."

# Set bucket policy to public read
# This allows direct access to files for presigned URLs
echo "Setting bucket policy to public read..."
/usr/bin/mc anonymous set public "local/$BUCKET_NAME" || true

# Configure CORS for presigned URL uploads
# Allows browser-based direct uploads to MinIO via presigned URLs
echo "Configuring CORS for presigned URLs..."
cat > /tmp/cors.json << EOF
[
  {
    "AllowedOrigins": [
      "http://localhost:3000",
      "http://localhost",
      "http://127.0.0.1:3000",
      "https://${DOMAIN}",
      "https://www.${DOMAIN}"
    ],
    "AllowedMethods": [
      "GET",
      "PUT",
      "POST",
      "DELETE",
      "HEAD",
      "OPTIONS"
    ],
    "AllowedHeaders": [
      "*"
    ],
    "ExposeHeaders": [
      "ETag",
      "x-amz-version-id"
    ],
    "MaxAgeSeconds": 3000
  }
]
EOF

# Apply CORS configuration
if /usr/bin/mc cors set local/$BUCKET_NAME /tmp/cors.json; then
    echo "✓ CORS: Configured for https://${DOMAIN} and localhost"
else
    echo "⚠️  CORS configuration failed – presigned URL uploads from browsers may not work"
fi

echo "✓ MinIO initialization complete!"
echo "✓ Bucket: $BUCKET_NAME"
echo "✓ Policy: Public read enabled"

