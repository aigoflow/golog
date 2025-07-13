package main

import (
	"database/sql"
	"testing"
	"time"
)

func TestNewEngine(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	if engine.db == nil {
		t.Error("Expected database connection to be initialized")
	}

	if engine.cache == nil {
		t.Error("Expected cache to be initialized")
	}
}

func TestDatabaseSchema(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	// Test sessions table exists
	var count int
	err := engine.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sessions'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query sessions table: %v", err)
	}
	if count != 1 {
		t.Error("Sessions table was not created")
	}

	// Test facts table exists
	err = engine.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='facts'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query facts table: %v", err)
	}
	if count != 1 {
		t.Error("Facts table was not created")
	}

	// Test rules table exists
	err = engine.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='rules'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query rules table: %v", err)
	}
	if count != 1 {
		t.Error("Rules table was not created")
	}
}

func TestSessionCRUD(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	// Test Create Session
	req := CreateSessionRequest{
		Name:        "test-session-crud",
		Description: "Testing CRUD operations",
	}
	session, err := engine.CreateSession(req)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == 0 {
		t.Error("Expected session ID to be set")
	}
	if session.Name != req.Name {
		t.Errorf("Expected session name '%s', got '%s'", req.Name, session.Name)
	}
	if session.Description != req.Description {
		t.Errorf("Expected session description '%s', got '%s'", req.Description, session.Description)
	}

	// Test Get Session
	retrieved, err := engine.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}
	if retrieved.ID != session.ID {
		t.Errorf("Expected session ID %d, got %d", session.ID, retrieved.ID)
	}
	if retrieved.Name != session.Name {
		t.Errorf("Expected session name '%s', got '%s'", session.Name, retrieved.Name)
	}

	// Test Get Session By Name
	retrievedByName, err := engine.GetSessionByName(session.Name)
	if err != nil {
		t.Fatalf("Failed to get session by name: %v", err)
	}
	if retrievedByName.ID != session.ID {
		t.Errorf("Expected session ID %d, got %d", session.ID, retrievedByName.ID)
	}

	// Test List Sessions
	sessions, err := engine.ListSessions()
	if err != nil {
		t.Fatalf("Failed to list sessions: %v", err)
	}
	if len(sessions) == 0 {
		t.Error("Expected at least one session in list")
	}

	found := false
	for _, s := range sessions {
		if s.ID == session.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created session not found in session list")
	}

	// Test Update Session Timestamp
	originalTime := session.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	err = engine.UpdateSessionTimestamp(session.ID)
	if err != nil {
		t.Fatalf("Failed to update session timestamp: %v", err)
	}

	updated, err := engine.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}
	if !updated.UpdatedAt.After(originalTime) {
		t.Error("Expected updated timestamp to be later than original")
	}

	// Test Delete Session
	err = engine.DeleteSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify session is deleted
	_, err = engine.GetSession(session.ID)
	if err == nil {
		t.Error("Expected error when getting deleted session")
	}
	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestSessionUniqueName(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	req := CreateSessionRequest{
		Name:        "unique-session",
		Description: "Testing unique constraint",
	}

	// Create first session
	_, err := engine.CreateSession(req)
	if err != nil {
		t.Fatalf("Failed to create first session: %v", err)
	}

	// Try to create second session with same name
	_, err = engine.CreateSession(req)
	if err == nil {
		t.Error("Expected error when creating session with duplicate name")
	}
}