# Database

Galaplate provides a robust database layer built on top of **GORM** with support for both **MySQL** and **PostgreSQL**. The framework includes automatic migration management, connection pooling, configurable GORM settings, and a comprehensive seeding system.

## Overview

The database system in Galaplate is designed for:
- **Multi-database support**: MySQL and PostgreSQL
- **Configurable GORM settings**: Fine-tune GORM behavior for different environments
- **Migration management**: Version-controlled schema changes via console commands
- **Connection pooling**: Optimized database connections
- **Seeding system**: Populate database with test data
- **Code generation**: Automated model and seeder creation via console commands
- **Developer-friendly CLI**: Powerful console commands for all database operations

## Configuration

Database configuration is managed through environment variables in your `.env` file:

```env
# Database Configuration
DB_CONNECTION=mysql          # or 'postgres'
DB_HOST=localhost
DB_PORT=3306                # 5432 for PostgreSQL
DB_DATABASE=galaplate
DB_USERNAME=root
DB_PASSWORD=password
```

### Supported Drivers

- **MySQL**: Primary support with optimized connection string
- **PostgreSQL**: Full support with SSL configuration options

### GORM Configuration

Galaplate allows you to customize GORM behavior through the bootstrap configuration system. This is particularly useful for different environments (development, testing, production).

#### Basic GORM Configuration

In your `main.go`, you can configure GORM settings:

```go
import (
    "time"
    "github.com/galaplate/core/bootstrap"
)

func main() {
    cfg := bootstrap.DefaultConfig()
    cfg.SetupRoutes = router.SetupRouter

    // Configure custom GORM settings
    cfg.DatabaseConfig = &bootstrap.DatabaseConfig{
        GormConfig: &bootstrap.GormConfig{
            DisableForeignKeyConstraintWhenMigrating: true,
            SlowThreshold:                            2 * time.Second,
            LogLevel:                                 "Info", // Silent, Error, Warn, Info
            IgnoreRecordNotFoundError:               true,
            ParameterizedQueries:                    true,
            Colorful:                                false, // Disable in production
        },
    }

    app := bootstrap.App(cfg)
    // ... rest of your app
}
```

#### Environment-Specific Configurations

**Development Environment:**
```go
cfg.DatabaseConfig = &bootstrap.DatabaseConfig{
    GormConfig: &bootstrap.GormConfig{
        DisableForeignKeyConstraintWhenMigrating: true,
        SlowThreshold:                            time.Second,
        LogLevel:                                 "Info",    // Detailed logging
        IgnoreRecordNotFoundError:               false,      // Show all errors
        ParameterizedQueries:                    true,
        Colorful:                                true,       // Pretty colors
    },
}
```

**Production Environment:**
```go
cfg.DatabaseConfig = &bootstrap.DatabaseConfig{
    GormConfig: &bootstrap.GormConfig{
        DisableForeignKeyConstraintWhenMigrating: true,
        SlowThreshold:                            5 * time.Second, // Higher threshold
        LogLevel:                                 "Error",         // Only errors
        IgnoreRecordNotFoundError:               true,
        ParameterizedQueries:                    true,
        Colorful:                                false,           // No colors in logs
    },
}
```

**Console Commands Configuration:**
```go
if len(os.Args) > 1 && os.Args[1] == "console" {
    bootstrap.InitWithConfig(&bootstrap.DatabaseConfig{
        GormConfig: &bootstrap.GormConfig{
            LogLevel: "Silent",  // Quiet for console operations
            Colorful: false,
        },
    })
    // ... rest of console setup
}
```

#### GORM Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `DisableForeignKeyConstraintWhenMigrating` | `bool` | `true` | Disable FK constraints during migration |
| `SlowThreshold` | `time.Duration` | `1 * time.Second` | Threshold for slow query logging |
| `LogLevel` | `string` | `"Warn"` | Log level: "Silent", "Error", "Warn", "Info" |
| `IgnoreRecordNotFoundError` | `bool` | `true` | Don't log "record not found" errors |
| `ParameterizedQueries` | `bool` | `true` | Use parameterized queries (recommended) |
| `Colorful` | `bool` | `true` | Enable colorful log output |

#### Using Default Configuration

If you want to use the default GORM configuration with custom database connection:

```go
cfg.DatabaseConfig = &bootstrap.DatabaseConfig{
    GormConfig: bootstrap.DefaultGormConfig(),
}
```

## Connection Setup

The database connection is established through the bootstrap system, which allows for configurable GORM settings. The core connection logic is in `core/database/connect.go`, but it's configured via the bootstrap layer.

