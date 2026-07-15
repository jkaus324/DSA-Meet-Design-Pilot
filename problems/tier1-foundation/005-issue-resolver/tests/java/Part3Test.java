// Issue Resolver — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testValidTransitionOpenToInprogress() {
        try {
            List<Issue> issues = {
            {200, "Bug", Category.TECHNICAL, Priority.HIGH, IssueState.OPEN, 0},
            };
            List<String> notifications;
            boolean ok = transition_issue(issues, 200, IssueState.IN_PROGRESS, notifications);
            boolean pass = ok == true
                && issues[0].state == IssueState.IN_PROGRESS
                && notifications.size() == 1
                && notifications[0] == "Issue 200: OPEN . IN_PROGRESS";
            System.out.println((pass ? "PASS" : "FAIL") + ": testValidTransitionOpenToInprogress");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testValidTransitionOpenToInprogress (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidTransitionRejected() {
        try {
            List<Issue> issues = {
            {201, "Bug", Category.TECHNICAL, Priority.HIGH, IssueState.OPEN, 0},
            };
            List<String> notifications;
            boolean ok = transition_issue(issues, 201, IssueState.CLOSED, notifications);
            boolean pass = ok == false
                && issues[0].state == IssueState.OPEN); // state unchanged
                && notifications.isEmpty()); // no notification;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidTransitionRejected");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidTransitionRejected (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFullLifecycle() {
        try {
            List<Issue> issues = {
            {202, "Payment fail", Category.BILLING, Priority.CRITICAL, IssueState.OPEN, 1},
            };
            List<String> notifications;
            boolean pass = transition_issue(issues, 202, IssueState.IN_PROGRESS, notifications)
                && transition_issue(issues, 202, IssueState.RESOLVED, notifications)
                && transition_issue(issues, 202, IssueState.CLOSED, notifications)
                && issues[0].state == IssueState.CLOSED
                && notifications.size() == 3;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFullLifecycle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFullLifecycle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPriorityHighestFirst() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            {1, "Bob",   0, {}},
            };
            List<Issue> issues = {
            {300, "Low prio",  Category.GENERAL, Priority.LOW,      IssueState.OPEN, -1},
            {301, "Critical",  Category.BILLING, Priority.CRITICAL, IssueState.OPEN, -1},
            {302, "Medium",    Category.TECHNICAL, Priority.MEDIUM, IssueState.OPEN, -1},
            };
            var first = assign_next_priority(agents, issues);
            boolean pass = first.id == 301); // CRITICAL picked first
                && first.assignedAgentId != -1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPriorityHighestFirst");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPriorityHighestFirst (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPriorityTiebreakById() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            };
            List<Issue> issues = {
            {403, "Third",  Category.GENERAL, Priority.HIGH, IssueState.OPEN, -1},
            {401, "First",  Category.GENERAL, Priority.HIGH, IssueState.OPEN, -1},
            {402, "Second", Category.GENERAL, Priority.HIGH, IssueState.OPEN, -1},
            };
            var first = assign_next_priority(agents, issues);
            boolean pass = first.id == 401); // lowest ID among HIGH priority;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPriorityTiebreakById");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPriorityTiebreakById (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPrioritySkipsAssigned() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {}},
            };
            List<Issue> issues = {
            {500, "Assigned",   Category.GENERAL, Priority.CRITICAL, IssueState.OPEN, 1},  // already assigned
            {501, "Unassigned", Category.GENERAL, Priority.LOW,      IssueState.OPEN, -1}, // unassigned
            };
            var result = assign_next_priority(agents, issues);
            boolean pass = result.id == 501); // skips 500 (already assigned;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPrioritySkipsAssigned");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPrioritySkipsAssigned (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testValidTransitionOpenToInprogress()) passed++;
        total++; if (testInvalidTransitionRejected()) passed++;
        total++; if (testFullLifecycle()) passed++;
        total++; if (testPriorityHighestFirst()) passed++;
        total++; if (testPriorityTiebreakById()) passed++;
        total++; if (testPrioritySkipsAssigned()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
