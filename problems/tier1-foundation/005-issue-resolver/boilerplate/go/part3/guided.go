package main

import "fmt"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Priority int

const (
	PriorityLow      Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

type IssueState int

const (
	IssueStateOpen       IssueState = iota
	IssueStateInProgress
	IssueStateResolved
	IssueStateClosed
)

type Category int

const (
	CategoryBilling   Category = iota
	CategoryTechnical
	CategoryGeneral
	CategoryAccount
)

type Issue struct {
	ID              int
	Description     string
	Cat             Category
	Prio            Priority
	State           IssueState
	AssignedAgentID int
}

type Agent struct {
	ID              int
	Name            string
	CurrentLoad     int
	Specializations []Category
}

// ─── Assignment Interface ───────────────────────────────────────────────────

type AssignmentStrategy interface {
	SelectAgent(agents []Agent, issue Issue) int
}

// ─── Existing Strategies ────────────────────────────────────────────────────
// TODO: Copy your Part 1 + Part 2 strategies here

type RoundRobinStrategy struct{ nextIndex int }

func (s *RoundRobinStrategy) SelectAgent(agents []Agent, issue Issue) int {
	return -1 // TODO: implement
}

type LeastLoadedStrategy struct{}

func (s *LeastLoadedStrategy) SelectAgent(agents []Agent, issue Issue) int {
	return -1 // TODO: implement
}

type SpecialistStrategy struct{ fallback LeastLoadedStrategy }

func (s *SpecialistStrategy) SelectAgent(agents []Agent, issue Issue) int {
	return -1 // TODO: implement
}

// ─── NEW: Observer Interface ────────────────────────────────────────────────
// HINT: An observer receives notification when an issue's state changes.
// It gets the issue ID, old state, and new state.

type IssueObserver interface {
	OnStateChange(issueID int, oldState, newState IssueState)
}

// ─── NEW: LoggingObserver ────────────────────────────────────────────────────
// HINT: Writes a formatted string to a notifications slice.
// Format: "Issue <id>: <OLD_STATE> -> <NEW_STATE>"

type LoggingObserver struct {
	log *[]string
}

func NewLoggingObserver(log *[]string) *LoggingObserver {
	return &LoggingObserver{log: log}
}

func stateName(s IssueState) string {
	switch s {
	case IssueStateOpen:
		return "OPEN"
	case IssueStateInProgress:
		return "IN_PROGRESS"
	case IssueStateResolved:
		return "RESOLVED"
	case IssueStateClosed:
		return "CLOSED"
	}
	return "UNKNOWN"
}

func (o *LoggingObserver) OnStateChange(issueID int, oldState, newState IssueState) {
	// TODO: Format and append notification string to o.log
	// Format: "Issue <id>: <OLD_STATE> -> <NEW_STATE>"
	_ = fmt.Sprintf
}

// ─── Resolver ───────────────────────────────────────────────────────────────

type IssueResolver struct {
	strategy  AssignmentStrategy
	observers []IssueObserver
}

func NewIssueResolver(s AssignmentStrategy) *IssueResolver {
	return &IssueResolver{strategy: s}
}

func (r *IssueResolver) SetStrategy(s AssignmentStrategy) { r.strategy = s }
func (r *IssueResolver) AddObserver(obs IssueObserver)    { r.observers = append(r.observers, obs) }

func (r *IssueResolver) Assign(agents []Agent, issues *[]Issue, issue Issue) Issue {
	// TODO: implement (same as Part 1)
	return issue
}

func (r *IssueResolver) GetAgentIssues(issues []Issue, agentID int) []Issue {
	// TODO: implement (same as Part 1)
	return nil
}

// HINT: Valid transitions: OPEN->IN_PROGRESS, IN_PROGRESS->RESOLVED, RESOLVED->CLOSED
// Return false for invalid transitions. Notify all observers on success.
func (r *IssueResolver) TransitionState(issues *[]Issue, issueID int, newState IssueState) bool {
	// TODO: implement
	return false
}

// HINT: Find the highest-priority unassigned (AssignedAgentID == -1) OPEN issue.
// Priority order: CRITICAL(3) > HIGH(2) > MEDIUM(1) > LOW(0)
// Tiebreak: lowest issue ID first.
func (r *IssueResolver) AssignNextPriority(agents []Agent, issues *[]Issue) Issue {
	// TODO: implement
	return Issue{ID: -1}
}

// ─── Test Entry Points ──────────────────────────────────────────────────────

var globalRoundRobin = &RoundRobinStrategy{}
var globalResolver = NewIssueResolver(globalRoundRobin)

func AssignIssue(agents []Agent, issues *[]Issue, issue Issue) Issue {
	return globalResolver.Assign(agents, issues, issue)
}

func GetAgentIssues(issues []Issue, agentID int) []Issue {
	return globalResolver.GetAgentIssues(issues, agentID)
}

func AssignLeastLoaded(agents []Agent, issues *[]Issue, issue Issue) Issue {
	s := &LeastLoadedStrategy{}
	return NewIssueResolver(s).Assign(agents, issues, issue)
}

func AssignBySpecialist(agents []Agent, issues *[]Issue, issue Issue) Issue {
	s := &SpecialistStrategy{}
	return NewIssueResolver(s).Assign(agents, issues, issue)
}

func TransitionIssue(issues *[]Issue, issueID int, newState IssueState, notifications *[]string) bool {
	logger := NewLoggingObserver(notifications)
	rr := &RoundRobinStrategy{}
	resolver := NewIssueResolver(rr)
	resolver.AddObserver(logger)
	return resolver.TransitionState(issues, issueID, newState)
}

func AssignNextPriority(agents []Agent, issues *[]Issue) Issue {
	rr := &RoundRobinStrategy{}
	resolver := NewIssueResolver(rr)
	return resolver.AssignNextPriority(agents, issues)
}
