package main

import (
	"testing"
)

func TestUnifyAtoms(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	atom1 := Atom("test")
	atom2 := Atom("test")
	atom3 := Atom("different")

	subst := make(Substitution)

	// Test unifying identical atoms
	result, ok := engine.unify(atom1, atom2, subst)
	if !ok {
		t.Error("Expected identical atoms to unify successfully")
	}
	if len(result) != 0 {
		t.Error("Expected no new bindings when unifying identical atoms")
	}

	// Test unifying different atoms
	_, ok = engine.unify(atom1, atom3, subst)
	if ok {
		t.Error("Expected different atoms to fail unification")
	}
}

func TestUnifyVariables(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	var1 := Variable("X")
	var2 := Variable("Y")
	atom := Atom("test")

	subst := make(Substitution)

	// Test unifying variable with atom
	result, ok := engine.unify(var1, atom, subst)
	if !ok {
		t.Error("Expected variable to unify with atom")
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 binding, got %d", len(result))
	}
	if result["X"].Value != "test" {
		t.Errorf("Expected X to be bound to 'test', got '%v'", result["X"].Value)
	}

	// Test unifying variable with variable
	result, ok = engine.unify(var1, var2, subst)
	if !ok {
		t.Error("Expected variables to unify")
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 binding, got %d", len(result))
	}
}

func TestUnifyCompounds(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	// parent(john, mary)
	compound1 := Compound("parent", []Term{Atom("john"), Atom("mary")})
	
	// parent(john, mary) - identical
	compound2 := Compound("parent", []Term{Atom("john"), Atom("mary")})
	
	// parent(X, mary) - with variable
	compound3 := Compound("parent", []Term{Variable("X"), Atom("mary")})
	
	// loves(john, mary) - different functor
	compound4 := Compound("loves", []Term{Atom("john"), Atom("mary")})
	
	// parent(john, mary, bob) - different arity
	compound5 := Compound("parent", []Term{Atom("john"), Atom("mary"), Atom("bob")})

	subst := make(Substitution)

	// Test unifying identical compounds
	result, ok := engine.unify(compound1, compound2, subst)
	if !ok {
		t.Error("Expected identical compounds to unify")
	}
	if len(result) != 0 {
		t.Error("Expected no new bindings for identical compounds")
	}

	// Test unifying compound with variable
	result, ok = engine.unify(compound1, compound3, subst)
	if !ok {
		t.Error("Expected compound to unify with compound containing variable")
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 binding, got %d", len(result))
	}
	if result["X"].Value != "john" {
		t.Errorf("Expected X to be bound to 'john', got '%v'", result["X"].Value)
	}

	// Test unifying compounds with different functors
	_, ok = engine.unify(compound1, compound4, subst)
	if ok {
		t.Error("Expected compounds with different functors to fail unification")
	}

	// Test unifying compounds with different arity
	_, ok = engine.unify(compound1, compound5, subst)
	if ok {
		t.Error("Expected compounds with different arity to fail unification")
	}
}

func TestUnifyNumbers(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	num1 := Number(42)
	num2 := Number(42)
	num3 := Number(43)

	subst := make(Substitution)

	// Test unifying identical numbers
	result, ok := engine.unify(num1, num2, subst)
	if !ok {
		t.Error("Expected identical numbers to unify")
	}
	if len(result) != 0 {
		t.Error("Expected no new bindings for identical numbers")
	}

	// Test unifying different numbers
	_, ok = engine.unify(num1, num3, subst)
	if ok {
		t.Error("Expected different numbers to fail unification")
	}
}

func TestOccursCheck(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	// Test X = f(X) - should fail occurs check
	varX := Variable("X")
	compound := Compound("f", []Term{varX})

	subst := make(Substitution)

	// This should fail due to occurs check
	_, ok := engine.unify(varX, compound, subst)
	if ok {
		t.Error("Expected occurs check to prevent infinite structure")
	}
}

func TestDeref(t *testing.T) {
	engine := setupTestEngine(t)
	defer teardownTestEngine(engine)

	varX := Variable("X")
	atom := Atom("test")

	subst := make(Substitution)
	subst["X"] = atom

	// Test dereferencing bound variable
	result := engine.deref(varX, subst)
	if result.Type != "atom" || result.Value != "test" {
		t.Errorf("Expected dereferenced variable to be atom 'test', got %v", result)
	}

	// Test dereferencing unbound variable
	varY := Variable("Y")
	result = engine.deref(varY, subst)
	if result.Type != "variable" || result.Value != "Y" {
		t.Errorf("Expected unbound variable to remain variable, got %v", result)
	}

	// Test dereferencing non-variable
	result = engine.deref(atom, subst)
	if result.Type != "atom" || result.Value != "test" {
		t.Errorf("Expected non-variable to remain unchanged, got %v", result)
	}
}