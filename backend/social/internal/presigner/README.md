# Presigner Service

Microservice für Presigned URL Generierung. Unterstützt sowohl **S3/MinIO** als auch **Google Cloud Storage (GCS)** über Umgebungsvariablen.

## Features

✅ **Multi-Provider** - Wählen Sie zwischen S3/MinIO oder GCS über `STORAGE_PROVIDER`
✅ **Isoliert** - Keine Storage-Dependencies in anderen Services
✅ **Konfigurierbar** - Bucket und Expiration über Env-Variablen
✅ **Bucket-agnostisch** - Requests können Bucket überschreiben

## Architektur

```
internal/
├── config/         # Provider-Auswahl und Initialization
├── provider/       # Abstraktes Interface, S3 und GCS Implementierungen
│   ├── provider.go # Interface
│   ├── s3.go       # S3/MinIO Implementierung
│   └── gcs.go      # Google Cloud Storage Implementierung
├── handler/        # HTTP Handler
├── service/        # Presigning Logic (Provider-agnostisch)
└── shared/         # Response Helpers
cmd/api/
└── main.go        # Einstiegspunkt
```

## API Endpoints

### Upload URL
```
POST /upload-url
Content-Type: application/json

{
  "key": "avatars/user-123.png",
  "bucket": "my-bucket"  # Optional, uses default if not provided
}

Response:
{
  "url": "https://...",
  "expiresIn": 900
}
```

### Download URL
```
POST /download-url
Content-Type: application/json

{
  "key": "avatars/user-123.png",
  "bucket": "my-bucket"  # Optional
}

Response:
{
  "url": "https://...",
  "expiresIn": 900
}
```

### Health
```
GET /health

Response:
{
  "status": "ok"
}
```

## Configuration

### S3/MinIO

```bash
# Wähle S3 Provider
export STORAGE_PROVIDER=s3
export STORAGE_BUCKET=trip-manager
export S3_ENDPOINT=http://minio:9000
export S3_REGION=us-east-1
export S3_ACCESS_KEY=minioadmin
export S3_SECRET_KEY=minioadmin
export S3_USE_SSL=false
export PRESIGNER_URL_EXPIRATION=15m
export PORT=8081

go run ./cmd/api/main.go
```

### Google Cloud Storage

```bash
# Wähle GCS Provider
export STORAGE_PROVIDER=gcs
export STORAGE_BUCKET=my-gcs-bucket
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
export PRESIGNER_URL_EXPIRATION=15m
export PORT=8081

go run ./cmd/api/main.go
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8081` | Server port |
| `STORAGE_PROVIDER` | `s3` | `s3` oder `gcs` |
| `STORAGE_BUCKET` | `trip-manager` | Default bucket name |
| `PRESIGNER_URL_EXPIRATION` | `15m` | URL expiration (e.g., "1h", "30m") |
| `LOG_LEVEL` | `info` | Log level |

### S3-spezifisch
| Variable | Default | Description |
|----------|---------|-------------|
| `S3_ENDPOINT` | `http://localhost:9000` | S3/MinIO endpoint |
| `S3_REGION` | `us-east-1` | AWS region |
| `S3_ACCESS_KEY` | `minioadmin` | Access key |
| `S3_SECRET_KEY` | `minioadmin` | Secret key |
| `S3_USE_SSL` | `true` | Use SSL |

### GCS-spezifisch
| Variable | Description |
|----------|-------------|
| `GOOGLE_APPLICATION_CREDENTIALS` | Path to GCP service account JSON |

## Docker

### S3/MinIO
```bash
docker build -t trip-manager-presigner .
docker run -p 8081:8081 \
  -e STORAGE_PROVIDER=s3 \
  -e S3_ENDPOINT=http://minio:9000 \
  -e STORAGE_BUCKET=trip-manager \
  trip-manager-presigner
```

### GCS
```bash
docker build -t trip-manager-presigner .
docker run -p 8081:8081 \
  -e STORAGE_PROVIDER=gcs \
  -e STORAGE_BUCKET=my-gcs-bucket \
  -v /path/to/credentials.json:/app/credentials.json \
  -e GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json \
  trip-manager-presigner
```

## Integration

### HTTP Client Usage

```go
// Upload URL
resp, err := http.Post(
  "http://presigner:8081/upload-url",
  "application/json",
  bytes.NewBuffer([]byte(`{"key":"avatars/user-123.png"}`)),
)

// Download URL
resp, err := http.Post(
  "http://presigner:8081/download-url",
  "application/json",
  bytes.NewBuffer([]byte(`{"key":"avatars/user-123.png"}`)),
)
```

## Vorteile

✅ **Storage-agnostisch** - Switch zwischen S3 und GCS ohne Code-Änderungen
✅ **Dependencies isoliert** - Andere Services brauchen nicht S3/GCS SDKs zu laden
✅ **Zentrale Konfiguration** - Alle Storage Settings an einem Ort
✅ **Skalierbar** - Kann independent skaliert werden
✅ **Sicherheit** - Services können nur URLs holen, nicht direkt auf Storage zugreifen
