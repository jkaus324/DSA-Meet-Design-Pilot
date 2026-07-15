// Rate Limiter — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testFreeTierLimit() {
        try {
            // Send 10 requests — all should pass
            for (int i = 0; i < 10; i++) {
            }
            // 11th should be rejected
            boolean pass = allow_request_for_tier(UserTier.FREE, Arrays.asList("free_user", 5000 + i, "/api/data")) == true
                && allow_request_for_tier(UserTier.FREE, Arrays.asList("free_user", 5010, "/api/data")) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFreeTierLimit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFreeTierLimit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testProTierLimit() {
        try {
            for (int i = 0; i < 100; i++) {
            }
            // 101st should be rejected
            boolean pass = allow_request_for_tier(UserTier.PRO, Arrays.asList("pro_user", 6000 + i, "/api/data")) == true
                && allow_request_for_tier(UserTier.PRO, Arrays.asList("pro_user", 6100, "/api/data")) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testProTierLimit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testProTierLimit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEnterpriseTierLimit() {
        try {
            for (int i = 0; i < 1000; i++) {
            }
            // 1001st should be rejected
            boolean pass = allow_request_for_tier(UserTier.ENTERPRISE, Arrays.asList("enterprise_user", 7000 + i, "/api/data")) == true
                && allow_request_for_tier(UserTier.ENTERPRISE, Arrays.asList("enterprise_user", 8000, "/api/data")) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEnterpriseTierLimit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEnterpriseTierLimit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTierIndependence() {
        try {
            // Free user hits limit at 10
            for (int i = 0; i < 10; i++) {
            allow_request_for_tier(UserTier.FREE,
            Arrays.asList("free_user_2", 9000 + i, "/api/x"));
            }
            // Pro user still has quota
            boolean pass = allow_request_for_tier(UserTier.FREE, Arrays.asList("free_user_2", 9010, "/api/x")) == false
                && allow_request_for_tier(UserTier.PRO, Arrays.asList("pro_user_2", 9010, "/api/x")) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTierIndependence");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTierIndependence (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testFreeTierLimit()) passed++;
        total++; if (testProTierLimit()) passed++;
        total++; if (testEnterpriseTierLimit()) passed++;
        total++; if (testTierIndependence()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
