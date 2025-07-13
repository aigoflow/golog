package main

import (
	"testing"
	"time"
)

func TestBuiltinUnification(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Test =(X, test)
	goal := Compound("=", []Term{Variable("X"), Atom("test")})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected = predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution, got %d", len(solutions))
	}
	if solutions[0]["X"].Value != "test" {
		t.Errorf("Expected X to be bound to 'test', got '%v'", solutions[0]["X"].Value)
	}

	// Test =(test, different) - should fail
	goal = Compound("=", []Term{Atom("test"), Atom("different")})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected = predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for failed unification, got %d", len(solutions))
	}
}

func TestBuiltinTypeChecking(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Test atom(test)
	goal := Compound("atom", []Term{Atom("test")})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected atom predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for atom check, got %d", len(solutions))
	}

	// Test atom(42) - should fail
	goal = Compound("atom", []Term{Number(42)})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected atom predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for failed atom check, got %d", len(solutions))
	}

	// Test var(X)
	goal = Compound("var", []Term{Variable("X")})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected var predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for var check, got %d", len(solutions))
	}

	// Test var(test) - should fail
	goal = Compound("var", []Term{Atom("test")})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected var predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for failed var check, got %d", len(solutions))
	}

	// Test number(42)
	goal = Compound("number", []Term{Number(42)})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected number predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for number check, got %d", len(solutions))
	}

	// Test number(test) - should fail
	goal = Compound("number", []Term{Atom("test")})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected number predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for failed number check, got %d", len(solutions))
	}
}

func TestBuiltinNow(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Test now(X)
	goal := Compound("now", []Term{Variable("X")})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected now predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for now, got %d", len(solutions))
	}
	
	if solutions[0]["X"].Type != "date" {
		t.Errorf("Expected X to be bound to date, got type '%s'", solutions[0]["X"].Type)
	}

	// Verify the date is valid RFC3339
	dateStr, ok := solutions[0]["X"].Value.(string)
	if !ok {
		t.Error("Expected date value to be string")
	}
	
	_, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t.Errorf("Expected valid RFC3339 date, got error: %v", err)
	}
}

func TestBuiltinDateComparison(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Create test dates
	earlier := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	later := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	
	earlierTerm := Date(earlier)
	laterTerm := Date(later)

	// Test date_before(earlier, later)
	goal := Compound("date_before", []Term{earlierTerm, laterTerm})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected date_before predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for date_before, got %d", len(solutions))
	}

	// Test date_before(later, earlier) - should fail
	goal = Compound("date_before", []Term{laterTerm, earlierTerm})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected date_before predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for failed date_before, got %d", len(solutions))
	}

	// Test date_after(later, earlier)
	goal = Compound("date_after", []Term{laterTerm, earlierTerm})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected date_after predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for date_after, got %d", len(solutions))
	}

	// Test days_between(earlier, later, X)
	goal = Compound("days_between", []Term{earlierTerm, laterTerm, Variable("X")})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected days_between predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for days_between, got %d", len(solutions))
	}
	
	if solutions[0]["X"].Type != "number" {
		t.Errorf("Expected X to be bound to number, got type '%s'", solutions[0]["X"].Type)
	}
	
	days, ok := solutions[0]["X"].Value.(float64)
	if !ok {
		t.Error("Expected days to be float64")
	}
	
	// Should be 364 days (2023 is not a leap year)
	expectedDays := 364.0
	if days != expectedDays {
		t.Errorf("Expected %f days, got %f", expectedDays, days)
	}
}

func TestUnhandledBuiltin(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Test non-existent builtin
	goal := Compound("nonexistent", []Term{Atom("test")})
	_, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if handled {
		t.Error("Expected non-existent predicate to not be handled")
	}

	// Test atom goal (not compound)
	goal = Atom("test")
	_, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if handled {
		t.Error("Expected atom goal to not be handled as builtin")
	}
}