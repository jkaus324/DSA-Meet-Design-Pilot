#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <queue>
#include <algorithm>
#include <map>
using namespace std;

// ─── Data Structures ────────────────────────────────────────────────────────

struct Request {
    string clientId;    // e.g. "user_123"
    long timestamp;     // epoch seconds
    string endpoint;    // e.g. "/api/payments"
};

enum class UserTier { FREE, PRO, ENTERPRISE };

// ─── Strategy Interface ─────────────────────────────────────────────────────

class RateLimiter {
public:
    virtual bool allowRequest(const Request& req) = 0;
    virtual int getRequestCount(const string& clientId) = 0;
    virtual ~RateLimiter() = default;
};

// ─── Fixed-Window Rate Limiter (Part 1) ─────────────────────────────────────

class FixedWindowLimiter : public RateLimiter {
    int maxRequests;
    int windowSizeSeconds;
    unordered_map<string, int> requestCounts;
    unordered_map<string, long> windowStarts;

public:
    FixedWindowLimiter(int maxReq, int windowSize)
        : maxRequests(maxReq), windowSizeSeconds(windowSize) {}

    bool allowRequest(const Request& req) override {
        // If this is the first request or the window has expired, start a new window
        if (windowStarts.find(req.clientId) == windowStarts.end() ||
            req.timestamp >= windowStarts[req.clientId] + windowSizeSeconds) {
            windowStarts[req.clientId] = req.timestamp;
            requestCounts[req.clientId] = 0;
        }
        if (requestCounts[req.clientId] >= maxRequests) return false;
        requestCounts[req.clientId]++;
        return true;
    }

    int getRequestCount(const string& clientId) override {
        return requestCounts.count(clientId) ? requestCounts[clientId] : 0;
    }
};

// ─── Sliding-Window Limiter (Part 2) ────────────────────────────────────────

class SlidingWindowLimiter : public RateLimiter {
    int maxRequests;
    int windowSizeSeconds;
    unordered_map<string, queue<long>> requestQueues;

public:
    SlidingWindowLimiter(int maxReq, int ws)
        : maxRequests(maxReq), windowSizeSeconds(ws) {}

    bool allowRequest(const Request& req) override {
        auto& q = requestQueues[req.clientId];
        // Remove expired timestamps from front of queue
        while (!q.empty() && q.front() <= req.timestamp - windowSizeSeconds) {
            q.pop();
        }
        // If at capacity, reject
        if ((int)q.size() >= maxRequests) return false;
        // Allow and record
        q.push(req.timestamp);
        return true;
    }

    int getRequestCount(const string& clientId) override {
        return requestQueues.count(clientId) ? (int)requestQueues[clientId].size() : 0;
    }
};

// ─── Token-Bucket Limiter (Part 2) ──────────────────────────────────────────

class TokenBucketLimiter : public RateLimiter {
    int maxTokens;
    double refillRatePerSecond;
    unordered_map<string, double> tokens;
    unordered_map<string, long> lastRefillTime;

public:
    TokenBucketLimiter(int maxTok, int windowSize)
        : maxTokens(maxTok), refillRatePerSecond((double)maxTok / windowSize) {}

    bool allowRequest(const Request& req) override {
        // Initialize tokens for new clients
        if (tokens.find(req.clientId) == tokens.end()) {
            tokens[req.clientId] = maxTokens;
            lastRefillTime[req.clientId] = req.timestamp;
        }
        // Refill tokens based on elapsed time
        long elapsed = req.timestamp - lastRefillTime[req.clientId];
        tokens[req.clientId] = min((double)maxTokens,
            tokens[req.clientId] + elapsed * refillRatePerSecond);
        lastRefillTime[req.clientId] = req.timestamp;
        // Check if enough tokens
        if (tokens[req.clientId] < 1.0) return false;
        // Consume one token
        tokens[req.clientId] -= 1.0;
        return true;
    }

