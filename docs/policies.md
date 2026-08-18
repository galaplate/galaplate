# Policies

Policies provide authorization logic that can be applied as middleware to routes.

## Creating Policies

Generate a policy:

```bash
go run main.go console make:policy AdminOnly
```

This creates `pkg/policies/admin_only_policy.go`:

```go
package policies

import (
    "context"
    "github.com/galaplate/core/policies"
    "github.com/gofiber/fiber/v2"
)

type AdminOnlyPolicy struct{}

func (p *AdminOnlyPolicy) Name() string {
    return "admin_only"
}

func (p *AdminOnlyPolicy) Evaluate(ctx context.Context, policyCtx *policies.PolicyContext) policies.PolicyResult {
    user := policyCtx.User
    if user == nil {
        return policies.PolicyResult{
            Allowed: false,
            Message: "Authentication required",
            Code:    fiber.StatusUnauthorized,
        }
    }

    return policies.PolicyResult{
        Allowed: true,
        Message: "Access granted",
        Code:    fiber.StatusOK,
    }
}

func init() {
    policies.GlobalPolicyManager.RegisterPolicy(&AdminOnlyPolicy{})
}
```

## Using Policies

Apply policies to routes:

```go
import "github.com/galaplate/core/policies"

app.Get("/admin", policies.WithPolicies("admin_only"), adminController.Index)
app.Get("/users", policies.WithPolicies("auth", "admin_only"), userController.Index)
```

If any policy fails, the middleware immediately returns the policy's error response:

```json
{
  "success": false,
  "message": "Authentication required"
}
```

## Policy Context

The `PolicyContext` contains:

| Field | Description |
|-------|-------------|
| `User` | The authenticated user (set by auth middleware) |
| `Request` | The Fiber context |
| `Resource` | The route path |
| `Action` | The HTTP method |
| `Data` | Additional data map |

## Direct Policy Usage

Use policy structs directly without registering by name:

```go
app.Get("/admin", policies.WithPoliciesDirect(&AdminOnlyPolicy{}), handler)
```