### Connection Architecture

```go
// Global database connection
var Connect *gorm.DB

// Default connection (uses environment variables)
func ConnectDB() {
    // Uses default GORM configuration
}

// Configurable connection (allows custom GORM settings)
func ConnectWithConfig(cfg *Config) {
    // Uses custom database and GORM configuration
}
```

### Connection Flow

1. **Bootstrap Configuration**: GORM settings are defined in `main.go`
2. **Configuration Translation**: Bootstrap config is converted to database config
3. **Connection Establishment**: Database connection is created with custom settings
4. **Global Assignment**: Connection is assigned to global `Connect` variable

### Connection Features

- **Configurable GORM behavior**: Customize logging, performance settings per environment
- **Automatic reconnection**: Handles connection drops gracefully
- **Query logging**: Configurable SQL query logging with slow query detection
- **Parameterized queries**: Protection against SQL injection by default
- **Connection pooling**: Efficient resource management
- **Multi-database support**: MySQL and PostgreSQL with optimized connection strings

### Connection Example

```go
// In main.go - Configure connection behavior
cfg := bootstrap.DefaultConfig()
cfg.DatabaseConfig = &bootstrap.DatabaseConfig{
    GormConfig: &bootstrap.GormConfig{
        SlowThreshold: 2 * time.Second,  // Log queries > 2s
        LogLevel:      "Warn",           // Only warnings and errors
        Colorful:      true,             // Colored output
    },
}

// Connection is automatically established when bootstrap.App(cfg) is called
app := bootstrap.App(cfg)
```

### Accessing the Database Connection

Once connected, use the global connection throughout your application:

```go
import "github.com/galaplate/core/database"

// In your controllers, services, etc.
func GetUsers() []models.User {
    var users []models.User
    database.Connect.Find(&users)
    return users
}
```

## Migrations

Galaplate provides a powerful migration system accessible through console commands, making database schema management simple and efficient. Migrations allow you to version control your database schema and deploy changes across different environments consistently.

### Migration System Overview

- **File-based migrations**: SQL migrations stored in `db/migrations/`
- **Bidirectional**: Each migration has both `up` and `down` operations
- **Atomic operations**: Migrations run within transactions
- **Status tracking**: Keep track of applied migrations
- **Cross-database support**: Works with both MySQL and PostgreSQL

### Migration Structure

Migrations are stored in `db/migrations/` with the following naming convention:
```
YYYYMMDDHHMMSS_migration_description.sql
```

#### Example Migration Files

**Creating a table (`20250830120000_create_users_table.sql`):**
```sql
-- migrate:up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email_verified_at TIMESTAMP NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(active);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- migrate:down
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

**Adding a column (`20250830130000_add_profile_to_users.sql`):**
```sql
-- migrate:up
ALTER TABLE users ADD COLUMN profile_picture VARCHAR(500);
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN last_login_at TIMESTAMP NULL;

-- migrate:down
ALTER TABLE users DROP COLUMN last_login_at;
ALTER TABLE users DROP COLUMN bio;
ALTER TABLE users DROP COLUMN profile_picture;
```

**Creating relationships (`20250830140000_create_posts_table.sql`):**
```sql
-- migrate:up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT false,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_published ON posts(published);

-- migrate:down
DROP INDEX IF EXISTS idx_posts_published;
DROP INDEX IF EXISTS idx_posts_slug;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;
```

### Migration Commands

Galaplate provides both modern **Console Commands** and traditional Make commands for database operations:

#### Console Commands

```bash
# Create a new migration
go run main.go console db:create create_users_table

# Run pending migrations
go run main.go console db:up

# Check migration status
go run main.go console db:status

# Rollback last migration
go run main.go console db:down

# Reset database (rollback all migrations and re-run)
go run main.go console db:reset

# Fresh migration (drop tables and run all migrations)
go run main.go console db:fresh

# Run database seeders
go run main.go console db:seed

# List all database commands
go run main.go console list | grep "db:"
```

### Migration Best Practices

#### **Development Workflow**
1. **Create descriptive migration names**: Use clear, action-oriented names
   - ✅ `create_users_table.sql`
   - ✅ `add_email_index_to_users.sql`
   - ❌ `update_table.sql`

2. **Test both directions**: Always test both `up` and `down` migrations
   ```bash
   # Test up migration
   go run main.go console db:up

   # Test down migration
   go run main.go console db:down

   # Test full cycle
   go run main.go console db:reset
   ```

3. **Use atomic operations**: Each migration should be self-contained
4. **Keep migrations small**: One logical change per migration file

#### **Production Considerations**
- **Always backup**: Create database backups before running migrations in production
- **Review performance impact**: Consider the impact on large datasets
- **Plan for rollbacks**: Ensure down migrations work correctly
- **Monitor execution**: Watch for long-running migrations that might lock tables

#### **Database-Specific Tips**

**PostgreSQL:**
```sql
-- Use IF NOT EXISTS for safer migrations
CREATE TABLE IF NOT EXISTS users (...);

