'use strict';

// Rate limiter — Strategy + Factory reference solution (JavaScript).

class Request {
  constructor(clientId, timestamp, endpoint) {
    this.clientId = clientId;
    this.timestamp = timestamp;
    this.endpoint = endpoint;
  }
}

class FixedWindowLimiter {
  constructor(maxRequests, windowSizeSeconds) {
    this.maxRequests = maxRequests;
    this.windowSizeSeconds = windowSizeSeconds;
    this.counts = new Map();
    this.starts = new Map();
  }
  allowRequest(req) {
    if (!this.starts.has(req.clientId) || req.timestamp >= this.starts.get(req.clientId) + this.windowSizeSeconds) {
      this.starts.set(req.clientId, req.timestamp);
      this.counts.set(req.clientId, 0);
    }
    if (this.counts.get(req.clientId) >= this.maxRequests) return false;
    this.counts.set(req.clientId, this.counts.get(req.clientId) + 1);
    return true;
  }
  getRequestCount(clientId) {
    return this.counts.has(clientId) ? this.counts.get(clientId) : 0;
  }
}

class SlidingWindowLimiter {
  constructor(maxRequests, windowSizeSeconds) {
    this.maxRequests = maxRequests;
    this.windowSizeSeconds = windowSizeSeconds;
    this.queues = new Map();
  }
  allowRequest(req) {
    if (!this.queues.has(req.clientId)) this.queues.set(req.clientId, []);
    const q = this.queues.get(req.clientId);
    while (q.length && q[0] <= req.timestamp - this.windowSizeSeconds) q.shift();
    if (q.length >= this.maxRequests) return false;
    q.push(req.timestamp);
    return true;
  }
  getRequestCount(clientId) {
    return this.queues.has(clientId) ? this.queues.get(clientId).length : 0;
  }
}

class TokenBucketLimiter {
  constructor(maxTokens, windowSize) {
    this.maxTokens = maxTokens;
    this.refillRate = maxTokens / windowSize;
    this.tokens = new Map();
    this.lastRefill = new Map();
  }
  allowRequest(req) {
    if (!this.tokens.has(req.clientId)) {
      this.tokens.set(req.clientId, this.maxTokens);
      this.lastRefill.set(req.clientId, req.timestamp);
    }
    const elapsed = req.timestamp - this.lastRefill.get(req.clientId);
    this.tokens.set(req.clientId, Math.min(this.maxTokens, this.tokens.get(req.clientId) + elapsed * this.refillRate));
    this.lastRefill.set(req.clientId, req.timestamp);
    if (this.tokens.get(req.clientId) < 1.0) return false;
    this.tokens.set(req.clientId, this.tokens.get(req.clientId) - 1.0);
    return true;
  }
  getRequestCount(clientId) {
    if (!this.tokens.has(clientId)) return 0;
    return this.maxTokens - Math.trunc(this.tokens.get(clientId));
  }
}

function create_limiter(algorithm, maxRequests, windowSize) {
  if (algorithm === 'fixed-window') return new FixedWindowLimiter(maxRequests, windowSize);
  if (algorithm === 'sliding-window') return new SlidingWindowLimiter(maxRequests, windowSize);
  if (algorithm === 'token-bucket') return new TokenBucketLimiter(maxRequests, windowSize);
  return null;
}

// ─── Module state ────────────────────────────────────────────────────────────

let _g_limiter = null;
let _g_strategy = new Map();
let _g_tier = new Map();

function reset_service() {
  _g_limiter = null;
  _g_strategy = new Map();
  _g_tier = new Map();
}

function init_limiter(maxRequests, windowSize) {
  _g_limiter = new FixedWindowLimiter(maxRequests, windowSize);
}

function allow_request(req) {
  if (_g_limiter === null) return false;
  return _g_limiter.allowRequest(req);
}

function allow_request_simple(clientId, timestamp, endpoint) {
  return allow_request(new Request(clientId, timestamp, endpoint));
}

function get_request_count(clientId) {
  if (_g_limiter === null) return 0;
  return _g_limiter.getRequestCount(clientId);
}

function allow_request_with_strategy(algorithm, req) {
  if (!_g_strategy.has(algorithm)) {
    _g_strategy.set(algorithm, create_limiter(algorithm, 100, 60));
  }
  if (_g_strategy.get(algorithm) === null) return false;
  return _g_strategy.get(algorithm).allowRequest(req);
}

function allow_request_with_strategy_simple(algorithm, clientId, timestamp, endpoint) {
  return allow_request_with_strategy(algorithm, new Request(clientId, timestamp, endpoint));
}

const _TIER_LIMITS = { FREE: 10, PRO: 100, ENTERPRISE: 1000 };

function allow_request_for_tier(tier, req) {
  if (!_g_tier.has(tier)) {
    const limit = Object.prototype.hasOwnProperty.call(_TIER_LIMITS, tier) ? _TIER_LIMITS[tier] : 10;
    _g_tier.set(tier, new SlidingWindowLimiter(limit, limit + 1));
  }
  return _g_tier.get(tier).allowRequest(req);
}

function allow_request_for_tier_str(tier, clientId, timestamp, endpoint) {
  return allow_request_for_tier(tier, new Request(clientId, timestamp, endpoint));
}

module.exports = {
  Request,
  FixedWindowLimiter,
  SlidingWindowLimiter,
  TokenBucketLimiter,
  create_limiter,
  reset_service,
  init_limiter,
  allow_request_simple,
  get_request_count,
  allow_request_with_strategy_simple,
  allow_request_for_tier_str,
};
