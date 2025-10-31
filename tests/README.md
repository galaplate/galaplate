# Database Testing with RefreshDatabase

This package provides Laravel-style `RefreshDatabase` functionality for Go tests in Galaplate.

## Usage Patterns

### 1. Using WithRefreshDatabase (Recommended)

This refreshes the database once per test suite, similar to Laravel's `RefreshDatabase` trait:

```go
package controllers

import (
    "testing"
    "github.com/galaplate/galaplate/tests"
    "github.com/stretchr/testify/suite"
)

type TestControllerSuite struct {
    tests.WithRefreshDatabase // Laravel-style RefreshDatabase
}

func (t *TestControllerSuite) SetupTest() {
    t.WithRefreshDatabase.SetupTest() // Automatically refreshes database
}

func (suite *TestControllerSuite) TestSomething() {
    // Your test code here
    // Database is clean and ready
}

func TestControllerSuiteRun(t *testing.T) {
    suite.Run(t, new(TestControllerSuite))
}
```

### 2. Using RefreshDatabaseBeforeEachTest

This refreshes the database before each individual test method:

```go
type TestControllerSuite struct {
    tests.RefreshDatabaseBeforeEachTest
}

func (t *TestControllerSuite) SetupTest() {
    t.RefreshDatabaseBeforeEachTest.SetupTest()
}
```

### 3. Without Database Refresh

For tests that don't need database refresh:

```go
type TestControllerSuite struct {
    tests.TestCase // No database refresh
}

func (t *TestControllerSuite) SetupTest() {
    t.TestCase.SetupTest()
}
```

## How It Works

The `RefreshDatabase` functionality:

1. Uses the existing `go run main.go console db:fresh --force` command
2. Loads test environment from `.env.test` file
3. Runs database migrations to ensure clean state
4. Integrates seamlessly with testify/suite

## Requirements

- `.env.testing` file with test database configuration
- `go run main.go console db:fresh` command available
- Test database properly configured (preferably SQLite in-memory for speed)

## Example .env.test

```env
APP_NAME=Galaplate
APP_ENV=testing
DB_CONNECTION=sqlite
DB_DATABASE=:memory:
```

This setup ensures your tests run against a clean database state, similar to Laravel's RefreshDatabase trait.
