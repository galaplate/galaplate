# Configuration

Galaplate uses YAML configuration files located in the `config/` directory. Environment variables are interpolated at load time.

## Config Files

| File | Purpose |
|------|---------|
| `config/app.yaml` | Application name, env, port, key |
| `config/database.yaml` | Database connections |
| `config/auth.yaml` | Authentication guards |
| `config/filesystems.yaml` | File storage drivers |

## Environment Variables

Galaplate loads `.env` variables into the shell. YAML config files reference them with `${VAR:default}` syntax.

### Application

```yaml
# config/app.yaml
name: Galaplate
env: ${APP_ENV:local}
debug: ${APP_DEBUG:false}
port: ${APP_PORT:8080}
key: ${APP_SECRET}
```

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `local` | Environment name |
| `APP_DEBUG` | `false` | Enable debug mode |
| `APP_PORT` | `8080` | HTTP server port |
| `APP_SECRET` | — | Encryption key (required) |

### Database

```yaml
# config/database.yaml
default: ${DB_CONNECTION:sqlite}

connections:
  sqlite:
    driver: sqlite
    database: ${DB_SQLITE_PATH:./storage/database.db}

  mysql:
    driver: mysql
    host: ${DB_HOST:localhost}
    port: ${DB_PORT:3306}
    database: ${DB_DATABASE:galaplate}
    username: ${DB_USERNAME:root}
    password: ${DB_PASSWORD:}
    pool_size: 10
    max_idle_connections: 5

  postgres:
    driver: postgres
    host: ${DB_HOST:localhost}
    port: ${DB_PORT:5432}
    database: ${DB_DATABASE:galaplate}
    username: ${DB_USERNAME:postgres}
    password: ${DB_PASSWORD:}
```

### Authentication

```yaml
# config/auth.yaml
default: jwt

guards:
  jwt:
    driver: jwt
    secret: ${JWT_SECRET:your-secret-key}
    expiration: ${JWT_EXPIRATION:86400}
    refresh_expiration: ${JWT_REFRESH_EXPIRATION:604800}
    algorithm: HS256
```

### File Storage

```yaml
# config/filesystems.yaml
default: ${FILESYSTEM_DRIVER:local}
max_size: ${FILESYSTEM_MAX_SIZE:10485760}

allowed_types:
  - image/jpeg
  - image/png
  - application/pdf

disks:
  local:
    driver: local
    path: ${FILESYSTEM_LOCAL_PATH:storage/app/uploads}

  s3:
    driver: s3
    region: ${AWS_REGION:}
    bucket: ${AWS_BUCKET:}
    key: ${AWS_ACCESS_KEY_ID:}
    secret: ${AWS_SECRET_ACCESS_KEY:}
```

## Reading Config in Code

```go
import "github.com/galaplate/core/config"

// Dot notation
name := config.ConfigString("app.name")
port := config.ConfigInt("app.port")
debug := config.ConfigBool("app.debug")

// Database config
host := config.ConfigString("database.connections.mysql.host")
```

## Custom Config Files

Add new YAML files to `config/`. They are automatically loaded and available via dot notation:

```yaml
# config/mail.yaml
default: smtp

drivers:
  smtp:
    host: ${SMTP_HOST:}
    port: ${SMTP_PORT:587}
```

```go
host := config.ConfigString("mail.drivers.smtp.host")
```
