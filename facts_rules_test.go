package main

import (
	"testing"
)

func TestAddAndLoadFacts(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	// Add a fact: parent(john, mary)
	fact := Fact{
		SessionID: sessionID,
		Predicate: Compound("parent", []Term{Atom("john"), Atom("mary")}),
	}

	err := engine.AddFact(fact)
	if err != nil {
		t.Fatalf("Failed to add fact: %v", err)
	}

	// Load facts matching parent(X, Y)
	goal := Compound("parent", []Term{Variable("X"), Variable("Y")})
	facts := engine.loadFacts(goal, sessionID)

	if len(facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(facts))
	}

	if facts[0].SessionID != sessionID {
		t.Errorf("Expected fact session ID %d, got %d", sessionID, facts[0].SessionID)
	}

	// Verify the loaded fact matches what we stored
	loadedPredicate := facts[0].Predicate
	if loadedPredicate.Type != "compound" || loadedPredicate.Value != "parent" {
		t.Errorf("Expected compound predicate 'parent', got %v", loadedPredicate)
	}

	if len(loadedPredicate.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(loadedPredicate.Args))
	}

	if loadedPredicate.Args[0].Value != "john" || loadedPredicate.Args[1].Value != "mary" {
		t.Errorf("Expected arguments john, mary, got %v, %v", 
			loadedPredicate.Args[0].Value, loadedPredicate.Args[1].Value)
	}
}

func TestAddAndLoadRules(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	// Add a rule: grandparent(X, Z) :- parent(X, Y), parent(Y, Z)
	rule := Rule{
		SessionID: sessionID,
		Head:      Compound("grandparent", []Term{Variable("X"), Variable("Z")}),
		Body: []Term{
			Compound("parent", []Term{Variable("X"), Variable("Y")}),
			Compound("parent", []Term{Variable("Y"), Variable("Z")}),
		},
	}

	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Load rules matching grandparent(X, Y)
	goal := Compound("grandparent", []Term{Variable("X"), Variable("Y")})
	rules := engine.loadRules(goal, sessionID)

	if len(rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(rules))
	}

	if rules[0].SessionID != sessionID {
		t.Errorf("Expected rule session ID %d, got %d", sessionID, rules[0].SessionID)
	}

	// Verify the loaded rule
	loadedRule := rules[0]
	if loadedRule.Head.Type != "compound" || loadedRule.Head.Value != "grandparent" {
		t.Errorf("Expected compound head 'grandparent', got %v", loadedRule.Head)
	}

	if len(loadedRule.Body) != 2 {
		t.Errorf("Expected 2 body goals, got %d", len(loadedRule.Body))
	}

	for i, bodyGoal := range loadedRule.Body {
		if bodyGoal.Type != "compound" || bodyGoal.Value != "parent" {
			t.Errorf("Expected body goal %d to be 'parent', got %v", i, bodyGoal)
		}
	}
}

func TestSessionIsolation(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	// Create two different sessions
	session1ID := createTestSession(t, engine)
	
	req2 := CreateSessionRequest{
		Name:        "test-session-2",
		Description: "Second test session",
	}
	session2, err := engine.CreateSession(req2)
	if err != nil {
		t.Fatalf("Failed to create second session: %v", err)
	}
	session2ID := session2.ID

	// Add facts to each session
	fact1 := Fact{
		SessionID: session1ID,
		Predicate: Compound("parent", []Term{Atom("john"), Atom("mary")}),
	}
	fact2 := Fact{
		SessionID: session2ID,
		Predicate: Compound("parent", []Term{Atom("bob"), Atom("alice")}),
	}

	err = engine.AddFact(fact1)
	if err != nil {
		t.Fatalf("Failed to add fact to session 1: %v", err)
	}

	err = engine.AddFact(fact2)
	if err != nil {
		t.Fatalf("Failed to add fact to session 2: %v", err)
	}

	// Query each session
	goal := Compound("parent", []Term{Variable("X"), Variable("Y")})

	facts1 := engine.loadFacts(goal, session1ID)
	facts2 := engine.loadFacts(goal, session2ID)

	// Each session should only see its own facts
	if len(facts1) != 1 {
		t.Errorf("Expected session 1 to have 1 fact, got %d", len(facts1))
	}
	if len(facts2) != 1 {
		t.Errorf("Expected session 2 to have 1 fact, got %d", len(facts2))
	}

	// Verify the facts are different
	if facts1[0].Predicate.Args[0].Value != "john" {
		t.Errorf("Expected session 1 fact to have 'john', got %v", facts1[0].Predicate.Args[0].Value)
	}
	if facts2[0].Predicate.Args[0].Value != "bob" {
		t.Errorf("Expected session 2 fact to have 'bob', got %v", facts2[0].Predicate.Args[0].Value)
	}
}

