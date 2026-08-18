# Quick Start

## Requirements

- Go 1.24+
- MySQL 8.0+, PostgreSQL 13+, or SQLite
- Make

## Install

```bash
git clone https://github.com/galaplate/galaplate.git
cd galaplate
go mod tidy
```

## Environment

```bash
cp .env.example .env
```

Edit `.env` with your database credentials. The framework uses YAML config files in `config/` that reference these environment variables.

## Database

```bash
# Run migrations
go run main.go console db:up

# Check status
go run main.go console db:status
```

## Start Server

```bash
# Development with hot reload
make dev

# Or run directly
go run main.go
```

The server starts on the port defined in `config/app.yaml` (default: 8080).

## Verify

```bash
curl http://localhost:8080/
# Hello world

curl http://localhost:8080/api/health
# {"status":"ok","message":"API is working"}
```

## Next Steps

- [Project Structure](/project-structure) — Understand the codebase layout
- [Console Commands](/console-commands) — Generate models, controllers, and more
- [Database](/database) — Write migrations and use the query builder
