# Routing

## Overview

Galaplate uses the Fiber framework for high-performance HTTP routing. All routes are defined in `router/router.go` and are organized for clarity, scalability, and maintainability. Routing connects HTTP endpoints to controller methods, applies middleware, and supports grouping for modular APIs.

- **Centralized routing** in `router/router.go`
- **RESTful conventions** for endpoints
- **Middleware** for authentication, CORS, etc.
- **Route groups** for API versioning and organization
- **Controller integration** for clean separation of logic

## Defining Routes

Routes are registered using Fiber's API:

```go
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello world")
})

app.Get("/logs", middleware.BasicAuth(), logController.ShowLogsPage)
```

## Route Methods

- `app.Get(path, handler)`
- `app.Post(path, handler)`
- `app.Put(path, handler)`
- `app.Delete(path, handler)`
- `app.Use(middleware)`

## Route Groups

Group related routes for modularity and versioning:

```go
api := app.Group("/api")
api.Get("/users", userController.GetUsers)
api.Post("/users", userController.CreateUser)
```

## Middleware

Apply middleware globally or to specific routes/groups:

```go
app.Use(cors.New()) // Global
app.Get("/logs", middleware.BasicAuth(), logController.ShowLogsPage) // Route-level
```

## Route Parameters and Query

```go
// Path parameter
app.Get("/api/users/:id", userController.GetUserByID)

// Query parameter
app.Get("/api/search", userController.SearchUsers)
```

## Controller Integration

Controllers are Go structs with handler methods, placed in `pkg/controllers/`.

```go
// pkg/controllers/user_controller.go
package controllers

import "github.com/gofiber/fiber/v2"

type UserController struct{}

func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
    // Fetch users from DB
    return ctx.JSON(fiber.Map{
        "success": true,
        "data": []string{"alice", "bob"},
    })
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
    // Parse and create user
    return ctx.Status(201).JSON(fiber.Map{
        "success": true,
        "message": "User created",
    })
}
```

Register controllers in `router/router.go`:

```go
var userController = controllers.UserController{}
app.Get("/api/users", userController.GetUsers)
```

## Example: Full Route Setup

```go
func SetupRouter(app *fiber.App) {
    app.Use(cors.New())
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello world")
    })
    app.Get("/logs", middleware.BasicAuth(), logController.ShowLogsPage)
    api := app.Group("/api")
    api.Get("/users", userController.GetUsers)
    api.Post("/users", userController.CreateUser)
}
```

## Best Practices

- Use route groups for versioning and organization
- Apply middleware at the group or route level as needed
- Keep controller logic in `pkg/controllers/`
- Use clear, RESTful route naming

## Next Steps

- **[Controllers](/console-commands)** - Learn about generating controllers
- **[Middleware](/middleware)** - Understand authentication and other middleware
- **[API Reference](/api-reference)** - See complete API endpoint examples
- **[Validation & DTOs](/validation-and-dto)** - Add request validation to your routes
