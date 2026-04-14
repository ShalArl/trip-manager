#!/bin/sh

# MinIO initialization script
# Sets up the bucket, policies, and CORS for presigned URL uploads

set -e

# Get environment variables (with defaults)
BUCKET_NAME="${S3_BUCKET:-trip-manager}"
MINIO_ENDPOINT="${S3_ENDPOINT:-http://localhost:9000}"
MINIO_ROOT_USER="${MINIO_ROOT_USER:-minioadmin}"
MINIO_ROOT_PASSWORD="${MINIO_ROOT_PASSWORD:-minioadmin}"

echo "Waiting for MinIO to start..."
sleep 5

# Configure mc (MinIO client)
echo "Configuring MinIO client..."
/usr/bin/mc config host add local "$MINIO_ENDPOINT" "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD" --api S3v4

# Create bucket if it doesn't exist
echo "Creating bucket '$BUCKET_NAME' if it doesn't exist..."
/usr/bin/mc mb "local/$BUCKET_NAME" --ignore-existing

# Set bucket policy to public read
# This allows direct access to files for presigned URLs
echo "Setting bucket policy to public read..."
/usr/bin/mc anonymous set public "local/$BUCKET_NAME"

# Configure CORS for presigned URL uploads
# Allows browser-based direct uploads to MinIO via presigned URLs
echo "Configuring CORS for presigned URLs..."
cat > /tmp/cors.json << 'EOF'
[
  {
    "AllowedOrigins": [
      "http://localhost:3000",
      "http://localhost",
      "http://127.0.0.1:3000"
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
/usr/bin/mc cors set local/$BUCKET_NAME /tmp/cors.json || true

# For production deployments, add domain-based CORS
# Uncomment and customize for your production domain:
# cat > /tmp/cors-prod.json << 'EOF'
# [
#   {
#     "AllowedOrigins": [
#       "https://yourdomain.com"
#     ],
#     "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS"],
#     "AllowedHeaders": ["*"],
#     "ExposeHeaders": ["ETag", "x-amz-version-id"],
#     "MaxAgeSeconds": 3000
#   }
# ]
# EOF
# /usr/bin/mc cors set local/$BUCKET_NAME /tmp/cors-prod.json

echo "✓ MinIO initialization complete!"
echo "✓ Bucket: $BUCKET_NAME"
echo "✓ Policy: Public read enabled"
echo "✓ CORS: Configured for presigned URLs"

