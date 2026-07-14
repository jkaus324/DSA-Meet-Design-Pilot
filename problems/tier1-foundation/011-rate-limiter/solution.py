"""Rate limiter — Strategy + Factory reference solution (Python)."""

from abc import ABC, abstractmethod
from collections import deque


class Request:
    def __init__(self, clientId, timestamp, endpoint):
        self.clientId = clientId
        self.timestamp = timestamp
        self.endpoint = endpoint


class UserTier:
    FREE = "FREE"
    PRO = "PRO"
    ENTERPRISE = "ENTERPRISE"


class RateLimiter(ABC):
    @abstractmethod
    def allowRequest(self, req):
        ...

    @abstractmethod
    def getRequestCount(self, clientId):
        ...


class FixedWindowLimiter(RateLimiter):
    def __init__(self, maxRequests, windowSizeSeconds):
        self.maxRequests = maxRequests
        self.windowSizeSeconds = windowSizeSeconds
        self.counts = {}
        self.starts = {}

    def allowRequest(self, req):
        if req.clientId not in self.starts or req.timestamp >= self.starts[req.clientId] + self.windowSizeSeconds:
            self.starts[req.clientId] = req.timestamp
            self.counts[req.clientId] = 0
        if self.counts[req.clientId] >= self.maxRequests:
            return False
        self.counts[req.clientId] += 1
        return True

    def getRequestCount(self, clientId):
        return self.counts.get(clientId, 0)


class SlidingWindowLimiter(RateLimiter):
    def __init__(self, maxRequests, windowSizeSeconds):
        self.maxRequests = maxRequests
        self.windowSizeSeconds = windowSizeSeconds
        self.queues = {}

    def allowRequest(self, req):
        q = self.queues.setdefault(req.clientId, deque())
        while q and q[0] <= req.timestamp - self.windowSizeSeconds:
            q.popleft()
        if len(q) >= self.maxRequests:
            return False
        q.append(req.timestamp)
        return True

    def getRequestCount(self, clientId):
        return len(self.queues.get(clientId, deque()))


class TokenBucketLimiter(RateLimiter):
    def __init__(self, maxTokens, windowSize):
        self.maxTokens = maxTokens
        self.refillRate = maxTokens / windowSize
        self.tokens = {}
        self.lastRefill = {}

    def allowRequest(self, req):
        if req.clientId not in self.tokens:
            self.tokens[req.clientId] = self.maxTokens
            self.lastRefill[req.clientId] = req.timestamp
        elapsed = req.timestamp - self.lastRefill[req.clientId]
        self.tokens[req.clientId] = min(self.maxTokens, self.tokens[req.clientId] + elapsed * self.refillRate)
        self.lastRefill[req.clientId] = req.timestamp
        if self.tokens[req.clientId] < 1.0:
            return False
        self.tokens[req.clientId] -= 1.0
        return True

    def getRequestCount(self, clientId):
        if clientId not in self.tokens:
            return 0
        return self.maxTokens - int(self.tokens[clientId])


def create_limiter(algorithm, maxRequests, windowSize):
    if algorithm == "fixed-window":
        return FixedWindowLimiter(maxRequests, windowSize)
    if algorithm == "sliding-window":
        return SlidingWindowLimiter(maxRequests, windowSize)
    if algorithm == "token-bucket":
        return TokenBucketLimiter(maxRequests, windowSize)
    return None


# ─── Module state ────────────────────────────────────────────────────────────

_g_limiter = None
_g_strategy = {}
_g_tier = {}


def reset_service():
    global _g_limiter, _g_strategy, _g_tier
    _g_limiter = None
    _g_strategy = {}
    _g_tier = {}


def init_limiter(maxRequests, windowSize):
    global _g_limiter
    _g_limiter = FixedWindowLimiter(maxRequests, windowSize)


def allow_request(req):
    if _g_limiter is None:
        return False
    return _g_limiter.allowRequest(req)


def allow_request_simple(clientId, timestamp, endpoint):
    return allow_request(Request(clientId, timestamp, endpoint))


def get_request_count(clientId):
    if _g_limiter is None:
        return 0
    return _g_limiter.getRequestCount(clientId)


def allow_request_with_strategy(algorithm, req):
    if algorithm not in _g_strategy:
        _g_strategy[algorithm] = create_limiter(algorithm, 100, 60)
    if _g_strategy[algorithm] is None:
        return False
    return _g_strategy[algorithm].allowRequest(req)


def allow_request_with_strategy_simple(algorithm, clientId, timestamp, endpoint):
    return allow_request_with_strategy(algorithm, Request(clientId, timestamp, endpoint))


_TIER_LIMITS = {"FREE": 10, "PRO": 100, "ENTERPRISE": 1000}


def allow_request_for_tier(tier, req):
    if tier not in _g_tier:
        limit = _TIER_LIMITS.get(tier, 10)
        _g_tier[tier] = SlidingWindowLimiter(limit, limit + 1)
    return _g_tier[tier].allowRequest(req)


def allow_request_for_tier_str(tier, clientId, timestamp, endpoint):
    return allow_request_for_tier(tier, Request(clientId, timestamp, endpoint))
