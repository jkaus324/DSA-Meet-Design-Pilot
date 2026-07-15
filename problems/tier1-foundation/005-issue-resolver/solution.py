"""Issue Resolver — Strategy + Observer reference solution (Python)."""

from abc import ABC, abstractmethod


class Priority:
    LOW = 0
    MEDIUM = 1
    HIGH = 2
    CRITICAL = 3


class IssueState:
    OPEN = "OPEN"
    IN_PROGRESS = "IN_PROGRESS"
    RESOLVED = "RESOLVED"
    CLOSED = "CLOSED"


class Category:
    BILLING = "BILLING"
    TECHNICAL = "TECHNICAL"
    GENERAL = "GENERAL"
    ACCOUNT = "ACCOUNT"


def state_name(s):
    return s


def prio_from_string(s):
    return {"LOW": Priority.LOW, "MEDIUM": Priority.MEDIUM, "HIGH": Priority.HIGH, "CRITICAL": Priority.CRITICAL}.get(s, Priority.LOW)


class Issue:
    def __init__(self, id, description, category, priority):
        self.id = id
        self.description = description
        self.category = category
        self.priority = priority
        self.state = IssueState.OPEN
        self.assignedAgentId = -1


class Agent:
    def __init__(self, id, name, specialization=None):
        self.id = id
        self.name = name
        self.currentLoad = 0
        self.specializations = [specialization] if specialization else []


class AssignmentStrategy(ABC):
    @abstractmethod
    def selectAgent(self, agents, issue):
        ...


class RoundRobinStrategy(AssignmentStrategy):
    def __init__(self):
        self.nextIndex = 0

    def selectAgent(self, agents, issue):
        if not agents:
            return -1
        idx = self.nextIndex % len(agents)
        self.nextIndex = (self.nextIndex + 1) % len(agents)
        return agents[idx].id


class LeastLoadedStrategy(AssignmentStrategy):
    def selectAgent(self, agents, issue):
        if not agents:
            return -1
        best = agents[0]
        for a in agents[1:]:
            if a.currentLoad < best.currentLoad or (a.currentLoad == best.currentLoad and a.id < best.id):
                best = a
        return best.id


class SpecialistStrategy(AssignmentStrategy):
    def __init__(self):
        self.fallback = LeastLoadedStrategy()

    def selectAgent(self, agents, issue):
        if not agents:
            return -1
        specialists = [a for a in agents if issue.category in a.specializations]
        if not specialists:
            return self.fallback.selectAgent(agents, issue)
        best = specialists[0]
        for a in specialists[1:]:
            if a.currentLoad < best.currentLoad or (a.currentLoad == best.currentLoad and a.id < best.id):
                best = a
        return best.id


class IssueObserver(ABC):
    @abstractmethod
    def onStateChange(self, issueId, oldState, newState):
        ...


class LoggingObserver(IssueObserver):
    def __init__(self, log):
        self.log = log

    def onStateChange(self, issueId, oldState, newState):
        self.log.append(f"Issue {issueId}: {oldState} -> {newState}")


class IssueResolver:
    def __init__(self, strategy):
        self.strategy = strategy
        self.observers = []

    def addObserver(self, obs):
        self.observers.append(obs)

    def assign(self, agents, issues, issue):
        agent_id = self.strategy.selectAgent(agents, issue)
        issue.assignedAgentId = agent_id
        for a in agents:
            if a.id == agent_id:
                a.currentLoad += 1
                break
        issues.append(issue)
        return issue

    def getAgentIssues(self, issues, agent_id):
        return [i for i in issues if i.assignedAgentId == agent_id]

    def transitionState(self, issues, issue_id, new_state):
        for issue in issues:
            if issue.id != issue_id:
                continue
            old = issue.state
            valid = (
                (old == IssueState.OPEN and new_state == IssueState.IN_PROGRESS)
                or (old == IssueState.IN_PROGRESS and new_state == IssueState.RESOLVED)
                or (old == IssueState.RESOLVED and new_state == IssueState.CLOSED)
            )
            if not valid:
                return False
            issue.state = new_state
            for obs in self.observers:
                obs.onStateChange(issue_id, old, new_state)
            return True
        return False


_g_agents = []
_g_issues = []
_g_issue_id = 0
_g_rr = None
_g_ir = None
_g_log = []
_g_logger = None


def reset_service():
    global _g_agents, _g_issues, _g_issue_id, _g_rr, _g_ir, _g_log, _g_logger
    _g_agents = []
    _g_issues = []
    _g_issue_id = 0
    _g_rr = RoundRobinStrategy()
    _g_ir = IssueResolver(_g_rr)
    _g_log = []
    _g_logger = LoggingObserver(_g_log)
    _g_ir.addObserver(_g_logger)


def ir_add_agent(id, name, specialization):
    _g_agents.append(Agent(id, name, specialization if specialization else None))


def ir_assign_issue_round_robin(description, category, priority):
    global _g_issue_id
    _g_issue_id += 1
    issue = Issue(_g_issue_id, description, category, prio_from_string(priority))
    return _g_ir.assign(_g_agents, _g_issues, issue).assignedAgentId


def ir_assign_issue_least_loaded(description, category, priority):
    global _g_issue_id
    _g_issue_id += 1
    issue = Issue(_g_issue_id, description, category, prio_from_string(priority))
    r = IssueResolver(LeastLoadedStrategy())
    return r.assign(_g_agents, _g_issues, issue).assignedAgentId


def ir_assign_issue_specialist(description, category, priority):
    global _g_issue_id
    _g_issue_id += 1
    issue = Issue(_g_issue_id, description, category, prio_from_string(priority))
    r = IssueResolver(SpecialistStrategy())
    return r.assign(_g_agents, _g_issues, issue).assignedAgentId


def ir_agent_issue_count(agent_id):
    return len(_g_ir.getAgentIssues(_g_issues, agent_id))


def ir_agent_load(agent_id):
    for a in _g_agents:
        if a.id == agent_id:
            return a.currentLoad
    return -1


def ir_transition(issue_id, new_state):
    return _g_ir.transitionState(_g_issues, issue_id, new_state)


def ir_get_issue_state(issue_id):
    for i in _g_issues:
        if i.id == issue_id:
            return i.state
    return ""


def ir_log_size():
    return len(_g_log)


def ir_log_entry(idx):
    if 0 <= idx < len(_g_log):
        return _g_log[idx]
    return ""