-- Use SERIAL for auto-incrementing IDs
id SERIAL PRIMARY KEY

-- Use proper JSON types
metadata JSONB
```

**MySQL:**
```sql
-- Use AUTO_INCREMENT for auto-incrementing IDs
id INT AUTO_INCREMENT PRIMARY KEY

-- Use proper JSON types (MySQL 5.7+)
metadata JSON
```

#### **Common Migration Patterns**

**Adding indexes safely:**
```sql
-- migrate:up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);

-- migrate:down
DROP INDEX IF EXISTS idx_users_email;
```

**Renaming columns:**
```sql
-- migrate:up
ALTER TABLE users RENAME COLUMN old_name TO new_name;

-- migrate:down
ALTER TABLE users RENAME COLUMN new_name TO old_name;
```

**Data migrations:**
```sql
-- migrate:up
UPDATE users SET status = 'active' WHERE status IS NULL;

-- migrate:down
UPDATE users SET status = NULL WHERE status = 'active';
```

## Models

Galaplate uses **GORM** models with struct tags for database mapping. Models are defined in `pkg/models/` and represent your database tables as Go structs.

> **💡 Learn More About GORM**
> For comprehensive GORM documentation, advanced features, and best practices, visit the official [GORM Documentation](https://gorm.io/docs/). This guide covers Galaplate-specific integration and common patterns.

### Model Structure

Each model should follow these conventions:
- **Struct tags**: Use `gorm` tags for database mapping and `json` tags for API serialization
- **Primary key**: Define ID field with `gorm:"primaryKey"`
- **Timestamps**: Include CreatedAt, UpdatedAt for automatic timestamp management
- **Soft deletes**: Add DeletedAt field for soft delete functionality

### Model Example

```go
package models

import (
    "encoding/json"
    "time"
    "gorm.io/gorm"
)

type JobState string

const (
    JobPending  JobState = "pending"
    JobStarted  JobState = "started"
    JobFinished JobState = "finished"
    JobFailed   JobState = "failed"
)

type Job struct {
    ID          uint            `gorm:"primaryKey" json:"id"`
    Type        string          `gorm:"type:varchar(255);not null" json:"type"`
    Payload     json.RawMessage `gorm:"type:text" json:"payload"`
    State       JobState        `gorm:"type:varchar(16);not null;check:state IN ('pending','started','finished','failed')" json:"state"`
    ErrorMsg    string          `gorm:"type:text" json:"error_msg"`
    Attempts    int             `gorm:"default:0" json:"attempts"`
    AvailableAt *time.Time      `json:"available_at"`

    // Standard GORM timestamps
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}

// TableName explicitly sets the table name
func (Job) TableName() string {
    return "jobs"
}

// Model hooks (optional)
func (j *Job) BeforeCreate(tx *gorm.DB) error {
    if j.State == "" {
        j.State = JobPending
    }
    return nil
}
```

### Common GORM Tags

| Tag | Example | Description |
|-----|---------|-------------|
| `primaryKey` | `gorm:"primaryKey"` | Marks field as primary key |
| `type` | `gorm:"type:varchar(255)"` | Specifies database column type |
| `not null` | `gorm:"not null"` | NOT NULL constraint |
| `unique` | `gorm:"unique"` | UNIQUE constraint |
| `index` | `gorm:"index"` | Creates database index |
| `default` | `gorm:"default:0"` | Sets default value |
| `check` | `gorm:"check:age > 0"` | CHECK constraint |
| `foreignKey` | `gorm:"foreignKey:UserID"` | Foreign key reference |

### Model Generation

Generate new models using Console Commands:

```bash
# Generate a new model
go run main.go console make:model User

# Generate a model with specific name
go run main.go console make:model ProductCategory

# List all make commands available
go run main.go console list | grep "make:"
```

**Generated model example:**
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

    // Add your fields here
}
```

This creates a new model file with:
- Proper struct definition
- GORM tags for database mapping
- JSON tags for API serialization
- Standard timestamps (CreatedAt, UpdatedAt, DeletedAt)
- Soft delete support

### GORM Features in Galaplate

Galaplate leverages these powerful GORM features:

