#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <queue>
#include <unordered_map>
using namespace std;

// ─── Data Structures ────────────────────────────────────────────────────────

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

// ─── Strategy Interface ─────────────────────────────────────────────────────

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
        log.push_back("Issue " + to_string(issueId) + ": "
                       + stateName(oldState) + " -> " + stateName(newState));
    }
};

// ─── IssueResolver ──────────────────────────────────────────────────────────

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
        for (auto& issue : issues) {
            if (issue.id != issueId) continue;
            IssueState old = issue.state;
            bool valid = false;
            if (old == IssueState::OPEN && newState == IssueState::IN_PROGRESS) valid = true;
            if (old == IssueState::IN_PROGRESS && newState == IssueState::RESOLVED) valid = true;
            if (old == IssueState::RESOLVED && newState == IssueState::CLOSED) valid = true;
            if (!valid) return false;
            issue.state = newState;
            for (auto* obs : observers) {
                obs->onStateChange(issueId, old, newState);
            }
            return true;
        }
        return false;
    }

    Issue assignNextPriority(vector<Agent>& agents, vector<Issue>& issues) {
        int bestIdx = -1;
        for (int i = 0; i < (int)issues.size(); i++) {
            if (issues[i].assignedAgentId != -1) continue;
            if (issues[i].state != IssueState::OPEN) continue;
            if (bestIdx == -1) { bestIdx = i; continue; }
            if ((int)issues[i].priority > (int)issues[bestIdx].priority) {
                bestIdx = i;
            } else if (issues[i].priority == issues[bestIdx].priority
                       && issues[i].id < issues[bestIdx].id) {
                bestIdx = i;
            }
        }
        if (bestIdx == -1) return {-1, "", Category::GENERAL, Priority::LOW, IssueState::OPEN, -1};
        Issue issue = issues[bestIdx];
        issues.erase(issues.begin() + bestIdx);
        return assign(agents, issues, issue);
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

// ─── Main ───────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    vector<Agent> agents = {
        {0, "Alice", 0, {Category::BILLING}},
        {1, "Bob",   0, {Category::TECHNICAL}},
        {2, "Carol", 0, {Category::GENERAL}},
    };
    vector<Issue> issues;

    Issue i1 = assign_issue(agents, issues, {1, "Can't pay", Category::BILLING, Priority::HIGH, IssueState::OPEN, -1});
    Issue i2 = assign_issue(agents, issues, {2, "App crash", Category::TECHNICAL, Priority::CRITICAL, IssueState::OPEN, -1});

    cout << "Issue " << i1.id << " assigned to agent " << i1.assignedAgentId << endl;
    cout << "Issue " << i2.id << " assigned to agent " << i2.assignedAgentId << endl;

    return 0;
}
#endif
