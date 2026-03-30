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

// ─── Concrete Strategy ──────────────────────────────────────────────────────
// TODO: Implement the selectAgent() method

class RoundRobinStrategy : public AssignmentStrategy {
    int nextIndex = 0;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        // TODO: Return the ID of the next agent in rotation
        // Cycle through agents: 0, 1, 2, 0, 1, 2, ...
        // Don't forget to handle wrapping around
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
        // Push issue into the issues vector and return it
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

#ifndef RUNNING_TESTS
int main() {
    cout << "Issue Resolver — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
