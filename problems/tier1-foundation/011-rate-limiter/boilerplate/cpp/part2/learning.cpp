#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <queue>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct Request {
    string clientId;
    long timestamp;
    string endpoint;
};

enum class UserTier { FREE, PRO, ENTERPRISE };

// ─── Strategy Interface ─────────────────────────────────────────────────────

class RateLimiter {
public:
    virtual bool allowRequest(const Request& req) = 0;
    virtual int getRequestCount(const string& clientId) = 0;
    virtual ~RateLimiter() = default;
};

// ─── Fixed-Window (complete from Part 1) ────────────────────────────────────

class FixedWindowLimiter : public RateLimiter {
    int maxRequests;
    int windowSizeSeconds;
    unordered_map<string, int> requestCounts;
    unordered_map<string, long> windowStarts;
    long getWindowStart(long ts) { return (ts / windowSizeSeconds) * windowSizeSeconds; }
public:
    FixedWindowLimiter(int maxReq, int ws) : maxRequests(maxReq), windowSizeSeconds(ws) {}
    bool allowRequest(const Request& req) override {
        long ws = getWindowStart(req.timestamp);
        if (windowStarts[req.clientId] != ws) { windowStarts[req.clientId] = ws; requestCounts[req.clientId] = 0; }
        if (requestCounts[req.clientId] >= maxRequests) return false;
        requestCounts[req.clientId]++;
        return true;
    }
    int getRequestCount(const string& cid) override { return requestCounts.count(cid) ? requestCounts[cid] : 0; }
};

// ─── Sliding-Window Limiter ─────────────────────────────────────────────────
// TODO: Implement allowRequest() and getRequestCount()

class SlidingWindowLimiter : public RateLimiter {
    int maxRequests;
    int windowSizeSeconds;
    unordered_map<string, queue<long>> requestQueues;
public:
    SlidingWindowLimiter(int maxReq, int ws) : maxRequests(maxReq), windowSizeSeconds(ws) {}

    bool allowRequest(const Request& req) override {
        // TODO: Remove expired timestamps from the front of the queue
        // TODO: If queue size >= maxRequests, return false
        // TODO: Push current timestamp and return true
        return false;
    }

    int getRequestCount(const string& clientId) override {
        // TODO: Return current queue size for this client
        return 0;
    }
};

// ─── Token-Bucket Limiter ───────────────────────────────────────────────────
// TODO: Implement allowRequest() and getRequestCount()

class TokenBucketLimiter : public RateLimiter {
    int maxTokens;
    double refillRatePerSecond;
    unordered_map<string, double> tokens;
    unordered_map<string, long> lastRefillTime;
public:
    TokenBucketLimiter(int maxTok, int windowSize)
        : maxTokens(maxTok), refillRatePerSecond((double)maxTok / windowSize) {}

    bool allowRequest(const Request& req) override {
        // TODO: Initialize tokens for new clients
        // TODO: Refill tokens based on elapsed time
        // TODO: If tokens < 1, return false
        // TODO: Consume 1 token and return true
        return false;
    }

    int getRequestCount(const string& clientId) override {
        // TODO: Return tokens consumed (maxTokens - remaining)
        return 0;
    }
};

// ─── Factory ────────────────────────────────────────────────────────────────

RateLimiter* create_limiter(const string& algorithm, int maxRequests, int windowSize) {
    if (algorithm == "fixed-window") return new FixedWindowLimiter(maxRequests, windowSize);
    if (algorithm == "sliding-window") return new SlidingWindowLimiter(maxRequests, windowSize);
    if (algorithm == "token-bucket") return new TokenBucketLimiter(maxRequests, windowSize);
    return nullptr;
}

// ─── Global Entry Points ────────────────────────────────────────────────────

static FixedWindowLimiter* g_limiter = nullptr;
static unordered_map<string, RateLimiter*> g_strategyLimiters;

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

bool allow_request_with_strategy(const string& algorithm, const Request& req) {
    if (g_strategyLimiters.find(algorithm) == g_strategyLimiters.end()) {
        g_strategyLimiters[algorithm] = create_limiter(algorithm, 100, 60);
    }
    return g_strategyLimiters[algorithm]->allowRequest(req);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Multiple algorithms — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
