import java.util.*;

// Data class (given).

// HINT: introduce an abstraction so new ranking rules don't change existing code.
public class Solution {
    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static void reset_service() {
        // TODO: write your solution
        // nothing to return
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static void init_limiter(int maxRequests, int windowSize) {
        // TODO: write your solution
        // nothing to return
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static boolean allow_request_simple(String clientId, int timestamp, String endpoint) {
        // TODO: write your solution
        return false;
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static int get_request_count(String clientId) {
        // TODO: write your solution
        return 0;
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static boolean allow_request_with_strategy_simple(String algorithm, String clientId, int timestamp, String endpoint) {
        // TODO: write your solution
        return false;
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static boolean allow_request_for_tier_str(String tier, String clientId, int timestamp, String endpoint) {
        // TODO: write your solution
        return false;
    }

}
