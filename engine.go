package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/oklog/ulid/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Engine struct {
	db    *sql.DB
	cache map[TableKey]TableEntry
}

func NewEngine(dbPath string) (*Engine, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createSchema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS facts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		predicate TEXT NOT NULL,
		data TEXT NOT NULL,
		FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS rules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		head_predicate TEXT NOT NULL,
		head_data TEXT NOT NULL,
		body_data TEXT NOT NULL,
		FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_fact_pred ON facts(predicate);
	CREATE INDEX IF NOT EXISTS idx_rule_pred ON rules(head_predicate);
	CREATE INDEX IF NOT EXISTS idx_fact_session ON facts(session_id);
	CREATE INDEX IF NOT EXISTS idx_rule_session ON rules(session_id);
	`

	if _, err = db.Exec(createSchema); err != nil {
		return nil, err
	}

	return &Engine{
		db:    db,
		cache: make(map[TableKey]TableEntry),
	}, nil
}

func (e *Engine) unify(t1, t2 Term, subst Substitution) (Substitution, bool) {
	t1 = e.deref(t1, subst)
	t2 = e.deref(t2, subst)

	if reflect.DeepEqual(t1, t2) {
		return subst, true
	}

	if t1.Type == "variable" {
		return e.bind(t1.Value.(string), t2, subst)
	}
	if t2.Type == "variable" {
		return e.bind(t2.Value.(string), t1, subst)
	}

	if t1.Type == "compound" && t2.Type == "compound" {
		if t1.Value != t2.Value || len(t1.Args) != len(t2.Args) {
			return subst, false
		}
		for i := range t1.Args {
			var ok bool
			if subst, ok = e.unify(t1.Args[i], t2.Args[i], subst); !ok {
				return subst, false
			}
		}
		return subst, true
	}

	return subst, false
}

func (e *Engine) deref(term Term, subst Substitution) Term {
	if term.Type == "variable" {
		if binding, exists := subst[term.Value.(string)]; exists {
			return e.deref(binding, subst)
		}
	}
	return term
}

func (e *Engine) bind(varName string, term Term, subst Substitution) (Substitution, bool) {
	if e.occursCheck(varName, term, subst) {
		return subst, false
	}

	newSubst := make(Substitution)
	for k, v := range subst {
		newSubst[k] = v
	}
	newSubst[varName] = term
	return newSubst, true
}

func (e *Engine) occursCheck(varName string, term Term, subst Substitution) bool {
	term = e.deref(term, subst)

	if term.Type == "variable" && term.Value == varName {
		return true
	}

	if term.Type == "compound" {
		for _, arg := range term.Args {
			if e.occursCheck(varName, arg, subst) {
				return true
			}
		}
	}

	return false
}

func (e *Engine) evalBuiltin(goal Term, subst Substitution, sessionID string) ([]Substitution, bool) {
	if goal.Type != "compound" {
		return nil, false
	}

	switch goal.Value {
	case "=":
		if len(goal.Args) == 2 {
			if newSubst, ok := e.unify(goal.Args[0], goal.Args[1], subst); ok {
				return []Substitution{newSubst}, true
			}
		}
		return []Substitution{}, true

	case "atom":
		if len(goal.Args) == 1 {
			arg := e.deref(goal.Args[0], subst)
			if arg.Type == "atom" {
				return []Substitution{subst}, true
			}
		}
		return []Substitution{}, true
	case "var":
		if len(goal.Args) == 1 {
			arg := e.deref(goal.Args[0], subst)
			if arg.Type == "variable" {
				return []Substitution{subst}, true
			}
		}
		return []Substitution{}, true
	case "number":
		if len(goal.Args) == 1 {
			arg := e.deref(goal.Args[0], subst)
			if arg.Type == "number" {
				return []Substitution{subst}, true
			}
		}
		return []Substitution{}, true

	case "count":
		return e.handleCount(goal, subst, sessionID)
	case "sum":
		return e.handleSum(goal, subst, sessionID)
	case "max":
		return e.handleMax(goal, subst, sessionID)
	case "min":
		return e.handleMin(goal, subst, sessionID)

	case "now":
		if len(goal.Args) == 1 {
			now := Date(time.Now())
			if newSubst, ok := e.unify(goal.Args[0], now, subst); ok {
				return []Substitution{newSubst}, true
			}
		}
		return []Substitution{}, true
	case "date_before":
		return e.handleDateBefore(goal, subst)
	case "date_after":
		return e.handleDateAfter(goal, subst)
	case "days_between":
		return e.handleDaysBetween(goal, subst)
	case "help":
		// Help predicate always succeeds (used by UI for command detection)
		return []Substitution{subst}, true
	}

	return nil, false
}

func (e *Engine) handleCount(goal Term, subst Substitution, sessionID string) ([]Substitution, bool) {
	if len(goal.Args) != 3 {
		return []Substitution{}, true
	}

	queryGoal := goal.Args[1]
	countVar := goal.Args[2]

	solutions := e.solve([]Term{queryGoal}, subst, sessionID)
	count := Number(float64(len(solutions)))

	if newSubst, ok := e.unify(countVar, count, subst); ok {
		return []Substitution{newSubst}, true
	}

	return []Substitution{}, true
}

func (e *Engine) handleSum(goal Term, subst Substitution, sessionID string) ([]Substitution, bool) {
	if len(goal.Args) != 3 {
		return []Substitution{}, true
	}

	template := goal.Args[0]
	queryGoal := goal.Args[1]
	sumVar := goal.Args[2]

	solutions := e.solve([]Term{queryGoal}, subst, sessionID)
	var total float64

	for _, sol := range solutions {
		instantiated := e.instantiate(template, sol)
		if instantiated.Type == "number" {
			if val, ok := instantiated.Value.(float64); ok {
				total += val
			}
		}
	}

	sum := Number(total)
	if newSubst, ok := e.unify(sumVar, sum, subst); ok {
		return []Substitution{newSubst}, true
	}

	return []Substitution{}, true
}

func (e *Engine) handleMax(goal Term, subst Substitution, sessionID string) ([]Substitution, bool) {
	if len(goal.Args) != 3 {
		return []Substitution{}, true
	}

	template := goal.Args[0]
	queryGoal := goal.Args[1]
	maxVar := goal.Args[2]

	solutions := e.solve([]Term{queryGoal}, subst, sessionID)
	if len(solutions) == 0 {
		return []Substitution{}, true
	}

	var max float64
	first := true

	for _, sol := range solutions {
		instantiated := e.instantiate(template, sol)
		if instantiated.Type == "number" {
			if val, ok := instantiated.Value.(float64); ok {
				if first || val > max {
					max = val
					first = false
				}
			}
		}
	}

	if !first {
		maxTerm := Number(max)
		if newSubst, ok := e.unify(maxVar, maxTerm, subst); ok {
			return []Substitution{newSubst}, true
		}
	}

	return []Substitution{}, true
}

func (e *Engine) handleMin(goal Term, subst Substitution, sessionID string) ([]Substitution, bool) {
	if len(goal.Args) != 3 {
		return []Substitution{}, true
	}

	template := goal.Args[0]
	queryGoal := goal.Args[1]
	minVar := goal.Args[2]

	solutions := e.solve([]Term{queryGoal}, subst, sessionID)
	if len(solutions) == 0 {
		return []Substitution{}, true
	}

	var min float64
	first := true

	for _, sol := range solutions {
		instantiated := e.instantiate(template, sol)
		if instantiated.Type == "number" {
			if val, ok := instantiated.Value.(float64); ok {
				if first || val < min {
					min = val
					first = false
				}
			}
		}
	}

	if !first {
		minTerm := Number(min)
		if newSubst, ok := e.unify(minVar, minTerm, subst); ok {
			return []Substitution{newSubst}, true
		}
	}

	return []Substitution{}, true
}

func (e *Engine) handleDateBefore(goal Term, subst Substitution) ([]Substitution, bool) {
	if len(goal.Args) != 2 {
		return []Substitution{}, true
	}

	date1 := e.deref(goal.Args[0], subst)
	date2 := e.deref(goal.Args[1], subst)

	if date1.Type == "date" && date2.Type == "date" {
		t1, err1 := time.Parse(time.RFC3339, date1.Value.(string))
		t2, err2 := time.Parse(time.RFC3339, date2.Value.(string))

		if err1 == nil && err2 == nil && t1.Before(t2) {
			return []Substitution{subst}, true
		}
	}

	return []Substitution{}, true
}

func (e *Engine) handleDateAfter(goal Term, subst Substitution) ([]Substitution, bool) {
	if len(goal.Args) != 2 {
		return []Substitution{}, true
	}

	date1 := e.deref(goal.Args[0], subst)
	date2 := e.deref(goal.Args[1], subst)

	if date1.Type == "date" && date2.Type == "date" {
		t1, err1 := time.Parse(time.RFC3339, date1.Value.(string))
		t2, err2 := time.Parse(time.RFC3339, date2.Value.(string))

		if err1 == nil && err2 == nil && t1.After(t2) {
			return []Substitution{subst}, true
		}
	}

	return []Substitution{}, true
}

func (e *Engine) handleDaysBetween(goal Term, subst Substitution) ([]Substitution, bool) {
	if len(goal.Args) != 3 {
		return []Substitution{}, true
	}

	date1 := e.deref(goal.Args[0], subst)
	date2 := e.deref(goal.Args[1], subst)

	if date1.Type == "date" && date2.Type == "date" {
		t1, err1 := time.Parse(time.RFC3339, date1.Value.(string))
		t2, err2 := time.Parse(time.RFC3339, date2.Value.(string))

		if err1 == nil && err2 == nil {
			days := Number(t2.Sub(t1).Hours() / 24)
			if newSubst, ok := e.unify(goal.Args[2], days, subst); ok {
				return []Substitution{newSubst}, true
			}
		}
	}

	return []Substitution{}, true
}

func (e *Engine) instantiate(term Term, subst Substitution) Term {
	term = e.deref(term, subst)

	if term.Type == "compound" {
		newArgs := make([]Term, len(term.Args))
		for i, arg := range term.Args {
			newArgs[i] = e.instantiate(arg, subst)
		}
		return Term{Type: term.Type, Value: term.Value, Args: newArgs}
	}

	return term
}

func (e *Engine) solve(goals []Term, subst Substitution, sessionID string) []Substitution {
	if len(goals) == 0 {
		return []Substitution{subst}
	}

	goal := goals[0]
	remaining := goals[1:]

	if solutions, handled := e.evalBuiltin(goal, subst, sessionID); handled {
		var allResults []Substitution
		for _, sol := range solutions {
			results := e.solve(remaining, sol, sessionID)
			allResults = append(allResults, results...)
		}
		return allResults
	}

	return e.solveUserDefined(goal, remaining, subst, sessionID)
}

func (e *Engine) solveUserDefined(goal Term, remaining []Term, subst Substitution, sessionID string) []Substitution {
	key := e.makeCacheKey(goal, sessionID)
	if entry, exists := e.cache[key]; exists && entry.Complete {
		var results []Substitution
		for _, cachedSubst := range entry.Solutions {
			merged := e.mergeSubstitutions(subst, cachedSubst)
			results = append(results, e.solve(remaining, merged, sessionID)...)
		}
		return results
	}

	var factSolutions []Substitution
	var allResults []Substitution

	// Handle facts
	for _, fact := range e.loadFacts(goal, sessionID) {
		if newSubst, ok := e.unify(goal, fact.Predicate, subst); ok {
			factSolutions = append(factSolutions, newSubst)
		}
	}

	// Handle rules (includes remaining goals in the rule processing)
	for _, rule := range e.loadRules(goal, sessionID) {
		renamedRule := e.renameVars(rule)
		if newSubst, ok := e.unify(goal, renamedRule.Head, subst); ok {
			newGoals := append(renamedRule.Body, remaining...)
			results := e.solve(newGoals, newSubst, sessionID)
			allResults = append(allResults, results...)
		}
	}

	// Only apply remaining goals to fact solutions
	for _, sol := range factSolutions {
		results := e.solve(remaining, sol, sessionID)
		allResults = append(allResults, results...)
	}

	// Cache only the solutions for this specific goal (not including remaining)
	e.cache[key] = TableEntry{Solutions: factSolutions, Complete: true}

	return allResults
}

func (e *Engine) makeCacheKey(goal Term, sessionID string) TableKey {
	argsJSON, _ := json.Marshal(goal.Args)
	return TableKey{
		Predicate: fmt.Sprintf("%v_%s", goal.Value, sessionID),
		Args:      string(argsJSON),
	}
}

func (e *Engine) mergeSubstitutions(s1, s2 Substitution) Substitution {
	merged := make(Substitution)
	for k, v := range s1 {
		merged[k] = v
	}
	for k, v := range s2 {
		merged[k] = v
	}
	return merged
}

var globalVarCounter int = 0

func (e *Engine) renameVars(rule Rule) Rule {
	// Create a mapping for variable renaming
	varMap := make(map[string]string)
	globalVarCounter++
	suffix := fmt.Sprintf("_%d", globalVarCounter)
	
	// Rename variables in the head
	renamedHead := e.renameTermVars(rule.Head, varMap, suffix)
	
	// Rename variables in the body
	renamedBody := make([]Term, len(rule.Body))
	for i, term := range rule.Body {
		renamedBody[i] = e.renameTermVars(term, varMap, suffix)
	}
	
	return Rule{
		ID:        rule.ID,
		SessionID: rule.SessionID,
		Head:      renamedHead,
		Body:      renamedBody,
	}
}

func (e *Engine) renameTermVars(term Term, varMap map[string]string, suffix string) Term {
	switch term.Type {
	case "variable":
		varName := term.Value.(string)
		if newName, exists := varMap[varName]; exists {
			return Variable(newName)
		}
		newName := varName + suffix
		varMap[varName] = newName
		return Variable(newName)
	case "compound":
		newArgs := make([]Term, len(term.Args))
		for i, arg := range term.Args {
			newArgs[i] = e.renameTermVars(arg, varMap, suffix)
		}
		return Compound(term.Value.(string), newArgs)
	default:
		// Atoms and numbers don't need renaming
		return term
	}
}

func (e *Engine) loadFacts(goal Term, sessionID string) []Fact {
	predicate := e.extractPredicate(goal)
	if predicate == "" {
		return nil
	}

	rows, err := e.db.Query("SELECT id, session_id, data FROM facts WHERE predicate = ? AND session_id = ?", predicate, sessionID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var facts []Fact
	for rows.Next() {
		var fact Fact
		var data string
		rows.Scan(&fact.ID, &fact.SessionID, &data)
		json.Unmarshal([]byte(data), &fact.Predicate)
		facts = append(facts, fact)
	}
	return facts
}

func (e *Engine) loadRules(goal Term, sessionID string) []Rule {
	predicate := e.extractPredicate(goal)
	if predicate == "" {
		return nil
	}

	rows, err := e.db.Query("SELECT id, session_id, head_data, body_data FROM rules WHERE head_predicate = ? AND session_id = ?", predicate, sessionID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		var headData, bodyData string
		rows.Scan(&rule.ID, &rule.SessionID, &headData, &bodyData)
		json.Unmarshal([]byte(headData), &rule.Head)
		json.Unmarshal([]byte(bodyData), &rule.Body)
		rules = append(rules, rule)
	}
	return rules
}

func (e *Engine) extractPredicate(goal Term) string {
	if goal.Type == "compound" {
		return goal.Value.(string)
	} else if goal.Type == "atom" {
		return goal.Value.(string)
	}
	return ""
}

func (e *Engine) AddFact(fact Fact) error {
	predicate := e.extractPredicate(fact.Predicate)
	data, err := json.Marshal(fact.Predicate)
	if err != nil {
		return err
	}

	_, err = e.db.Exec("INSERT INTO facts (session_id, predicate, data) VALUES (?, ?, ?)", 
		fact.SessionID, predicate, string(data))
	return err
}

func (e *Engine) AddRule(rule Rule) error {
	predicate := e.extractPredicate(rule.Head)
	headData, err := json.Marshal(rule.Head)
	if err != nil {
		return err
	}
	bodyData, err := json.Marshal(rule.Body)
	if err != nil {
		return err
	}

	_, err = e.db.Exec("INSERT INTO rules (session_id, head_predicate, head_data, body_data) VALUES (?, ?, ?, ?)",
		rule.SessionID, predicate, string(headData), string(bodyData))
	return err
}

func (e *Engine) Query(query Query, sessionID string) QueryResult {
	solutions := e.solve(query.Goals, make(Substitution), sessionID)

	var result QueryResult
	if len(solutions) == 0 {
		result.Solutions = []Solution{{Success: false}}
	} else {
		for _, subst := range solutions {
			// Only include bindings for variables that appeared in the original query
			cleanedBindings := e.extractQueryBindings(query.Goals, subst)
			result.Solutions = append(result.Solutions, Solution{
				Bindings: cleanedBindings,
				Success:  true,
			})
		}
	}

	return result
}

func (e *Engine) extractQueryBindings(goals []Term, subst Substitution) Substitution {
	queryVars := make(map[string]bool)
	// Collect all variables from the query
	for _, goal := range goals {
		e.collectVars(goal, queryVars)
	}
	
	// Create a new substitution with only query variables, fully dereferenced
	cleaned := make(Substitution)
	for varName := range queryVars {
		if val, exists := subst[varName]; exists {
			cleaned[varName] = e.deref(val, subst)
		}
	}
	return cleaned
}

func (e *Engine) collectVars(term Term, vars map[string]bool) {
	switch term.Type {
	case "variable":
		vars[term.Value.(string)] = true
	case "compound":
		for _, arg := range term.Args {
			e.collectVars(arg, vars)
		}
	}
}

func (e *Engine) ClearCache() {
	e.cache = make(map[TableKey]TableEntry)
}

func (e *Engine) CreateSession(req CreateSessionRequest) (*Session, error) {
	now := time.Now()
	
	// Generate ULID
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(now), entropy).String()
	
	_, err := e.db.Exec("INSERT INTO sessions (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		id, req.Name, req.Description, now, now)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (e *Engine) GetSession(id string) (*Session, error) {
	var session Session
	err := e.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM sessions WHERE id = ?", id).Scan(
		&session.ID, &session.Name, &session.Description, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (e *Engine) GetSessionByName(name string) (*Session, error) {
	var session Session
	err := e.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM sessions WHERE name = ?", name).Scan(
		&session.ID, &session.Name, &session.Description, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (e *Engine) ListSessions() ([]Session, error) {
	rows, err := e.db.Query("SELECT id, name, description, created_at, updated_at FROM sessions ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.Name, &session.Description, &session.CreatedAt, &session.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (e *Engine) DeleteSession(id string) error {
	_, err := e.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

func (e *Engine) UpdateSessionTimestamp(sessionID string) error {
	_, err := e.db.Exec("UPDATE sessions SET updated_at = ? WHERE id = ?", time.Now(), sessionID)
	return err
}

func (e *Engine) Close() error {
	return e.db.Close()
}