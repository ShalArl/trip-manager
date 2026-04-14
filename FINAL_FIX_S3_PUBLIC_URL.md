# Final Fix: S3_PUBLIC_URL mit Bucket-Namen

## Das Problem war:

Es gibt **3 verschiedene URLs** und jede muss richtig konfiguriert sein:

### 1. **S3_ENDPOINT** (intern)
```
http://minio:9000
```
→ Nur für Backend-interne Kommunikation

### 2. **S3_PUBLIC_URL** (öffentlich, mit Bucket!)
```
https://travel-nugget.duckdns.org/minio/trip-manager
```
→ Muss den **Bucket-Namen** enthalten!
→ Wird vom Backend für presigned URLs und öffentliche URLs verwendet

### 3. **NEXT_PUBLIC_S3_PUBLIC_URL** (Frontend)
```
https://travel-nugget.duckdns.org/minio/trip-manager
```
→ Wird nur noch fallback verwendet (sollte nicht nötig sein)

## Fixes Applied:

### 1. **backend/internal/storage/s3.go** ✅
Fixed presigned URL generation - entfernt doppelte Bucket-Namen-einfügung:
```go
// FALSCH:
fileURL := fmt.Sprintf("%s/%s/%s?%s", s.publicURL, s.bucket, fileName, sig)
// Ergebnis: https://domain/minio/trip-manager/trip-manager/avatars/...

// RICHTIG:
fileURL := fmt.Sprintf("%s/%s?%s", s.publicURL, fileName, sig)
// Ergebnis: https://domain/minio/trip-manager/avatars/...
```

### 2. **frontend/lib/api/uploads.ts** ✅
Fixed public URL extraction - nutzt presigned URL ohne Signatur statt hardcoded URL zu konstruieren:
```typescript
// FALSCH:
const publicUrl = `${NEXT_PUBLIC_S3_PUBLIC_URL}/avatars/${userId}${ext}`;

// RICHTIG:
const publicUrl = presignedUrl.split('?')[0];  // Entfernt Query-Parameter
// Ergebnis: presigned URL ohne Signatur = öffentliche URL
```

### 3. **.env.example** ✅
Aktualisiert mit klarer Dokumentation:
```bash
S3_PUBLIC_URL=https://travel-nugget.duckdns.org/minio/trip-manager  # Mit Bucket-Name!
```

## Was du jetzt tun musst:

### 1. Code Pushen
```bash
cd ~/dev/cloud/trip-manager
git add -A
git commit -m "fix: S3_PUBLIC_URL format with bucket name and presigned URL generation"
git push origin main
```

GitHub Actions wird triggert → neues Backend-Image bauen

### 2. Auf Server Deployen
```bash
cd ~/app/cloud
docker-compose pull
docker-compose down
docker-compose up -d
sleep 30
```

### 3. Testen

**Test 1: Presigned URL korrekt?**
```bash
# Upload neuer Avatar
# Browser Console sollte zeigen:
[uploadAvatar] Got presigned URL...
# Die URL sollte sein: https://travel-nugget.duckdns.org/minio/trip-manager/avatars/...?X-Amz-...
```

**Test 2: Public URL korrekt?**
```bash
# Nach Upload
# Browser Console sollte zeigen:
[uploadAvatar] Avatar uploaded successfully: https://travel-nugget.duckdns.org/minio/trip-manager/avatars/...
# OHNE Query-Parameter!
```

**Test 3: GET funktioniert?**
```bash
curl -v https://travel-nugget.duckdns.org/minio/trip-manager/avatars/a41cc4c9-32f1-4579-8a8b-45925f69b242.png
# HTTP/2 200 ✅
```

**Test 4: Kein Mixed Content?**
- Browser Console (F12) sollte KEINE Mixed Content Fehler zeigen
- Avatar sollte angezeigt werden

## Zusammenfassung der Änderungen:

| Datei | Änderung | Grund |
|-------|----------|-------|
| `backend/internal/storage/s3.go` | Presigned URL: `%s/%s?%s` statt `%s/%s/%s?%s` | Bucket-Name darf nicht doppelt sein |
| `frontend/lib/api/uploads.ts` | Public URL: `presignedUrl.split('?')[0]` | Extrahiert wahre URL ohne falsche Konstruktion |
| `.env.example` | Dokumentation: Bucket-Name in URL | Klarheit für zukünftige Deployments |

## Das war's!

Der Avatar-Upload sollte jetzt mit korrekten URLs funktionieren:
- PUT Presigned: `https://domain/minio/trip-manager/avatars/...?signature` ✅
- GET Public: `https://domain/minio/trip-manager/avatars/...` ✅

