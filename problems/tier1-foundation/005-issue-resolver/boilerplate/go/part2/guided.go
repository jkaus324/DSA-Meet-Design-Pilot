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

// ─── Existing Strategy ──────────────────────────────────────────────────────
// TODO: Copy your Part 1 round-robin strategy here (or extend it)

type RoundRobinStrategy struct {
	nextIndex int
}

func (s *RoundRobinStrategy) SelectAgent(agents []Agent, issue Issue) int {
	return -1 // TODO: implement
}

// ─── NEW: LeastLoadedStrategy ────────────────────────────────────────────────
// HINT: Find the agent with the lowest CurrentLoad.
// Break ties by lowest agent ID.

type LeastLoadedStrategy struct{}

func (s *LeastLoadedStrategy) SelectAgent(agents []Agent, issue Issue) int {
	return -1 // TODO: implement
}

// ─── NEW: SpecialistStrategy ─────────────────────────────────────────────────
// HINT: Find agents whose Specializations include the issue's category.
// Among specialists, pick the one with lowest CurrentLoad.
// If no specialist exists, fall back to least-loaded.

type SpecialistStrategy struct {
	fallback LeastLoadedStrategy
}

func (s *SpecialistStrategy) SelectAgent(agents []Agent, issue Issue) int {
	return -1 // TODO: implement
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
	// TODO: implement
	return issue
}

func (r *IssueResolver) GetAgentIssues(issues []Issue, agentID int) []Issue {
	// TODO: implement
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
