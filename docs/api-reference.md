# API Reference

## Base URL

```
http://localhost:8080
```

## Response Format

### Success

```json
{
  "success": true,
  "data": {},
  "message": "..."
}
```

### Error

```json
{
  "success": false,
  "message": "...",
  "error": "..."
}
```

### Validation Error

```json
{
  "success": false,
  "status": 422,
  "message": "Field validation failed",
  "errors": {
    "field": "error message"
  }
}
```

## Authentication

### Register

```http
POST /api/register
Content-Type: application/json

{
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

Response `201`:

```json
{
  "success": true,
  "data": {
    "user": { "id": 1, "username": "alice", "email": "alice@example.com" },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "message": "User registered successfully"
}
```

### Login

```http
POST /api/login
Content-Type: application/json

{
  "email": "alice@example.com",
  "password": "secret123"
}
```

Response `200`:

```json
{
  "success": true,
  "data": {
    "user": { "id": 1, "username": "alice", "email": "alice@example.com" },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "message": "Login successful"
}
```

### Profile

```http
GET /api/profile
Authorization: Bearer {token}
```

Response `200`:

```json
{
  "success": true,
  "message": "Profile data",
  "data": { "id": 1, "username": "alice" }
}
```

## Test Endpoints

### Health Check

```http
GET /api/health
```

Response `200`:

```json
{ "status": "ok", "message": "API is working" }
```

### Create Test Data

```http
POST /api/test
Content-Type: application/json

{ "name": "Sample", "description": "A test item" }
```

Response `201`:

```json
{
  "success": true,
  "message": "Test data created successfully",
  "data": { "id": 1, "name": "Sample", "description": "A test item" }
}
```

### Get Test Data

```http
GET /api/test/{id}
```

Response `200`:

```json
{
  "success": true,
  "data": { "id": "1", "name": "Test Item", "description": "This is a test item" }
}
```

## Log Viewer (Admin)

The log viewer provides an HTML interface for browsing application logs.

### View Logs

```http
GET /admin/logs
```

### Export Logs

```http
GET /admin/logs/export?file=app.2025-01-15.log&format=json
GET /admin/logs/export?file=app.2025-01-15.log&format=csv
```

### Cleanup Old Logs

```http
POST /admin/logs/cleanup?days=30
```

Response:

```json
{
  "success": true,
  "deleted_count": 12,
  "total_size_mb": 4.5,
  "cutoff_date": "2024-12-16",
  "retention_days": 30
}
```

### Log Statistics

```http
GET /admin/logs/stats
```

Response:

```json
{
  "success": true,
  "total_files": 15,
  "total_size_mb": 8.2,
  "oldest_date": "2025-01-01",
  "newest_date": "2025-01-15"
}
```

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 422 | Unprocessable Entity |
| 500 | Internal Server Error |
