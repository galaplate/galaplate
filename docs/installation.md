# Installation

## Requirements

- Go 1.24 or higher
- MySQL 8.0+, PostgreSQL 13+, or SQLite
- Make (optional, for convenience commands)

## Clone

```bash
git clone https://github.com/galaplate/galaplate.git
cd galaplate
go mod tidy
```

## Database Setup

### MySQL

```sql
CREATE DATABASE galaplate CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### PostgreSQL

```sql
CREATE DATABASE galaplate;
```

### SQLite

No setup required. A database file is created automatically at `storage/database.db`.

## Environment File

```bash
cp .env.example .env
```

Required variables:

| Variable | Description |
|----------|-------------|
| `APP_SECRET` | Application encryption key |
| `DB_CONNECTION` | `mysql`, `postgres`, or `sqlite` |
| `DB_DATABASE` | Database name |
| `DB_USERNAME` | Database user |
| `DB_PASSWORD` | Database password |
| `JWT_SECRET` | Secret for JWT token signing |

## Generate App Key

Generate a secure key for `APP_SECRET` and `JWT_SECRET`:

```bash
openssl rand -base64 32
```

## Run Migrations

```bash
go run main.go console db:up
```

## Start Development

```bash
make dev
```

This starts the server with hot reload via `reflex`.

## Common Issues

**Port already in use:**

Change `APP_PORT` in `.env`.

**Database connection refused:**

Verify the database server is running and credentials in `.env` are correct.

**Missing dependencies:**

```bash
go mod tidy
```
