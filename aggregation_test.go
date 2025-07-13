package main

import (
	"testing"
)

func setupAggregationTestData(t *testing.T, engine *Engine, sessionID int) {
	// Add some test facts: score(john, 85), score(mary, 92), score(bob, 78)
	facts := []Fact{
		{SessionID: sessionID, Predicate: Compound("score", []Term{Atom("john"), Number(85)})},
		{SessionID: sessionID, Predicate: Compound("score", []Term{Atom("mary"), Number(92)})},
		{SessionID: sessionID, Predicate: Compound("score", []Term{Atom("bob"), Number(78)})},
	}

	for _, fact := range facts {
		err := engine.AddFact(fact)
		if err != nil {
			t.Fatalf("Failed to add test fact: %v", err)
		}
	}
}

func TestBuiltinCount(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)
	setupAggregationTestData(t, engine, sessionID)

	subst := make(Substitution)

	// Test count(_, score(X, Y), Count)
	template := Variable("_")
	queryGoal := Compound("score", []Term{Variable("X"), Variable("Y")})
	countVar := Variable("Count")
	
	goal := Compound("count", []Term{template, queryGoal, countVar})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected count predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for count, got %d", len(solutions))
	}
	
	if solutions[0]["Count"].Type != "number" {
		t.Errorf("Expected Count to be number, got type '%s'", solutions[0]["Count"].Type)
	}
	
	count, ok := solutions[0]["Count"].Value.(float64)
	if !ok {
		t.Error("Expected count value to be float64")
	}
	
	if count != 3.0 {
		t.Errorf("Expected count of 3, got %f", count)
	}
}

func TestBuiltinSum(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)
	setupAggregationTestData(t, engine, sessionID)

	subst := make(Substitution)

	// Test sum(Y, score(X, Y), Total)
	template := Variable("Y")
	queryGoal := Compound("score", []Term{Variable("X"), Variable("Y")})
	sumVar := Variable("Total")
	
	goal := Compound("sum", []Term{template, queryGoal, sumVar})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected sum predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for sum, got %d", len(solutions))
	}
	
	if solutions[0]["Total"].Type != "number" {
		t.Errorf("Expected Total to be number, got type '%s'", solutions[0]["Total"].Type)
	}
	
	total, ok := solutions[0]["Total"].Value.(float64)
	if !ok {
		t.Error("Expected sum value to be float64")
	}
	
	expectedSum := 85.0 + 92.0 + 78.0 // 255
	if total != expectedSum {
		t.Errorf("Expected sum of %f, got %f", expectedSum, total)
	}
}

func TestBuiltinMax(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)
	setupAggregationTestData(t, engine, sessionID)

	subst := make(Substitution)

	// Test max(Y, score(X, Y), Maximum)
	template := Variable("Y")
	queryGoal := Compound("score", []Term{Variable("X"), Variable("Y")})
	maxVar := Variable("Maximum")
	
	goal := Compound("max", []Term{template, queryGoal, maxVar})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected max predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for max, got %d", len(solutions))
	}
	
	if solutions[0]["Maximum"].Type != "number" {
		t.Errorf("Expected Maximum to be number, got type '%s'", solutions[0]["Maximum"].Type)
	}
	
	max, ok := solutions[0]["Maximum"].Value.(float64)
	if !ok {
		t.Error("Expected max value to be float64")
	}
	
	if max != 92.0 {
		t.Errorf("Expected max of 92, got %f", max)
	}
}

func TestBuiltinMin(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)
	setupAggregationTestData(t, engine, sessionID)

	subst := make(Substitution)

	// Test min(Y, score(X, Y), Minimum)
	template := Variable("Y")
	queryGoal := Compound("score", []Term{Variable("X"), Variable("Y")})
	minVar := Variable("Minimum")
	
	goal := Compound("min", []Term{template, queryGoal, minVar})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected min predicate to be handled")
	}
	if len(solutions) != 1 {
		t.Errorf("Expected 1 solution for min, got %d", len(solutions))
	}
	
	if solutions[0]["Minimum"].Type != "number" {
		t.Errorf("Expected Minimum to be number, got type '%s'", solutions[0]["Minimum"].Type)
	}
	
	min, ok := solutions[0]["Minimum"].Value.(float64)
	if !ok {
		t.Error("Expected min value to be float64")
	}
	
	if min != 78.0 {
		t.Errorf("Expected min of 78, got %f", min)
	}
}

func TestAggregationWithNoSolutions(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Test max with no matching facts
	template := Variable("Y")
	queryGoal := Compound("nonexistent", []Term{Variable("X"), Variable("Y")})
	maxVar := Variable("Maximum")
	
	goal := Compound("max", []Term{template, queryGoal, maxVar})
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected max predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for max with no data, got %d", len(solutions))
	}

	// Test min with no matching facts
	goal = Compound("min", []Term{template, queryGoal, Variable("Minimum")})
	solutions, handled = engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected min predicate to be handled")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for min with no data, got %d", len(solutions))
	}
}

func TestAggregationWrongArity(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	subst := make(Substitution)

	// Test count with wrong number of arguments
	goal := Compound("count", []Term{Variable("X")}) // Should have 3 args
	solutions, handled := engine.evalBuiltin(goal, subst, sessionID)
	
	if !handled {
		t.Error("Expected count predicate to be handled even with wrong arity")
	}
	if len(solutions) != 0 {
		t.Errorf("Expected 0 solutions for count with wrong arity, got %d", len(solutions))
	}
}