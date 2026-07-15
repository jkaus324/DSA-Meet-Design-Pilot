// Data class (given — do not modify).

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function reset_service() {
  // TODO: implement this
  return null;
}

function ir_add_agent(id, name, specialization) {
  // TODO: implement this
  return null;
}

function ir_assign_issue_round_robin(description, category, priority) {
  // TODO: implement this
  return null;
}

function ir_agent_issue_count(agentId) {
  // TODO: implement this
  return null;
}

function ir_agent_load(agentId) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { reset_service, ir_add_agent, ir_assign_issue_round_robin, ir_agent_issue_count, ir_agent_load };
