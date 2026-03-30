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

// ─── Existing Strategy ──────────────────────────────────────────────────────
// TODO: Copy your Part 1 round-robin strategy here (or extend it)

class RoundRobinStrategy : public AssignmentStrategy {
    int nextIndex = 0;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        return -1; // TODO: implement
    }
};

// ─── NEW: Least-Loaded Strategy ─────────────────────────────────────────────
// HINT: Find the agent with the lowest currentLoad.
// Break ties by lowest agent ID.

class LeastLoadedStrategy : public AssignmentStrategy {
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        return -1; // TODO: implement
    }
};

// ─── NEW: Specialist Strategy ───────────────────────────────────────────────
// HINT: Find agents whose specializations include the issue's category.
// Among specialists, pick the one with lowest currentLoad.
// If no specialist exists, fall back to least-loaded.

class SpecialistStrategy : public AssignmentStrategy {
    LeastLoadedStrategy fallback;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        return -1; // TODO: implement
    }
};

// ─── Resolver ───────────────────────────────────────────────────────────────

class IssueResolver {
    AssignmentStrategy* strategy;
public:
    IssueResolver(AssignmentStrategy* s) : strategy(s) {}
    void setStrategy(AssignmentStrategy* s) { strategy = s; }

    Issue assign(vector<Agent>& agents, vector<Issue>& issues, Issue issue) {
        // TODO: implement
        return issue;
    }

    vector<Issue> getAgentIssues(const vector<Issue>& issues, int agentId) {
        // TODO: implement
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
    cout << "Part 2: Multiple strategies — implement the TODOs above." << endl;
    return 0;
}
#endif
