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

// ─── Assignment Interface ───────────────────────────────────────────────────
// HINT: This interface lets you swap assignment logic at runtime.
// What method signature would let you select an agent for an issue?

type AssignmentStrategy interface {
	// HINT: returns the ID of the selected agent
	SelectAgent(agents []Agent, issue Issue) int
}

// ─── Concrete Strategy ──────────────────────────────────────────────────────
// TODO: Implement a round-robin strategy:
//   - Track which agent was assigned last
//   - Cycle through agents in order (0, 1, 2, 0, 1, 2, ...)
//   - Return the selected agent's ID

// ─── Resolver ───────────────────────────────────────────────────────────────
// TODO: Implement an IssueResolver struct that:
//   - Accepts any AssignmentStrategy
//   - Has an Assign() method that assigns an issue to the selected agent
//   - Has a GetAgentIssues() method to retrieve issues for a given agent
//   - Does NOT know about specific assignment algorithms

// type IssueResolver struct { ... }
// func (r *IssueResolver) Assign(agents []Agent, issues *[]Issue, issue Issue) Issue
// func (r *IssueResolver) GetAgentIssues(issues []Issue, agentID int) []Issue

// ─── Test Entry Points (must exist for tests to compile) ────────────────────
// Your solution must provide these functions:

func AssignIssue(agents []Agent, issues *[]Issue, issue Issue) Issue {
	return issue // TODO: use IssueResolver with RoundRobinStrategy
}

func GetAgentIssues(issues []Issue, agentID int) []Issue {
	return nil // TODO: use IssueResolver
}
