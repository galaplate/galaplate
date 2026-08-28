# Database

Galaplate uses GORM with a schema builder and Go-based migrations.

## Configuration

See `config/database.yaml`. Supported drivers: `mysql`, `postgres`, `sqlite`.

## Connection

The database connection is established automatically during bootstrap. Access it via:

```go
import "github.com/galaplate/core/database"

database.Connect.Find(&users)
```

### Context-Aware Queries

In HTTP controllers and middleware, pass the request context to DB calls so queries respect cancellation and timeouts:

```go
func (c *UserController) Index(ctx fiber.Ctx) error {
    db := database.Connect.WithContext(ctx.Context())
    db.Find(&users)
    return ctx.JSON(users)
}
```

In background jobs and schedulers, use the context provided by the worker:

```go
func (j SendEmail) Handle(ctx context.Context, payload json.RawMessage) error {
    db := database.Connect.WithContext(ctx)
    db.Find(&users)
    return nil
}
```

## Migrations

Migrations are Go files that use a Blueprint schema builder.

### Creating Migrations

```bash
go run main.go console db:create create_posts_table
```

This generates a migration file in `db/migrations/`:

```go
package migrations

import "github.com/galaplate/core/database"

type Migration1234567890 struct {
    database.BaseMigration
}

func init() {
    database.Register(&Migration1234567890{
        BaseMigration: database.BaseMigration{
            Name:      "create_posts_table",
            Timestamp: 1234567890,
        },
    })
}

func (m *Migration1234567890) Up(schema *database.Schema) error {
    return schema.Create("posts", func(table *database.Blueprint) {
        table.ID()
        table.String("title").NotNullable()
        table.String("slug").Unique().NotNullable()
        table.Text("content")
        table.Boolean("published").Default(false)
        table.Timestamps()
    })
}

func (m *Migration1234567890) Down(schema *database.Schema) error {
    return schema.DropIfExists("posts")
}
```

### Running Migrations

```bash
go run main.go console db:up
go run main.go console db:status
go run main.go console db:down
go run main.go console db:fresh
go run main.go console db:reset
```

### Blueprint Methods

#### Columns

| Method | Description |
|--------|-------------|
| `ID()` | Auto-incrementing primary key |
| `String(name, length...)` | VARCHAR |
| `Text(name)` | TEXT |
| `Integer(name)` | INT |
| `BigInteger(name)` | BIGINT |
| `Boolean(name)` | BOOLEAN / TINYINT(1) |
| `JSON(name)` | JSON / JSONB |
| `Timestamp(name)` | TIMESTAMP |
| `Timestamps()` | `created_at` and `updated_at` |
| `Decimal(name, precision, scale)` | DECIMAL |
| `UUID(name)` | UUID |
| `Enum(name, values[])` | ENUM |
| `Date(name)` | DATE |
| `DateTime(name)` | DATETIME |
| `Float(name)` | FLOAT |
| `Double(name)` | DOUBLE |
| `Char(name, length...)` | CHAR |
| `Blob(name)` | BLOB |
| `TinyInt(name)` | TINYINT |
| `SmallInt(name)` | SMALLINT |
| `MediumInt(name)` | MEDIUMINT |
| `TinyText(name)` | TINYTEXT |
| `MediumText(name)` | MEDIUMTEXT |
| `LongText(name)` | LONGTEXT |

#### Modifiers

| Method | Description |
|--------|-------------|
| `.NotNullable()` | NOT NULL |
| `.Nullable()` | NULL |
| `.Default(value)` | Default value |
| `.Unique()` | UNIQUE constraint |
| `.Unsigned()` | UNSIGNED (MySQL) |
| `.Comment(text)` | Column comment (MySQL) |

#### Indexes

```go
table.Index([]string{"user_id"})
table.UniqueIndex([]string{"email"})
table.Primary([]string{"id", "tenant_id"})
```

#### Foreign Keys

```go
table.Foreign("user_id").References("id").On("users").OnDelete("CASCADE").Finish()
```

#### Altering Tables

```go
// Add columns
schema.Table("users", func(table *database.Blueprint) {
    table.String("avatar").Nullable()
    table.Index([]string{"created_at"})
})

// Modify columns
schema.Table("users", func(table *database.Blueprint) {
    table.Modify("name").String("name", 100).NotNullable()
})

// Drop columns
schema.Table("users", func(table *database.Blueprint) {
    table.DropColumn("temporary_field")
    table.DropIndex("idx_users_old")
})
```

## Models

Generate models with:

```bash
go run main.go console make:model Post
```

Example model:

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `gorm:"size:255;not null" json:"title"`
    Slug      string         `gorm:"size:255;uniqueIndex" json:"slug"`
    Content   string         `gorm:"type:text" json:"content"`
    Published bool           `gorm:"default:false" json:"published"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
```

## Querying

Galaplate uses GORM for queries. See the [GORM documentation](https://gorm.io/docs/) for full query capabilities.

```go
var post Post
database.Connect.First(&post, 1)

database.Connect.Where("published = ?", true).Find(&posts)

database.Connect.Create(&post)

database.Connect.Model(&post).Update("published", true)

database.Connect.Delete(&post)
```

## Pagination

```go
import "github.com/galaplate/core/supports"

result := database.Connect.Scopes(supports.Paginate(page, pageSize)).Find(&posts)
```

## Seeders

Generate a seeder:

```bash
go run main.go console make:seeder UserSeeder
```

Implement the `Seeder` interface:

```go
package seeders

import (
    "gorm.io/gorm"
    "github.com/galaplate/core/database/seeders"
    "github.com/galaplate/galaplate/pkg/models"
)

type UserSeeder struct{}

func (s *UserSeeder) Seed(db *gorm.DB) error {
    users := []models.User{
        {Username: "admin", Email: "admin@example.com"},
    }
    for _, user := range users {
        db.FirstOrCreate(&user, models.User{Email: user.Email})
    }
    return nil
}

func init() {
    seeders.RegisterSeeder("userseeder", &UserSeeder{})
}
```

Run seeders:

```bash
go run main.go console db:seed
```
