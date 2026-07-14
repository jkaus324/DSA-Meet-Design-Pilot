# Data class (given — do not modify).

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def reset_service():
    # TODO: implement this
    return None

def ir_add_agent(id, name, specialization):
    # TODO: implement this
    return None

def ir_assign_issue_round_robin(description, category, priority):
    # TODO: implement this
    return None

def ir_agent_issue_count(agentId):
    # TODO: implement this
    return None

def ir_agent_load(agentId):
    # TODO: implement this
    return None

def ir_assign_issue_least_loaded(description, category, priority):
    # TODO: implement this
    return None

def ir_assign_issue_specialist(description, category, priority):
    # TODO: implement this
    return None

def ir_transition(issueId, newState):
    # TODO: implement this
    return None

def ir_get_issue_state(issueId):
    # TODO: implement this
    return None

def ir_log_size():
    # TODO: implement this
    return None

def ir_log_entry(idx):
    # TODO: implement this
    return None
