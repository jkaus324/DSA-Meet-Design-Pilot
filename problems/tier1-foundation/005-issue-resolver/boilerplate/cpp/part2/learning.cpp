#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

enum class Priority { LOW, MEDIUM, HIGH, CRITICAL };
enum class IssueState { OPEN, IN_PROGRESS, RESOLVED, CLOSED };
enum class Category { BILLING, TECHNICAL, GENERAL, ACCOUNT };

struct Issue {
    int id;
    string description;
    Category category;
    Priority priority;
    IssueState state;
    int assignedAgentId;
};

struct Agent {
    int id;
    string name;
    int currentLoad;
    vector<Category> specializations;
};

// ─── Assignment Interface ───────────────────────────────────────────────────

class AssignmentStrategy {
public:
    virtual int selectAgent(vector<Agent>& agents, const Issue& issue) = 0;
    virtual ~AssignmentStrategy() = default;
};

// ─── Concrete Strategies ────────────────────────────────────────────────────

class RoundRobinStrategy : public AssignmentStrategy {
    int nextIndex = 0;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        // TODO: Return ID of next agent in rotation, wrapping around
        return -1;
    }
};

class LeastLoadedStrategy : public AssignmentStrategy {
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        // TODO: Find agent with lowest currentLoad
        // Break ties by lowest agent ID
        return -1;
    }
};

class SpecialistStrategy : public AssignmentStrategy {
    LeastLoadedStrategy fallback;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        // TODO: Find agents whose specializations include issue.category
        // Among specialists, pick least-loaded (tiebreak by ID)
        // If no specialist found, delegate to fallback.selectAgent()
        return -1;
    }
};

// ─── Resolver ───────────────────────────────────────────────────────────────

class IssueResolver {
    AssignmentStrategy* strategy;
public:
    IssueResolver(AssignmentStrategy* s) : strategy(s) {}
    void setStrategy(AssignmentStrategy* s) { strategy = s; }

    Issue assign(vector<Agent>& agents, vector<Issue>& issues, Issue issue) {
        // TODO: Use strategy->selectAgent() to pick an agent
        // Set issue.assignedAgentId, increment agent's currentLoad
        // Push issue into issues vector and return it
        return issue;
    }

    vector<Issue> getAgentIssues(const vector<Issue>& issues, int agentId) {
        // TODO: Return all issues assigned to the given agentId
        return {};
    }
};

// ─── Test Entry Points ──────────────────────────────────────────────────────

static RoundRobinStrategy globalRoundRobin;
static IssueResolver globalResolver(&globalRoundRobin);

Issue assign_issue(vector<Agent>& agents, vector<Issue>& issues, Issue issue) {
    return globalResolver.assign(agents, issues, issue);
}

vector<Issue> get_agent_issues(const vector<Issue>& issues, int agentId) {
    return globalResolver.getAgentIssues(issues, agentId);
}

Issue assign_least_loaded(vector<Agent>& agents, vector<Issue>& issues, Issue issue) {
    LeastLoadedStrategy s;
    IssueResolver resolver(&s);
    return resolver.assign(agents, issues, issue);
}

Issue assign_by_specialist(vector<Agent>& agents, vector<Issue>& issues, Issue issue) {
    SpecialistStrategy s;
    IssueResolver resolver(&s);
    return resolver.assign(agents, issues, issue);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Multiple strategies — all scaffolding provided, implement selectAgent() methods." << endl;
    return 0;
}
#endif
