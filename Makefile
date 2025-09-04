# Galaplate Makefile - Development and Database Management Tools

# Default target
.DEFAULT_GOAL := help

# Build the application
build:
	@echo "🔨 Building application..."
	go build -o server main.go

# Run the application (builds first)
run: build
	@echo "🚀 Starting server..."
	./server

# Watch for changes and auto-reload (requires reflex: go install github.com/cespare/reflex@latest)
watch:
	@echo "👀 Watching for changes..."
	reflex -s -r '\.go$$' make run

# Development server with hot reload
dev: watch

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f server

# Format Go code
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./tests/...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests in the tests directory only
# Install development dependencies
install-deps:
	@echo "📦 Installing development dependencies..."
	go install github.com/cespare/reflex@latest
	go install github.com/amacneil/dbmate@latest
	npm install -g dotenv-cli

# Tidy go modules
tidy:
	@echo "🔧 Tidying go modules..."
	go mod tidy

# Show help
help:
	@echo "🚀 Galaplate Development Commands"
	@echo ""
	@echo "📋 Available commands:"
	@echo "  build                       Build the application"
	@echo "  run                         Build and run the application"
	@echo "  dev/watch                   Start development server with hot reload"
	@echo "  clean                       Clean build artifacts"
	@echo "  fmt                         Format Go code"
	@echo "  test                        Run tests"
	@echo "  install-deps                Install development dependencies"
	@echo "  tidy                        Tidy go modules"
	@echo ""
	@echo "🖥️  Console System:"
	@echo "  console                     Show console command help"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make dev                    # Start development server"
	@echo "  make test-coverage          # Run tests with coverage"
	@echo "  go run main.go console list          # List all console commands"
	@echo "  go run main.go console make:model User   # Create User model"
	@echo "  go run main.go console db:create         # Create migration"
	@echo "  go run main.go console db:up             # Run migrations"

# Console command runner
.PHONY: console
console:
	@echo "🖥️  Galaplate Console System"
	@echo ""
	@echo "All database and code generation commands have been migrated to Go!"
	@echo ""
	@echo "📋 Quick commands:"
	@echo "  go run main.go console list              # List all available commands"
	@echo ""
	@echo "🗄️  Database management:"
	@echo "  go run main.go console db:create         # Create migration"
	@echo "  go run main.go console db:up             # Run migrations"
	@echo "  go run main.go console db:down           # Rollback migration"
	@echo "  go run main.go console db:status         # Migration status"
	@echo "  go run main.go console db:seed           # Run seeders"
	@echo ""
	@echo "🏗️  Code generation:"
	@echo "  go run main.go console make:model User   # Create model"
	@echo "  go run main.go console make:dto UserDTO  # Create DTO"
	@echo "  go run main.go console make:job EmailJob # Create job"
	@echo "  go run main.go console make:cron Daily   # Create cron"
	@echo ""
	@echo "💡 For complete list: go run main.go console list"

# Prevent make from trying to build console command arguments as targets
list:
	@:

make\:model make\:dto make\:job make\:cron make\:seeder make\:command:
	@:

example\:demo:
	@:
