# Quick Start

Get up and running with Galaplate in just a few minutes!

## Prerequisites

Before you begin, make sure you have:

- **Go 1.22+** installed on your system
- **MySQL** or **PostgreSQL** database server
- **Git** for version control
- **Make** utility (usually pre-installed on Unix systems)

## Installation Methods

### Method 1: Using Galaplate CLI (Recommended)

The fastest way to get started with the new CLI approach:

```bash
# Install the Galaplate CLI
curl -s https://raw.githubusercontent.com/galaplate/cli/main/install.sh -o /tmp/install.sh && chmod +x /tmp/install.sh && sudo /tmp/install.sh

# Create your first project
galaplate new my-awesome-api
cd my-awesome-api
```

### Method 2: Manual Installation

```bash
# Download and install CLI manually
wget https://github.com/galaplate/cli/releases/download/v0.1.0/galaplate-linux-amd64.tar.gz
tar -xzf galaplate-linux-amd64.tar.gz
sudo mv galaplate /usr/local/bin/

# Create your project
galaplate new my-awesome-api
cd my-awesome-api
```

## Create Your First Project

The Galaplate CLI will generate a complete project structure for you:

```bash
# Create a new API project (default template)
galaplate new my-awesome-api

# Or create with specific options
galaplate new my-fullstack-app --template=full --db=postgres

# Navigate to your project
cd my-awesome-api

# Install dependencies (automatically done by CLI)
go mod tidy
```

The CLI generates:
- Complete project structure with all necessary files
- Database models and migrations
- Console command system
- Environment configuration
- Hot reload setup

## Environment Configuration

1. **Copy the environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit your `.env` file:**
   ```env
   APP_NAME=MyAwesomeAPI
   APP_ENV=local
   APP_DEBUG=true
   APP_URL=http://localhost
   APP_PORT=8080
   APP_SECRET=your-super-secret-key-here

   # Database Configuration
   DB_CONNECTION=mysql
   DB_HOST=localhost
   DB_PORT=3306
   DB_DATABASE=my_awesome_api
   DB_USERNAME=root
   DB_PASSWORD=your_password

   # Basic Auth for Admin Endpoints
   BASIC_AUTH_USERNAME=admin
   BASIC_AUTH_PASSWORD=secure_password
   ```

3. **Generate a secure secret key:**
   ```bash
   # Generate a random 32-character secret
   openssl rand -base64 32
   ```

## Database Setup

### Create Database

<!-- tabs:start -->

#### **MySQL**

```sql
CREATE DATABASE my_awesome_api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### **PostgreSQL**

```sql
CREATE DATABASE my_awesome_api;
```

<!-- tabs:end -->

### Run Migrations

```bash
# View all available console commands
go run main.go console list

# Run database migrations
go run main.go console db:up

# Check migration status
go run main.go console db:status

# Seed the database (optional)
go run main.go console db:seed
```

## Start Development Server

### Option 1: Hot Reload Development (Recommended)

```bash
# Start development server with hot reload
make dev
```

This will:
- Install `reflex` if not already installed
- Watch for file changes
- Automatically rebuild and restart the server

### Option 2: Manual Build and Run

```bash
# Build and run
make run

# Or build separately
make build
./server
```

## Verify Installation

Once your server is running, you can test it:

### 1. Basic Health Check

```bash
curl http://localhost:8080/
```

**Expected Response:**
```
Hello world
```

### 2. Access Logs (Admin Endpoint)

```bash
curl -u admin:secure_password http://localhost:8080/logs
```

This should return an HTML page showing your application logs.

### 3. Check Server Logs

Your application logs will be available in:
```
storage/logs/app.YYYY-MM-DD.log
```

## Next Steps

Congratulations! 🎉 Your Galaplate application is now running. Here's what you can do next:

### 1. Explore the Codebase
- **[Project Structure](/project-structure)** - Understand how the code is organized
- **[Configuration](/configuration)** - Learn about all configuration options

### 2. Build Your API
- **[DTOs & Validation](/validation-and-dto)** - Create data models and validation
- **[Database](/database)** - Work with databases and migrations
- **[Routing](/routings)** - Define API endpoints and routes

### 3. Add Features
- **[Policies & Security](/policies)** - Implement authentication and security
- **[Background Tasks](/background-tasks)** - Process tasks asynchronously
- **[Console Commands](/console-commands)** - Generate boilerplate code

### 4. Development Tools
- **[Testing](/testing)** - Write and run tests
- **[API Reference](/api-reference)** - Explore available endpoints

## Common Issues

### Port Already in Use

If port 8080 is already in use:

```bash
# Change the port in your .env file
APP_PORT=3000

# Or set it temporarily
APP_PORT=3000 make run
```

### Database Connection Issues

1. **Check your database is running:**
   ```bash
   # MySQL
   brew services start mysql
   # or
   sudo systemctl start mysql

   # PostgreSQL
   brew services start postgresql
   # or
   sudo systemctl start postgresql
   ```

2. **Verify connection details in `.env`**

3. **Test database connection:**
   ```bash
   make db-connect
   ```

### Missing Dependencies

If you encounter missing dependencies:

```bash
# Install all development dependencies
make install-deps

# Tidy go modules
make tidy
```

## Development Commands

Here are the most commonly used commands during development:

```bash
# Development
make dev              # Start with hot reload
make run              # Build and run once
make build            # Build binary only
make clean            # Clean build artifacts

# Database (Console Commands)
go run main.go console db:up        # Run migrations
go run main.go console db:down      # Rollback migration
go run main.go console db:status    # Check migration status
go run main.go console db:fresh     # Reset and migrate
go run main.go console db:seed      # Run database seeders

# Code Quality
make fmt              # Format code
make test             # Run tests

# Code Generation (Console Commands)
go run main.go console make:model User      # Generate new model
go run main.go console make:dto UserDto     # Generate new DTO
go run main.go console make:job ProcessData # Generate background job
go run main.go console make:seeder UserSeeder # Generate database seeder
go run main.go console make:controller UserController # Generate controller

# Console System
go run main.go console list         # List all available commands
go run main.go console example      # Run example command
go run main.go console interactive  # Interactive demo
```

## Getting Help

- **[Full Documentation](/)**
- **[API Reference](/api-reference)**
- **[Examples](/examples)**
- **[CLI Repository](https://github.com/galaplate/cli)**
- **[GitHub Issues](https://github.com/galaplate/cli/issues)**

---

**You're all set!** 🚀 Start building your amazing API with Galaplate!
