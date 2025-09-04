# DTOs & Validation

## Overview

Galaplate uses Data Transfer Objects (DTOs) to define the structure of request and response payloads, and leverages `go-playground/validator` for robust request validation. This ensures type safety, clear API contracts, and reliable error handling.

- **DTOs**: Structs for input/output data with console command generation
- **Validation**: Tag-based, automatic, and extensible
- **Error Handling**: Consistent validation error responses
- **Code Generation**: Streamlined DTO creation via console commands

## Defining DTOs

### Generating DTOs with Console Commands

Use the console command system to quickly generate new DTOs:

```bash
# Generate a new DTO
go run main.go console make:dto CreateUserDto

# Generate DTO with specific name
go run main.go console make:dto ProductDto

# Generate multiple DTOs
go run main.go console make:dto UserLoginDto
go run main.go console make:dto UpdateProfileDto

# List all available make commands
go run main.go console list | grep "make:"
```

**Generated DTO example:**
```go
// pkg/dto/create_user_dto.go
package dto

type CreateUserDto struct {
    // Add your fields here with validation tags
    // Example:
    // Email    string `json:"email" validate:"required,email"`
    // Password string `json:"password" validate:"required,min=8"`
}
```

### Manual DTO Definition

DTOs are Go structs, typically placed in `pkg/dto/` or alongside controllers:

```go
// pkg/dto/create_user_dto.go
package dto

type CreateUserDto struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
}
```

## Using DTOs in Controllers

Bind and validate DTOs in your handler:

```go
func (c *UserController) Register(ctx *fiber.Ctx) error {
    var dto CreateUserDto  // Using generated DTO
    if err := ctx.BodyParser(&dto); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    if err := validator.New().Struct(dto); err != nil {
        return ctx.Status(422).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
    }
    // Proceed with valid data
    user := models.User{
        Email: dto.Email,
        Name:  dto.Name,
        // Map other fields...
    }
    // Save user logic...
}
```

## Validation Tags

- `required` — Field must be present
- `email` — Must be a valid email
- `min`, `max` — Length or value constraints
- `len` — Exact length
- `oneof` — Must match one of the listed values

**Example:**
```go
// Generated using: go run main.go console make:dto ProductDto
type ProductDto struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Description string  `json:"description" validate:"required,min=10,max=500"`
    CategoryID  uint    `json:"category_id" validate:"required,gt=0"`
}
```

## Custom Validation

You can register custom validation functions:

```go
validate := validator.New()
validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
    return len(fl.Field().String()) >= 3
})
```

## Validation Error Response

Validation errors return a 422 status and details:

```json
{
  "error": "Validation failed",
  "details": {
    "email": ["email is required", "email must be valid"],
    "password": ["password is required", "password must be at least 8 characters"]
  }
}
```

## Development Workflow

1. **Generate DTO structure**:
   ```bash
   go run main.go console make:dto UserRegistrationDto
   ```

2. **Add fields and validation tags**:
   ```go
   type UserRegistrationDto struct {
       Email           string `json:"email" validate:"required,email"`
       Password        string `json:"password" validate:"required,min=8,max=100"`
       ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
       Name            string `json:"name" validate:"required,min=2,max=50"`
       Age             int    `json:"age" validate:"required,gte=18,lte=120"`
   }
   ```

3. **Use in controllers** with proper error handling and validation

## Best Practices

### DTO Development with Console Commands

- **Use code generation**: Always start with `go run main.go console make:dto` for consistency
- **Follow naming conventions**: Use descriptive DTO names like `CreateUserDto`, `UpdateProfileDto`
- **One DTO per file**: Each DTO should have its own file in `pkg/dto/`

### General DTO Best Practices

- Define DTOs for all request/response bodies
- Use validation tags for all fields
- Return clear error messages for invalid input
- Keep DTOs and models separate
- Use consistent naming patterns (e.g., `CreateUserDto`, `UpdateUserDto`)

## Console Commands Reference

```bash
# DTO Development
go run main.go console make:dto <DtoName>    # Generate new DTO
go run main.go console list                  # List all available commands
```

## Common DTO Patterns

### Request DTOs
```bash
# Generate request DTOs
go run main.go console make:dto CreateUserDto
go run main.go console make:dto UpdateUserDto
go run main.go console make:dto LoginDto
```

### Response DTOs
```bash
# Generate response DTOs
go run main.go console make:dto UserResponseDto
go run main.go console make:dto ApiResponseDto
```

## Next Steps

- **[Console Commands](/console-commands)** - Learn about all available generation commands
- **[Testing](/testing)** - Write tests for your DTOs and validation logic
- **[API Reference](/api-reference)** - Build endpoints using your DTOs
- **[Database](/database)** - Create models that work with your DTOs
