// Notification System — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testCriticalReachesAll() {
        try {
            User u1 = Arrays.asList("u1", "u1@test.com", "+1-555-0001", {"email")};
            Map<String, String> prefs = {Arrays.asList("*", "critical")};
            // Should not throw
            notify("CRITICAL: System down", "critical", {u1}, prefs);
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testCriticalReachesAll");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCriticalReachesAll (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPromotionalFiltered() {
        try {
            User u1 = Arrays.asList("u1", "u1@test.com", "+1-555-0001", {"email")};
            Map<String, String> prefs = {Arrays.asList("*", "info")};
            // Promotional should be filtered out — no exception expected
            notify("50% off sale!", "promotional", {u1}, prefs);
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testPromotionalFiltered");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPromotionalFiltered (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyPrefsAllowAll() {
        try {
            User u1 = Arrays.asList("u1", "u1@test.com", "+1-555-0001", {"email")};
            Map<String, String> emptyPrefs;
            notify("Informational update", "info", {u1}, emptyPrefs);
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyPrefsAllowAll");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyPrefsAllowAll (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testCriticalReachesAll()) passed++;
        total++; if (testPromotionalFiltered()) passed++;
        total++; if (testEmptyPrefsAllowAll()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
