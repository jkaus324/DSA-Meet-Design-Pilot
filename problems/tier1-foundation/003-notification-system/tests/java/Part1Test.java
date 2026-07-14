// Notification System — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testNotifyNoThrow() {
        try {
            User u1 = Arrays.asList("user1", "u1@test.com", "+91-9000000001", {"email")};
            User u2 = Arrays.asList("user2", "u2@test.com", "+91-9000000002", {"sms", "push")};
            User u3 = Arrays.asList("user3", "u3@test.com", "+91-9000000003", {"email", "sms")};
            // notify should not throw
            notify("Order shipped", {u1, u2, u3});
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testNotifyNoThrow");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNotifyNoThrow (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyUserList() {
        try {
            // If we can capture output, verify only subscribed channels receive
            // For simplicity, we test that an empty subscriber list causes no crash
            List<User> empty;
            notify("Event", empty);
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyUserList");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyUserList (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultiChannelUser() {
        try {
            User u = Arrays.asList("u1", "u1@test.com", "+1-555-0001", {"email", "sms", "push")};
            // Should not throw even with multiple channels
            notify("Flash sale", {u});
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultiChannelUser");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultiChannelUser (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testNotifyNoThrow()) passed++;
        total++; if (testEmptyUserList()) passed++;
        total++; if (testMultiChannelUser()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
