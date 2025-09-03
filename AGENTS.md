# Galaplate Agent Context

This guide provides essential dev context for automated agents working in this repository (Galaplate — a Go-based REST API boilerplate).

## 1. Build, Lint, and Test
- **Build**: `make build` or `go build -o server main.go`
- **Run**: `make run` (builds then runs)
- **Dev (hot reload)**: `make dev` (needs reflex)
- **Format**: `make fmt` or `go fmt ./...`
- **Test all**: `make test` or `go test ./...`
- **Test (single file)**: `go test ./path/to/pkg/file_test.go`
- **Test Coverage**: `make test-coverage`
- **Linting**: Uses `go fmt` (strict gofmt style); no golangci-lint config by default

## 2. Console Command System
Galaplate features a powerful console command system for development tasks:

### Database Operations
- **Run migrations**: `go run main.go console db:up`
- **Check migration status**: `go run main.go console db:status`
- **Rollback migration**: `go run main.go console db:down`
- **Fresh migration**: `go run main.go console db:fresh`
- **Create migration**: `go run main.go console db:create <name>`
- **Run seeders**: `go run main.go console db:seed`

### Code Generation
- **Generate model**: `go run main.go console make:model <ModelName>`
- **Generate DTO**: `go run main.go console make:dto <DtoName>`
- **Generate job**: `go run main.go console make:job <JobName>`
- **Generate seeder**: `go run main.go console make:seeder <SeederName>`
- **Generate cron job**: `go run main.go console make:cron <CronName>`

### Utility Commands
- **List all commands**: `go run main.go console list`
- **Run example**: `go run main.go console example`
- **Interactive demo**: `go run main.go console interactive`

### Alternative Make Commands (Legacy)
- **Database**: `make db-up`, `make db-down`, `make db-status`, `make db-fresh`
- **Development**: `make dev`, `make run`, `make build`, `make clean`

## 3. Code Style Guidelines
- **Imports**: Standard library, then third-party, then local imports (grouped, no blank lines between groups)
- **Formatting**: Always use `gofmt`; tabs not spaces
- **Types**: Prefer explicit struct field types; use `interface{}` or `any` sparingly (replace where modern Go suggests)
- **Naming**: `CamelCase` for structs/interfaces, `camelCase` for vars, `ALL_CAPS` for constants; exported symbols are Capitalized
- **Error handling**: Return Go `error` type; wrap or enrich errors with context (e.g., `fmt.Errorf`); log errors with `logs.Error`/`logs.Fatal`
- **JSON struct tags**: Always used for API inputs/outputs
- **Tests**: Standard `_test.go` files; assertions via Go testing or testify if added
- **Env/config**: Managed via `env/` package and `.env` file

## 4. Project Structure
- App code is under `cmd/`, `pkg/`, `domains/`, `router/`
- DB/migrations/seed in `db/`, support utility code in `pkg/supports/`
- Console commands in `pkg/console/commands/`

## 5. Custom Command Development
To create new console commands:
1. Create command file in `pkg/console/commands/`
2. Implement `Command` interface with `GetSignature()`, `GetDescription()`, `Execute(args []string) error`
3. Register in `pkg/console/commands.go` via `RegisterCommands()`

No Cursor or Copilot rules present. If new lint, type, or style guidelines are added, update this file accordingly.