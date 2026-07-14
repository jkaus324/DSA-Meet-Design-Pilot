# Problem 005 — Customer Issue Resolution System

**Tier:** 1 (Foundation) | **Pattern:** Strategy + Observer | **DSA:** HashMap + Priority Queue
**Companies:** PhonePe, Flipkart | **Time:** 45 minutes

---

## Problem Statement

You're building the backend for a customer support platform. When customers report issues, the system must **assign each issue to an available support agent** and track its lifecycle from creation to closure.

Different teams use different assignment policies:
- A **round-robin** assigner cycles through agents in order
- A **least-loaded** assigner picks the agent with the fewest open issues
- A **category-specialist** assigner routes issues to agents who specialize in that category

As issues progress through states (Open, InProgress, Resolved, Closed), stakeholders must be **notified** of every state transition. High-priority issues must be processed before low-priority ones.

**Your task:** Design and implement an `IssueResolver` system that assigns issues to agents using pluggable strategies, tracks issue state transitions with observer notifications, and prioritizes high-priority issues using a priority queue.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What varies here? The *assignment algorithm* varies. The *resolver itself* stays the same.
2. If you used `if-else` inside `assign_issue()`, what happens when a 4th assignment policy is added? You modify existing code — violating Open/Closed Principle.
3. How does the Strategy pattern solve this? Each assignment policy becomes a separate class implementing a common interface.
4. When an issue changes state, who needs to know? Multiple observers — the assigned agent, a logger, an escalation system. How do you decouple the issue from its observers?

**The key insight:** Assignment policy is a Strategy; state change notification is an Observer. These two patterns work together — Strategy decides *who* gets the issue, Observer decides *who hears about* state changes.

---

## Data Structures

```cpp
enum class Priority { LOW, MEDIUM, HIGH, CRITICAL };
enum class IssueState { OPEN, IN_PROGRESS, RESOLVED, CLOSED };
enum class Category { BILLING, TECHNICAL, GENERAL, ACCOUNT };

struct Issue {
    int id;
    std::string description;
    Category category;
    Priority priority;
    IssueState state;       // starts as OPEN
    int assignedAgentId;    // -1 if unassigned
};

struct Agent {
    int id;
    std::string name;
    int currentLoad;        // number of currently assigned open issues
    std::vector<Category> specializations;  // categories this agent handles
};
```

---

## Base Requirement — Round-robin agent assignment

Implement an `IssueResolver` that assigns incoming issues to agents using a round-robin policy. The first issue goes to agent 0, the second to agent 1, and so on, wrapping around when all agents have been visited.

When an issue is assigned, its `assignedAgentId` is set and the agent's `currentLoad` increments by 1.

**Entry points (tests will call these):**
```cpp
Issue assign_issue(vector<Agent>& agents, vector<Issue>& issues, Issue issue);
vector<Issue> get_agent_issues(const vector<Issue>& issues, int agentId);
```

**What to implement:**
```cpp
class AssignmentStrategy {
public:
    virtual int selectAgent(vector<Agent>& agents, const Issue& issue) = 0;
    virtual ~AssignmentStrategy() = default;
};

class RoundRobinStrategy : public AssignmentStrategy { ... };

class IssueResolver {
public:
    IssueResolver(AssignmentStrategy* strategy);
    void setStrategy(AssignmentStrategy* strategy);
    Issue assign(vector<Agent>& agents, vector<Issue>& issues, Issue issue);
    vector<Issue> getAgentIssues(const vector<Issue>& issues, int agentId);
};
```

**Constraints:**
- `1 <= agents.size() <= 100`
- `1 <= issue.id <= 10000`
- Agent IDs are unique and zero-indexed
- Round-robin wraps around: after the last agent, the next assignment goes back to agent 0

---

## Extension 1 — Multiple assignment strategies

The ops team now wants **pluggable assignment policies**. In addition to round-robin, support:

| Strategy | Rule |
|----------|------|
| Round Robin | Cycle through agents in order |
| Least Loaded | Assign to the agent with the lowest `currentLoad`. Break ties by lowest agent ID. |
| Category Specialist | Assign to the agent whose `specializations` include the issue's `category`. Among specialists, pick the one with lowest `currentLoad`. If no specialist exists, fall back to least-loaded. |

**Design goal:** Adding a 4th strategy must require **zero changes** to `IssueResolver` itself.

**New entry points:**
```cpp
Issue assign_least_loaded(vector<Agent>& agents, vector<Issue>& issues, Issue issue);
Issue assign_by_specialist(vector<Agent>& agents, vector<Issue>& issues, Issue issue);
```

**Design challenge:** How do you add new assignment policies without modifying the resolver? Can you swap strategies at runtime?

---

## Extension 2 — Issue state tracking + priority

The product team now requires:

1. **State machine:** Issues transition through `OPEN -> IN_PROGRESS -> RESOLVED -> CLOSED`. Invalid transitions (e.g., `OPEN -> CLOSED`) must be rejected.

2. **Observer notifications:** When an issue's state changes, all registered observers are notified. An observer receives the issue ID, old state, and new state.

3. **Priority queue:** High-priority issues must be processed first. When multiple unassigned issues are waiting, `assign_next()` should pick the one with the highest priority (CRITICAL > HIGH > MEDIUM > LOW). Break ties by lowest issue ID.

**New entry points:**
```cpp
bool transition_issue(vector<Issue>& issues, int issueId,
                      IssueState newState, vector<string>& notifications);
Issue assign_next_priority(vector<Agent>& agents, vector<Issue>& issues);
```

**What to implement:**
```cpp
class IssueObserver {
public:
    virtual void onStateChange(int issueId, IssueState oldState, IssueState newState) = 0;
    virtual ~IssueObserver() = default;
};

class LoggingObserver : public IssueObserver { ... };  // logs to notifications vector
```

**Constraints:**
- Valid transitions: OPEN->IN_PROGRESS, IN_PROGRESS->RESOLVED, RESOLVED->CLOSED
- `transition_issue()` returns `false` for invalid transitions and does not modify state
- Priority ordering: CRITICAL(3) > HIGH(2) > MEDIUM(1) > LOW(0)
- Tie-breaking for equal priority: lowest issue ID first

---

## Running Tests

```bash
./run-tests.sh 005-issue-resolver cpp
```
