# BFF Code Cleanup - Presigned URL Migration

## Status: Frontend ✅ Complete | Backend ✅ Complete

## What Changed

### Frontend (DONE)
✅ **Removed old BFF upload code:**
- `updateMeWithAvatar()` function removed from `frontend/lib/api/auth.ts`
- Multipart form data upload removed from `ProfileSettings.tsx`

✅ **Created new presigned URL API:**
- `frontend/lib/api/uploads.ts` with:
  - `getPresignedUrl()` - Request presigned URL from backend
  - `uploadToPresignedUrl()` - Upload directly to MinIO
  - `uploadAvatar()` - Convenience wrapper

✅ **Updated ProfileSettings.tsx:**
- New flow: Get URL → Upload → Update user profile
- Simpler, cleaner code
- Better error handling

## Backend Cleanup (DONE)

✅ **Verified existing implementation:**
1. ✅ `backend/internal/handler/multipart.go` - Already removed (no longer exists)
2. ✅ `backend/internal/service/user_service.go` - Already clean (no UpdateUserWithAvatar method)
3. ✅ `backend/internal/handler/auth.go` - Already simplified (JSON-only, no multipart)

✅ **Completed new implementation:**
1. ✅ Added presigned URL route: `POST /api/uploads/presigned`
2. ✅ Created `backend/internal/handler/uploads.go` with `GetPresignedURLHandler`
3. ✅ Infrastructure support via `GeneratePresignedURL` in `presigned_url.go`
4. ✅ Removed unused import from s3.go
5. ✅ Backend builds and compiles successfully

## Backend Cleanup TODO

### 1. Remove multipart.go helper
- File: `backend/internal/handler/multipart.go`
- Function to remove: `ParseMultipartFormWithFile()`
- Reason: No longer needed for presigned URLs

### 2. Simplify handler/auth.go
- Remove multipart handling from `UpdateMeHandler`
- Remove file parsing logic
- Keep only JSON PUT for user updates
- Remove calls to `UpdateUserWithAvatar()`

### 3. Simplify user_service.go
- Remove `UpdateUserWithAvatar()` method
- Keep only `UpdateUser()` method
- Remove avatar file handling logic
- Service now only handles URL storage

### 4. Update Tests
- Remove multipart upload tests
- Add presigned URL generation tests
- Add integration tests for the new flow

## New Architecture

```
OLD (BFF):
┌─────────┐     ┌─────────┐     ┌────────┐
│ Frontend│ --→ │ Backend │ --→ │ MinIO  │
└─────────┘     └─────────┘     └────────┘
  (Upload)     (Route + Upload)  (Store)

NEW (Presigned URLs):
┌─────────┐     ┌─────────┐
│ Frontend│ --→ │ Backend │  (1. Get URL)
└─────────┘     └─────────┘

       ↓ (2. Direct Upload)

    ┌────────┐
    │ MinIO  │
    └────────┘

┌─────────┐     ┌─────────┐  (3. Store URL in DB)
│ Frontend│ --→ │ Backend │
└─────────┘     └─────────┘
```

## Benefits

✅ **Performance**: Files never go through backend
✅ **Scalability**: Backend handles only metadata operations
✅ **Security**: Presigned URLs expire, no credentials exposed
✅ **Simplicity**: No multipart handling needed
✅ **Reliability**: Direct S3/MinIO connection more reliable

## Files to Cleanup

### High Priority (Remove) - ✅ COMPLETED
- [x] `backend/internal/handler/multipart.go` - ✅ Already removed
- [x] `backend/internal/service/user_service.go` - ✅ Already clean (no UpdateUserWithAvatar method)
- [x] `backend/internal/handler/auth.go` - ✅ Already simplified (JSON-only, no multipart)

### Medium Priority (Update) - ✅ COMPLETED
- [x] `backend/cmd/api/main.go` - ✅ Added route: `r.Post("/uploads/presigned", GetPresignedURLHandler(application))`
- [x] `backend/internal/handler/uploads.go` - ✅ Already exists with GetPresignedURLHandler
- [x] `backend/internal/infrastructure/media_service.go` - ✅ Contains presigned URL support via GeneratePresignedURL
- [x] `backend/internal/infrastructure/presigned_url.go` - ✅ Full presigned URL generation implementation
- [x] `backend/internal/storage/s3.go` - ✅ Removed unused `manager` import

### Verification Tasks - ✅ COMPLETED
- ✅ Presigned URL endpoint registered in routing
- ✅ All type signatures correct (time.Duration for expiration)
- ✅ Backend compiles successfully with no errors

### Low Priority (Nice to have)
- [ ] Documentation - Update deployment docs
- [ ] Comments - Remove old BFF references
- [ ] Migration guide - For any dependent systems

## Testing Checklist

```bash
# 1. Get presigned URL
curl -X POST http://localhost:8000/api/uploads/presigned \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"fileName": "avatar.jpg", "mediaType": "avatar"}'

# 2. Upload directly to MinIO using presigned URL
curl -X PUT "{presignedUrl}" \
  -H "Content-Type: image/jpeg" \
  --data-binary @avatar.jpg

# 3. Update user profile with avatar URL
curl -X PUT http://localhost:8000/api/users/me \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "email": "john@example.com",
    "bio": "Hello",
    "avatarUrl": "http://minio:9000/trip-manager/avatars/user-id.jpg"
  }'

# 4. Verify avatar displays correctly
# - Check browser Network tab for direct MinIO requests
# - Verify avatar appears in Settings and Navbar
```

## Notes

- Presigned URLs are valid for 15 minutes
- Only authenticated users can request presigned URLs
- Backend still validates final state
- Works with AWS S3, MinIO, GCS, etc.
- This is a breaking change - all clients must use new flow

