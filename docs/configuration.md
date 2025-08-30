# Configuration

Galaplate uses environment variables for configuration, making it easy to deploy across different environments while keeping sensitive data secure.

## Environment Variables

All configuration is managed through environment variables defined in your `.env` file.

### Application Settings

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_NAME` | string | `Galaplate` | Application name used in logs and UI |
| `APP_ENV` | string | `local` | Environment: `local`, `staging`, `production` |
| `APP_DEBUG` | boolean | `true` | Enable debug mode and verbose logging |
| `APP_URL` | string | `http://localhost` | Base URL for the application |
| `APP_PORT` | string | `8080` | Port number for the HTTP server |
| `APP_SCREET` | string | **required** | Secret key for JWT and encryption |

### Database Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `DB_CONNECTION` | string | `mysql` | Database driver: `mysql` or `postgres` |
| `DB_HOST` | string | `localhost` | Database server hostname |
| `DB_PORT` | string | `3306` | Database server port |
| `DB_DATABASE` | string | **required** | Database name |
| `DB_USERNAME` | string | **required** | Database username |
| `DB_PASSWORD` | string | | Database password |

### Authentication

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `BASIC_AUTH_USERNAME` | string | **required** | Username for admin endpoints |
| `BASIC_AUTH_PASSWORD` | string | **required** | Password for admin endpoints |

## Environment Files

### `.env` File

Create your `.env` file from the template:

```bash
cp .env.example .env
```

**Example `.env` file:**
```env
# Application
APP_NAME=MyAwesomeAPI
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost
APP_PORT=8080
APP_SCREET=super-secret-key-change-this-in-production

# Database
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=my_awesome_api
DB_USERNAME=root
DB_PASSWORD=my_secure_password

# Authentication
BASIC_AUTH_USERNAME=admin
BASIC_AUTH_PASSWORD=secure_admin_password
```

---

## Next Steps

- **[Database](/database)** - Configure database connections
- **[Console Commands](/console-commands)** - Learn about code generation
- **[Project Structure](/project-structure)** - Understand the codebase
- **[Quick Start](/quick-start)** - Get your application running
