# Routing

Routes are defined in `router/router.go` using Fiber's API.

## Basic Routes

```go
package router

import (
    "github.com/gofiber/fiber/v2"
    "github.com/galaplate/galaplate/pkg/controllers"
    "github.com/galaplate/galaplate/pkg/middleware"
)

func SetupRouter(app *fiber.App) {
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello world")
    })
}
```

## Route Methods

```go
app.Get("/users", userController.Index)
app.Post("/users", userController.Store)
app.Put("/users/:id", userController.Update)
app.Delete("/users/:id", userController.Destroy)
```

## Route Parameters

```go
func (c *UserController) Show(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    return ctx.JSON(fiber.Map{"id": id})
}
```

## Query Parameters

```go
page := ctx.Query("page", "1")
search := ctx.Query("search")
```

## Route Groups

```go
api := app.Group("/api")
api.Get("/users", userController.Index)
api.Post("/users", userController.Store)

admin := app.Group("/admin")
admin.Use(middleware.BasicAuth())
admin.Get("/logs", logController.Index)
```

## Middleware

Apply middleware globally, to groups, or to individual routes:

```go
// Global
app.Use(cors.New())

// Group
api := app.Group("/api", middleware.JWTAuth())

// Route
app.Get("/profile", middleware.JWTAuth(), profileController.Show)
```

## Controllers

Controllers are structs with handler methods:

```go
package controllers

import "github.com/gofiber/fiber/v2"

type UserController struct{}

func (c *UserController) Index(ctx *fiber.Ctx) error {
    return ctx.JSON(fiber.Map{
        "success": true,
        "data":    []string{"alice", "bob"},
    })
}

var UserControllerInstance = &UserController{}
```

Register in `router/router.go`:

```go
var userController = controllers.UserControllerInstance
app.Get("/api/users", userController.Index)
```

## Error Responses

Fiber's default error handler returns JSON. Galaplate bootstraps a custom error handler that returns:

```json
{
  "success": false,
  "message": "error message",
  "error":   "error message",
  "status":  500
}
```
