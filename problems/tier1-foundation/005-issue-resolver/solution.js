// Issue Resolver — Strategy + Observer reference solution (JavaScript).

const Priority = { LOW: 0, MEDIUM: 1, HIGH: 2, CRITICAL: 3 };

const IssueState = {
  OPEN: "OPEN",
  IN_PROGRESS: "IN_PROGRESS",
  RESOLVED: "RESOLVED",
  CLOSED: "CLOSED",
};

const Category = {
  BILLING: "BILLING",
  TECHNICAL: "TECHNICAL",
  GENERAL: "GENERAL",
  ACCOUNT: "ACCOUNT",
};

function state_name(s) {
  return s;
}

function prio_from_string(s) {
  const m = { LOW: Priority.LOW, MEDIUM: Priority.MEDIUM, HIGH: Priority.HIGH, CRITICAL: Priority.CRITICAL };
  return Object.prototype.hasOwnProperty.call(m, s) ? m[s] : Priority.LOW;
}

class Issue {
  constructor(id, description, category, priority) {
    this.id = id;
    this.description = description;
    this.category = category;
    this.priority = priority;
    this.state = IssueState.OPEN;
    this.assignedAgentId = -1;
  }
}

class Agent {
  constructor(id, name, specialization = null) {
    this.id = id;
    this.name = name;
    this.currentLoad = 0;
    this.specializations = specialization ? [specialization] : [];
  }
}

class AssignmentStrategy {
  selectAgent(agents, issue) { throw new Error('not implemented'); }
}

class RoundRobinStrategy extends AssignmentStrategy {
  constructor() {
    super();
    this.nextIndex = 0;
  }
  selectAgent(agents, issue) {
    if (agents.length === 0) return -1;
    const idx = this.nextIndex % agents.length;
    this.nextIndex = (this.nextIndex + 1) % agents.length;
    return agents[idx].id;
  }
}

class LeastLoadedStrategy extends AssignmentStrategy {
  selectAgent(agents, issue) {
    if (agents.length === 0) return -1;
    let best = agents[0];
    for (let i = 1; i < agents.length; i++) {
      const a = agents[i];
      if (a.currentLoad < best.currentLoad || (a.currentLoad === best.currentLoad && a.id < best.id)) {
        best = a;
      }
    }
    return best.id;
  }
}

class SpecialistStrategy extends AssignmentStrategy {
  constructor() {
    super();
    this.fallback = new LeastLoadedStrategy();
  }
  selectAgent(agents, issue) {
    if (agents.length === 0) return -1;
    const specialists = agents.filter(a => a.specializations.includes(issue.category));
    if (specialists.length === 0) {
      return this.fallback.selectAgent(agents, issue);
    }
    let best = specialists[0];
    for (let i = 1; i < specialists.length; i++) {
      const a = specialists[i];
      if (a.currentLoad < best.currentLoad || (a.currentLoad === best.currentLoad && a.id < best.id)) {
        best = a;
      }
    }
    return best.id;
  }
}

class IssueObserver {
  onStateChange(issueId, oldState, newState) { throw new Error('not implemented'); }
}

class LoggingObserver extends IssueObserver {
  constructor(log) {
    super();
    this.log = log;
  }
  onStateChange(issueId, oldState, newState) {
    this.log.push(`Issue ${issueId}: ${oldState} -> ${newState}`);
  }
}

class IssueResolver {
  constructor(strategy) {
    this.strategy = strategy;
    this.observers = [];
  }
  addObserver(obs) {
    this.observers.push(obs);
  }
  assign(agents, issues, issue) {
    const agentId = this.strategy.selectAgent(agents, issue);
    issue.assignedAgentId = agentId;
    for (const a of agents) {
      if (a.id === agentId) {
        a.currentLoad += 1;
        break;
      }
    }
    issues.push(issue);
    return issue;
  }
  getAgentIssues(issues, agentId) {
    return issues.filter(i => i.assignedAgentId === agentId);
  }
  transitionState(issues, issueId, newState) {
    for (const issue of issues) {
      if (issue.id !== issueId) continue;
      const old = issue.state;
      const valid =
        (old === IssueState.OPEN && newState === IssueState.IN_PROGRESS) ||
        (old === IssueState.IN_PROGRESS && newState === IssueState.RESOLVED) ||
        (old === IssueState.RESOLVED && newState === IssueState.CLOSED);
      if (!valid) return false;
      issue.state = newState;
      for (const obs of this.observers) {
        obs.onStateChange(issueId, old, newState);
      }
      return true;
    }
    return false;
  }
}

let _g_agents = [];
let _g_issues = [];
let _g_issue_id = 0;
let _g_rr = null;
let _g_ir = null;
let _g_log = [];
let _g_logger = null;

function reset_service() {
  _g_agents = [];
  _g_issues = [];
  _g_issue_id = 0;
  _g_rr = new RoundRobinStrategy();
  _g_ir = new IssueResolver(_g_rr);
  _g_log = [];
  _g_logger = new LoggingObserver(_g_log);
  _g_ir.addObserver(_g_logger);
}

function ir_add_agent(id, name, specialization) {
  _g_agents.push(new Agent(id, name, specialization ? specialization : null));
}

function ir_assign_issue_round_robin(description, category, priority) {
  _g_issue_id += 1;
  const issue = new Issue(_g_issue_id, description, category, prio_from_string(priority));
  return _g_ir.assign(_g_agents, _g_issues, issue).assignedAgentId;
}

function ir_assign_issue_least_loaded(description, category, priority) {
  _g_issue_id += 1;
  const issue = new Issue(_g_issue_id, description, category, prio_from_string(priority));
  const r = new IssueResolver(new LeastLoadedStrategy());
  return r.assign(_g_agents, _g_issues, issue).assignedAgentId;
}

function ir_assign_issue_specialist(description, category, priority) {
  _g_issue_id += 1;
  const issue = new Issue(_g_issue_id, description, category, prio_from_string(priority));
  const r = new IssueResolver(new SpecialistStrategy());
  return r.assign(_g_agents, _g_issues, issue).assignedAgentId;
}

function ir_agent_issue_count(agentId) {
  return _g_ir.getAgentIssues(_g_issues, agentId).length;
}

function ir_agent_load(agentId) {
  for (const a of _g_agents) {
    if (a.id === agentId) return a.currentLoad;
  }
  return -1;
}

function ir_transition(issueId, newState) {
  return _g_ir.transitionState(_g_issues, issueId, newState);
}

function ir_get_issue_state(issueId) {
  for (const i of _g_issues) {
    if (i.id === issueId) return i.state;
  }
  return "";
}

function ir_log_size() {
  return _g_log.length;
}

function ir_log_entry(idx) {
  if (idx >= 0 && idx < _g_log.length) return _g_log[idx];
  return "";
}

module.exports = {
  Issue,
  Agent,
  AssignmentStrategy,
  RoundRobinStrategy,
  LeastLoadedStrategy,
  SpecialistStrategy,
  IssueObserver,
  LoggingObserver,
  IssueResolver,
  reset_service,
  ir_add_agent,
  ir_assign_issue_round_robin,
  ir_assign_issue_least_loaded,
  ir_assign_issue_specialist,
  ir_agent_issue_count,
  ir_agent_load,
  ir_transition,
  ir_get_issue_state,
  ir_log_size,
  ir_log_entry,
};
