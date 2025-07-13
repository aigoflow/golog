# Prolog Engine Makefile
# Provides convenient targets for different modes and operations

# Default values
BINARY_NAME=golog
HOST=localhost
PORT=3000
DB_FILE=prolog.db

# Build targets
.PHONY: build clean test lint fmt vet deps

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✅ Build complete"

clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f $(DB_FILE)
	@rm -f .env
	@echo "✅ Cleanup complete"

test:
	@echo "🧪 Running tests..."
	@go test -v
	@echo "✅ All tests passed"

test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

lint:
	@echo "🔍 Running linter..."
	@go vet ./...
	@echo "✅ Linting complete"

fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Formatting complete"

vet:
	@echo "🔬 Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete"

deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependencies updated"

# Development targets
.PHONY: dev dev-watch dev-clean

dev: build
	@echo "🚀 Starting development server..."
	@./$(BINARY_NAME)

dev-watch:
	@echo "👀 Starting with file watching (requires 'entr')..."
	@echo "Install with: brew install entr (macOS) or apt-get install entr (Ubuntu)"
	@find . -name "*.go" | entr -r make dev

dev-clean: clean build dev

# Production modes
.PHONY: api-only ui-basic ui-protected server

api-only: build
	@echo "🔌 Starting API-only mode..."
	@echo "HOST=$(HOST)" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo ""
	@echo "📡 API Mode - No Web UI"
	@echo "🌐 Available at: http://$(HOST):$(PORT)/api/v1/"
	@echo ""
	@./$(BINARY_NAME)

ui-basic: build
	@echo "🖥️  Starting with basic UI..."
	@echo "HOST=$(HOST)" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "ENABLE_UI=true" >> .env
	@echo ""
	@echo "🌐 Basic UI Mode"
	@echo "🖥️  Web UI at: http://$(HOST):$(PORT)/ui"
	@echo "📡 API at: http://$(HOST):$(PORT)/api/v1/"
	@echo ""
	@./$(BINARY_NAME)

ui-protected: build
	@echo "🔒 Starting with password-protected UI..."
	@echo "HOST=$(HOST)" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "ENABLE_UI=true" >> .env
	@echo "UI_PASSWORD=admin123" >> .env
	@echo ""
	@echo "🔒 Protected UI Mode"
	@echo "🖥️  Web UI at: http://$(HOST):$(PORT)/ui (Password: admin123)"
	@echo "📡 API at: http://$(HOST):$(PORT)/api/v1/"
	@echo ""
	@./$(BINARY_NAME)

server: build
	@echo "🏭 Starting production server..."
	@echo "HOST=0.0.0.0" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "ENABLE_UI=true" >> .env
	@echo "UI_PASSWORD=$$(openssl rand -base64 12)" >> .env
	@echo ""
	@echo "🏭 Production Server Mode"
	@echo "🌐 Listening on all interfaces: http://0.0.0.0:$(PORT)"
	@echo "🔒 UI Password: $$(grep UI_PASSWORD .env | cut -d= -f2)"
	@echo ""
	@./$(BINARY_NAME)

# Secure modes
.PHONY: api-secure ui-secure full-secure

api-secure: build
	@echo "🔐 Starting API with security..."
	@echo "HOST=$(HOST)" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "API_KEY=$$(openssl rand -base64 32)" >> .env
	@echo ""
	@echo "🔐 Secure API Mode"
	@echo "📡 API at: http://$(HOST):$(PORT)/api/v1/"
	@echo "🔑 API Key: $$(grep API_KEY .env | cut -d= -f2)"
	@echo ""
	@./$(BINARY_NAME)

ui-secure: build
	@echo "🔐 Starting secure UI mode..."
	@echo "HOST=$(HOST)" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "ENABLE_UI=true" >> .env
	@echo "UI_PASSWORD=$$(openssl rand -base64 12)" >> .env
	@echo "API_KEY=$$(openssl rand -base64 32)" >> .env
	@echo ""
	@echo "🔐 Secure UI Mode"
	@echo "🖥️  Web UI at: http://$(HOST):$(PORT)/ui"
	@echo "🔒 UI Password: $$(grep UI_PASSWORD .env | cut -d= -f2)"
	@echo "🔑 API Key: $$(grep API_KEY .env | cut -d= -f2)"
	@echo ""
	@./$(BINARY_NAME)

full-secure: build
	@echo "🛡️  Starting maximum security mode..."
	@echo "HOST=127.0.0.1" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "ENABLE_UI=true" >> .env
	@echo "UI_PASSWORD=$$(openssl rand -base64 16)" >> .env
	@echo "API_KEY=$$(openssl rand -base64 32)" >> .env
	@echo ""
	@echo "🛡️  Maximum Security Mode"
	@echo "🖥️  Web UI at: http://127.0.0.1:$(PORT)/ui"
	@echo "🔒 UI Password: $$(grep UI_PASSWORD .env | cut -d= -f2)"
	@echo "🔑 API Key: $$(grep API_KEY .env | cut -d= -f2)"
	@echo "⚠️  Only accessible from localhost"
	@echo ""
	@./$(BINARY_NAME)

# Custom configuration
.PHONY: custom config

custom: build
	@echo "⚙️  Starting with custom configuration..."
	@if [ ! -f .env ]; then \
		echo "❌ No .env file found. Run 'make config' first or copy from .env.example"; \
		exit 1; \
	fi
	@echo "📋 Using configuration from .env file:"
	@cat .env | sed 's/^/   /'
	@echo ""
	@./$(BINARY_NAME)