#### **Associations & Relationships**
```go
type User struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    Name     string    `gorm:"type:varchar(255);not null" json:"name"`
    Email    string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
    Profile  Profile   `gorm:"foreignKey:UserID" json:"profile"`        // HasOne
    Posts    []Post    `gorm:"foreignKey:UserID" json:"posts"`          // HasMany
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    UserID uint   `gorm:"not null" json:"user_id"`
    Bio    string `gorm:"type:text" json:"bio"`
    User   User   `gorm:"references:ID" json:"user"` // BelongsTo
}
```

#### **Model Hooks**
```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Hash password, set defaults, etc.
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email, create profile, etc.
    return nil
}
```

#### **Scopes for Reusable Queries**
```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func RecentUsers(db *gorm.DB) *gorm.DB {
    return db.Where("created_at > ?", time.Now().AddDate(0, -1, 0))
}

// Usage
var users []User
database.Connect.Scopes(ActiveUsers, RecentUsers).Find(&users)
```

#### **Soft Deletes**
```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `json:"name"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Soft delete (sets deleted_at timestamp)
database.Connect.Delete(&user)

// Find with deleted records
database.Connect.Unscoped().Find(&users)

// Permanently delete
database.Connect.Unscoped().Delete(&user)
```

> **🔗 Advanced GORM Features**
> For advanced features like custom data types, complex queries, transactions, and performance optimization, refer to the [GORM Documentation](https://gorm.io/docs/).

## Database Seeding

The seeding system allows you to populate your database with test or initial data. Galaplate provides console commands to easily generate and run seeders.

### Seeder Structure

Seeders are located in `db/seeders/` and implement the `Seeder` interface:

```go
func init() {
    seeders.RegisterSeeder("userseeder", &UserSeeder{})
}
```

**Note**: The seeder is automatically registered and will be executed when running `go run main.go console db:seed`.

### Creating Seeders

Generate a new seeder using Console Commands:

```bash
# Create a new seeder
go run main.go console make:seeder UserSeeder

# Create multiple seeders
go run main.go console make:seeder ProductSeeder
go run main.go console make:seeder CategorySeeder
```

This creates a seeder file in `db/seeders/` with the following structure:

```go
package seeders

import (
    "gorm.io/gorm"
    "github.com/galaplate/galaplate/pkg/models"
    "github.com/galaplate/core/seeders"
)

type UserSeeder struct{}

func (s *UserSeeder) Seed(db *gorm.DB) error {
    users := []models.User{
        {
            Name:  "Admin User",
            Email: "admin@example.com",
        },
        // Add more users...
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

### Running Seeders

Use Console Commands to run database seeders:

```bash
go run main.go console db:seed

go run main.go console db:fresh  # Drop tables and re-migrate
```

### Seeder Best Practices

- **Idempotent operations**: Use `FirstOrCreate` to avoid duplicates
- **Environment-specific**: Different data for development/staging/production
- **Dependency order**: Seed related data in proper order
- **Error handling**: Graceful handling of seeding failures

## Database Utilities

### Connection Pooling

GORM automatically handles connection pooling with sensible defaults:

```go
// Configure connection pool (optional)
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Query Optimization

- **Preloading**: Load related data efficiently
- **Indexing**: Proper database indexes for common queries
- **Pagination**: Built-in pagination utilities
- **Raw SQL**: Option to use raw SQL for complex queries

```go
// Preloading example
var users []models.User
db.Preload("Profile").Find(&users)

// Pagination example
result := db.Scopes(utils.Paginate(page, pageSize)).Find(&users)
```

### Database Debugging

Enable SQL query logging for development:

```go
// Enable detailed logging
gormConfig.Logger = logger.Default.LogMode(logger.Info)
```

## Security Considerations

- **Parameterized queries**: All queries use parameter binding
- **Connection encryption**: SSL/TLS support for production
- **Access control**: Database user with minimal required permissions
- **Environment isolation**: Separate databases for different environments

## Performance Tips

1. **Use appropriate indexes**: Add indexes for frequently queried fields
2. **Optimize N+1 queries**: Use preloading and joins
3. **Connection pooling**: Configure pool size based on application load
4. **Query analysis**: Monitor slow queries and optimize
5. **Database maintenance**: Regular VACUUM/OPTIMIZE operations

## Next Steps

- **[Console Commands](/console-commands)** - Learn about database management commands
- **[Models & DTOs](/validation-and-dto)** - Create data models and transfer objects
- **[API Reference](/api-reference)** - Build APIs using your database models
- **[Background Tasks](/background-tasks)** - Process data asynchronously
