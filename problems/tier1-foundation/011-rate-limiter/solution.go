// Rate limiter — Strategy + Factory reference solution (Go).
package main

type Request struct {
	clientId  string
	timestamp int
	endpoint  string
}

type RateLimiter interface {
	allowRequest(req Request) bool
	getRequestCount(clientId string) int
}

type FixedWindowLimiter struct {
	maxRequests       int
	windowSizeSeconds int
	counts            map[string]int
	starts            map[string]int
}

func newFixedWindowLimiter(maxRequests, windowSize int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		maxRequests:       maxRequests,
		windowSizeSeconds: windowSize,
		counts:            map[string]int{},
		starts:            map[string]int{},
	}
}

func (l *FixedWindowLimiter) allowRequest(req Request) bool {
	start, ok := l.starts[req.clientId]
	if !ok || req.timestamp >= start+l.windowSizeSeconds {
		l.starts[req.clientId] = req.timestamp
		l.counts[req.clientId] = 0
	}
	if l.counts[req.clientId] >= l.maxRequests {
		return false
	}
	l.counts[req.clientId]++
	return true
}

func (l *FixedWindowLimiter) getRequestCount(clientId string) int {
	return l.counts[clientId]
}

type SlidingWindowLimiter struct {
	maxRequests       int
	windowSizeSeconds int
	queues            map[string][]int
}

func newSlidingWindowLimiter(maxRequests, windowSize int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		maxRequests:       maxRequests,
		windowSizeSeconds: windowSize,
		queues:            map[string][]int{},
	}
}

func (l *SlidingWindowLimiter) allowRequest(req Request) bool {
	q := l.queues[req.clientId]
	for len(q) > 0 && q[0] <= req.timestamp-l.windowSizeSeconds {
		q = q[1:]
	}
	if len(q) >= l.maxRequests {
		l.queues[req.clientId] = q
		return false
	}
	q = append(q, req.timestamp)
	l.queues[req.clientId] = q
	return true
}

func (l *SlidingWindowLimiter) getRequestCount(clientId string) int {
	return len(l.queues[clientId])
}

type TokenBucketLimiter struct {
	maxTokens  int
	refillRate float64
	tokens     map[string]float64
	lastRefill map[string]int
	seen       map[string]bool
}

func newTokenBucketLimiter(maxTokens, windowSize int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		maxTokens:  maxTokens,
		refillRate: float64(maxTokens) / float64(windowSize),
		tokens:     map[string]float64{},
		lastRefill: map[string]int{},
		seen:       map[string]bool{},
	}
}

func (l *TokenBucketLimiter) allowRequest(req Request) bool {
	if !l.seen[req.clientId] {
		l.seen[req.clientId] = true
		l.tokens[req.clientId] = float64(l.maxTokens)
		l.lastRefill[req.clientId] = req.timestamp
	}
	elapsed := req.timestamp - l.lastRefill[req.clientId]
	tok := l.tokens[req.clientId] + float64(elapsed)*l.refillRate
	if tok > float64(l.maxTokens) {
		tok = float64(l.maxTokens)
	}
	l.tokens[req.clientId] = tok
	l.lastRefill[req.clientId] = req.timestamp
	if l.tokens[req.clientId] < 1.0 {
		return false
	}
	l.tokens[req.clientId] -= 1.0
	return true
}

func (l *TokenBucketLimiter) getRequestCount(clientId string) int {
	if !l.seen[clientId] {
		return 0
	}
	return l.maxTokens - int(l.tokens[clientId])
}

func createLimiter(algorithm string, maxRequests, windowSize int) RateLimiter {
	switch algorithm {
	case "fixed-window":
		return newFixedWindowLimiter(maxRequests, windowSize)
	case "sliding-window":
		return newSlidingWindowLimiter(maxRequests, windowSize)
	case "token-bucket":
		return newTokenBucketLimiter(maxRequests, windowSize)
	}
	return nil
}

// ─── Module state ────────────────────────────────────────────────────────────

var gLimiter RateLimiter
var gStrategy map[string]RateLimiter
var gTier map[string]RateLimiter

func reset_service() {
	gLimiter = nil
	gStrategy = map[string]RateLimiter{}
	gTier = map[string]RateLimiter{}
}

func init_limiter(maxRequests int, windowSize int) {
	gLimiter = newFixedWindowLimiter(maxRequests, windowSize)
}

func allow_request_simple(clientId string, timestamp int, endpoint string) bool {
	if gLimiter == nil {
		return false
	}
	return gLimiter.allowRequest(Request{clientId, timestamp, endpoint})
}

func get_request_count(clientId string) int {
	if gLimiter == nil {
		return 0
	}
	return gLimiter.getRequestCount(clientId)
}

func allow_request_with_strategy_simple(algorithm string, clientId string, timestamp int, endpoint string) bool {
	if gStrategy == nil {
		gStrategy = map[string]RateLimiter{}
	}
	if _, ok := gStrategy[algorithm]; !ok {
		gStrategy[algorithm] = createLimiter(algorithm, 100, 60)
	}
	if gStrategy[algorithm] == nil {
		return false
	}
	return gStrategy[algorithm].allowRequest(Request{clientId, timestamp, endpoint})
}

var tierLimits = map[string]int{"FREE": 10, "PRO": 100, "ENTERPRISE": 1000}

func allow_request_for_tier_str(tier string, clientId string, timestamp int, endpoint string) bool {
	if gTier == nil {
		gTier = map[string]RateLimiter{}
	}
	if _, ok := gTier[tier]; !ok {
		limit, has := tierLimits[tier]
		if !has {
			limit = 10
		}
		gTier[tier] = newSlidingWindowLimiter(limit, limit+1)
	}
	return gTier[tier].allowRequest(Request{clientId, timestamp, endpoint})
}
