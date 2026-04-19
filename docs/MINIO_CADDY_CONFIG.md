# MinIO Configuration Guide

## Architecture Overview

**UPDATED: Presigned URL Architecture (NEW)**

The application now supports **direct uploads to MinIO via presigned URLs**:

```
Frontend Presigned URL Upload Flow:
  1. Frontend requests presigned URL: POST /api/uploads/presigned
  2. Backend generates signed URL via MinIO S3 API
  3. Frontend receives URL valid for 15 minutes
  4. Frontend uploads directly to MinIO using presigned URL
  5. Frontend updates user profile: PUT /api/users/me with avatarUrl
  6. Backend stores URL in database
  7. Frontend displays image from returned URL
```

**This means MinIO MUST be publicly accessible!**

---

## MinIO Exposure Strategy

### Development (Local with Caddy)
```
Frontend → Caddy (http://localhost/minio) → MinIO
```

### Production (Via Caddy Reverse Proxy)
```
Frontend → Caddy (https://{DOMAIN}/minio) → MinIO (internal)
```

**MinIO never needs to be directly exposed** - Caddy reverse proxy handles public access

---

## Configuration

### 1. Caddyfile Configuration

The Caddyfile now includes a MinIO reverse proxy route:

```caddyfile
$DOMAIN {
    # API Reverse Proxy
    handle /api/* {
        reverse_proxy backend:8000
    }
    
    # MinIO Reverse Proxy for Presigned URLs
    handle /minio/* {
        reverse_proxy minio:9000 {
            header_up Host {host}
            header_up X-Forwarded-Proto {scheme}
            header_up X-Forwarded-For {remote}
        }
    }
    
    # Frontend Reverse Proxy
    handle {
        reverse_proxy frontend:3000
    }
}
```

**Key points:**
- `/minio/*` routes all MinIO requests through Caddy
- Header forwarding preserves client information
- Works for both HTTP (dev) and HTTPS (production)

---

### 2. Environment Configuration

#### Development (Local MinIO via Caddy):
```dotenv
S3_ENDPOINT=http://minio:9000          # Internal: Backend talks to MinIO
S3_PUBLIC_URL=http://localhost/minio    # Public: Frontend talks via Caddy
S3_BUCKET=trip-manager
S3_REGION=us-east-1
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_USE_SSL=false
```

#### Production (MinIO via Caddy with HTTPS):
```dotenv
S3_ENDPOINT=http://minio:9000                    # Internal: Backend talks to MinIO
S3_PUBLIC_URL=https://{DOMAIN}/minio             # Public: Frontend via HTTPS Caddy
S3_BUCKET=trip-manager
S3_REGION=us-east-1
S3_ACCESS_KEY=your-key
S3_SECRET_KEY=your-secret
S3_USE_SSL=false                                 # SSL between Backend→MinIO not needed
```

#### Production (AWS S3 instead of MinIO):
```dotenv
S3_ENDPOINT=https://s3.amazonaws.com
S3_PUBLIC_URL=https://your-bucket.s3.amazonaws.com
S3_BUCKET=your-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=your-aws-key
S3_SECRET_KEY=your-aws-secret
S3_USE_SSL=true
```

---

### 3. MinIO CORS Configuration

MinIO needs CORS headers for presigned URL uploads to work from browsers.

The `minio-init.sh` script automatically configures CORS:

```bash
mc cors set local/$BUCKET_NAME cors-config.json
```

**CORS allows:**
- ✅ Preflight OPTIONS requests (browsers require this)
- ✅ PUT/POST requests for uploads
- ✅ Cross-origin headers

For production, update the CORS origins in `docker/minio-init.sh`:

```json
{
  "AllowedOrigins": [
    "https://yourdomain.com"
  ],
  "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS"],
  "AllowedHeaders": ["*"],
  "ExposeHeaders": ["ETag", "x-amz-version-id"],
  "MaxAgeSeconds": 3000
}
```

---

## Upload Flow Details

### Step 1: Request Presigned URL
```bash
POST /api/uploads/presigned
Authorization: Bearer {token}
Content-Type: application/json

{
  "fileName": "avatar.jpg",
  "mediaType": "avatar"
}
```

**Response:**
```json
{
  "presignedUrl": "http://localhost/minio/trip-manager/avatars/user-id.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&...",
  "expiresIn": 900
}
```

### Step 2: Direct Upload to MinIO
```bash
PUT {presignedUrl}
Content-Type: image/jpeg

[binary file data]
```

**Response (200 OK):**
- File stored in MinIO
- No body needed

### Step 3: Update User Profile
```bash
PUT /api/users/me
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "John Doe",
  "avatarUrl": "http://localhost/minio/trip-manager/avatars/user-id.jpg"
}
```

### Step 4: Backend Stores URL
```json
{
  "id": "user-123",
  "name": "John Doe",
  "avatarUrl": "http://localhost/minio/trip-manager/avatars/user-id.jpg",
  ...
}
```

### Step 5: Frontend Displays Image
```html
<img src={avatarUrl} alt="Avatar" />
```

---

## Security Considerations

### ✅ Presigned URLs Are Secure Because:
- **Time-limited**: Expire after 15 minutes
- **Signed**: HMAC-SHA256 signature prevents tampering
- **User-specific**: Only the authenticated user receives the URL
- **Action-specific**: Only allows PUT (upload), not GET or DELETE

### ✅ MinIO Public Access Is Safe Because:
- **Presigned URLs only**: Direct bucket access is not allowed
- **Caddy reverse proxy**: Can add rate limiting/WAF if needed
- **CORS configured**: Only allows expected origins
- **Bucket policy**: Public read-only, no write permissions

---

## Troubleshooting

### Images not loading in browser
1. Check S3_PUBLIC_URL is correct: `echo $S3_PUBLIC_URL`
2. Test Caddy proxy: `curl http://localhost/minio/`
3. Check backend logs: `docker-compose logs backend`
4. Verify CORS headers: `curl -I http://localhost/minio/`

### Presigned URL upload fails (403 Forbidden)
1. Check CORS configuration: `docker-compose logs minio-init`
2. Verify MinIO bucket exists: `docker-compose exec minio mc ls local/`
3. Check bucket policy: `docker-compose exec minio mc anonymous get local/trip-manager`
4. Test CORS preflight: 
   ```bash
   curl -X OPTIONS -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: PUT" \
     http://localhost/minio/trip-manager/avatars/test.jpg
   ```

### MinIO can't be reached
1. Check Caddy is running: `docker-compose ps caddy`
2. Check MinIO is running: `docker-compose ps minio`
3. Check network: `docker-compose exec backend curl http://minio:9000`
4. Check Caddy logs: `docker-compose logs caddy`

### S3_PUBLIC_URL issues
- **Development**: Must be `http://localhost/minio` or `http://localhost:9000`
- **Production**: Must be `https://{DOMAIN}/minio`
- **Never use**: `http://minio:9000` for public URLs (internal only)

---

## Switching from MinIO to AWS S3

To switch from MinIO to AWS S3:

```bash
# 1. Update .env
S3_ENDPOINT=https://s3.amazonaws.com
S3_BUCKET=your-bucket-name
S3_PUBLIC_URL=https://your-bucket.s3.amazonaws.com
S3_ACCESS_KEY=your-aws-key
S3_SECRET_KEY=your-aws-secret
S3_USE_SSL=true

# 2. Restart application
docker-compose down
docker-compose up -d
```

**No code changes needed** - just environment variables!

---

## Summary

| Component | Development | Production |
|-----------|-------------|-----------|
| **S3_ENDPOINT** | http://minio:9000 | http://minio:9000 |
| **S3_PUBLIC_URL** | http://localhost/minio | https://{DOMAIN}/minio |
| **Caddy Route** | /minio/* → minio:9000 | /minio/* → minio:9000 |
| **SSL** | No | Yes (Caddy handles it) |
| **MinIO Exposure** | Via Caddy on /minio | Via Caddy on /minio |
| **CORS** | Configured for localhost | Configured for {DOMAIN} |

---

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────┐
│ Frontend (Browser)                                        │
├──────────────────────────────────────────────────────────┤
│                                                          │
│ 1. POST /api/uploads/presigned  →  Caddy:/api  →  Backend
│    ↓ Get presigned URL (expires in 15 min)              │
│                                                          │
│ 2. PUT {presignedUrl}  →  Caddy:/minio  →  MinIO
│    ↓ Direct file upload                                  │
│                                                          │
│ 3. PUT /api/users/me with avatarUrl  →  Backend
│    ↓ Update profile                                      │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**Benefits:**
- ✅ File upload bypasses backend (fast & scalable)
- ✅ MinIO never exposes credentials
- ✅ Presigned URLs provide temporary access
- ✅ Caddy handles HTTPS & proxying
- ✅ Easy switch between MinIO and AWS S3


