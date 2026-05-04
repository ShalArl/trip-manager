# Social Service

Microservice für Social Features (Likes und Comments) des Trip Manager Systems.

## Struktur

```
internal/
├── like/              # Like-spezifischer Code
│   ├── types.go       # Konstanten und TargetType
│   ├── repository.go  # Firestore Repository Interface und Impl
│   ├── service.go     # Business Logic für Likes
│   └── handler.go     # HTTP Handler Funktionen
├── comment/           # Comment-spezifischer Code
│   ├── types.go       # Comment Datenstrukturen
│   ├── repository.go  # Firestore Repository Interface und Impl
│   ├── service.go     # Business Logic für Comments
│   └── handler.go     # HTTP Handler Funktionen
├── shared/            # Gemeinsame Utilities
│   ├── errors.go      # Standard Error Definitionen
│   └── types.go       # Response Helper (RespondJSON, RespondError)
├── middleware/        # Auth und andere Middleware
│   └── auth.go        # Authorization und UserID Extraction
└── config/            # Konfiguration
    └── config.go      # Config Loading und Firestore Connection
cmd/api/
└── main.go           # Einstiegspunkt der Anwendung
```

## Architektur

Jedes Feature (Like, Comment) folgt dem gleichen Pattern:
1. **Repository**: Datenbankzugriff (Firestore)
2. **Service**: Business Logic und Validierung
3. **Handler**: HTTP Endpoints und Request/Response Handling
4. **Types**: DTOs und Domain Models

## API Endpoints

### Likes
```
GET    /entities/{entityId}/likes       # Get like count and user's like status
POST   /entities/{entityId}/likes       # Like an entity (requires auth)
DELETE /entities/{entityId}/likes       # Unlike an entity (requires auth)
```

### Comments
```
GET    /entities/{entityId}/comments    # List comments for entity
POST   /entities/{entityId}/comments    # Create comment (requires auth)
PUT    /comments/{commentId}            # Update comment (requires auth)
DELETE /comments/{commentId}            # Delete comment (requires auth)
```

### Health
```
GET    /health                          # Health check
```

## Building and Running

```bash
# Build
go build -o social ./cmd/api

# Run
./social

# With environment variables
PORT=8080 FIRESTORE_PROJECT_ID=trip-manager-local ./social
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `FIRESTORE_PROJECT_ID`: GCP Firestore Project ID (default: trip-manager-local)
- `FIRESTORE_EMULATOR_HOST`: Firestore Emulator host (optional, for local development)
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to GCP credentials file (optional)

## Docker

```bash
docker build -t trip-manager-social .
docker run -p 8080:8080 trip-manager-social
```

## Development

### Running with Firestore Emulator

```bash
# Start Firestore emulator
firebase emulators:start --only firestore

# Run service
FIRESTORE_EMULATOR_HOST=localhost:8080 PORT=8080 go run ./cmd/api/main.go
```

### Testing

```bash
go test ./...
```

