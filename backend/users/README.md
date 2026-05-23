# Users Service

Manages user provisioning and profile management for Trip Manager.

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/users/provision | required | Provision user after Firebase signup |
| GET | /api/users/me | required | Get current user profile |
| PUT | /api/users/me | required | Update current user profile |
| GET | /health | - | Health check |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8001 | Service port |
| DATABASE_URL | - | PostgreSQL connection string |
| AUTH_SERVICE_URL | http://localhost:8080 | Auth service URL |
| FIREBASE_PROJECT_ID | trip-manager-local | Firebase project ID |

## Usage

### Provision a user

```go
POST /api/users/provision
Authorization: Bearer <firebase-token>

// Optional body
{
  "name": "John Doe"
}

// Response 200 (already exists) or 201 (created)
{
  "id": "uuid",
  "email": "john@example.com",
  "name": "John Doe",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Get current user

```go
GET /api/users/me
Authorization: Bearer <firebase-token>
```

### Update profile

```go
PUT /api/users/me
Authorization: Bearer <firebase-token>

{
  "name": "New Name",
  "bio": "My bio",
  "avatarKey": "avatars/user-id.jpg"
}
```

## Run locally

```bash
go run ./cmd/api
```