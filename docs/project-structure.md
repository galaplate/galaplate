# Project Structure

```
galaplate/
├── config/                 # YAML configuration files
│   ├── app.yaml
│   ├── database.yaml
│   ├── auth.yaml
│   └── filesystems.yaml
├── console/
│   └── kernel.go          # Custom console command registration
├── db/
│   ├── migrations/        # Go-based migration files
│   └── seeders/           # Database seeders
├── docs/                  # Documentation (Docsify)
├── pkg/
│   ├── controllers/       # HTTP request handlers
│   ├── dto/               # Request/response structs
│   ├── jobs/              # Background job handlers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # GORM models
│   ├── policies/          # Authorization policies
│   └── scheduler/         # Cron task handlers
├── router/
│   └── router.go          # Route definitions
├── storage/
│   ├── logs/              # Application logs
│   └── app/               # File uploads
├── templates/             # HTML templates
├── tests/                 # Test suites
├── main.go                # Application entry point
├── go.mod
├── Makefile
└── .env
```

## Entry Point

`main.go` bootstraps the application:

```go
func main() {
    app := bootstrap.NewApp(withSetupRoutes)

    if len(os.Args) > 1 && os.Args[1] == "console" {
        // Run console commands
        kernel := console.NewKernel()
        pkgConsole.RegisterCommands(kernel)
        kernel.Run(os.Args)
        return
    }

    // Start HTTP server
    port := config.ConfigString("app.port")
    app.Listen(":" + port)
}
```

## Where to Put Code

| Task | Location |
|------|----------|
| Add a route | `router/router.go` |
| Add a controller | `pkg/controllers/` |
| Add a model | `pkg/models/` |
| Add a migration | `db/migrations/` |
| Add a seeder | `db/seeders/` |
| Add middleware | `pkg/middleware/` |
| Add a background job | `pkg/jobs/` |
| Add a cron task | `pkg/scheduler/` |
| Add a policy | `pkg/policies/` |
| Add a DTO | `pkg/dto/` |
