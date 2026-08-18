# Galaplate

A Go web framework inspired by Laravel. Built on Fiber and GORM with migrations, background jobs, task scheduling, and a powerful console command system.

## Requirements

- Go 1.24+
- MySQL 8.0+, PostgreSQL 13+, or SQLite

## Quick Start

```bash
git clone https://github.com/galaplate/galaplate.git
cd galaplate
cp .env.example .env
# Edit .env with your database credentials
go run main.go console db:up
make dev
```

## Documentation

Full documentation is in the [docs/](docs/) directory:

- [Quick Start](docs/quick-start.md)
- [Installation](docs/installation.md)
- [Configuration](docs/configuration.md)
- [Project Structure](docs/project-structure.md)
- [Console Commands](docs/console-commands.md)
- [Database](docs/database.md)
- [Routing](docs/routings.md)
- [Validation & DTOs](docs/validation-and-dto.md)
- [Background Tasks](docs/background-tasks.md)
- [Policies](docs/policies.md)
- [File Storage](docs/file-storage.md)
- [Testing](docs/testing.md)
- [API Reference](docs/api-reference.md)

## License

MIT License
