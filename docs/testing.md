# Testing

Galaplate provides a `TestCase` built on [testify/suite](https://github.com/stretchr/testify) for feature testing HTTP endpoints.

## Test Structure

Create test files in `tests/feature/`:

```go
package controllers

import (
    "net/http"
    "strings"
    "testing"

    "github.com/galaplate/galaplate/tests"
    "github.com/stretchr/testify/suite"
)

type UserControllerSuite struct {
    tests.TestCase
}

func TestUserController(t *testing.T) {
    suite.Run(t, new(UserControllerSuite))
}
```

## HTTP Testing

### GET Request

```go
func (s *UserControllerSuite) TestIndex() {
    req, _ := http.NewRequest("GET", "/api/users", nil)
    resp, err := s.App.Test(req)

    s.NoError(err)
    s.Equal(200, resp.StatusCode)
}
```

### POST Request

```go
func (s *UserControllerSuite) TestStore() {
    payload := strings.NewReader(`{"name":"Alice"}`)
    req, _ := http.NewRequest("POST", "/api/users", payload)
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.App.Test(req)

    s.NoError(err)
    s.Equal(201, resp.StatusCode)
}
```

### Authenticated Request

```go
func (s *UserControllerSuite) TestProfile() {
    req, _ := http.NewRequest("GET", "/api/profile", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := s.App.Test(req)
    s.NoError(err)
    s.Equal(200, resp.StatusCode)
}
```

## Database Testing

Access the database in tests:

```go
func (s *UserControllerSuite) TestDatabaseCount() {
    var count int64
    s.DB.Model(&models.User{}).Count(&count)
    s.Equal(int64(0), count)
}
```

## Refresh Database

Enable database refresh for a suite:

```go
type UserControllerSuite struct {
    tests.WithRefreshDatabase
}
```

This runs `db:fresh` before the suite starts.

To refresh before each test:

```go
type UserControllerSuite struct {
    tests.RefreshDatabaseBeforeEachTest
}
```

## Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Specific package
go test ./tests/feature/pkg/controllers/...

# Verbose
go test -v ./...
```

## Test Environment

Tests load `.env.testing` by default. The `TestCase` automatically:

- Sets `APP_ENV=testing`
- Locates the project root
- Boots the Fiber app with routes
- Provides a database connection

## Custom Bootstrap

Override bootstrap settings per test:

```go
func (s *UserControllerSuite) SetupTest() {
    s.Config = tests.DefaultTestConfig()
    s.Config.RefreshDatabase = true
    s.Config.SetupRoutes = router.SetupRouter
    s.TestCase.SetupTest()
}
```
