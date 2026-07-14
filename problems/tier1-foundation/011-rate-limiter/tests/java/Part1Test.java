// Rate Limiter — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testAllowWithinLimit() {
        try {
            init_limiter(3, 60); // 3 requests per 60 seconds
            Request r1Arrays.asList("client_A", 1000, "/api/search");
            Request r2Arrays.asList("client_A", 1001, "/api/search");
            Request r3Arrays.asList("client_A", 1002, "/api/search");
            boolean pass = allow_request(r1) == true
                && allow_request(r2) == true
                && allow_request(r3) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAllowWithinLimit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAllowWithinLimit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRejectOverLimit() {
        try {
            init_limiter(3, 60);
            boolean pass = allow_request(Arrays.asList("client_B", 2000, "/api/pay")) == true
                && allow_request(Arrays.asList("client_B", 2001, "/api/pay")) == true
                && allow_request(Arrays.asList("client_B", 2002, "/api/pay")) == true
                && allow_request(Arrays.asList("client_B", 2003, "/api/pay")) == false); // 4th rejected
                && allow_request(Arrays.asList("client_B", 2004, "/api/pay")) == false); // 5th rejected;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRejectOverLimit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRejectOverLimit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIndependentClientLimits() {
        try {
            init_limiter(2, 60);
            boolean pass = allow_request(Arrays.asList("alice", 3000, "/api/x")) == true
                && allow_request(Arrays.asList("alice", 3001, "/api/x")) == true
                && allow_request(Arrays.asList("alice", 3002, "/api/x")) == false); // alice exhausted
                && allow_request(Arrays.asList("bob",   3003, "/api/x")) == true);  // bob still has quota
                && allow_request(Arrays.asList("bob",   3004, "/api/x")) == true
                && allow_request(Arrays.asList("bob",   3005, "/api/x")) == false); // bob exhausted;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIndependentClientLimits");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIndependentClientLimits (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testWindowReset() {
        try {
            init_limiter(2, 60);
            // New window starts at 1060
            boolean pass = allow_request(Arrays.asList("client_C", 1000, "/api/y")) == true
                && allow_request(Arrays.asList("client_C", 1020, "/api/y")) == true
                && allow_request(Arrays.asList("client_C", 1040, "/api/y")) == false); // limit hit in window [960, 1020
                && allow_request(Arrays.asList("client_C", 1060, "/api/y")) == true);  // new window
                && allow_request(Arrays.asList("client_C", 1080, "/api/y")) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testWindowReset");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testWindowReset (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testGetRequestCount() {
        try {
            init_limiter(5, 60);
            allow_request(Arrays.asList("new_client", 5000, "/api/z"));
            allow_request(Arrays.asList("new_client", 5001, "/api/z"));
            boolean pass = get_request_count("new_client") == 0
                && get_request_count("new_client") == 1
                && get_request_count("new_client") == 2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testGetRequestCount");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testGetRequestCount (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testAllowWithinLimit()) passed++;
        total++; if (testRejectOverLimit()) passed++;
        total++; if (testIndependentClientLimits()) passed++;
        total++; if (testWindowReset()) passed++;
        total++; if (testGetRequestCount()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
