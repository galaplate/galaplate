# Console Commands

Galaplate includes a console command system for code generation and database management.

## Usage

```bash
go run main.go console <command> [arguments]
```

## Database Commands

| Command | Description |
|---------|-------------|
| `db:up` | Run pending migrations |
| `db:down` | Rollback last batch |
| `db:status` | Show migration status |
| `db:fresh` | Drop all tables and re-run migrations |
| `db:reset` | Rollback all and re-run |
| `db:seed` | Run database seeders |
| `db:create <name>` | Create a new migration file |

## Generator Commands

| Command | Description | Output |
|---------|-------------|--------|
| `make:model <Name>` | Generate a GORM model | `pkg/models/name.go` |
| `make:dto <Name>` | Generate a DTO struct | `pkg/dto/name.go` |
| `make:job <Name>` | Generate a queue job | `pkg/jobs/name.go` |
| `make:cron <Name>` | Generate a cron task | `pkg/scheduler/name.go` |
| `make:seeder <Name>` | Generate a seeder | `db/seeders/name.go` |
| `make:factory <Name>` | Generate a model factory | `db/factories/name.go` |
| `make:policy <Name>` | Generate a policy | `pkg/policies/name.go` |

## Custom Commands

Register custom commands in `console/kernel.go`:

```go
package console

import "github.com/galaplate/core/console"

func RegisterCommands(kernel *console.Kernel) {
    kernel.Register(&MyCustomCommand{})
}
```

A command implements the `Command` interface:

```go
type Command interface {
    GetSignature() string
    GetDescription() string
    Execute(args []string) error
}
```

Example:

```go
package commands

import "fmt"

type GreetCommand struct{}

func (c *GreetCommand) GetSignature() string {
    return "greet"
}

func (c *GreetCommand) GetDescription() string {
    return "Print a greeting"
}

func (c *GreetCommand) Execute(args []string) error {
    name := "world"
    if len(args) > 0 {
        name = args[0]
    }
    fmt.Printf("Hello, %s!\n", name)
    return nil
}
```
