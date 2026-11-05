# Service Providers

Service providers are the central place to wire an application or a package into Galaplate, the same idea as Laravel's service providers. The boilerplate ships three of them in `providers/`, and `main.go` loads them with a blank import:

```go
import _ "github.com/galaplate/galaplate/providers"
```

| Provider | Responsibility |
|----------|----------------|
| `AppServiceProvider` | Middleware aliases and app-wide setup |
| `RouteServiceProvider` | Mounts `router.SetupRouter` |
| `ConsoleServiceProvider` | Registers custom console commands |

## Writing a Provider

Generate one with:

```bash
go run main.go console make:provider Billing --routes --commands
```

This writes `providers/billing_service_provider.go`. `--routes` and `--commands` add `Routes()` and `Commands()` stubs; leave them off for a provider that only needs `Register()`/`Boot()`. Because `providers` is already imported by `main.go` and the tests, the new provider is picked up with no further wiring.

A provider implements `Register()` and `Boot()` and adds itself in `init()`:

```go
package providers

import "github.com/galaplate/core/provider"

type BillingServiceProvider struct{}

func init() { provider.Add(&BillingServiceProvider{}) }

func (p *BillingServiceProvider) Register() {}

func (p *BillingServiceProvider) Boot() {}
```

`Register()` runs for every provider before any `Boot()`, so `Boot()` can rely on what other providers registered. Both run once, after config, database and file storage are initialized.

## Routes

Implement `Routes(app *fiber.App)` to mount routes. Provider routes are mounted after `AppConfig.SetupRoutes`, and they show up in `route:list` and in tests.

```go
func (p *BillingServiceProvider) Routes(app *fiber.App) {
    billing := app.Group("/api/billing", provider.Middleware("auth"))
    billing.Get("/invoices", controllers.InvoiceControllerInstance.Index)
}
```

## Console Commands

Implement `Commands() []provider.Command`:

```go
func (p *ConsoleServiceProvider) Commands() []provider.Command {
    return []provider.Command{
        &commands.SendReportCommand{},
    }
}
```

## Middleware Aliases

Packages cannot import your application's middleware, so the application names them in `AppServiceProvider`:

```go
func (p *AppServiceProvider) Register() {
    provider.AliasMiddleware("auth", middleware.JWTAuth())
}
```

A package then uses `provider.Middleware("auth")`. A missing alias panics at startup. Inside your own application you can keep importing `pkg/middleware` directly.

## Packages

A provider can live in its own Go module. Keep the implementation under `internal/` so only the provider is importable, and install it with a blank import:

```go
import _ "github.com/acme/billing"
```

Its routes and commands are registered without touching `router/router.go` or `main.go`. See [Creating a Package](/packages) for the full guide.
