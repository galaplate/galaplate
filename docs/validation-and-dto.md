# Validation & DTOs

Galaplate uses `go-playground/validator` for struct validation. DTOs (Data Transfer Objects) define the shape of request and response payloads.

## Creating DTOs

Generate a DTO:

```bash
go run main.go console make:dto CreateUserRequest
```

This creates `pkg/dto/create_user_request.go`:

```go
package dto

import (
    "github.com/galaplate/core/fiber"
    "github.com/galaplate/core/supports"
)

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func (s *CreateUserRequest) Validate(c fiber.Ctx) (*CreateUserRequest, error) {
    if err := supports.NewValidator(c).Validate(s); err != nil {
        return nil, err
    }
    return s, nil
}
```

## Validation in Controllers

```go
func (c *UserController) Store(ctx fiber.Ctx) error {
    req, err := new(dto.CreateUserRequest).Validate(ctx)
    if err != nil {
        return err // Returns 422 with validation errors
    }

    user := models.User{
        Username: req.Username,
        Email:    req.Email,
    }
    database.Connect.Create(&user)

    return ctx.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    user,
    })
}
```

## Validation Tags

| Tag | Description |
|-----|-------------|
| `required` | Field must be present |
| `email` | Must be a valid email |
| `min` | Minimum length or value |
| `max` | Maximum length or value |
| `gte` | Greater than or equal to |
| `lte` | Less than or equal to |
| `eqfield` | Must equal another field |
| `confirmation` | Custom tag for field confirmation |

## Custom Validation

Register custom validators in an `init()` function:

```go
import (
    "github.com/galaplate/core/supports"
    "github.com/go-playground/validator/v10"
)

func init() {
    supports.RegisterValidation("username", func(fl validator.FieldLevel) bool {
        return len(fl.Field().String()) >= 3
    })
}
```

## Validation Response

Failed validation returns a 422 response:

```json
{
  "success": false,
  "status": 422,
  "message": "Field validation for 'email' failed on the 'email' tag",
  "errors": {
    "email": "Field validation for 'email' failed on the 'email' tag",
    "password": "Field validation for 'password' failed on the 'min' tag"
  }
}
```
