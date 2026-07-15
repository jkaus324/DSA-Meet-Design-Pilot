// Issue Resolver — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testRoundRobinAssignment() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            {1, "Bob",   0, {}},
            {2, "Carol", 0, {}},
            };
            List<Issue> issues;
            var i1 = assign_issue(agents, issues, {1, "Issue A", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            var i2 = assign_issue(agents, issues, {2, "Issue B", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            var i3 = assign_issue(agents, issues, {3, "Issue C", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            boolean pass = i1.assignedAgentId == 0); // first agent
                && i2.assignedAgentId == 1); // second agent
                && i3.assignedAgentId == 2); // third agent;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRoundRobinAssignment");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRoundRobinAssignment (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRoundRobinWrap() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            {1, "Bob",   0, {}},
            };
            List<Issue> issues;
            // Reset global round-robin by creating fresh resolver
            RoundRobinStrategy rr = new RoundRobinStrategy();
            IssueResolver resolver = new IssueResolver( rr);
            var i1 = resolver.assign(agents, issues, {10, "A", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            var i2 = resolver.assign(agents, issues, {11, "B", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            var i3 = resolver.assign(agents, issues, {12, "C", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            boolean pass = i1.assignedAgentId == 0
                && i2.assignedAgentId == 1
                && i3.assignedAgentId == 0); // wraps back to first;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRoundRobinWrap");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRoundRobinWrap (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testGetAgentIssues() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            {1, "Bob",   0, {}},
            };
            List<Issue> issues;
            RoundRobinStrategy rr = new RoundRobinStrategy();
            IssueResolver resolver = new IssueResolver( rr);
            resolver.assign(agents, issues, {20, "X", Category.BILLING, Priority.HIGH, IssueState.OPEN, -1});
            resolver.assign(agents, issues, {21, "Y", Category.TECHNICAL, Priority.LOW, IssueState.OPEN, -1});
            resolver.assign(agents, issues, {22, "Z", Category.GENERAL, Priority.MEDIUM, IssueState.OPEN, -1});
            var aliceIssues = resolver.getAgentIssues(issues, 0);
            var bobIssues = resolver.getAgentIssues(issues, 1);
            boolean pass = aliceIssues.size() == 2); // issues 20, 22
                && bobIssues.size() == 1);   // issue 21
                && aliceIssues[0].id == 20
                && aliceIssues[1].id == 22
                && bobIssues[0].id == 21;
            System.out.println((pass ? "PASS" : "FAIL") + ": testGetAgentIssues");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testGetAgentIssues (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAgentLoadIncrement() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            {1, "Bob",   0, {}},
            };
            List<Issue> issues;
            RoundRobinStrategy rr = new RoundRobinStrategy();
            IssueResolver resolver = new IssueResolver( rr);
            resolver.assign(agents, issues, {30, "A", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            resolver.assign(agents, issues, {31, "B", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            resolver.assign(agents, issues, {32, "C", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            boolean pass = agents[0].currentLoad == 2); // got issues 30, 32
                && agents[1].currentLoad == 1); // got issue 31;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAgentLoadIncrement");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAgentLoadIncrement (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyIssues() {
        try {
            List<Issue> issues;
            RoundRobinStrategy rr = new RoundRobinStrategy();
            IssueResolver resolver = new IssueResolver( rr);
            var result = resolver.getAgentIssues(issues, 0);
            boolean pass = result.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyIssues");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyIssues (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testRoundRobinAssignment()) passed++;
        total++; if (testRoundRobinWrap()) passed++;
        total++; if (testGetAgentIssues()) passed++;
        total++; if (testAgentLoadIncrement()) passed++;
        total++; if (testEmptyIssues()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
