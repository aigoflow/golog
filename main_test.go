package main

import (
	"testing"
)

// setupTestEngine creates a new engine with an in-memory SQLite database for testing
func setupTestEngine(t *testing.T) *Engine {
	engine, err := NewEngine(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test engine: %v", err)
	}
	return engine
}

// teardownTestEngine cleans up the test engine
func teardownTestEngine(engine *Engine) {
	if engine != nil {
		engine.Close()
	}
}

// createTestSession creates a test session and returns its ID
func createTestSession(t *testing.T, engine *Engine) int {
	req := CreateSessionRequest{
		Name:        "test-session",
		Description: "Test session for unit tests",
	}
	session, err := engine.CreateSession(req)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	return session.ID
}