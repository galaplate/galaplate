# API Reference

This section provides comprehensive documentation for all available API endpoints in Galaplate.

## Base URL

```
http://localhost:8080
```

## Authentication

Galaplate supports multiple authentication methods:

### Basic Authentication
Used for admin endpoints like logs viewer.

```http
Authorization: Basic <base64(username:password)>
```

### JWT Authentication
For protected API endpoints (when implemented).

```http
Authorization: Bearer <jwt_token>
```

## Response Format

All API responses follow a consistent JSON format:

### Success Response
```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error description",
  "message": "Human-readable error message"
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200  | OK - Request successful |
| 201  | Created - Resource created successfully |
| 400  | Bad Request - Invalid request data |
| 401  | Unauthorized - Authentication required |
| 403  | Forbidden - Access denied |
| 404  | Not Found - Resource not found |
| 422  | Unprocessable Entity - Validation failed |
| 500  | Internal Server Error - Server error |

## Endpoints

### Health Check

#### GET /

Basic health check endpoint to verify the server is running.

**Request:**
```http
GET /
```

**Response:**
```
Hello world
```

**Example:**
```bash
curl http://localhost:8080/
```

---

### Logs Viewer

#### GET /logs

Admin endpoint to view application logs. Requires basic authentication.

**Request:**
```http
GET /logs?file=<filename>
Authorization: Basic <credentials>
```

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| file | string | No | Specific log file to view |

**Response:**
Returns an HTML page with log viewer interface.

**Example:**
```bash
# View logs page
curl -u admin:password http://localhost:8080/logs

# View specific log file
curl -u admin:password http://localhost:8080/logs?file=app.2025-06-24.log
```

**Log File Format:**
- Files are stored in `storage/logs/`
- Named with pattern: `app.YYYY-MM-DD.log`
- JSON formatted log entries
- Automatic daily rotation

---

## Data Models

### Job Model

Represents background jobs in the queue system.

```go
type Job struct {
    ID          uint            `json:"id"`
    Type        string          `json:"type"`
    Payload     json.RawMessage `json:"payload"`
    State       JobState        `json:"state"`
    ErrorMsg    string          `json:"error_msg"`
    Attempts    int             `json:"attempts"`
    AvailableAt time.Time       `json:"available_at"`
    CreatedAt   time.Time       `json:"created_at"`
    StartedAt   *time.Time      `json:"started_at"`
    FinishedAt  *time.Time      `json:"finished_at"`
}
```

**Job States:**
- `pending` - Job is waiting to be processed
- `started` - Job is currently being processed
- `finished` - Job completed successfully
- `failed` - Job failed after max attempts

### Log File Model

Represents log files in the logs viewer.

```go
type LogFile struct {
    Name         string    `json:"name"`
    Size         int64     `json:"size"`
    ModifiedTime time.Time `json:"modified_time"`
    Content      string    `json:"content,omitempty"`
}
```

---

## Error Handling

Galaplate implements comprehensive error handling with structured error responses.

### Global Error Handler

All errors are processed through a global error handler that:
- Logs errors for debugging
- Returns consistent error responses
- Handles different error types appropriately

### Validation Errors

When request validation fails, the API returns detailed validation errors:

```json
{
  "success": false,
  "error": "Validation failed",
  "message": "The request data is invalid",
  "details": {
    "field_name": ["Field is required", "Field must be valid email"]
  }
}
```

### Custom Error Types

```go
type GlobalErrorHandlerResp struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error"`
    Status  int    `json:"status"`
}
```

---

## Rate Limiting

Currently, Galaplate doesn't implement rate limiting by default, but it can be easily added using Fiber middleware:

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
}))
```

---

## CORS Configuration

CORS is enabled by default for all origins. To customize:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://example.com",
    AllowHeaders: "Origin, Content-Type, Accept",
}))
```

---

## Request/Response Examples

### Successful Request

```bash
curl -X GET http://localhost:8080/ \
  -H "Content-Type: application/json"
```

**Response:**
```
Hello world
```

### Authentication Required

```bash
curl -X GET http://localhost:8080/logs \
  -H "Content-Type: application/json"
```

**Response:**
```json
{
  "success": false,
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

### With Authentication

```bash
curl -X GET http://localhost:8080/logs \
  -H "Content-Type: application/json" \
  -u admin:password
```

**Response:**
```html
<!DOCTYPE html>
<html>
<!-- Log viewer HTML page -->
</html>
```

---

## SDK Examples

### Go Client Example

```go
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {
    // Health check
    resp, err := http.Get("http://localhost:8080/")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body)) // "Hello world"
}
```

### JavaScript/Node.js Example

```javascript
// Health check
fetch('http://localhost:8080/')
  .then(response => response.text())
  .then(data => console.log(data)); // "Hello world"

// Logs with authentication
fetch('http://localhost:8080/logs', {
  headers: {
    'Authorization': 'Basic ' + btoa('admin:password')
  }
})
  .then(response => response.text())
  .then(html => console.log(html));
```

### Python Example

```python
import requests
from requests.auth import HTTPBasicAuth

# Health check
response = requests.get('http://localhost:8080/')
print(response.text)  # "Hello world"

# Logs with authentication
response = requests.get(
    'http://localhost:8080/logs',
    auth=HTTPBasicAuth('admin', 'password')
)
print(response.text)  # HTML content
```

---

## Testing the API

### Using curl

```bash
# Health check
curl -v http://localhost:8080/

# Logs viewer
curl -v -u admin:password http://localhost:8080/logs

# Check specific log file
curl -v -u admin:password "http://localhost:8080/logs?file=app.2025-06-24.log"
```

### Using HTTPie

```bash
# Health check
http GET localhost:8080/

# Logs viewer
http GET localhost:8080/logs --auth admin:password

# Check specific log file
http GET localhost:8080/logs file==app.2025-06-24.log --auth admin:password
```

### Using Postman

1. **Health Check:**
   - Method: GET
   - URL: `http://localhost:8080/`

2. **Logs Viewer:**
   - Method: GET
   - URL: `http://localhost:8080/logs`
   - Authorization: Basic Auth (username: admin, password: your_password)

---

## Extending the API

### Adding New Endpoints

1. **Create a controller:**
   ```go
   // pkg/controllers/user_controller.go
   func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
       // Implementation
   }
   ```

2. **Register routes:**
   ```go
   // router/router.go
   app.Get("/api/users", userController.GetUsers)
   ```

3. **Add middleware if needed:**
   ```go
   app.Get("/api/users", middleware.Auth(), userController.GetUsers)
   ```

### Adding Authentication

```go
// middleware/auth.go
func JWTAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // JWT validation logic
        return c.Next()
    }
}
```

---

## Next Steps

- **[Console Commands](/console-commands)** - Generate controllers and other components
- **[Routing](/routings)** - Learn about route definitions and middleware
- **[Validation & DTOs](/validation-and-dto)** - Add request validation
- **[Database](/database)** - Work with data models and migrations