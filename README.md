# 🧠 Prolog Engine

A clean, modern Prolog engine with REST API and optional web UI, built in Go.

## Features

✅ **Core Prolog Engine**
- Unification & Backtracking
- Tabling/Memoization 
- Built-in predicates (=, atom, var, number, now, date functions)
- Aggregation functions (count, sum, max, min)

✅ **Session Management**
- SQLite persistence
- Named sessions with descriptions
- Session-isolated facts and rules

✅ **REST API**
- Complete CRUD operations for sessions, facts, rules
- Query execution with JSON responses
- Optional API key protection

✅ **Web UI** (Optional)
- Browser-based terminal emulator
- Interactive help sidebar with examples
- Session management interface
- Optional password protection

## Quick Start

### 1. Basic Usage
```bash
# Build and start with basic UI
make ui-basic

# API only (no UI)
make api-only

# With password protection
make ui-protected
```

### 2. Development
```bash
# Development mode with auto-rebuild
make dev

# Run tests
make test

# Clean build
make clean build
```

### 3. Production
```bash
# Production server (all interfaces, random password)
make server

# Maximum security (localhost only, random keys)
make full-secure
```

## Available Make Targets

### 📋 Basic Commands
- `make build` - Build the application
- `make clean` - Clean build artifacts and database  
- `make test` - Run all tests
- `make deps` - Update dependencies

### 🚀 Development
- `make dev` - Start development server
- `make dev-watch` - Start with file watching (requires `entr`)
- `make dev-clean` - Clean build and start fresh

### 🌐 Server Modes
- `make api-only` - API only, no UI
- `make ui-basic` - Basic UI without password
- `make ui-protected` - UI with password (admin123)
- `make server` - Production mode (0.0.0.0, random password)

### 🔒 Secure Modes  
- `make api-secure` - API with random API key
- `make ui-secure` - UI with random password and API key
- `make full-secure` - Maximum security (localhost only)

### ⚙️ Custom Configuration
- `make config` - Create .env from template
- `make custom` - Start with custom .env configuration

### 🎭 Demo & Testing
- `make demo` - Start demo mode
- `make demo-data` - Demo with sample data  
- `make demo-reset` - Reset demo environment

### 🛠️ Utilities
- `make status` - Show current configuration
- `make backup` - Backup database
- `make help` - Show all commands

## Configuration

The application can be configured via environment variables or `.env` file:

```bash
# Server configuration
HOST=localhost          # Default: localhost
PORT=8080               # Default: 8080

# Security (optional)
API_KEY=your-secret-key # Protects API routes
UI_PASSWORD=admin123    # Protects web UI

# Features (optional)  
ENABLE_UI=true          # Enable web interface
```

### Custom Configuration

```bash
# Create config template
make config

# Edit .env file with your settings
nano .env

# Start with custom config
make custom
```

## API Usage

### Sessions
```bash
# Create session
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"Test session"}'

# List sessions  
curl http://localhost:8080/api/v1/sessions
```

### Facts & Rules
```bash
# Add fact
curl -X POST http://localhost:8080/api/v1/sessions/1/facts \
  -H "Content-Type: application/json" \
  -d '{"predicate":{"type":"compound","value":"parent","args":[{"type":"atom","value":"tom"},{"type":"atom","value":"bob"}]}}'

# Add rule
curl -X POST http://localhost:8080/api/v1/sessions/1/rules \
  -H "Content-Type: application/json" \
  -d '{"head":{"type":"compound","value":"grandparent","args":[...]},"body":[...]}'
```

### Queries
```bash
# Execute query
curl -X POST http://localhost:8080/api/v1/sessions/1/query \
  -H "Content-Type: application/json" \
  -d '{"goals":[{"type":"compound","value":"parent","args":[{"type":"variable","value":"X"},{"type":"atom","value":"bob"}]}]}'
```

## Web UI Usage

1. Start with UI enabled: `make ui-basic`
2. Open browser to `http://localhost:8080/ui`
3. Create or select a session
4. Use the terminal to interact with Prolog:

```prolog
% Add facts
parent(tom, bob).
parent(bob, alice).

% Add rules  
grandparent(X, Z) :- parent(X, Y), parent(Y, Z).

% Query
?- grandparent(tom, X)
```

### Terminal Commands
- `help` - Show available commands
- `clear` - Clear terminal
- `sessions` - List all sessions
- Use ↑/↓ arrows for command history

## Examples

### Family Relationships
```prolog
% Facts
parent(tom, bob).
parent(bob, alice).
parent(alice, charlie).

% Rules
grandparent(X, Z) :- parent(X, Y), parent(Y, Z).
ancestor(X, Y) :- parent(X, Y).
ancestor(X, Y) :- parent(X, Z), ancestor(Z, Y).

% Queries  
?- grandparent(tom, X)     % Find grandchildren of tom
?- ancestor(tom, charlie)  % Check if tom is ancestor of charlie
```

### Aggregation
```prolog  
% Score facts
score(alice, 95).
score(bob, 87).
score(charlie, 92).

% Queries
?- count(_, score(X, Y), N)           % Count all scores
?- sum(Score, score(X, Score), Total) % Sum all scores  
?- max(Score, score(X, Score), Max)   % Find highest score
```

### Date/Time
```prolog
% Get current time
?- now(X)

% Date comparisons
?- date_before(date("2023-01-01"), date("2023-12-31"))
?- days_between(date("2023-01-01"), date("2023-01-31"), Days)
```

## Development

### Prerequisites
- Go 1.21+
- Make
- SQLite3

### Building from Source
```bash
git clone <repository>
cd golog
make deps
make build
make test
```

### File Watching (Optional)
Install `entr` for automatic rebuilds:
```bash
# macOS
brew install entr

# Ubuntu/Debian  
sudo apt-get install entr

# Then use
make dev-watch
```

## Architecture

- **types.go** - Core data structures and types
- **engine.go** - Prolog engine with SQLite persistence  
- **handlers.go** - HTTP handlers and middleware
- **templates.go** - Embedded HTML templates and JavaScript
- **main.go** - Application entry point

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]