func TestExtractPredicate(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	// Test compound term
	compound := Compound("parent", []Term{Atom("john"), Atom("mary")})
	predicate := engine.extractPredicate(compound)
	if predicate != "parent" {
		t.Errorf("Expected predicate 'parent', got '%s'", predicate)
	}

	// Test atom term
	atom := Atom("test")
	predicate = engine.extractPredicate(atom)
	if predicate != "test" {
		t.Errorf("Expected predicate 'test', got '%s'", predicate)
	}

	// Test variable term (should return empty)
	variable := Variable("X")
	predicate = engine.extractPredicate(variable)
	if predicate != "" {
		t.Errorf("Expected empty predicate for variable, got '%s'", predicate)
	}

	// Test number term (should return empty)
	number := Number(42)
	predicate = engine.extractPredicate(number)
	if predicate != "" {
		t.Errorf("Expected empty predicate for number, got '%s'", predicate)
	}
}

func TestQueryExecution(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)
	sessionID := createTestSession(t, engine)

	// Add some facts
	facts := []Fact{
		{SessionID: sessionID, Predicate: Compound("parent", []Term{Atom("john"), Atom("mary")})},
		{SessionID: sessionID, Predicate: Compound("parent", []Term{Atom("mary"), Atom("bob")})},
		{SessionID: sessionID, Predicate: Compound("parent", []Term{Atom("alice"), Atom("carol")})},
	}

	for _, fact := range facts {
		err := engine.AddFact(fact)
		if err != nil {
			t.Fatalf("Failed to add fact: %v", err)
		}
	}

	// Query: parent(X, mary)
	query := Query{
		Goals: []Term{
			Compound("parent", []Term{Variable("X"), Atom("mary")}),
		},
	}

	result := engine.Query(query, sessionID)

	// Should find one solution: X = john
	if len(result.Solutions) != 1 {
		t.Errorf("Expected 1 solution, got %d", len(result.Solutions))
	}

	solution := result.Solutions[0]
	if !solution.Success {
		t.Error("Expected successful solution")
	}

	if solution.Bindings["X"].Value != "john" {
		t.Errorf("Expected X to be bound to 'john', got '%v'", solution.Bindings["X"].Value)
	}

	// Query: parent(X, Y) - should find all facts
	query = Query{
		Goals: []Term{
			Compound("parent", []Term{Variable("X"), Variable("Y")}),
		},
	}

	result = engine.Query(query, sessionID)

	// Should find three solutions
	if len(result.Solutions) != 3 {
		t.Errorf("Expected 3 solutions, got %d", len(result.Solutions))
	}

	for _, solution := range result.Solutions {
		if !solution.Success {
			t.Error("Expected all solutions to be successful")
		}
	}

	// Query: parent(nonexistent, Y) - should find no solutions
	query = Query{
		Goals: []Term{
			Compound("parent", []Term{Atom("nonexistent"), Variable("Y")}),
		},
	}

	result = engine.Query(query, sessionID)

	// Should find one unsuccessful solution
	if len(result.Solutions) != 1 {
		t.Errorf("Expected 1 solution, got %d", len(result.Solutions))
	}

	if result.Solutions[0].Success {
		t.Error("Expected unsuccessful solution for non-matching query")
	}
}