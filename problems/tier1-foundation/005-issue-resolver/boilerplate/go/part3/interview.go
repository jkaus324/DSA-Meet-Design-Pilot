package main

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

// ─── NEW in Extension 2 ─────────────────────────────────────────────────────
//
// The product team now requires:
//   1. State machine: OPEN -> IN_PROGRESS -> RESOLVED -> CLOSED
//      Invalid transitions must be rejected (return false).
//   2. Notifications: When state changes, all registered observers are notified.
//   3. Priority queue: AssignNextPriority() picks the highest-priority
//      unassigned issue first. Tiebreak by lowest issue ID.
//
// Think about:
//   - How do you decouple notifications from the state machine?
//   - What interface lets you add new observers without modifying the resolver?
//   - How do you efficiently find the highest-priority unassigned issue?
//
// Entry points (must exist for tests):
//   func AssignIssue(agents []Agent, issues *[]Issue, issue Issue) Issue
//   func GetAgentIssues(issues []Issue, agentID int) []Issue
//   func AssignLeastLoaded(agents []Agent, issues *[]Issue, issue Issue) Issue
//   func AssignBySpecialist(agents []Agent, issues *[]Issue, issue Issue) Issue
//   func TransitionIssue(issues *[]Issue, issueID int, newState IssueState, notifications *[]string) bool
//   func AssignNextPriority(agents []Agent, issues *[]Issue) Issue
//
// ─────────────────────────────────────────────────────────────────────────────
