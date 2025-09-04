# Testing Guide

This guide explains how to write and run tests in Galaplate using the framework's built-in testing infrastructure.

## Quick Start

### Running Tests

```bash
# Run all tests
make test
# or
go test ./...

# Run tests with coverage
make test-coverage

# Run specific test file
go test ./tests/feature/pkg/controllers/test_controller_test.go

# Run tests in verbose mode
go test -v ./...
```

## Test Structure

Galaplate uses Go's native testing framework enhanced with [testify/suite](https://github.com/stretchr/testify) for better organization and setup/teardown capabilities.

### Base Test Case

All feature tests should embed the `tests.TestCase` struct which provides:

- **App**: Configured Fiber application instance
- **DB**: GORM database connection (if needed)
- **Automatic setup**: Environment configuration and app initialization

## Writing Feature Tests

### 1. Test File Structure

Create test files following the pattern: `tests/feature/pkg/[module]/[controller]_test.go`

```go
package controllers

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"
    "testing"

    "github.com/galaplate/galaplate/tests"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type YourControllerSuite struct {
    tests.TestCase
}

func (t *YourControllerSuite) SetupTest() {
    t.TestCase.SetupTest()
    // Add any additional setup here
}

// Test methods go here...

func TestYourControllerSuiteRun(t *testing.T) {
    suite.Run(t, new(YourControllerSuite))
}
```

### 2. HTTP Request Testing

#### GET Requests

```go
func (suite *YourControllerSuite) TestGetEndpoint() {
    t := suite.T()

    req, err := http.NewRequest("GET", "/api/your-endpoint", nil)
    assert.NoError(t, err)

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(200, resp.StatusCode)
}
```

#### POST Requests with JSON

```go
func (suite *YourControllerSuite) TestCreateResource() {
    t := suite.T()

    payload := strings.NewReader(`{"name": "Test Resource"}`)
    req, err := http.NewRequest("POST", "/api/resources", payload)
    assert.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(201, resp.StatusCode)
}
```

#### Testing with Path Parameters

```go
func (suite *YourControllerSuite) TestGetResourceById() {
    t := suite.T()

    req, err := http.NewRequest("GET", "/api/resources/123", nil)
    assert.NoError(t, err)

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(200, resp.StatusCode)
}
```

### 3. Response Validation

#### Basic Status Code Validation

```go
suite.Equal(200, resp.StatusCode)
suite.Equal(201, resp.StatusCode)
suite.Equal(422, resp.StatusCode) // Validation errors
suite.Equal(500, resp.StatusCode) // Server errors
```

#### JSON Response Validation

```go
func (suite *YourControllerSuite) TestResponseStructure() {
    // ... make request ...

    body, err := io.ReadAll(resp.Body)
    suite.NoError(err)

    var response map[string]interface{}
    err = json.Unmarshal(body, &response)
    suite.NoError(err)

    // Validate response structure
    suite.Equal(true, response["success"])
    suite.Equal("Expected message", response["message"])
    
    // Validate nested data
    data := response["data"].(map[string]interface{})
    suite.Equal("expected_value", data["field_name"])
    suite.Equal(float64(123), data["numeric_field"]) // JSON numbers are float64
}
```

### 4. Common Test Scenarios

#### Testing Validation Errors

```go
func (suite *YourControllerSuite) TestValidationError() {
    t := suite.T()

    // Send invalid data (empty required field)
    payload := strings.NewReader(`{"name": ""}`)
    req, err := http.NewRequest("POST", "/api/resources", payload)
    assert.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(422, resp.StatusCode)
}
```

#### Testing Invalid JSON

```go
func (suite *YourControllerSuite) TestInvalidJSON() {
    t := suite.T()

    payload := strings.NewReader(`{"invalid": json}`)
    req, err := http.NewRequest("POST", "/api/resources", payload)
    assert.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(500, resp.StatusCode)
}
```

#### Testing Missing Content-Type

```go
func (suite *YourControllerSuite) TestMissingContentType() {
    t := suite.T()

    payload := strings.NewReader(`{"name": "test"}`)
    req, err := http.NewRequest("POST", "/api/resources", payload)
    assert.NoError(t, err)
    // Intentionally omit Content-Type header

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(422, resp.StatusCode)
}
```

