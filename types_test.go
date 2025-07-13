package main

import (
	"testing"
	"time"
)

func TestAtom(t *testing.T) {
	atom := Atom("test")
	if atom.Type != "atom" {
		t.Errorf("Expected type 'atom', got '%s'", atom.Type)
	}
	if atom.Value != "test" {
		t.Errorf("Expected value 'test', got '%v'", atom.Value)
	}
	if len(atom.Args) != 0 {
		t.Errorf("Expected no args, got %d", len(atom.Args))
	}
}

func TestVariable(t *testing.T) {
	variable := Variable("X")
	if variable.Type != "variable" {
		t.Errorf("Expected type 'variable', got '%s'", variable.Type)
	}
	if variable.Value != "X" {
		t.Errorf("Expected value 'X', got '%v'", variable.Value)
	}
	if len(variable.Args) != 0 {
		t.Errorf("Expected no args, got %d", len(variable.Args))
	}
}

func TestCompound(t *testing.T) {
	args := []Term{Atom("john"), Atom("mary")}
	compound := Compound("parent", args)
	
	if compound.Type != "compound" {
		t.Errorf("Expected type 'compound', got '%s'", compound.Type)
	}
	if compound.Value != "parent" {
		t.Errorf("Expected value 'parent', got '%v'", compound.Value)
	}
	if len(compound.Args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(compound.Args))
	}
	if compound.Args[0].Value != "john" {
		t.Errorf("Expected first arg 'john', got '%v'", compound.Args[0].Value)
	}
	if compound.Args[1].Value != "mary" {
		t.Errorf("Expected second arg 'mary', got '%v'", compound.Args[1].Value)
	}
}

func TestNumber(t *testing.T) {
	number := Number(42.5)
	if number.Type != "number" {
		t.Errorf("Expected type 'number', got '%s'", number.Type)
	}
	if number.Value != 42.5 {
		t.Errorf("Expected value 42.5, got '%v'", number.Value)
	}
	if len(number.Args) != 0 {
		t.Errorf("Expected no args, got %d", len(number.Args))
	}
}

func TestDate(t *testing.T) {
	now := time.Now()
	dateTerm := Date(now)
	
	if dateTerm.Type != "date" {
		t.Errorf("Expected type 'date', got '%s'", dateTerm.Type)
	}
	
	expectedValue := now.Format(time.RFC3339)
	if dateTerm.Value != expectedValue {
		t.Errorf("Expected value '%s', got '%v'", expectedValue, dateTerm.Value)
	}
	if len(dateTerm.Args) != 0 {
		t.Errorf("Expected no args, got %d", len(dateTerm.Args))
	}
}