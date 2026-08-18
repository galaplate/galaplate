# Introduction

Galaplate is a Go web framework inspired by Laravel. It provides a structured, convention-driven approach to building REST APIs using Fiber, GORM, and a powerful console command system.

## What is Galaplate?

Galaplate combines modern Go idioms with developer experience patterns from Laravel. It includes everything you need to build production APIs: routing, database migrations, background jobs, task scheduling, request validation, authentication, and file storage.

## Stack

| Component | Library |
|-----------|---------|
| HTTP Framework | [Fiber v2](https://gofiber.io/) |
| ORM | [GORM](https://gorm.io/) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Scheduling | [robfig/cron](https://github.com/robfig/cron) |
| Logging | Structured JSON with rotation |
| Databases | MySQL, PostgreSQL, SQLite |

## Core vs Application

Galaplate is split into two parts:

- **`github.com/galaplate/core`** — The framework library. Provides bootstrap, database, queue, scheduler, policies, logger, file storage, config, console commands, and testing utilities.
- **`github.com/galaplate/galaplate`** — The application boilerplate. Uses core to build a working API project.

You work within the application boilerplate. The core package is updated independently.

## Quick Start

```bash
git clone https://github.com/galaplate/galaplate.git
cd galaplate
cp .env.example .env
# Edit .env with your database credentials
go run main.go console db:up
make dev
```

## Available Documentation

- [Quick Start](/quick-start)
- [Installation](/installation)
- [Configuration](/configuration)
- [Project Structure](/project-structure)
- [Console Commands](/console-commands)
- [Database](/database)
- [Routing](/routings)
- [Validation & DTOs](/validation-and-dto)
- [Background Tasks](/background-tasks)
- [Policies](/policies)
- [File Storage](/file-storage)
- [Testing](/testing)
- [API Reference](/api-reference)
