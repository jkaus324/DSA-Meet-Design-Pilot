// Issue Resolver — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testLeastLoadedBasic() {
        try {
            List<Agent> agents = {
            {0, "Alice", 3, {}},
            {1, "Bob",   1, {}},
            {2, "Carol", 2, {}},
            };
            List<Issue> issues;
            var result = assign_least_loaded(agents, issues,
            {100, "Help", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            boolean pass = result.assignedAgentId == 1); // Bob has lowest load (1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLeastLoadedBasic");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLeastLoadedBasic (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLeastLoadedTiebreak() {
        try {
            List<Agent> agents = {
            {0, "Alice", 2, {}},
            {1, "Bob",   2, {}},
            {2, "Carol", 2, {}},
            };
            List<Issue> issues;
            var result = assign_least_loaded(agents, issues,
            {101, "Tie", Category.GENERAL, Priority.LOW, IssueState.OPEN, -1});
            boolean pass = result.assignedAgentId == 0); // Alice wins tie (lowest ID;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLeastLoadedTiebreak");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLeastLoadedTiebreak (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSpecialistMatch() {
        try {
            List<Agent> agents = {
            {0, "Alice", 0, {Category.BILLING}},
            {1, "Bob",   0, {Category.TECHNICAL}},
            {2, "Carol", 0, {Category.GENERAL}},
            };
            List<Issue> issues;
            var result = assign_by_specialist(agents, issues,
            {102, "Tech issue", Category.TECHNICAL, Priority.HIGH, IssueState.OPEN, -1});
            boolean pass = result.assignedAgentId == 1); // Bob specializes in TECHNICAL;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSpecialistMatch");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSpecialistMatch (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSpecialistFallback() {
        try {
            List<Agent> agents = {
            {0, "Alice", 3, {Category.BILLING}},
            {1, "Bob",   1, {Category.BILLING}},
            {2, "Carol", 2, {Category.BILLING}},
            };
            List<Issue> issues;
            var result = assign_by_specialist(agents, issues,
            {103, "Account issue", Category.ACCOUNT, Priority.MEDIUM, IssueState.OPEN, -1});
            // No ACCOUNT specialist — falls back to least-loaded: Bob (load=1)
            boolean pass = result.assignedAgentId == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSpecialistFallback");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSpecialistFallback (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSpecialistLeastLoaded() {
        try {
            List<Agent> agents = {
            {0, "Alice", 5, {Category.BILLING, Category.TECHNICAL}},
            {1, "Bob",   2, {Category.TECHNICAL}},
            {2, "Carol", 8, {Category.TECHNICAL, Category.GENERAL}},
            };
            List<Issue> issues;
            var result = assign_by_specialist(agents, issues,
            {104, "Server down", Category.TECHNICAL, Priority.CRITICAL, IssueState.OPEN, -1});
            boolean pass = result.assignedAgentId == 1); // Bob is least-loaded TECHNICAL specialist;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSpecialistLeastLoaded");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSpecialistLeastLoaded (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testLeastLoadedBasic()) passed++;
        total++; if (testLeastLoadedTiebreak()) passed++;
        total++; if (testSpecialistMatch()) passed++;
        total++; if (testSpecialistFallback()) passed++;
        total++; if (testSpecialistLeastLoaded()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
