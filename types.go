package main

import "time"

type Term struct {
	Type  string      `json:"type"` // "atom", "variable", "compound", "list", "date", "number"
	Value interface{} `json:"value"`
	Args  []Term      `json:"args,omitempty"`
}

type Session struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSessionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type Fact struct {
	ID        int    `json:"id,omitempty"`
	SessionID string `json:"session_id"`
	Predicate Term   `json:"predicate"`
}

type Rule struct {
	ID        int    `json:"id,omitempty"`
	SessionID string `json:"session_id"`
	Head      Term   `json:"head"`
	Body      []Term `json:"body"`
}

type Query struct {
	Goals []Term `json:"goals"`
}

type Substitution map[string]Term

type Solution struct {
	Bindings Substitution `json:"bindings"`
	Success  bool         `json:"success"`
}

type QueryResult struct {
	Solutions []Solution `json:"solutions"`
}

type TableKey struct {
	Predicate string
	Args      string
}

type TableEntry struct {
	Solutions []Substitution
	Complete  bool
}

func Atom(value string) Term {
	return Term{Type: "atom", Value: value}
}

func Variable(name string) Term {
	return Term{Type: "variable", Value: name}
}

func Compound(functor string, args []Term) Term {
	return Term{Type: "compound", Value: functor, Args: args}
}

func Number(n float64) Term {
	return Term{Type: "number", Value: n}
}

func Date(t time.Time) Term {
	return Term{Type: "date", Value: t.Format(time.RFC3339)}
}