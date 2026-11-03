# Galaplate Makefile - Development and Database Management Tools

# Default target
.DEFAULT_GOAL := help

# Check for forbidden direct gorm imports in boilerplate code
# GORM should only be imported in core; boilerplate should use core/database aliases.
check-gorm:
	@found=$$(grep -rn 'gorm\.io/gorm' --include="*.go" . | grep -v '/vendor/' || true); \
	if [ -n "$$found" ]; then \
		echo "❌ Forbidden: direct gorm import detected in boilerplate code:"; \
		echo "$$found"; \
		echo ""; \
		echo "Use github.com/galaplate/core/database instead (database.DeletedAt, database.ErrRecordNotFound)."; \
		exit 1; \
	fi

# Build the application
build: check-gorm
	@echo "🔨 Building application..."
	go build -o server main.go

# Run the application (builds first)
run: build
	@echo "🚀 Starting server..."
	./server

# Watch for changes and auto-reload (requires reflex: go install github.com/cespare/reflex@latest)
watch: export LOG_STDOUT=true
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
test: check-gorm
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
	npm install -g dotenv-cli --force || echo "dotenv-cli already installed or npm not available"

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
	@echo "  console <cmd> [args]        Build once, then run a console command (e.g. make console db:up)"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make dev                    # Start development server"
	@echo "  make test-coverage          # Run tests with coverage"
	@echo "  make console                # List all console commands"
	@echo "  make console make:model User   # Create User model"
	@echo "  make console db:create         # Create migration"
	@echo "  make console db:up             # Run migrations"

# Console command runner: builds ./server once (incremental), then runs the
# binary instead of `go run`, which relinks the whole app on every call.
#   make console                      # list commands
#   make console db:up
#   make console make:model User
#   make console ARGS="db:seed --force"   # use ARGS for flags make would eat
ifeq (console,$(firstword $(MAKECMDGOALS)))
  CONSOLE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # Turn the extra words into no-op targets so make doesn't try to build them.
  $(eval $(subst :,\:,$(CONSOLE_ARGS)):;@:)
endif

.PHONY: console check-gorm
console: build
	@./server console $(CONSOLE_ARGS) $(ARGS)
