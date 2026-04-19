# Presigned URL Implementation Plan

## Overview
Migration from BFF (Backend-for-Frontend) file upload approach to presigned URLs for direct browser-to-S3/MinIO uploads.

## New Upload Flow

```
User Upload Avatar:
1. Frontend: POST /api/uploads/presigned
   Request: { fileName: "avatar.jpg", mediaType: "avatar" }
   Response: { presignedUrl: "https://s3.../...", expiresIn: 900 }

2. Frontend: PUT to presignedUrl directly (multipart/form-data)
   No authorization needed - URL is signed with credentials

3. MinIO/S3: Stores file, returns 200 OK

4. Frontend: PUT /api/users/me
   Request: { name, email, bio }
   (No avatarFile needed anymore - URL will be stored by backend)

5. Backend: Checks if presigned upload was completed, stores URL in DB

6. Frontend: Updates UserContext with new avatarUrl
```

## Files Modified/Created

### Backend

1. **`backend/internal/storage/s3.go`**
   - Added: `GeneratePresignedURL()` method
   - Uses `s3.PresignPutObject` for upload URLs (valid for 15 mins)
   - Works with AWS S3, MinIO, and S3-compatible services

2. **`backend/internal/infrastructure/presigned_url.go`** (NEW)
   - `GeneratePresignedURL()` method in MediaService
   - Path generation logic (same as UploadImage)
   - Error handling and validation

3. **`backend/internal/handler/uploads.go`** (NEW)
   - `GeneratePresignedURLHandler`: POST /api/uploads/presigned
   - Parses request (fileName, mediaType)
   - Calls MediaService to generate URL
   - Returns presignedUrl + expiresIn

4. **`backend/cmd/api/main.go`**
   - Add route: `r.Post("/uploads/presigned", handler.GeneratePresignedURLHandler(application))`
   - Inside authenticated API routes

5. **`backend/internal/service/user_service.go`**
   - Modify `UpdateUser()`: Check if avatar URL provided, validate it's on S3

### Frontend

1. **`frontend/lib/api/uploads.ts`** (NEW)
   - `getPresignedUrl()`: POST /api/uploads/presigned
   - `uploadToPresignedUrl()`: PUT to presigned URL

2. **`frontend/components/settings/ProfileSettings.tsx`**
   - New flow:
     1. User selects file
     2. Get presigned URL from backend
     3. Upload directly to MinIO
     4. Update user profile (without file, just metadata)

## Advantages Over BFF Approach

✅ **Better Performance**
- Files don't go through backend
- Reduced bandwidth usage
- Faster uploads

✅ **Better Scalability**
- Backend doesn't bottleneck on uploads
- Direct S3/MinIO connection

✅ **Security**
- Presigned URLs expire after 15 minutes
- No need to expose S3 credentials to frontend
- Backend still validates everything

✅ **Simpler Architecture**
- Frontend handles uploading and can retry
- Backend just validates final state
- MinIO/S3 handles storage

## API Endpoints

### Generate Presigned URL
```
POST /api/uploads/presigned
Authorization: Bearer {token}

Request:
{
  "fileName": "avatar.jpg",
  "mediaType": "avatar"
}

Response:
{
  "presignedUrl": "https://minio.../trip-manager/avatars/...",
  "expiresIn": 900
}
```

### Upload Using Presigned URL
```
PUT {presignedUrl}
(No authorization needed - signed into URL)
Content-Type: application/octet-stream

Body: {file binary}
```

### Update User Profile
```
PUT /api/users/me
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
  "name": "John Doe",
  "email": "john@example.com",
  "bio": "...",
  "avatarUrl": "http://minio:9000/trip-manager/avatars/user-id.jpg"
}
```

## Implementation Status

- [ ] Add `GeneratePresignedURL()` to S3Storage
- [ ] Create `presigned_url.go` in MediaService
- [ ] Create `uploads.go` handler
- [ ] Register routes in main.go
- [ ] Create frontend `uploads.ts` API file
- [ ] Update ProfileSettings component
- [ ] Test with MinIO
- [ ] Documentation update

## Testing

```bash
# 1. Get presigned URL
curl -X POST http://localhost:8000/api/uploads/presigned \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"fileName": "test.jpg", "mediaType": "avatar"}'

# 2. Upload using presigned URL
curl -X PUT "{presignedUrl}" \
  -H "Content-Type: image/jpeg" \
  --data-binary @/path/to/image.jpg

# 3. Update user with new avatar URL
curl -X PUT http://localhost:8000/api/users/me \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "email": "john@example.com",
    "bio": "Hello",
    "avatarUrl": "http://minio:9000/trip-manager/avatars/user-id.jpg"
  }'
```

## Notes

- Presigned URLs are valid for 15 minutes
- Only authenticated users can request presigned URLs
- Backend still validates that files end up on S3
- Path generation is identical to existing UploadImage logic
- Works with AWS S3, MinIO, and any S3-compatible service

