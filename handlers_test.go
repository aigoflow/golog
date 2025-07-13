package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *Engine) {
	gin.SetMode(gin.TestMode)
	engine := setupTestEngine(t)
	router := engine.setupRoutes()
	return router, engine
}

func TestCreateSessionHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	req := CreateSessionRequest{
		Name:        "test-api-session",
		Description: "Testing API session creation",
	}

	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/sessions", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var session Session
	err := json.Unmarshal(w.Body.Bytes(), &session)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if session.Name != req.Name {
		t.Errorf("Expected session name '%s', got '%s'", req.Name, session.Name)
	}
	if session.Description != req.Description {
		t.Errorf("Expected session description '%s', got '%s'", req.Description, session.Description)
	}
}

func TestListSessionsHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	// Create a test session first
	createTestSession(t, engine)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/sessions", nil)

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]Session
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	sessions, exists := response["sessions"]
	if !exists {
		t.Error("Expected 'sessions' key in response")
	}

	if len(sessions) == 0 {
		t.Error("Expected at least one session")
	}
}

func TestGetSessionHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	sessionID := createTestSession(t, engine)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/sessions/"+strconv.Itoa(sessionID), nil)

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var session Session
	err := json.Unmarshal(w.Body.Bytes(), &session)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if session.ID != sessionID {
		t.Errorf("Expected session ID %d, got %d", sessionID, session.ID)
	}
}

func TestGetSessionHandlerNotFound(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/sessions/99999", nil)

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteSessionHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	sessionID := createTestSession(t, engine)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/v1/sessions/"+strconv.Itoa(sessionID), nil)

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify session is deleted
	_, err := engine.GetSession(sessionID)
	if err == nil {
		t.Error("Expected error when getting deleted session")
	}
}

func TestAddFactHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	sessionID := createTestSession(t, engine)

	fact := Fact{
		Predicate: Compound("parent", []Term{Atom("john"), Atom("mary")}),
	}

	jsonData, _ := json.Marshal(fact)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/sessions/"+strconv.Itoa(sessionID)+"/facts", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify fact was added by querying
	goal := Compound("parent", []Term{Variable("X"), Variable("Y")})
	facts := engine.loadFacts(goal, sessionID)
	if len(facts) != 1 {
		t.Errorf("Expected 1 fact to be stored, got %d", len(facts))
	}
}

func TestAddRuleHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	sessionID := createTestSession(t, engine)

	rule := Rule{
		Head: Compound("grandparent", []Term{Variable("X"), Variable("Z")}),
		Body: []Term{
			Compound("parent", []Term{Variable("X"), Variable("Y")}),
			Compound("parent", []Term{Variable("Y"), Variable("Z")}),
		},
	}

	jsonData, _ := json.Marshal(rule)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/sessions/"+strconv.Itoa(sessionID)+"/rules", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify rule was added by querying
	goal := Compound("grandparent", []Term{Variable("X"), Variable("Y")})
	rules := engine.loadRules(goal, sessionID)
	if len(rules) != 1 {
		t.Errorf("Expected 1 rule to be stored, got %d", len(rules))
	}
}

func TestQueryHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	sessionID := createTestSession(t, engine)

	// Add a fact first
	fact := Fact{
		SessionID: sessionID,
		Predicate: Compound("parent", []Term{Atom("john"), Atom("mary")}),
	}
	err := engine.AddFact(fact)
	if err != nil {
		t.Fatalf("Failed to add fact: %v", err)
	}

	query := Query{
		Goals: []Term{
			Compound("parent", []Term{Variable("X"), Atom("mary")}),
		},
	}

	jsonData, _ := json.Marshal(query)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/sessions/"+strconv.Itoa(sessionID)+"/query", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result QueryResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result.Solutions) != 1 {
		t.Errorf("Expected 1 solution, got %d", len(result.Solutions))
	}

	if !result.Solutions[0].Success {
		t.Error("Expected successful solution")
	}

	if result.Solutions[0].Bindings["X"].Value != "john" {
		t.Errorf("Expected X to be bound to 'john', got '%v'", result.Solutions[0].Bindings["X"].Value)
	}
}

func TestClearCacheHandler(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/cache/clear", nil)

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "cache cleared" {
		t.Errorf("Expected 'cache cleared' status, got '%s'", response["status"])
	}
}

func TestInvalidSessionID(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	// Test with invalid session ID for facts
	fact := Fact{
		Predicate: Compound("parent", []Term{Atom("john"), Atom("mary")}),
	}

	jsonData, _ := json.Marshal(fact)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/sessions/invalid/facts", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMalformedJSON(t *testing.T) {
	router, engine := setupTestRouter(t)
	defer teardownTestEngine(engine)

	sessionID := createTestSession(t, engine)

	// Send malformed JSON
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/sessions/"+strconv.Itoa(sessionID)+"/facts", bytes.NewBuffer([]byte("{invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}