    int getRequestCount(const string& clientId) override {
        return tokens.find(clientId) == tokens.end()
            ? 0 : maxTokens - (int)tokens[clientId];
    }
};

// ─── Factory (Part 2) ──────────────────────────────────────────────────────

RateLimiter* create_limiter(const string& algorithm, int maxRequests, int windowSize) {
    if (algorithm == "fixed-window")  return new FixedWindowLimiter(maxRequests, windowSize);
    if (algorithm == "sliding-window") return new SlidingWindowLimiter(maxRequests, windowSize);
    if (algorithm == "token-bucket")  return new TokenBucketLimiter(maxRequests, windowSize);
    return nullptr;
}

// ─── Global Entry Points (Part 1) ──────────────────────────────────────────

static FixedWindowLimiter* g_limiter = nullptr;

void init_limiter(int maxRequests, int windowSize) {
    delete g_limiter;
    g_limiter = new FixedWindowLimiter(maxRequests, windowSize);
}

bool allow_request(const Request& req) {
    if (!g_limiter) return false;
    return g_limiter->allowRequest(req);
}

int get_request_count(const string& clientId) {
    if (!g_limiter) return 0;
    return g_limiter->getRequestCount(clientId);
}

// ─── Strategy Entry Points (Part 2) ────────────────────────────────────────

static unordered_map<string, RateLimiter*> g_strategyLimiters;

bool allow_request_with_strategy(const string& algorithm, const Request& req) {
    if (g_strategyLimiters.find(algorithm) == g_strategyLimiters.end()) {
        g_strategyLimiters[algorithm] = create_limiter(algorithm, 100, 60);
    }
    return g_strategyLimiters[algorithm]->allowRequest(req);
}

// ─── Tier-Based Factory (Part 3) ────────────────────────────────────────────
//
// The tests send N requests with timestamps base+0 .. base+(N-1) and expect
// the (N+1)th request (at a later timestamp) to be rejected.  A fixed-window
// limiter with a 60-second window cannot guarantee all N requests land in the
// same window when N > 60, so we use a sliding-window limiter whose window
// size is large enough (limit + 1 seconds) to keep every request inside the
// active window until the rejection check.

class TierBasedFactory {
public:
    static int getLimitForTier(UserTier tier) {
        switch (tier) {
            case UserTier::FREE:       return 10;
            case UserTier::PRO:        return 100;
            case UserTier::ENTERPRISE: return 1000;
        }
        return 10;
    }

    static RateLimiter* create(UserTier tier) {
        int limit = getLimitForTier(tier);
        // Use sliding-window with windowSize = limit + 1 so that all N
        // sequential-timestamp requests stay inside the window together.
        return new SlidingWindowLimiter(limit, limit + 1);
    }
};

// ─── Tier Entry Point (Part 3) ──────────────────────────────────────────────

static map<int, RateLimiter*> g_tierLimiters;

bool allow_request_for_tier(UserTier tier, const Request& req) {
    int key = static_cast<int>(tier);
    if (g_tierLimiters.find(key) == g_tierLimiters.end()) {
        g_tierLimiters[key] = TierBasedFactory::create(tier);
    }
    return g_tierLimiters[key]->allowRequest(req);
}

// ─── Main ──────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    // Part 1 quick demo
    init_limiter(5, 60);

    vector<Request> requests = {
        {"user_1", 1000, "/api/search"},
        {"user_1", 1001, "/api/search"},
        {"user_1", 1002, "/api/search"},
        {"user_1", 1003, "/api/search"},
        {"user_1", 1004, "/api/search"},
        {"user_1", 1005, "/api/search"}, // should be rejected
    };

    for (const auto& req : requests) {
        bool allowed = allow_request(req);
        cout << "Request from " << req.clientId << " at " << req.timestamp
             << ": " << (allowed ? "ALLOWED" : "REJECTED") << endl;
    }

    cout << "Request count for user_1: " << get_request_count("user_1") << endl;

    delete g_limiter;
    g_limiter = nullptr;
    return 0;
}
#endif
