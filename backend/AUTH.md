# Authentication Setup

## Environment Variables

Configure the following environment variables to run the application:

```bash
# Database connection
DATABASE_URL=postgres://postgres:postgres@localhost:5432/trip_manager?sslmode=disable

# Server configuration
SERVER_PORT=8080
ENVIRONMENT=development

# JWT token secret (change this in production!)
JWT_SECRET=your-secret-key-change-in-production
```

## JWT Token

The API uses JWT (JSON Web Tokens) for authentication:

- **Token Duration**: 15 minutes
- **Authentication Header**: `Authorization: Bearer <token>`
- **Signing Method**: HS256 (HMAC with SHA-256)

## Authentication Flow

### 1. Register a new user
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "SecurePass123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "createdAt": "2024-03-28T10:00:00Z",
    "updatedAt": "2024-03-28T10:00:00Z"
  }
}
```

### 2. Login with email and password
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123"
}

Response: (same as register)
```

### 3. Use the token for authenticated requests
```bash
GET /api/users/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "createdAt": "2024-03-28T10:00:00Z",
  "updatedAt": "2024-03-28T10:00:00Z"
}
```

### 4. Change password
```bash
PUT /api/users/me/password
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "currentPassword": "SecurePass123",
  "newPassword": "NewSecurePass456"
}

Response: 204 No Content
```

## Password Requirements

- Minimum 8 characters
- Must contain at least one uppercase letter
- Must contain at least one lowercase letter
- Must contain at least one digit

## Protected Routes

All routes except `/auth/register` and `/auth/login` require a valid JWT token in the `Authorization` header.

If the token is missing or invalid, the API will return:
```
401 Unauthorized
```

## Token Claims

The JWT token contains the following claims:
- `userId`: The unique identifier of the user
- `email`: The user's email address
- `name`: The user's display name
- `exp`: Token expiration time
- `iat`: Token issued at time
- `nbf`: Token not valid before time

## Security Notes

- **Never** commit your `JWT_SECRET` to version control
- Use a strong, randomly generated secret for production
- Rotate your JWT_SECRET periodically
- Implement HTTPS in production to protect tokens in transit
- Consider implementing token refresh mechanisms for long-running sessions

