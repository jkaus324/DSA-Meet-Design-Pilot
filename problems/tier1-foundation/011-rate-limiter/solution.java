// Rate Limiter — Solution (Java, Strategy + Factory)
import java.util.*;

class RLRequest {
    public String clientId;
    public long timestamp;
    public String endpoint;
    public RLRequest(String clientId, long timestamp, String endpoint) {
        this.clientId = clientId; this.timestamp = timestamp; this.endpoint = endpoint;
    }
}

interface RateLimiter {
    boolean allowRequest(RLRequest req);
    int getRequestCount(String clientId);
}

class FixedWindowLimiter implements RateLimiter {
    private final int maxRequests;
    private final long windowSizeSeconds;
    private final Map<String, Integer> counts = new HashMap<>();
    private final Map<String, Long> starts = new HashMap<>();

    public FixedWindowLimiter(int maxRequests, long windowSizeSeconds) {
        this.maxRequests = maxRequests; this.windowSizeSeconds = windowSizeSeconds;
    }

    public boolean allowRequest(RLRequest req) {
        Long st = starts.get(req.clientId);
        if (st == null || req.timestamp >= st + windowSizeSeconds) {
            starts.put(req.clientId, req.timestamp);
            counts.put(req.clientId, 0);
        }
        int cnt = counts.getOrDefault(req.clientId, 0);
        if (cnt >= maxRequests) return false;
        counts.put(req.clientId, cnt + 1);
        return true;
    }

    public int getRequestCount(String clientId) {
        return counts.getOrDefault(clientId, 0);
    }
}

class SlidingWindowLimiter implements RateLimiter {
    private final int maxRequests;
    private final long windowSizeSeconds;
    private final Map<String, Deque<Long>> queues = new HashMap<>();

    public SlidingWindowLimiter(int maxRequests, long windowSizeSeconds) {
        this.maxRequests = maxRequests; this.windowSizeSeconds = windowSizeSeconds;
    }

    public boolean allowRequest(RLRequest req) {
        Deque<Long> q = queues.computeIfAbsent(req.clientId, k -> new ArrayDeque<>());
        while (!q.isEmpty() && q.peekFirst() <= req.timestamp - windowSizeSeconds) {
            q.pollFirst();
        }
        if (q.size() >= maxRequests) return false;
        q.addLast(req.timestamp);
        return true;
    }

    public int getRequestCount(String clientId) {
        Deque<Long> q = queues.get(clientId);
        return q == null ? 0 : q.size();
    }
}

class TokenBucketLimiter implements RateLimiter {
    private final int maxTokens;
    private final double refillRate;
    private final Map<String, Double> tokens = new HashMap<>();
    private final Map<String, Long> lastRefill = new HashMap<>();

    public TokenBucketLimiter(int maxTokens, long windowSize) {
        this.maxTokens = maxTokens;
        this.refillRate = (double) maxTokens / windowSize;
    }

    public boolean allowRequest(RLRequest req) {
        if (!tokens.containsKey(req.clientId)) {
            tokens.put(req.clientId, (double) maxTokens);
            lastRefill.put(req.clientId, req.timestamp);
        }
        long elapsed = req.timestamp - lastRefill.get(req.clientId);
        double t = Math.min(maxTokens, tokens.get(req.clientId) + elapsed * refillRate);
        tokens.put(req.clientId, t);
        lastRefill.put(req.clientId, req.timestamp);
        if (t < 1.0) return false;
        tokens.put(req.clientId, t - 1.0);
        return true;
    }

    public int getRequestCount(String clientId) {
        if (!tokens.containsKey(clientId)) return 0;
        return maxTokens - (int) (double) tokens.get(clientId);
    }
}

public class Solution {
    private static RateLimiter limiter = null;
    private static Map<String, RateLimiter> strategies = new HashMap<>();
    private static Map<String, RateLimiter> tiers = new HashMap<>();

    private static final Map<String, Integer> TIER_LIMITS = new HashMap<>();
    static {
        TIER_LIMITS.put("FREE", 10);
        TIER_LIMITS.put("PRO", 100);
        TIER_LIMITS.put("ENTERPRISE", 1000);
    }

    private static RateLimiter createLimiter(String algorithm, int maxRequests, long windowSize) {
        switch (algorithm) {
            case "fixed-window": return new FixedWindowLimiter(maxRequests, windowSize);
            case "sliding-window": return new SlidingWindowLimiter(maxRequests, windowSize);
            case "token-bucket": return new TokenBucketLimiter(maxRequests, windowSize);
            default: return null;
        }
    }

    public static void reset_service() {
        limiter = null;
        strategies = new HashMap<>();
        tiers = new HashMap<>();
    }

    public static void init_limiter(int maxRequests, int windowSize) {
        limiter = new FixedWindowLimiter(maxRequests, windowSize);
    }

    public static boolean allow_request_simple(String clientId, int timestamp, String endpoint) {
        if (limiter == null) return false;
        return limiter.allowRequest(new RLRequest(clientId, timestamp, endpoint));
    }

    public static int get_request_count(String clientId) {
        if (limiter == null) return 0;
        return limiter.getRequestCount(clientId);
    }

    public static boolean allow_request_with_strategy_simple(String algorithm, String clientId, int timestamp, String endpoint) {
        RateLimiter rl = strategies.get(algorithm);
        if (rl == null) {
            rl = createLimiter(algorithm, 100, 60);
            if (rl == null) return false;
            strategies.put(algorithm, rl);
        }
        return rl.allowRequest(new RLRequest(clientId, timestamp, endpoint));
    }

    public static boolean allow_request_for_tier_str(String tier, String clientId, int timestamp, String endpoint) {
        RateLimiter rl = tiers.get(tier);
        if (rl == null) {
            int limit = TIER_LIMITS.getOrDefault(tier, 10);
            rl = new SlidingWindowLimiter(limit, limit + 1);
            tiers.put(tier, rl);
        }
        return rl.allowRequest(new RLRequest(clientId, timestamp, endpoint));
    }
}
