# Creating a Package

A Galaplate package (plugin) is a separate Go module that ships routes, console commands, migrations, jobs, cron tasks or policies, and installs into an application with a single blank import. It is built on [Service Providers](/providers): the package registers a provider in `init()`, and the application picks it up when the package is imported.

## Quick Start

From your application, scaffold a package and link it in one step:

```bash
go run main.go console make:package blog --link
```

This creates `../galaplate-blog` (module `github.com/<you>/galaplate-blog`, derived from your app's module path) with a working provider, an example route (`GET /api/blog/ping`), an example command (`blog:hello`), a migrations package, a provider test, a `Makefile` and a `README.md`. With `--link` it also:

- adds `require` and `replace` directives for the package to your app's `go.mod`
- adds `_ "github.com/<you>/galaplate-blog"` to `providers/packages.go` (created if missing)

Then:

```bash
cd ../galaplate-blog && go mod tidy && make test
cd - && go run main.go console blog:hello
```

| Option | Default | Description |
|--------|---------|-------------|
| `--module <path>` | sibling of the app's module, e.g. `github.com/acme/galaplate-blog` | Go module path |
| `--dir <path>` | `../galaplate-<name>` | Where to create the package; must be empty or missing |
| `--link` | off (asked in interactive mode) | Wire the package into the current app for local development |

Run it without arguments to be prompted for each value. The package's `go.mod` requires the same `github.com/galaplate/core` version and `go` directive as your application.

The rest of this page explains what the generated files do.

## How It Works

Everything in Galaplate that a package can contribute is registered from `init()`, so importing the package is enough to install it:

| Contribution | How the package registers it | Picked up by |
|--------------|------------------------------|--------------|
| Setup code | `provider.Add(&Provider{})` with `Register()` / `Boot()` | `bootstrap.NewApp()` |
| HTTP routes | `Routes(app *fiber.App)` on the provider | `provider.MountRoutes`, after the app's own routes |
| Console commands | `Commands() []provider.Command` on the provider | `console.NewKernel()` |
| Migrations | `database.Register(...)` | `db:up`, `db:status`, ... |
| Background jobs | `queue.RegisterJob(...)` | Queue workers |
| Cron tasks | `scheduler.RegisterScheduler(...)` | Scheduler |
| Policies | `policies.GlobalPolicyManager.RegisterPolicy(...)` | `policies.WithPolicies(...)` |

The boot order is:

1. Config, database and file storage are initialized.
2. `Register()` runs on every provider (application and packages).
3. `Boot()` runs on every provider.
4. The application's routes are mounted, then every provider's `Routes()`.

## Package Layout

`make:package blog` generates:

```
galaplate-blog/                  # module github.com/acme/galaplate-blog
├── go.mod                       # requires github.com/galaplate/core
├── provider.go                  # the only public entry point
├── provider_test.go             # boots the provider and hits its routes
├── Makefile                     # test, fmt
├── README.md                    # install steps, routes, commands, aliases
└── internal/
    ├── commands/
    │   └── hello_command.go     # blog:hello
    ├── controllers/
    │   └── ping_controller.go   # GET /api/blog/ping
    └── migrations/
        └── migrations.go        # add your migrations here
```

As the package grows, add `internal/models/`, `internal/jobs/` and so on next to these.

Keep the implementation under `internal/` so the provider is the only importable package. The application cannot reach into the package's internals, and you can change them freely.

Depend only on `github.com/galaplate/core`, never on the application module. Anything the package needs from the application, such as auth middleware, comes through [middleware aliases](#using-application-middleware).

## The Provider

```go
// provider.go
package blog

import (
	"github.com/galaplate/core/fiber"
	"github.com/galaplate/core/provider"
	"github.com/acme/galaplate-blog/internal/commands"
	"github.com/acme/galaplate-blog/internal/controllers"
	_ "github.com/acme/galaplate-blog/internal/migrations"
)

type BlogServiceProvider struct{}

func init() { provider.Add(&BlogServiceProvider{}) }

func (p *BlogServiceProvider) Register() {}

func (p *BlogServiceProvider) Boot() {}

func (p *BlogServiceProvider) Routes(app *fiber.App) {
	blog := app.Group("/api/blog")
	blog.Get("/posts", controllers.PostControllerInstance.Index)
	blog.Post("/posts", provider.Middleware("auth"), controllers.PostControllerInstance.Store)
}

func (p *BlogServiceProvider) Commands() []provider.Command {
	return []provider.Command{
		&commands.SyncPostsCommand{},
	}
}
```

`Routes()` and `Commands()` are optional; implement only the ones you need.

## Console Commands

A command implements `provider.Command`, which has the same methods as a console command in the application:

```go
package commands

import "fmt"

type SyncPostsCommand struct{}

func (c *SyncPostsCommand) GetSignature() string   { return "blog:sync" }
func (c *SyncPostsCommand) GetDescription() string { return "Sync blog posts from the remote feed" }

func (c *SyncPostsCommand) Execute(args []string) error {
	fmt.Println("Syncing posts...")
	return nil
}
```

Prefix signatures with the package name (`blog:sync`) so they don't clash with the application's commands or with other packages.

Commands from packages always get the database and file storage initialized before they run.

## Migrations

Package migrations use the same format as the application's (see [Database](/database)) and are imported by the provider for their `init()` side effect:

```go
package migrations

import "github.com/galaplate/core/database"

type Migration1767000000 struct {
	database.BaseMigration
}

func init() {
	database.Register(&Migration1767000000{
		BaseMigration: database.BaseMigration{
			Name:      "blog_create_posts_table",
			Timestamp: 1767000000,
		},
	})
}

func (m *Migration1767000000) Up(schema *database.Schema) error {
	return schema.Create("blog_posts", func(table *database.Blueprint) {
		table.ID()
		table.String("title").NotNullable()
		table.Text("body")
		table.Timestamps()
	})
}

func (m *Migration1767000000) Down(schema *database.Schema) error {
	return schema.DropIfExists("blog_posts")
}
```

Migrations from every package and the application run together, ordered by timestamp. A migration counts as applied when a row named `<timestamp>_<name>` exists in the migrations table, so:

- Prefix migration names with the package name (`blog_create_posts_table`) so they can never match a migration in the application.
- Prefix table names (`blog_posts`) to avoid colliding with application tables.

## Jobs, Cron Tasks and Policies

These self-register exactly as they do in an application (see [Background Tasks](/background-tasks) and [Policies](/policies)). Put them in an `internal/` package and blank-import it from `provider.go`, like the migrations above. Prefix job types and policy names with the package name (`blog.publish_post`) to keep them unique.

## Using Application Middleware

A package cannot import the application, so it asks for middleware by name:

```go
blog.Post("/posts", provider.Middleware("auth"), handler)
```

The application supplies the alias in its `AppServiceProvider`:

```go
func (p *AppServiceProvider) Register() {
	provider.AliasMiddleware("auth", middleware.JWTAuth())
}
```

`provider.Middleware` panics when the alias is missing, so a misconfigured app fails at startup instead of on a request. Document every alias your package needs in its README.

> Only call `provider.Middleware` from `Routes()` or `Boot()`. Every provider's `Register()` has finished by then, so all aliases exist. Inside `Register()` the alias may not be set yet, because the order in which packages' `init()` functions run (and so the order of providers) is up to Go.

## Database Access

Like applications, packages use the shared connection from `github.com/galaplate/core/database`:

```go
func (c *PostController) Index(ctx fiber.Ctx) error {
	var posts []models.Post
	if err := database.Connect.WithContext(ctx.Context()).Find(&posts).Error; err != nil {
		return err
	}
	return ctx.JSON(posts)
}
```

## Configuration

There is no publish step for packages yet, so a package can't copy config files into the application. Read namespaced keys and fall back to sensible defaults:

```go
perPage := config.ConfigInt("blog.per_page")
if perPage == 0 {
	perPage = 15
}
```

The application can then add a `config/blog.yaml` to override them.

## Installing a Package

`make:package --link` does this for a package you are developing locally. For a published package, in the application:

```bash
go get github.com/acme/galaplate-blog
```

Then add a blank import. The simplest place is a file in `providers/`, which both `main.go` and `tests/test_case.go` already import, so the package is loaded for the server, the console and the tests:

```go
// providers/packages.go
package providers

import (
	_ "github.com/acme/galaplate-blog"
)
```

Then check that everything was picked up:

```bash
go run main.go console list         # shows blog:sync
go run main.go console route:list   # shows /api/blog/posts
go run main.go console db:status    # shows blog_create_posts_table
go run main.go console db:up
```

## Developing a Package Locally

`make:package --link` sets this up for you. To do it by hand, point the application at your working copy with a `replace` directive, the same way this repository uses a local `core`:

```
// galaplate/go.mod
replace github.com/acme/galaplate-blog => ../galaplate-blog
```

Changes to the package are then picked up on the next build. Remove the `replace` before you release the application.

## Testing a Package

Provider wiring can be tested in the package itself with `provider.Reset()` and a bare Fiber app:

```go
func TestRoutesAreMounted(t *testing.T) {
	provider.Reset()
	provider.Add(&blog.BlogServiceProvider{})
	provider.AliasMiddleware("auth", func(c fiber.Ctx) error { return c.Next() })
	provider.Boot()

	app := fiber.New()
	provider.MountRoutes(app)

	resp, err := app.Test(httptest.NewRequest("GET", "/api/blog/posts", nil))
	// ...
}
```

To test against a real database, install the package into an application and write feature tests there (see [Testing](/testing)).

## Checklist

- [ ] Module depends only on `github.com/galaplate/core`
- [ ] Implementation is under `internal/`; `provider.go` is the only public package
- [ ] Commands, migrations, tables, job types and policy names are prefixed with the package name
- [ ] `provider.Middleware` is only called from `Routes()` or `Boot()`
- [ ] README lists the middleware aliases and config keys the package uses
