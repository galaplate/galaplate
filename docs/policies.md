# Policies

Galaplate includes a robust policy system that allows you to implement authorization, validation, rate limiting, and other access control mechanisms in your application. Policies are reusable middleware components that can be applied to routes and handlers.

## Overview

The policy system in Galaplate is built around the `Policy` interface, which provides a standardized way to implement authorization and validation logic. Each policy evaluates a request context and returns a result indicating whether the action should be allowed or denied.

## Policy Interface

All policies must implement the `Policy` interface:

```go
type Policy interface {
    Name() string
    Evaluate(ctx context.Context, policyCtx *PolicyContext) PolicyResult
}
```

### Policy Context

The `PolicyContext` struct contains all the information needed for policy evaluation:

```go
type PolicyContext struct {
    User     any            // Current authenticated user
    Request  *fiber.Ctx     // Fiber context with request data
    Resource string         // Resource being accessed (usually the route path)
    Action   string         // Action being performed (HTTP method)
    Data     map[string]any // Additional data for policy evaluation
}
```

### Policy Result

Policies return a `PolicyResult` that indicates the outcome:

```go
type PolicyResult struct {
    Allowed bool   // Whether the action is allowed
    Message string // Human-readable message
    Code    int    // HTTP status code to return if denied
}
```

## Using Policies as Middleware

Policies can be easily applied as Fiber middleware:

```go
// Apply single policy
app.Get("/admin", policies.WithPolicies("admin_only"), handler)

// Apply multiple policies
app.Post("/api/users", policies.WithPolicies("auth", "rate_limit"), createUser)

// Or using the struct directly
app.Use("/api", policies.WithPoliciesDirect(new(pkgPolicies.AdminOnlyPolicy)))
```

## Creating Policies

### Basic Policy

Here's how to create a custom policy:

```go
package policies

import (
    "context"
    "github.com/galaplate/core/policies"
    "github.com/gofiber/fiber/v2"
)

type AdminOnlyPolicy struct {
    // Add configuration fields here
}

func NewAdminOnlyPolicy() *AdminOnlyPolicy {
    return &AdminOnlyPolicy{}
}

func (p *AdminOnlyPolicy) Name() string {
    return "admin_only"
}

func (p *AdminOnlyPolicy) Evaluate(ctx context.Context, policyCtx *policies.PolicyContext) policies.PolicyResult {
    // Check if user is authenticated
    if policyCtx.User == nil {
        return policies.PolicyResult{
            Allowed: false,
            Message: "Authentication required",
            Code:    fiber.StatusUnauthorized,
        }
    }

    // Check if user is admin (example logic)
    user, ok := policyCtx.User.(map[string]interface{})
    if !ok {
        return policies.PolicyResult{
            Allowed: false,
            Message: "Invalid user data",
            Code:    fiber.StatusInternalServerError,
        }
    }

    if role, exists := user["role"]; !exists || role != "admin" {
        return policies.PolicyResult{
            Allowed: false,
            Message: "Admin access required",
            Code:    fiber.StatusForbidden,
        }
    }

    return policies.PolicyResult{
        Allowed: true,
        Message: "Admin access granted",
        Code:    fiber.StatusOK,
    }
}

// Register the policy automatically
func init() {
    policies.GlobalPolicyManager.RegisterPolicy(NewAdminOnlyPolicy())
}
```

### Advanced Policy with Configuration

```go
type ResourceOwnerPolicy struct {
    resourceIDParam string
    userIDField     string
}

func NewResourceOwnerPolicy(resourceIDParam, userIDField string) *ResourceOwnerPolicy {
    return &ResourceOwnerPolicy{
        resourceIDParam: resourceIDParam,
        userIDField:     userIDField,
    }
}

func (p *ResourceOwnerPolicy) Name() string {
    return "resource_owner"
}

func (p *ResourceOwnerPolicy) Evaluate(ctx context.Context, policyCtx *policies.PolicyContext) policies.PolicyResult {
    // Get resource ID from URL parameters
    resourceID := policyCtx.Request.Params(p.resourceIDParam)
    if resourceID == "" {
        return policies.PolicyResult{
            Allowed: false,
            Message: "Resource ID not found",
            Code:    fiber.StatusBadRequest,
        }
    }

    // Get user ID from user context
    user, ok := policyCtx.User.(map[string]interface{})
    if !ok {
        return policies.PolicyResult{
            Allowed: false,
            Message: "User not authenticated",
            Code:    fiber.StatusUnauthorized,
        }
    }

    userID, exists := user[p.userIDField]
    if !exists {
        return policies.PolicyResult{
            Allowed: false,
            Message: "User ID not found",
            Code:    fiber.StatusInternalServerError,
        }
    }

    // Check if user owns the resource (simplified example)
    // In real implementation, you'd query the database
    if userID != resourceID {
        return policies.PolicyResult{
            Allowed: false,
            Message: "Access denied: not resource owner",
            Code:    fiber.StatusForbidden,
        }
    }

    return policies.PolicyResult{
        Allowed: true,
        Message: "Resource owner access granted",
        Code:    fiber.StatusOK,
    }
}
```

