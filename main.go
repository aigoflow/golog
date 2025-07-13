package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (optional)
	godotenv.Load()

	engine, err := NewEngine("prolog.db")
	if err != nil {
		log.Fatal("Failed to initialize engine:", err)
	}
	defer engine.Close()

	r := engine.setupRoutes()

	fmt.Println("🧠 Clean Prolog Engine v2.0")
	fmt.Println("===========================")
	fmt.Println("✅ Unification & Backtracking")
	fmt.Println("✅ Tabling/Memoization")
	fmt.Println("✅ Aggregation (count, sum, max, min)")
	fmt.Println("✅ Date/Time Reasoning")
	fmt.Println("✅ Clean Architecture")
	fmt.Println("✅ Gin Web Framework")
	fmt.Println("✅ Session Management")
	fmt.Println("✅ SQLite Persistence")
	fmt.Println("\nAPI Endpoints:")
	fmt.Println("Session Management:")
	fmt.Println("  POST /api/v1/sessions - Create a new session")
	fmt.Println("  GET  /api/v1/sessions - List all sessions")
	fmt.Println("  GET  /api/v1/sessions/:id - Get session details")
	fmt.Println("  DEL  /api/v1/sessions/:id - Delete a session")
	fmt.Println("\nProlog Operations (session-scoped):")
	fmt.Println("  POST /api/v1/sessions/:sessionId/facts - Add a fact")
	fmt.Println("  POST /api/v1/sessions/:sessionId/rules - Add a rule")  
	fmt.Println("  POST /api/v1/sessions/:sessionId/query - Execute a query")
	fmt.Println("\nUtilities:")
	fmt.Println("  POST /api/v1/cache/clear - Clear cache")
	
	// Show UI information if enabled
	if os.Getenv("ENABLE_UI") == "true" {
		fmt.Println("\nWeb UI:")
		fmt.Println("  GET  /ui - Interactive Prolog REPL (browser-based)")
		if os.Getenv("UI_PASSWORD") != "" {
			fmt.Println("       🔒 Password protected")
		}
	}
	// Get host and port from environment or use defaults
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	address := host + ":" + port
	fmt.Printf("\nListening on %s\n", address)

	if err := r.Run(address); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}