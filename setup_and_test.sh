#!/bin/bash

# Setup and Test Runner for Prolog Server
# This script sets up dependencies and runs tests

set -e

echo "🔧 Setting up Prolog Server"
echo "============================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if jq is installed (for pretty JSON in demo)
if ! command -v jq &> /dev/null; then
    echo "⚠️  jq not found. Installing for pretty JSON output..."
    if command -v apt &> /dev/null; then
        sudo apt install jq
    elif command -v brew &> /dev/null; then
        brew install jq
    else
        echo "Please install jq manually for better demo output"
    fi
fi

# Initialize Go module if needed
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init prolog-server
fi

# Install dependencies
echo "📥 Installing dependencies..."
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3

# Make test scripts executable
chmod +x test_prolog.sh 2>/dev/null || true
chmod +x demo_prolog.sh 2>/dev/null || true

echo "✅ Setup complete!"
echo ""
echo "🚀 To run the server:"
echo "   go run main.go"
echo ""
echo "🧪 To run tests (in another terminal):"
echo "   ./test_prolog.sh"
echo ""
echo "🎭 To run demo (in another terminal):"
echo "   ./demo_prolog.sh"
echo ""
echo "📖 Example manual test:"
echo '   curl -X POST http://localhost:8080/facts \'
echo '     -H "Content-Type: application/json" \'
echo '     -d '"'"'{"predicate":{"type":"atom","value":"hello"}}'"'"

# Optionally start server in background for testing
read -p "🔥 Start server in background for testing? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🚀 Starting server in background..."
    go run main.go &
    SERVER_PID=$!
    echo "Server PID: $SERVER_PID"
    
    # Wait for server to start
    sleep 2
    
    echo "🧪 Running tests..."
    if [ -f "test_prolog.sh" ]; then
        ./test_prolog.sh
    else
        echo "test_prolog.sh not found - running basic test"
        curl -s http://localhost:8080/facts > /dev/null && echo "✅ Server is responding"
    fi
    
    echo "🎭 Running demo..."
    if [ -f "demo_prolog.sh" ]; then
        ./demo_prolog.sh
    fi
    
    echo ""
    echo "🛑 To stop the server:"
    echo "   kill $SERVER_PID"
    echo "   or use: pkill -f 'go run main.go'"
fi