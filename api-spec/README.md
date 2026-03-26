# Trip Manager API Specification

This directory contains the OpenAPI specification for the Trip Manager API.

## Overview

The Trip Manager API provides endpoints for managing travel trips, locations, and activities. It allows users to:

- **Manage Trips**: Create, read, update, and delete trips
- **Manage Locations**: Add destinations to trips with geographic information
- **Plan Activities**: Schedule and organize activities within trips
- **Track Status**: Monitor trip progress and activity details

## API Documentation

### Specification Format

- **Format**: OpenAPI 3.0.0
- **File**: `openapi.yaml`

### Main Resources

#### Trips
CRUD operations for trip management.
- `GET /trips` - List all trips
- `POST /trips` - Create a new trip
- `GET /trips/{tripId}` - Get trip details
- `PUT /trips/{tripId}` - Update trip
- `DELETE /trips/{tripId}` - Delete trip

#### Locations
Manage destinations within a trip.
- `GET /trips/{tripId}/locations` - List trip locations
- `POST /trips/{tripId}/locations` - Add location to trip
- `DELETE /trips/{tripId}/locations/{locationId}` - Remove location

#### Activities
Plan and organize activities for each location.
- `GET /trips/{tripId}/activities` - List trip activities
- `POST /trips/{tripId}/activities` - Add activity
- `PUT /trips/{tripId}/activities/{activityId}` - Update activity
- `DELETE /trips/{tripId}/activities/{activityId}` - Delete activity

## Viewing the Spec

### Online (Swagger UI / ReDoc)
You can view the spec in an interactive editor:
1. Copy the content of `openapi.yaml`
2. Paste it at [Swagger Editor](https://editor.swagger.io/) or [ReDoc](https://redocly.github.io/redoc/)

### Locally
Install a tool like [Prism](https://stoplight.io/open-source/prism) for local previewing:
```sh
npm install -g @stoplight/prism-cli
prism mock openapi.yaml
```

## Response Format

All responses follow a consistent JSON structure:

**Success (2xx)**:
```json
{
  "data": { ... },
  "timestamp": "2026-03-23T10:00:00Z"
}
```

**Error (4xx/5xx)**:
```json
{
  "error": "error_code",
  "message": "Human-readable error message"
}
```

## Versioning

- **Current Version**: 1.0.0
- Breaking changes will increment the major version number.


# API Definition (OpenAPI)