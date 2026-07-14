// Issue Resolver — Strategy + Observer reference solution (Go).
package main

import "fmt"

const (
	prioLow      = 0
	prioMedium   = 1
	prioHigh     = 2
	prioCritical = 3
)

const (
	stOpen       = "OPEN"
	stInProgress = "IN_PROGRESS"
	stResolved   = "RESOLVED"
	stClosed     = "CLOSED"
)

func prioFromString(s string) int {
	switch s {
	case "LOW":
		return prioLow
	case "MEDIUM":
		return prioMedium
	case "HIGH":
		return prioHigh
	case "CRITICAL":
		return prioCritical
	}
	return prioLow
}

type issue struct {
	id              int
	description     string
	category        string
	priority        int
	state           string
	assignedAgentID int
}

type agent struct {
	id              int
	name            string
	currentLoad     int
	specializations []string
}

func (a *agent) hasSpecialization(cat string) bool {
	for _, s := range a.specializations {
		if s == cat {
			return true
		}
	}
	return false
}

type assignmentStrategy interface {
	selectAgent(agents []*agent, iss *issue) int
}

type roundRobinStrategy struct {
	nextIndex int
}

func (r *roundRobinStrategy) selectAgent(agents []*agent, iss *issue) int {
	if len(agents) == 0 {
		return -1
	}
	idx := r.nextIndex % len(agents)
	r.nextIndex = (r.nextIndex + 1) % len(agents)
	return agents[idx].id
}

type leastLoadedStrategy struct{}

func (leastLoadedStrategy) selectAgent(agents []*agent, iss *issue) int {
	if len(agents) == 0 {
		return -1
	}
	best := agents[0]
	for _, a := range agents[1:] {
		if a.currentLoad < best.currentLoad || (a.currentLoad == best.currentLoad && a.id < best.id) {
			best = a
		}
	}
	return best.id
}

type specialistStrategy struct {
	fallback leastLoadedStrategy
}

func (s specialistStrategy) selectAgent(agents []*agent, iss *issue) int {
	if len(agents) == 0 {
		return -1
	}
	var specialists []*agent
	for _, a := range agents {
		if a.hasSpecialization(iss.category) {
			specialists = append(specialists, a)
		}
	}
	if len(specialists) == 0 {
		return s.fallback.selectAgent(agents, iss)
	}
	best := specialists[0]
	for _, a := range specialists[1:] {
		if a.currentLoad < best.currentLoad || (a.currentLoad == best.currentLoad && a.id < best.id) {
			best = a
		}
	}
	return best.id
}

type issueObserver interface {
	onStateChange(issueID int, oldState, newState string)
}

type loggingObserver struct {
	log *[]string
}

func (l loggingObserver) onStateChange(issueID int, oldState, newState string) {
	*l.log = append(*l.log, fmt.Sprintf("Issue %d: %s -> %s", issueID, oldState, newState))
}

type issueResolver struct {
	strategy  assignmentStrategy
	observers []issueObserver
}

func (ir *issueResolver) addObserver(obs issueObserver) {
	ir.observers = append(ir.observers, obs)
}

func (ir *issueResolver) assign(agents []*agent, issues *[]*issue, iss *issue) *issue {
	agentID := ir.strategy.selectAgent(agents, iss)
	iss.assignedAgentID = agentID
	for _, a := range agents {
		if a.id == agentID {
			a.currentLoad++
			break
		}
	}
	*issues = append(*issues, iss)
	return iss
}

func (ir *issueResolver) getAgentIssues(issues []*issue, agentID int) int {
	count := 0
	for _, i := range issues {
		if i.assignedAgentID == agentID {
			count++
		}
	}
	return count
}

func (ir *issueResolver) transitionState(issues []*issue, issueID int, newState string) bool {
	for _, iss := range issues {
		if iss.id != issueID {
			continue
		}
		old := iss.state
		valid := (old == stOpen && newState == stInProgress) ||
			(old == stInProgress && newState == stResolved) ||
			(old == stResolved && newState == stClosed)
		if !valid {
			return false
		}
		iss.state = newState
		for _, obs := range ir.observers {
			obs.onStateChange(issueID, old, newState)
		}
		return true
	}
	return false
}

var (
	gAgents  []*agent
	gIssues  []*issue
	gIssueID int
	gIR      *issueResolver
	gLog     []string
)

func reset_service() {
	gAgents = nil
	gIssues = nil
	gIssueID = 0
	gLog = nil
	gIR = &issueResolver{strategy: &roundRobinStrategy{}}
	gIR.addObserver(loggingObserver{log: &gLog})
}

func ir_add_agent(id int, name, specialization string) {
	a := &agent{id: id, name: name}
	if specialization != "" {
		a.specializations = []string{specialization}
	}
	gAgents = append(gAgents, a)
}

func ir_assign_issue_round_robin(description, category, priority string) int {
	gIssueID++
	iss := &issue{id: gIssueID, description: description, category: category, priority: prioFromString(priority), state: stOpen, assignedAgentID: -1}
	return gIR.assign(gAgents, &gIssues, iss).assignedAgentID
}

func ir_assign_issue_least_loaded(description, category, priority string) int {
	gIssueID++
	iss := &issue{id: gIssueID, description: description, category: category, priority: prioFromString(priority), state: stOpen, assignedAgentID: -1}
	r := &issueResolver{strategy: leastLoadedStrategy{}}
	return r.assign(gAgents, &gIssues, iss).assignedAgentID
}

func ir_assign_issue_specialist(description, category, priority string) int {
	gIssueID++
	iss := &issue{id: gIssueID, description: description, category: category, priority: prioFromString(priority), state: stOpen, assignedAgentID: -1}
	r := &issueResolver{strategy: specialistStrategy{}}
	return r.assign(gAgents, &gIssues, iss).assignedAgentID
}

func ir_agent_issue_count(agentID int) int {
	return gIR.getAgentIssues(gIssues, agentID)
}

func ir_agent_load(agentID int) int {
	for _, a := range gAgents {
		if a.id == agentID {
			return a.currentLoad
		}
	}
	return -1
}

func ir_transition(issueID int, newState string) bool {
	return gIR.transitionState(gIssues, issueID, newState)
}

func ir_get_issue_state(issueID int) string {
	for _, i := range gIssues {
		if i.id == issueID {
			return i.state
		}
	}
	return ""
}

func ir_log_size() int {
	return len(gLog)
}

func ir_log_entry(idx int) string {
	if idx >= 0 && idx < len(gLog) {
		return gLog[idx]
	}
	return ""
}
