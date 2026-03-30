#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <queue>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

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

// ─── Fixed-Window Rate Limiter ──────────────────────────────────────────────
// TODO: Implement the allowRequest() and getRequestCount() methods

class FixedWindowLimiter : public RateLimiter {
    int maxRequests;
    int windowSizeSeconds;
    unordered_map<string, int> requestCounts;
    unordered_map<string, long> windowStarts;

    long getWindowStart(long timestamp) {
        return (timestamp / windowSizeSeconds) * windowSizeSeconds;
    }

public:
    FixedWindowLimiter(int maxReq, int windowSize)
        : maxRequests(maxReq), windowSizeSeconds(windowSize) {}

    bool allowRequest(const Request& req) override {
        // TODO: Check if we're in a new window (reset count if so)
        // TODO: If count >= maxRequests, return false
        // TODO: Increment count and return true
        return false;
    }

    int getRequestCount(const string& clientId) override {
        // TODO: Return the current request count for this client
        return 0;
    }
};

// ─── Global Entry Points ────────────────────────────────────────────────────

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

#ifndef RUNNING_TESTS
int main() {
    cout << "Rate Limiter — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
