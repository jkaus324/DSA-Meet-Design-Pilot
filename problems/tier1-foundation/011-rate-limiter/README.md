# Problem 011 — API Rate Limiter

**Tier:** 1 (Foundation) | **Pattern:** Strategy + Factory | **DSA:** Queue + HashMap
**Companies:** Amazon, Razorpay, Uber | **Time:** 45 minutes

---

## Problem Statement

You're building the API gateway for a high-traffic platform. Every incoming request must be checked against a **rate limiter** before it reaches the backend. If a client has exceeded their allowed number of requests in a given time window, the request is **rejected**.

Different endpoints may use different rate-limiting algorithms, and different user tiers (free, pro, enterprise) have different rate limits.

**Your task:** Design and implement a `RateLimiter` system that can enforce request limits using multiple algorithms and adapt to user tiers — with new algorithms addable without modifying existing code.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What varies here? The *rate-limiting algorithm* varies (fixed-window, sliding-window, token-bucket). The *enforcement logic* stays the same.
2. If you used `if-else` inside `allow_request()` to pick the algorithm, what happens when a 4th algorithm is added? You modify existing code — violating Open/Closed Principle.
3. How does the Strategy pattern solve this? Each algorithm becomes a separate class implementing a common interface.
4. How does a Factory help? It creates the right limiter based on endpoint configuration, so the caller doesn't need to know which algorithm to use.

**The key insight:** A rate-limiting algorithm is a strategy. The factory encapsulates the creation decision. Together, they give you a system where adding a new algorithm or a new user tier requires zero changes to existing code.

---

## Data Structures

```cpp
struct Request {
    std::string clientId;   // "user_123", "api_key_abc"
    long timestamp;         // epoch seconds
    std::string endpoint;   // "/api/payments", "/api/search"
};

enum class UserTier {
    FREE,        // 10 requests per minute
    PRO,         // 100 requests per minute
    ENTERPRISE   // 1000 requests per minute
};
```

---

## Base Requirement — Fixed-window rate limiting

Implement a fixed-window rate limiter that tracks how many requests each client has made in the current time window. If a client exceeds the limit, reject the request.

**Fixed-window algorithm:**
- Divide time into fixed windows (e.g., 60-second windows)
- Count requests per client in the current window
- If count >= limit, reject
- When a new window starts, reset the count

**Entry points (tests will call these):**
```cpp
bool allow_request(const Request& req);
int get_request_count(const std::string& clientId);
```

**What to implement:**
```cpp
class RateLimiter {
public:
    virtual bool allowRequest(const Request& req) = 0;
    virtual int getRequestCount(const std::string& clientId) = 0;
    virtual ~RateLimiter() = default;
};

class FixedWindowLimiter : public RateLimiter {
    int maxRequests;        // max allowed per window
    int windowSizeSeconds;  // window duration
    // TODO: track per-client request counts and window start times
public:
    FixedWindowLimiter(int maxReq, int windowSize);
    bool allowRequest(const Request& req) override;
    int getRequestCount(const std::string& clientId) override;
};
```

---

## Extension 1 — Multiple algorithms per endpoint

The platform team wants different endpoints to use different rate-limiting algorithms:
- `/api/search` uses **sliding-window** (smoother traffic shaping)
- `/api/payments` uses **token-bucket** (allows short bursts)
- `/api/users` uses **fixed-window** (simple and sufficient)

**Sliding-window algorithm:**
- Count requests in a rolling window ending at the current timestamp
- If count >= limit, reject
- Uses a queue to track individual request timestamps per client

**Token-bucket algorithm:**
- Each client has a bucket of tokens (starts full)
- Each request consumes one token
- Tokens replenish at a fixed rate (e.g., 1 token per second)
- If bucket is empty, reject

**New entry points:**
```cpp
RateLimiter* create_limiter(const std::string& algorithm, int maxRequests, int windowSize);
bool allow_request_with_strategy(const std::string& algorithm, const Request& req);
```

**Design challenge:** How does the Factory pattern help here? The caller just says "I need a limiter for /api/search" and gets the right algorithm without knowing the details.

---

## Extension 2 — User tier-based rate limits

Different user tiers get different rate limits:

| Tier | Requests per Minute |
|------|-------------------|
| FREE | 10 |
| PRO | 100 |
| ENTERPRISE | 1000 |

The factory now takes a user tier and creates a limiter with the appropriate limits.

**New entry point:**
```cpp
bool allow_request_for_tier(UserTier tier, const Request& req);
```

**Design challenge:** How do you combine tier-based limits with per-endpoint algorithm selection? Can your Factory handle both dimensions (algorithm + tier)?

---

## Running Tests

```bash
./run-tests.sh 011-rate-limiter cpp
```