config:
	@echo "⚙️  Creating custom configuration..."
	@if [ -f .env ]; then \
		echo "⚠️  .env file already exists. Backing up to .env.backup"; \
		cp .env .env.backup; \
	fi
	@cp .env.example .env
	@echo "✅ Configuration template copied to .env"
	@echo "📝 Edit .env file with your preferred settings, then run 'make custom'"

# Demo and testing modes
.PHONY: demo demo-data demo-reset

demo: ui-basic
	@echo "🎭 Demo mode started!"

demo-data: build
	@echo "📊 Setting up demo with sample data..."
	@echo "HOST=$(HOST)" > .env
	@echo "PORT=$(PORT)" >> .env
	@echo "ENABLE_UI=true" >> .env
	@echo ""
	@echo "🎭 Demo Mode with Sample Data"
	@echo "🖥️  Web UI at: http://$(HOST):$(PORT)/ui"
	@echo "📚 Sample data will be available in the UI"
	@echo ""
	@./$(BINARY_NAME) &
	@sleep 2
	@echo "🔄 Loading sample data..."
	@curl -s -X POST http://$(HOST):$(PORT)/api/v1/sessions \
		-H "Content-Type: application/json" \
		-d '{"name":"family-demo","description":"Family relationships demo"}' > /dev/null || true
	@echo "✅ Demo data loaded!"
	@echo "🎯 Try these examples in the UI:"
	@echo "   parent(tom, bob)."
	@echo "   parent(bob, alice)."
	@echo "   grandparent(X,Z) :- parent(X,Y), parent(Y,Z)."
	@echo "   ?- grandparent(tom, X)"
	@wait

demo-reset: clean
	@echo "🔄 Resetting demo environment..."
	@rm -f $(DB_FILE)
	@echo "✅ Demo reset complete"

# Utility targets
.PHONY: help status logs backup restore

help:
	@echo "🧠 Prolog Engine Makefile Commands"
	@echo "=================================="
	@echo ""
	@echo "📋 Basic Commands:"
	@echo "   build         - Build the application"
	@echo "   clean         - Clean build artifacts and database"
	@echo "   test          - Run all tests"
	@echo "   test-coverage - Run tests with coverage report"
	@echo "   deps          - Update dependencies"
	@echo ""
	@echo "🚀 Development:"
	@echo "   dev           - Start development server"
	@echo "   dev-watch     - Start with file watching (requires entr)"
	@echo "   dev-clean     - Clean build and start fresh"
	@echo ""
	@echo "🌐 Server Modes:"
	@echo "   api-only      - API only, no UI"
	@echo "   ui-basic      - Basic UI without password"
	@echo "   ui-protected  - UI with password (admin123)"
	@echo "   server        - Production mode (0.0.0.0, random password)"
	@echo ""
	@echo "🔒 Secure Modes:"
	@echo "   api-secure    - API with random API key"
	@echo "   ui-secure     - UI with random password and API key"
	@echo "   full-secure   - Maximum security (localhost only)"
	@echo ""
	@echo "⚙️  Custom Configuration:"
	@echo "   config        - Create .env from template"
	@echo "   custom        - Start with custom .env configuration"
	@echo ""
	@echo "🎭 Demo & Testing:"
	@echo "   demo          - Start demo mode"
	@echo "   demo-data     - Demo with sample data"
	@echo "   demo-reset    - Reset demo environment"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "   status        - Show current configuration"
	@echo "   logs          - Show recent activity (if running)"
	@echo "   backup        - Backup database"
	@echo "   restore       - Restore database"
	@echo ""
	@echo "💡 Examples:"
	@echo "   make ui-basic HOST=192.168.1.100 PORT=9000"
	@echo "   make server PORT=3000"
	@echo "   make custom    # after editing .env file"

status:
	@echo "📊 Current Status"
	@echo "================"
	@echo "Binary: $(BINARY_NAME) $$(if [ -f $(BINARY_NAME) ]; then echo '✅'; else echo '❌ (run make build)'; fi)"
	@echo "Database: $(DB_FILE) $$(if [ -f $(DB_FILE) ]; then echo '✅'; else echo '❌ (will be created on first run)'; fi)"
	@echo "Config: .env $$(if [ -f .env ]; then echo '✅'; else echo '❌ (using defaults)'; fi)"
	@if [ -f .env ]; then \
		echo ""; \
		echo "📋 Current Configuration:"; \
		cat .env | sed 's/^/   /'; \
	fi
	@echo ""
	@echo "🔗 Process: $$(pgrep -f $(BINARY_NAME) | wc -l | xargs) instance(s) running"

logs:
	@echo "📄 Recent Activity"
	@echo "=================="
	@if [ -f $(DB_FILE) ]; then \
		sqlite3 $(DB_FILE) "SELECT name, description, created_at FROM sessions ORDER BY created_at DESC LIMIT 5;" 2>/dev/null || echo "No session data available"; \
	else \
		echo "No database file found"; \
	fi

backup:
	@echo "💾 Creating backup..."
	@if [ -f $(DB_FILE) ]; then \
		cp $(DB_FILE) $(DB_FILE).backup.$$(date +%Y%m%d_%H%M%S); \
		echo "✅ Database backed up to $(DB_FILE).backup.$$(date +%Y%m%d_%H%M%S)"; \
	else \
		echo "❌ No database file to backup"; \
	fi

restore:
	@echo "🔄 Available backups:"
	@ls -la $(DB_FILE).backup.* 2>/dev/null || echo "No backups found"
	@echo ""
	@echo "To restore, run: cp $(DB_FILE).backup.TIMESTAMP $(DB_FILE)"

# Default target
.DEFAULT_GOAL := help

# Variables can be overridden from command line
# Example: make ui-basic HOST=192.168.1.100 PORT=9000