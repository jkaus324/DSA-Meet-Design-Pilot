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

// ─── NEW in Extension 1 ─────────────────────────────────────────────────────
//
// The ops team now wants PLUGGABLE assignment strategies:
// round-robin, least-loaded, and category-specialist.
//
// Think about:
//   - How do you add new assignment algorithms without modifying the resolver?
//   - What if the specialist strategy needs a fallback?
//   - Can you swap strategies at runtime?
//
// Entry points (must exist for tests):
//   func AssignIssue(agents []Agent, issues *[]Issue, issue Issue) Issue
//   func GetAgentIssues(issues []Issue, agentID int) []Issue
//   func AssignLeastLoaded(agents []Agent, issues *[]Issue, issue Issue) Issue
//   func AssignBySpecialist(agents []Agent, issues *[]Issue, issue Issue) Issue
//
// ─────────────────────────────────────────────────────────────────────────────
