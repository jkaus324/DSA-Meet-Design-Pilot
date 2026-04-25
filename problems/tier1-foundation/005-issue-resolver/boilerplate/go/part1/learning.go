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

type AssignmentStrategy interface {
	SelectAgent(agents []Agent, issue Issue) int
}

// ─── Concrete Strategy ──────────────────────────────────────────────────────
// TODO: Implement the SelectAgent() method

type RoundRobinStrategy struct {
	nextIndex int
}

func (s *RoundRobinStrategy) SelectAgent(agents []Agent, issue Issue) int {
	// TODO: Return the ID of the next agent in rotation
	// Cycle through agents: 0, 1, 2, 0, 1, 2, ...
	// Don't forget to handle wrapping around
	return -1
}

// ─── Resolver ───────────────────────────────────────────────────────────────

type IssueResolver struct {
	strategy AssignmentStrategy
}

func NewIssueResolver(s AssignmentStrategy) *IssueResolver {
	return &IssueResolver{strategy: s}
}

func (r *IssueResolver) SetStrategy(s AssignmentStrategy) {
	r.strategy = s
}

func (r *IssueResolver) Assign(agents []Agent, issues *[]Issue, issue Issue) Issue {
	// TODO: Use strategy.SelectAgent() to pick an agent
	// Set issue.AssignedAgentID, increment agent's CurrentLoad
	// Append issue to issues slice and return it
	return issue
}

func (r *IssueResolver) GetAgentIssues(issues []Issue, agentID int) []Issue {
	// TODO: Return all issues assigned to the given agentID
	return nil
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
