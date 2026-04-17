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

// ─── Concrete Strategies ────────────────────────────────────────────────────

type RoundRobinStrategy struct {
	nextIndex int
}

func (s *RoundRobinStrategy) SelectAgent(agents []Agent, issue Issue) int {
	if len(agents) == 0 {
		return -1
	}
	idx := s.nextIndex % len(agents)
	s.nextIndex = (s.nextIndex + 1) % len(agents)
	return agents[idx].ID
}

type LeastLoadedStrategy struct{}

func (s *LeastLoadedStrategy) SelectAgent(agents []Agent, issue Issue) int {
	if len(agents) == 0 {
		return -1
	}
	bestIdx := 0
	for i := 1; i < len(agents); i++ {
		if agents[i].CurrentLoad < agents[bestIdx].CurrentLoad ||
			(agents[i].CurrentLoad == agents[bestIdx].CurrentLoad && agents[i].ID < agents[bestIdx].ID) {
			bestIdx = i
		}
	}
	return agents[bestIdx].ID
}

type SpecialistStrategy struct {
	fallback LeastLoadedStrategy
}

func (s *SpecialistStrategy) SelectAgent(agents []Agent, issue Issue) int {
	if len(agents) == 0 {
		return -1
	}
	bestIdx := -1
	for i, agent := range agents {
		isSpecialist := false
		for _, cat := range agent.Specializations {
			if cat == issue.Cat {
				isSpecialist = true
				break
			}
		}
		if !isSpecialist {
			continue
		}
		if bestIdx == -1 ||
			agents[i].CurrentLoad < agents[bestIdx].CurrentLoad ||
			(agents[i].CurrentLoad == agents[bestIdx].CurrentLoad && agents[i].ID < agents[bestIdx].ID) {
			bestIdx = i
		}
	}
	if bestIdx == -1 {
		return s.fallback.SelectAgent(agents, issue)
	}
	return agents[bestIdx].ID
}

// ─── Observer Interface ─────────────────────────────────────────────────────

type IssueObserver interface {
	OnStateChange(issueID int, oldState, newState IssueState)
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

// ─── LoggingObserver ─────────────────────────────────────────────────────────

type LoggingObserver struct {
	log *[]string
}

func NewLoggingObserver(log *[]string) *LoggingObserver {
	return &LoggingObserver{log: log}
}

func (o *LoggingObserver) OnStateChange(issueID int, oldState, newState IssueState) {
	// TODO: Push formatted string to log
	// Format: "Issue <id>: <OLD_STATE> -> <NEW_STATE>"
	msg := fmt.Sprintf("Issue %d: %s -> %s", issueID, stateName(oldState), stateName(newState))
	*o.log = append(*o.log, msg)
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
	agentID := r.strategy.SelectAgent(agents, issue)
	issue.AssignedAgentID = agentID
	issue.State = IssueStateOpen
	for i := range agents {
		if agents[i].ID == agentID {
			agents[i].CurrentLoad++
			break
		}
	}
	*issues = append(*issues, issue)
	return issue
}

func (r *IssueResolver) GetAgentIssues(issues []Issue, agentID int) []Issue {
	var result []Issue
	for _, issue := range issues {
		if issue.AssignedAgentID == agentID {
			result = append(result, issue)
		}
	}
	return result
}

func (r *IssueResolver) TransitionState(issues *[]Issue, issueID int, newState IssueState) bool {
	// TODO: Find the issue by ID
	// Validate the transition: OPEN->IN_PROGRESS, IN_PROGRESS->RESOLVED, RESOLVED->CLOSED
	// If invalid, return false without modifying state
	// If valid, update state and notify all observers
	return false
}

func (r *IssueResolver) AssignNextPriority(agents []Agent, issues *[]Issue) Issue {
	// TODO: Find the highest-priority unassigned OPEN issue
	// Priority: CRITICAL(3) > HIGH(2) > MEDIUM(1) > LOW(0)
	// Tiebreak: lowest issue ID first
	// Remove it from issues, then call r.Assign() to assign it
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