## Best Practices

### 1. Test Organization

- **Group related tests** in the same test suite
- **Use descriptive test method names** that explain what is being tested
- **Follow the AAA pattern**: Arrange, Act, Assert

### 2. Test Isolation

```go
func (suite *YourControllerSuite) SetupTest() {
    suite.TestCase.SetupTest()
    // Each test gets a fresh app instance
    // Add any test-specific setup here
}
```

### 3. Assertion Best Practices

```go
// Use suite assertions for better error messages
suite.Equal(expected, actual)
suite.NoError(err)
suite.Contains(haystack, needle)

// Use assert for simple validations in setup
assert.NoError(t, err)
```

### 4. Testing Authenticated Endpoints

```go
func (suite *YourControllerSuite) TestAuthenticatedEndpoint() {
    t := suite.T()

    req, err := http.NewRequest("GET", "/api/protected", nil)
    assert.NoError(t, err)
    req.Header.Set("Authorization", "Bearer your-test-token")

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(200, resp.StatusCode)
}
```

## Test Environment

The testing framework automatically:

1. Sets `APP_ENV=testing`
2. Changes working directory to project root
3. Initializes the Fiber application with all routes
4. Provides access to database connection (if needed)

### Environment Variables

Create a `.env.testing` file for test-specific configuration:

```env
APP_ENV=testing
DB_CONNECTION=sqlite
DB_DATABASE=:memory:
LOG_LEVEL=error
```

## Database Testing

For tests that require database interaction:

```go
func (suite *YourControllerSuite) TestDatabaseOperation() {
    // Access database through suite.DB
    var count int64
    suite.DB.Model(&YourModel{}).Count(&count)
    suite.Equal(int64(0), count)
}
```

## Running Tests in CI/CD

The framework is designed to work well in CI environments:

```yaml
# Example GitHub Actions step
- name: Run Tests
  run: |
    make test-coverage
```

## Debugging Tests

```bash
# Run with verbose output
go test -v ./tests/feature/pkg/controllers/

# Run specific test method
go test -v -run TestCanCreateTestData ./tests/feature/pkg/controllers/

# Run tests with race detection
go test -race ./...
```

## Example Complete Test Suite

```go
package controllers

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"
    "testing"

    "github.com/galaplate/galaplate/tests"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type ProductControllerSuite struct {
    tests.TestCase
}

func (t *ProductControllerSuite) SetupTest() {
    t.TestCase.SetupTest()
}

func (suite *ProductControllerSuite) TestListProducts() {
    req, err := http.NewRequest("GET", "/api/products", nil)
    assert.NoError(suite.T(), err)

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(200, resp.StatusCode)
}

func (suite *ProductControllerSuite) TestCreateProduct() {
    payload := strings.NewReader(`{"name": "Test Product", "price": 99.99}`)
    req, err := http.NewRequest("POST", "/api/products", payload)
    assert.NoError(suite.T(), err)
    req.Header.Set("Content-Type", "application/json")

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(201, resp.StatusCode)

    body, err := io.ReadAll(resp.Body)
    suite.NoError(err)

    var response map[string]interface{}
    err = json.Unmarshal(body, &response)
    suite.NoError(err)

    suite.Equal(true, response["success"])
    data := response["data"].(map[string]interface{})
    suite.Equal("Test Product", data["name"])
    suite.Equal(99.99, data["price"])
}

func (suite *ProductControllerSuite) TestValidationError() {
    payload := strings.NewReader(`{"name": ""}`)
    req, err := http.NewRequest("POST", "/api/products", payload)
    assert.NoError(suite.T(), err)
    req.Header.Set("Content-Type", "application/json")

    resp, err := suite.App.Test(req)
    suite.NoError(err)
    suite.Equal(422, resp.StatusCode)
}

func TestProductControllerSuiteRun(t *testing.T) {
    suite.Run(t, new(ProductControllerSuite))
}
```

## Next Steps

- **[DTOs & Validation](/validation-and-dto)** - Create DTOs for your test data structures
- **[Database](/database)** - Set up test database configurations
- **[Console Commands](/console-commands)** - Generate test files and components
- **[API Reference](/api-reference)** - Test your API endpoints comprehensively
```