# Installation

This guide covers different ways to install and set up Galaplate for your development environment.

## System Requirements

### Minimum Requirements

- **Go**: Version 1.22 or higher
- **Operating System**: Linux, macOS, or Windows
- **Memory**: 512MB RAM minimum (2GB recommended)
- **Disk Space**: 100MB for Galaplate + dependencies

### Database Requirements

Choose one of the following databases:

- **MySQL**: Version 5.7 or higher (8.0+ recommended)
- **PostgreSQL**: Version 12 or higher (14+ recommended)

### Optional Tools

- **Make**: For using the provided Makefile commands
- **Docker**: For containerized development (optional)
- **Git**: For version control

## Installation Methods

### Method 1: Automated Installation Script

The easiest way to get started with Galaplate:

```bash
curl -s https://raw.githubusercontent.com/galaplate/cli/main/install.sh -o /tmp/install.sh && chmod +x /tmp/install.sh && sudo /tmp/install.sh
```

This script will:
- Download the latest Galaplate CLI release
- Install the CLI tool to `/usr/local/bin`
- Set up necessary permissions
- Verify the installation

#### Alternative: Install to Custom Directory

```bash
curl -s https://raw.githubusercontent.com/galaplate/cli/main/install.sh -o /tmp/install.sh && chmod +x /tmp/install.sh && sudo /tmp/install.sh -d ~/.local/bin
```

#### Install Specific Version

```bash
curl -s https://raw.githubusercontent.com/galaplate/cli/main/install.sh -o /tmp/install.sh && chmod +x /tmp/install.sh && sudo /tmp/install.sh -v v0.1.0
```

### Method 2: Manual Download

