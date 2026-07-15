// Rate Limiter — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testFactoryFixedWindow() {
        try {
            RateLimiter limiter = create_limiter("fixed-window", 3, 60);
            // limiter = null; // GC handles deallocation
            boolean pass = limiter != null
                && limiter.allowRequest(Arrays.asList("user_1", 1000, "/api/a")) == true
                && limiter.allowRequest(Arrays.asList("user_1", 1001, "/api/a")) == true
                && limiter.allowRequest(Arrays.asList("user_1", 1002, "/api/a")) == true
                && limiter.allowRequest(Arrays.asList("user_1", 1003, "/api/a")) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFactoryFixedWindow");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFactoryFixedWindow (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSlidingWindow() {
        try {
            RateLimiter limiter = create_limiter("sliding-window", 3, 60);
            // After first request expires from window
            // limiter = null; // GC handles deallocation
            boolean pass = limiter != null
                && limiter.allowRequest(Arrays.asList("user_2", 1000, "/api/b")) == true
                && limiter.allowRequest(Arrays.asList("user_2", 1030, "/api/b")) == true
                && limiter.allowRequest(Arrays.asList("user_2", 1050, "/api/b")) == true
                && limiter.allowRequest(Arrays.asList("user_2", 1055, "/api/b")) == false); // 4th in 60s window
                && limiter.allowRequest(Arrays.asList("user_2", 1061, "/api/b")) == true); // 1000 is now outside [1001,1061];
            System.out.println((pass ? "PASS" : "FAIL") + ": testSlidingWindow");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSlidingWindow (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTokenBucket() {
        try {
            RateLimiter limiter = create_limiter("token-bucket", 3, 60);
            // Burst: use all 3 tokens immediately
            // Wait for tokens to refill (rate = 3/60 = 0.05/sec, need 20sec for 1 token)
            // limiter = null; // GC handles deallocation
            boolean pass = limiter != null
                && limiter.allowRequest(Arrays.asList("user_3", 1000, "/api/c")) == true
                && limiter.allowRequest(Arrays.asList("user_3", 1000, "/api/c")) == true
                && limiter.allowRequest(Arrays.asList("user_3", 1000, "/api/c")) == true
                && limiter.allowRequest(Arrays.asList("user_3", 1000, "/api/c")) == false); // empty
                && limiter.allowRequest(Arrays.asList("user_3", 1020, "/api/c")) == true); // 1 token refilled
                && limiter.allowRequest(Arrays.asList("user_3", 1020, "/api/c")) == false); // empty again;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTokenBucket");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTokenBucket (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFactoryUnknownAlgorithm() {
        try {
            RateLimiter limiter = create_limiter("unknown-algo", 10, 60);
            boolean pass = limiter == null;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFactoryUnknownAlgorithm");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFactoryUnknownAlgorithm (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAllowRequestWithStrategy() {
        try {
            // Reset state for strategy-based requests
            boolean pass = allow_request_with_strategy("fixed-window", Arrays.asList("user_4", 2000, "/api/d")) == true
                && allow_request_with_strategy("fixed-window", Arrays.asList("user_4", 2001, "/api/d")) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAllowRequestWithStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAllowRequestWithStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testFactoryFixedWindow()) passed++;
        total++; if (testSlidingWindow()) passed++;
        total++; if (testTokenBucket()) passed++;
        total++; if (testFactoryUnknownAlgorithm()) passed++;
        total++; if (testAllowRequestWithStrategy()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
