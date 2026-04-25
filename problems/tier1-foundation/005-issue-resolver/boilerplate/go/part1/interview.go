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

// ─── Your Design Starts Here ────────────────────────────────────────────────
//
// Design and implement an IssueResolver that:
//   1. Assigns issues to agents using round-robin rotation
//   2. Allows new assignment strategies to be added WITHOUT modifying
//      the resolver itself
//
// Think about:
//   - What abstraction lets you swap assignment logic at runtime?
//   - How would you add a 4th assignment policy with zero changes
//     to existing code?
//   - What happens when Extension 1 (multiple strategies) is added?
//
// Entry points (must exist for tests):
//   func AssignIssue(agents []Agent, issues *[]Issue, issue Issue) Issue
//   func GetAgentIssues(issues []Issue, agentID int) []Issue
//
// ─────────────────────────────────────────────────────────────────────────────