## Generating Policy Files

Galaplate provides a console command to generate policy templates:

```bash
# Generate a new policy
go run main.go console make:policy MyCustomPolicy

# This creates a new policy file with boilerplate code
```

## Best Practices

### 1. Keep Policies Focused
Each policy should have a single responsibility. Create separate policies for authentication, authorization, rate limiting, etc.

### 2. Use Meaningful Names
Policy names should clearly indicate their purpose:
- `auth` - Authentication required
- `admin_only` - Admin access required
- `rate_limit` - Rate limiting
- `owner_only` - Resource ownership required

### 3. Provide Clear Error Messages
Always return descriptive error messages that help developers and users understand why access was denied.

### 4. Handle Edge Cases
Consider scenarios like missing user data, invalid request format, or database errors.

### 5. Performance Considerations
- Cache expensive operations (like database queries) when possible
- Use context timeouts for long-running operations
- Clean up resources properly

### 6. Testing Policies
Create comprehensive tests for your policies:

```go
func TestAdminOnlyPolicy(t *testing.T) {
    policy := NewAdminOnlyPolicy()

    // Test with no user
    result := policy.Evaluate(context.Background(), &policies.PolicyContext{
        User: nil,
    })
    assert.False(t, result.Allowed)
    assert.Equal(t, fiber.StatusUnauthorized, result.Code)

    // Test with non-admin user
    result = policy.Evaluate(context.Background(), &policies.PolicyContext{
        User: map[string]interface{}{"role": "user"},
    })
    assert.False(t, result.Allowed)
    assert.Equal(t, fiber.StatusForbidden, result.Code)

    // Test with admin user
    result = policy.Evaluate(context.Background(), &policies.PolicyContext{
        User: map[string]interface{}{"role": "admin"},
    })
    assert.True(t, result.Allowed)
}
```

## Integration with Routes

### Domain-based Routes
Policies integrate seamlessly with Galaplate's domain structure:

```go
// In domains/user/routes/routes.go
func SetupUserRoutes(app *fiber.App) {
    api := app.Group("/api/users")

    // Public endpoints
    api.Post("/register", handlers.Register)
    api.Post("/login", handlers.Login)

    // Protected endpoints
    api.Get("/profile",
        policies.WithPolicies("auth"),
        handlers.GetProfile)

    api.Put("/profile",
        policies.WithPolicies("auth"),
        handlers.UpdateProfile)

    // Admin only endpoints
    api.Get("/",
        policies.WithPolicies("auth", "admin_only"),
        handlers.ListUsers)

    api.Delete("/:id",
        policies.WithPolicies("auth", "admin_only", "rate_limit"),
        handlers.DeleteUser)
}
```

### Middleware Chain
Policies are executed in the order they are specified. If any policy fails, the request is rejected and subsequent policies are not executed.

```go
// This chain executes: auth -> admin_only -> rate_limit
api.Post("/admin/action",
    policies.WithPolicies("auth", "admin_only", "rate_limit"),
    handlers.AdminAction)
```

## Error Handling

When a policy fails, the middleware automatically returns a JSON response:

```json
{
    "success": false,
    "message": "Admin access required"
}
```

## Debugging Policies

You can add logging to understand policy execution:

```go
func (p *MyPolicy) Evaluate(ctx context.Context, policyCtx *policies.PolicyContext) policies.PolicyResult {
    log.Printf("Evaluating %s policy for user: %v, resource: %s, action: %s",
        p.Name(), policyCtx.User, policyCtx.Resource, policyCtx.Action)

    // Policy logic here...
}
```

This comprehensive policy system provides the flexibility to implement any authorization or validation logic your application needs while maintaining clean, reusable, and testable code.

## Next Steps

- **[Console Commands](/console-commands)** - Learn how to generate policy files using commands
- **[Routings](/routings)** - Apply policies to your API routes and domain handlers
- **[Validation & DTOs](/validation-and-dto)** - Combine policies with request validation
- **[API Reference](/api-reference)** - See complete examples of secured endpoints
