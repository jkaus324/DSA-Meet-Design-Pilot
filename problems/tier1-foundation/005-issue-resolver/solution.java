// Issue Resolver — Solution (Java, Strategy + Observer)
import java.util.*;

class IRIssue {
    public int id;
    public String description;
    public String category;
    public int priority;
    public String state = "OPEN";
    public int assignedAgentId = -1;

    public IRIssue(int id, String description, String category, int priority) {
        this.id = id; this.description = description;
        this.category = category; this.priority = priority;
    }
}

class IRAgent {
    public int id;
    public String name;
    public int currentLoad = 0;
    public List<String> specializations = new ArrayList<>();

    public IRAgent(int id, String name, String specialization) {
        this.id = id;
        this.name = name;
        if (specialization != null && !specialization.isEmpty()) {
            this.specializations.add(specialization);
        }
    }
}

interface AssignmentStrategy {
    int selectAgent(List<IRAgent> agents, IRIssue issue);
}

class RoundRobinStrategy implements AssignmentStrategy {
    private int nextIndex = 0;
    public int selectAgent(List<IRAgent> agents, IRIssue issue) {
        if (agents.isEmpty()) return -1;
        int idx = nextIndex % agents.size();
        nextIndex = (nextIndex + 1) % agents.size();
        return agents.get(idx).id;
    }
}

class LeastLoadedStrategy implements AssignmentStrategy {
    public int selectAgent(List<IRAgent> agents, IRIssue issue) {
        if (agents.isEmpty()) return -1;
        IRAgent best = agents.get(0);
        for (int i = 1; i < agents.size(); i++) {
            IRAgent a = agents.get(i);
            if (a.currentLoad < best.currentLoad ||
                (a.currentLoad == best.currentLoad && a.id < best.id)) {
                best = a;
            }
        }
        return best.id;
    }
}

class SpecialistStrategy implements AssignmentStrategy {
    private final LeastLoadedStrategy fallback = new LeastLoadedStrategy();
    public int selectAgent(List<IRAgent> agents, IRIssue issue) {
        if (agents.isEmpty()) return -1;
        List<IRAgent> specialists = new ArrayList<>();
        for (IRAgent a : agents) {
            if (a.specializations.contains(issue.category)) specialists.add(a);
        }
        if (specialists.isEmpty()) return fallback.selectAgent(agents, issue);
        IRAgent best = specialists.get(0);
        for (int i = 1; i < specialists.size(); i++) {
            IRAgent a = specialists.get(i);
            if (a.currentLoad < best.currentLoad ||
                (a.currentLoad == best.currentLoad && a.id < best.id)) {
                best = a;
            }
        }
        return best.id;
    }
}

interface IssueObserver {
    void onStateChange(int issueId, String oldState, String newState);
}

class LoggingObserver implements IssueObserver {
    private final List<String> log;
    public LoggingObserver(List<String> log) { this.log = log; }
    public void onStateChange(int issueId, String oldState, String newState) {
        log.add("Issue " + issueId + ": " + oldState + " -> " + newState);
    }
}

class IssueResolver {
    private AssignmentStrategy strategy;
    private final List<IssueObserver> observers = new ArrayList<>();

    public IssueResolver(AssignmentStrategy strategy) { this.strategy = strategy; }

    public void addObserver(IssueObserver obs) { observers.add(obs); }

    public IRIssue assign(List<IRAgent> agents, List<IRIssue> issues, IRIssue issue) {
        int agentId = strategy.selectAgent(agents, issue);
        issue.assignedAgentId = agentId;
        for (IRAgent a : agents) {
            if (a.id == agentId) { a.currentLoad += 1; break; }
        }
        issues.add(issue);
        return issue;
    }

    public List<IRIssue> getAgentIssues(List<IRIssue> issues, int agentId) {
        List<IRIssue> out = new ArrayList<>();
        for (IRIssue i : issues) if (i.assignedAgentId == agentId) out.add(i);
        return out;
    }

    public boolean transitionState(List<IRIssue> issues, int issueId, String newState) {
        for (IRIssue issue : issues) {
            if (issue.id != issueId) continue;
            String old = issue.state;
            boolean valid =
                (old.equals("OPEN") && newState.equals("IN_PROGRESS")) ||
                (old.equals("IN_PROGRESS") && newState.equals("RESOLVED")) ||
                (old.equals("RESOLVED") && newState.equals("CLOSED"));
            if (!valid) return false;
            issue.state = newState;
            for (IssueObserver obs : observers) obs.onStateChange(issueId, old, newState);
            return true;
        }
        return false;
    }
}

public class Solution {
    private static List<IRAgent> agents = new ArrayList<>();
    private static List<IRIssue> issues = new ArrayList<>();
    private static int nextIssueId = 0;
    private static RoundRobinStrategy rr = new RoundRobinStrategy();
    private static IssueResolver ir = new IssueResolver(rr);
    private static List<String> log = new ArrayList<>();
    private static LoggingObserver logger = new LoggingObserver(log);

    static { ir.addObserver(logger); }

    private static int prioFromString(String s) {
        switch (s) {
            case "MEDIUM": return 1;
            case "HIGH": return 2;
            case "CRITICAL": return 3;
            case "LOW":
            default: return 0;
        }
    }

    public static void reset_service() {
        agents = new ArrayList<>();
        issues = new ArrayList<>();
        nextIssueId = 0;
        rr = new RoundRobinStrategy();
        ir = new IssueResolver(rr);
        log = new ArrayList<>();
        logger = new LoggingObserver(log);
        ir.addObserver(logger);
    }

    public static void ir_add_agent(int id, String name, String specialization) {
        agents.add(new IRAgent(id, name, specialization));
    }

    public static int ir_assign_issue_round_robin(String description, String category, String priority) {
        nextIssueId += 1;
        IRIssue issue = new IRIssue(nextIssueId, description, category, prioFromString(priority));
        return ir.assign(agents, issues, issue).assignedAgentId;
    }

    public static int ir_assign_issue_least_loaded(String description, String category, String priority) {
        nextIssueId += 1;
        IRIssue issue = new IRIssue(nextIssueId, description, category, prioFromString(priority));
        IssueResolver r = new IssueResolver(new LeastLoadedStrategy());
        return r.assign(agents, issues, issue).assignedAgentId;
    }

    public static int ir_assign_issue_specialist(String description, String category, String priority) {
        nextIssueId += 1;
        IRIssue issue = new IRIssue(nextIssueId, description, category, prioFromString(priority));
        IssueResolver r = new IssueResolver(new SpecialistStrategy());
        return r.assign(agents, issues, issue).assignedAgentId;
    }

    public static int ir_agent_issue_count(int agentId) {
        return ir.getAgentIssues(issues, agentId).size();
    }

    public static int ir_agent_load(int agentId) {
        for (IRAgent a : agents) if (a.id == agentId) return a.currentLoad;
        return -1;
    }

    public static boolean ir_transition(int issueId, String newState) {
        return ir.transitionState(issues, issueId, newState);
    }

    public static String ir_get_issue_state(int issueId) {
        for (IRIssue i : issues) if (i.id == issueId) return i.state;
        return "";
    }

    public static int ir_log_size() { return log.size(); }

    public static String ir_log_entry(int idx) {
        if (idx >= 0 && idx < log.size()) return log.get(idx);
        return "";
    }
}
