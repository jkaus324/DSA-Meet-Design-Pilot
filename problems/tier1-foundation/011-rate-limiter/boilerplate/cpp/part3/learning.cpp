#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <queue>
#include <algorithm>
#include <map>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

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

// ─── All Algorithms (complete from Parts 1-2) ───────────────────────────────

class FixedWindowLimiter : public RateLimiter {
    int maxRequests; int windowSizeSeconds;
    unordered_map<string, int> requestCounts;
    unordered_map<string, long> windowStarts;
    long getWindowStart(long ts) { return (ts / windowSizeSeconds) * windowSizeSeconds; }
public:
    FixedWindowLimiter(int maxReq, int ws) : maxRequests(maxReq), windowSizeSeconds(ws) {}
    bool allowRequest(const Request& req) override {
        long ws = getWindowStart(req.timestamp);
        if (windowStarts[req.clientId] != ws) { windowStarts[req.clientId] = ws; requestCounts[req.clientId] = 0; }
        if (requestCounts[req.clientId] >= maxRequests) return false;
        requestCounts[req.clientId]++; return true;
    }
    int getRequestCount(const string& cid) override { return requestCounts.count(cid) ? requestCounts[cid] : 0; }
};

class SlidingWindowLimiter : public RateLimiter {
    int maxRequests; int windowSizeSeconds;
    unordered_map<string, queue<long>> requestQueues;
public:
    SlidingWindowLimiter(int maxReq, int ws) : maxRequests(maxReq), windowSizeSeconds(ws) {}
    bool allowRequest(const Request& req) override {
        auto& q = requestQueues[req.clientId];
        while (!q.empty() && q.front() <= req.timestamp - windowSizeSeconds) q.pop();
        if ((int)q.size() >= maxRequests) return false;
        q.push(req.timestamp); return true;
    }
    int getRequestCount(const string& cid) override { return requestQueues.count(cid) ? (int)requestQueues[cid].size() : 0; }
};

class TokenBucketLimiter : public RateLimiter {
    int maxTokens; double refillRate;
    unordered_map<string, double> tokens;
    unordered_map<string, long> lastRefill;
public:
    TokenBucketLimiter(int maxTok, int ws) : maxTokens(maxTok), refillRate((double)maxTok / ws) {}
    bool allowRequest(const Request& req) override {
        if (tokens.find(req.clientId) == tokens.end()) { tokens[req.clientId] = maxTokens; lastRefill[req.clientId] = req.timestamp; }
        long elapsed = req.timestamp - lastRefill[req.clientId];
        tokens[req.clientId] = min((double)maxTokens, tokens[req.clientId] + elapsed * refillRate);
        lastRefill[req.clientId] = req.timestamp;
        if (tokens[req.clientId] < 1.0) return false;
        tokens[req.clientId] -= 1.0; return true;
    }
    int getRequestCount(const string& cid) override { return tokens.find(cid) == tokens.end() ? 0 : maxTokens - (int)tokens[cid]; }
};

// ─── Factory ────────────────────────────────────────────────────────────────

RateLimiter* create_limiter(const string& algo, int maxReq, int ws) {
    if (algo == "fixed-window") return new FixedWindowLimiter(maxReq, ws);
    if (algo == "sliding-window") return new SlidingWindowLimiter(maxReq, ws);
    if (algo == "token-bucket") return new TokenBucketLimiter(maxReq, ws);
    return nullptr;
}

// ─── Global Entry Points (Parts 1-2) ────────────────────────────────────────

static FixedWindowLimiter* g_limiter = nullptr;
static unordered_map<string, RateLimiter*> g_strategyLimiters;

void init_limiter(int maxRequests, int windowSize) { delete g_limiter; g_limiter = new FixedWindowLimiter(maxRequests, windowSize); }
bool allow_request(const Request& req) { return g_limiter ? g_limiter->allowRequest(req) : false; }
int get_request_count(const string& cid) { return g_limiter ? g_limiter->getRequestCount(cid) : 0; }
bool allow_request_with_strategy(const string& algo, const Request& req) {
    if (g_strategyLimiters.find(algo) == g_strategyLimiters.end()) g_strategyLimiters[algo] = create_limiter(algo, 100, 60);
    return g_strategyLimiters[algo]->allowRequest(req);
}

// ─── Tier-Based Factory ─────────────────────────────────────────────────────

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

    static RateLimiter* create(UserTier tier, const string& algo = "fixed-window") {
        return create_limiter(algo, getLimitForTier(tier), 60);
    }
};

// ─── Tier Entry Point ───────────────────────────────────────────────────────
// TODO: Implement allow_request_for_tier()
// The limiter for each tier should persist across calls (use a static map)

static map<int, RateLimiter*> g_tierLimiters;

bool allow_request_for_tier(UserTier tier, const Request& req) {
    int key = static_cast<int>(tier);
    if (g_tierLimiters.find(key) == g_tierLimiters.end()) {
        g_tierLimiters[key] = TierBasedFactory::create(tier);
    }
    // TODO: Use the tier-specific limiter to allow or reject the request
    return false;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Tier-based rate limits — implement the TODO above, then run tests." << endl;
    return 0;
}
#endif