Download the binary for your platform from the [releases page](https://github.com/galaplate/cli/releases).

```bash
# Example for Linux amd64
wget https://github.com/galaplate/cli/releases/download/v0.1.0/galaplate-linux-amd64.tar.gz
tar -xzf galaplate-linux-amd64.tar.gz
sudo mv galaplate /usr/local/bin/
chmod +x /usr/local/bin/galaplate
```

### Method 3: Build from Source

For developers who want the latest features:

```bash
# Clone the CLI repository
git clone https://github.com/galaplate/cli.git
cd cli

# Build the CLI tool
go build -o galaplate ./cmd/galaplate

# Install to your PATH
sudo mv galaplate /usr/local/bin/
```

## Verify Installation

Check that Galaplate CLI is installed correctly:

```bash
galaplate version
```

You should see output similar to:
```
Galaplate CLI v0.1.0
```

## Create Your First Project

After installation, create your first Galaplate project:

```bash
# Create a basic API project
galaplate new my-api-project

# Create a full-stack project with MySQL
galaplate new my-fullstack-app --template=full --db=mysql

# Create a microservice
galaplate new my-microservice --template=micro
```

This will generate a complete project structure with:
- REST API endpoints
- Database models and migrations
- Console command system
- Built-in code generators

## Development Dependencies

Install additional tools for the best development experience:

### Essential Tools

```bash
# Hot reload tool for development
go install github.com/cespare/reflex@latest

# Environment CLI (requires Node.js)
npm install -g dotenv-cli
```

### Built-in Project Features

Each generated Galaplate project includes:

- **Database Integration**: PostgreSQL and MySQL support with migrations
- **Code Generators**: Built-in console commands for models, controllers, jobs, DTOs
- **Hot Reload**: Development server with automatic restart
- **Environment Management**: `.env` file support

### Using Console Commands

After creating a project, navigate to it and use the built-in console system:

```bash
cd my-api-project

# View all available console commands
go run main.go console list

# Database operations
go run main.go console db:up          # Run migrations
go run main.go console db:down        # Rollback migrations
go run main.go console db:status      # Check migration status
go run main.go console db:fresh       # Reset and run all migrations

# Code generation
go run main.go console make:model User
go run main.go console make:controller UserController
go run main.go console make:job EmailNotification
go run main.go console make:dto UserCreateRequest
```

## Database Setup

### MySQL Installation

<!-- tabs:start -->

#### **macOS (Homebrew)**

```bash
# Install MySQL
brew install mysql

# Start MySQL service
brew services start mysql

# Secure installation (optional but recommended)
mysql_secure_installation
```

#### **Ubuntu/Debian**

```bash
# Update package index
sudo apt update

# Install MySQL
sudo apt install mysql-server

# Start MySQL service
sudo systemctl start mysql
sudo systemctl enable mysql

# Secure installation
sudo mysql_secure_installation
```

#### **CentOS/RHEL**

```bash
# Install MySQL repository
sudo yum install mysql-server

# Start MySQL service
sudo systemctl start mysqld
sudo systemctl enable mysqld

# Get temporary root password
sudo grep 'temporary password' /var/log/mysqld.log

# Secure installation
mysql_secure_installation
```

#### **Windows**

1. Download MySQL installer from [mysql.com](https://dev.mysql.com/downloads/installer/)
2. Run the installer and follow the setup wizard
3. Choose "Developer Default" for a complete installation
4. Set a root password during installation

<!-- tabs:end -->

### PostgreSQL Installation

<!-- tabs:start -->

#### **macOS (Homebrew)**

```bash
# Install PostgreSQL
brew install postgresql

# Start PostgreSQL service
brew services start postgresql

# Create a database user (optional)
createuser --interactive
```

#### **Ubuntu/Debian**

```bash
# Update package index
sudo apt update

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Switch to postgres user
sudo -u postgres psql
```

#### **CentOS/RHEL**

```bash
# Install PostgreSQL
sudo yum install postgresql-server postgresql-contrib

# Initialize database
sudo postgresql-setup initdb

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### **Windows**

1. Download PostgreSQL installer from [postgresql.org](https://www.postgresql.org/download/windows/)
2. Run the installer and follow the setup wizard
3. Remember the password you set for the postgres user

<!-- tabs:end -->

## IDE and Editor Setup

### Visual Studio Code

Recommended extensions for Go development:

```bash
# Install VS Code extensions
code --install-extension golang.go
code --install-extension ms-vscode.vscode-json
code --install-extension bradlc.vscode-tailwindcss
```

### GoLand

GoLand by JetBrains provides excellent Go support out of the box. No additional setup required.

### Vim/Neovim

For Vim users, consider these plugins:
- **vim-go**: Comprehensive Go support
- **coc.nvim**: Language server support
- **nerdtree**: File explorer

## Environment Setup

### Git Configuration

Set up Git hooks for better development workflow:

```bash
# In your Galaplate project directory
git config core.hooksPath .githooks
chmod +x .githooks/*
```

## Troubleshooting

### Common Installation Issues

#### Go Not Found

```bash
# Check if Go is installed
go version

# If not installed, download from https://golang.org/dl/
# Or use package manager:
# macOS: brew install go
# Ubuntu: sudo apt install golang-go
```

#### Permission Denied

```bash
# Fix permissions for global installation
sudo chown -R $(whoami) $(go env GOPATH)
sudo chown -R $(whoami) $(go env GOROOT)
```

#### PATH Issues

```bash
# Add Go binary path to your shell profile
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Database Connection Issues

#### MySQL Connection Refused

```bash
# Check if MySQL is running
sudo systemctl status mysql

# Start MySQL if not running
sudo systemctl start mysql

# Check MySQL port
netstat -tlnp | grep :3306
```

#### PostgreSQL Connection Issues

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check PostgreSQL port
netstat -tlnp | grep :5432

# Connect to PostgreSQL
sudo -u postgres psql
```

### Performance Optimization

#### Go Module Proxy

Speed up dependency downloads:

```bash
# Set Go module proxy
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
```

#### Build Cache

Enable Go build cache for faster builds:

```bash
# Check build cache location
go env GOCACHE

# Clean build cache if needed
go clean -cache
```

## Next Steps

After successful installation:

1. **Create Your Project**: Use `galaplate new my-project` to generate your first application
2. **[Quick Start](/quick-start)** - Learn about running and developing your project
3. **[Configuration](/configuration)** - Configure your database and environment
4. **[Project Structure](/project-structure)** - Understand the generated codebase
5. **[Console Commands](/console-commands)** - Explore the built-in code generators

## Getting Help

If you encounter issues during installation:

- **[CLI Repository](https://github.com/galaplate/cli)** - Galaplate CLI on GitHub
- **[Issues](https://github.com/galaplate/cli/issues)** - Report bugs and request features
- **[Releases](https://github.com/galaplate/cli/releases)** - Download specific versions
- **[Documentation](/)** - Browse the full documentation

---

**Installation complete!** 🎉 You're ready to start building with Galaplate!
