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

// ─── Existing Strategies ────────────────────────────────────────────────────
// TODO: Copy your Part 1 + Part 2 strategies here

class RoundRobinStrategy : public AssignmentStrategy {
    int nextIndex = 0;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        return -1; // TODO: implement
    }
};

class LeastLoadedStrategy : public AssignmentStrategy {
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        return -1; // TODO: implement
    }
};

class SpecialistStrategy : public AssignmentStrategy {
    LeastLoadedStrategy fallback;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        return -1; // TODO: implement
    }
};

// ─── NEW: Observer Interface ────────────────────────────────────────────────
// HINT: An observer receives notification when an issue's state changes.
// It gets the issue ID, old state, and new state.

class IssueObserver {
public:
    virtual void onStateChange(int issueId, IssueState oldState, IssueState newState) = 0;
    virtual ~IssueObserver() = default;
};

// ─── NEW: Logging Observer ──────────────────────────────────────────────────
// HINT: Writes a formatted string to a notifications vector.
// Format: "Issue <id>: <OLD_STATE> -> <NEW_STATE>"

class LoggingObserver : public IssueObserver {
    vector<string>& log;
public:
    LoggingObserver(vector<string>& logRef) : log(logRef) {}
    void onStateChange(int issueId, IssueState oldState, IssueState newState) override {
        // TODO: Format and push notification string to log
    }
};

// ─── Resolver ───────────────────────────────────────────────────────────────

class IssueResolver {
    AssignmentStrategy* strategy;
    vector<IssueObserver*> observers;
public:
    IssueResolver(AssignmentStrategy* s) : strategy(s) {}
    void setStrategy(AssignmentStrategy* s) { strategy = s; }
    void addObserver(IssueObserver* obs) { observers.push_back(obs); }

    Issue assign(vector<Agent>& agents, vector<Issue>& issues, Issue issue) {
        // TODO: implement (same as Part 1)
        return issue;
    }

    vector<Issue> getAgentIssues(const vector<Issue>& issues, int agentId) {
        // TODO: implement (same as Part 1)
        return {};
    }

    // HINT: Valid transitions: OPEN->IN_PROGRESS, IN_PROGRESS->RESOLVED, RESOLVED->CLOSED
    // Return false for invalid transitions. Notify all observers on success.
    bool transitionState(vector<Issue>& issues, int issueId, IssueState newState) {
        // TODO: implement
        return false;
    }

    // HINT: Find the highest-priority unassigned (assignedAgentId == -1) OPEN issue.
    // Priority order: CRITICAL > HIGH > MEDIUM > LOW
    // Tiebreak: lowest issue ID first.
    Issue assignNextPriority(vector<Agent>& agents, vector<Issue>& issues) {
        // TODO: implement
        return {-1, "", Category::GENERAL, Priority::LOW, IssueState::OPEN, -1};
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

bool transition_issue(vector<Issue>& issues, int issueId,
                      IssueState newState, vector<string>& notifications) {
    LoggingObserver logger(notifications);
    RoundRobinStrategy rr;
    IssueResolver resolver(&rr);
    resolver.addObserver(&logger);
    return resolver.transitionState(issues, issueId, newState);
}

Issue assign_next_priority(vector<Agent>& agents, vector<Issue>& issues) {
    RoundRobinStrategy rr;
    IssueResolver resolver(&rr);
    return resolver.assignNextPriority(agents, issues);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: State tracking + priority — implement the TODOs above." << endl;
    return 0;
}
#endif
