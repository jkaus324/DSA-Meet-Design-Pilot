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

// ─── Concrete Strategies ────────────────────────────────────────────────────

type RoundRobinStrategy struct {
	nextIndex int
}

func (s *RoundRobinStrategy) SelectAgent(agents []Agent, issue Issue) int {
	// TODO: Return ID of next agent in rotation, wrapping around
	return -1
}

type LeastLoadedStrategy struct{}

func (s *LeastLoadedStrategy) SelectAgent(agents []Agent, issue Issue) int {
	// TODO: Find agent with lowest CurrentLoad
	// Break ties by lowest agent ID
	return -1
}

type SpecialistStrategy struct {
	fallback LeastLoadedStrategy
}

func (s *SpecialistStrategy) SelectAgent(agents []Agent, issue Issue) int {
	// TODO: Find agents whose Specializations include issue.Cat
	// Among specialists, pick least-loaded (tiebreak by ID)
	// If no specialist found, delegate to s.fallback.SelectAgent()
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

func AssignLeastLoaded(agents []Agent, issues *[]Issue, issue Issue) Issue {
	s := &LeastLoadedStrategy{}
	return NewIssueResolver(s).Assign(agents, issues, issue)
}

func AssignBySpecialist(agents []Agent, issues *[]Issue, issue Issue) Issue {
	s := &SpecialistStrategy{}
	return NewIssueResolver(s).Assign(agents, issues, issue)
}
