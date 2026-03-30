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
        if (agents.empty()) return -1;
        int idx = nextIndex % agents.size();
        nextIndex = (nextIndex + 1) % agents.size();
        return agents[idx].id;
    }
};

class LeastLoadedStrategy : public AssignmentStrategy {
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        if (agents.empty()) return -1;
        int bestIdx = 0;
        for (int i = 1; i < (int)agents.size(); i++) {
            if (agents[i].currentLoad < agents[bestIdx].currentLoad) {
                bestIdx = i;
            } else if (agents[i].currentLoad == agents[bestIdx].currentLoad
                       && agents[i].id < agents[bestIdx].id) {
                bestIdx = i;
            }
        }
        return agents[bestIdx].id;
    }
};

class SpecialistStrategy : public AssignmentStrategy {
    LeastLoadedStrategy fallback;
public:
    int selectAgent(vector<Agent>& agents, const Issue& issue) override {
        if (agents.empty()) return -1;
        int bestIdx = -1;
        for (int i = 0; i < (int)agents.size(); i++) {
            bool isSpecialist = false;
            for (auto& cat : agents[i].specializations) {
                if (cat == issue.category) { isSpecialist = true; break; }
            }
            if (!isSpecialist) continue;
            if (bestIdx == -1
                || agents[i].currentLoad < agents[bestIdx].currentLoad
                || (agents[i].currentLoad == agents[bestIdx].currentLoad
                    && agents[i].id < agents[bestIdx].id)) {
                bestIdx = i;
            }
        }
        if (bestIdx == -1) return fallback.selectAgent(agents, issue);
        return agents[bestIdx].id;
    }
};

// ─── Observer Interface ─────────────────────────────────────────────────────

class IssueObserver {
public:
    virtual void onStateChange(int issueId, IssueState oldState, IssueState newState) = 0;
    virtual ~IssueObserver() = default;
};

string stateName(IssueState s) {
    switch (s) {
        case IssueState::OPEN: return "OPEN";
        case IssueState::IN_PROGRESS: return "IN_PROGRESS";
        case IssueState::RESOLVED: return "RESOLVED";
        case IssueState::CLOSED: return "CLOSED";
    }
    return "UNKNOWN";
}

class LoggingObserver : public IssueObserver {
    vector<string>& log;
public:
    LoggingObserver(vector<string>& logRef) : log(logRef) {}
    void onStateChange(int issueId, IssueState oldState, IssueState newState) override {
        // TODO: Push formatted string to log
        // Format: "Issue <id>: <OLD_STATE> -> <NEW_STATE>"
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
        int agentId = strategy->selectAgent(agents, issue);
        issue.assignedAgentId = agentId;
        issue.state = IssueState::OPEN;
        for (auto& agent : agents) {
            if (agent.id == agentId) { agent.currentLoad++; break; }
        }
        issues.push_back(issue);
        return issue;
    }

    vector<Issue> getAgentIssues(const vector<Issue>& issues, int agentId) {
        vector<Issue> result;
        for (auto& issue : issues) {
            if (issue.assignedAgentId == agentId) result.push_back(issue);
        }
        return result;
    }

    bool transitionState(vector<Issue>& issues, int issueId, IssueState newState) {
        // TODO: Find the issue by ID
        // Validate the transition: OPEN->IN_PROGRESS, IN_PROGRESS->RESOLVED, RESOLVED->CLOSED
        // If invalid, return false without modifying state
        // If valid, update state and notify all observers
        return false;
    }

    Issue assignNextPriority(vector<Agent>& agents, vector<Issue>& issues) {
        // TODO: Find the highest-priority unassigned OPEN issue
        // Priority: CRITICAL(3) > HIGH(2) > MEDIUM(1) > LOW(0)
        // Tiebreak: lowest issue ID first
        // Remove it from issues, then call assign() to assign it
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
    cout << "Part 3: State tracking + priority — full scaffolding provided." << endl;
    return 0;
}
#endif
