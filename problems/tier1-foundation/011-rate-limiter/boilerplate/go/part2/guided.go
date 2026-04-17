package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type Request struct {
	ClientId  string
	Timestamp int64
	Endpoint  string
}

type UserTier int

const (
	FREE       UserTier = iota
	PRO        UserTier = iota
	ENTERPRISE UserTier = iota
)

// ─── Strategy Interface ───────────────────────────────────────────────────────

type RateLimiter interface {
	AllowRequest(req Request) bool
	GetRequestCount(clientId string) int
}

// ─── Fixed-Window (provided from Part 1) ──────────────────────────────────────

type FixedWindowLimiter struct {
	maxRequests       int
	windowSizeSeconds int64
	requestCounts     map[string]int
	windowStarts      map[string]int64
}

func NewFixedWindowLimiter(maxReq int, windowSize int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		maxRequests:       maxReq,
		windowSizeSeconds: int64(windowSize),
		requestCounts:     make(map[string]int),
		windowStarts:      make(map[string]int64),
	}
}

func (l *FixedWindowLimiter) getWindowStart(ts int64) int64 {
	return (ts / l.windowSizeSeconds) * l.windowSizeSeconds
}

func (l *FixedWindowLimiter) AllowRequest(req Request) bool {
	ws := l.getWindowStart(req.Timestamp)
	if l.windowStarts[req.ClientId] != ws {
		l.windowStarts[req.ClientId] = ws
		l.requestCounts[req.ClientId] = 0
	}
	if l.requestCounts[req.ClientId] >= l.maxRequests {
		return false
	}
	l.requestCounts[req.ClientId]++
	return true
}

func (l *FixedWindowLimiter) GetRequestCount(clientId string) int {
	return l.requestCounts[clientId]
}

// ─── Sliding-Window Limiter ───────────────────────────────────────────────────
// TODO: Implement using a slice of timestamps per client
// HINT: Remove expired timestamps (older than windowSize seconds)
//       If slice length >= maxRequests, reject

// type SlidingWindowLimiter struct { ... }

// ─── Token-Bucket Limiter ─────────────────────────────────────────────────────
// TODO: Implement using a token count and last-refill timestamp per client
// HINT: tokens refill at rate = maxTokens / windowSize per second
//       On each request, refill based on elapsed time, then consume 1 token

// type TokenBucketLimiter struct { ... }

// ─── Factory ──────────────────────────────────────────────────────────────────
// TODO: Create a factory that returns the right limiter based on algorithm name
//   "fixed-window"   -> FixedWindowLimiter
//   "sliding-window" -> SlidingWindowLimiter
//   "token-bucket"   -> TokenBucketLimiter
//   unknown          -> nil

func CreateLimiter(algorithm string, maxRequests int, windowSize int) RateLimiter {
	if algorithm == "fixed-window" {
		return NewFixedWindowLimiter(maxRequests, windowSize)
	}
	// TODO: add sliding-window and token-bucket cases
	return nil
}

// ─── Global Entry Points ──────────────────────────────────────────────────────

var gLimiter RateLimiter
var gStrategyLimiters = make(map[string]RateLimiter)

func InitLimiter(maxRequests int, windowSize int) {
	gLimiter = NewFixedWindowLimiter(maxRequests, windowSize)
}

func AllowRequest(req Request) bool {
	if gLimiter == nil {
		return false
	}
	return gLimiter.AllowRequest(req)
}

func GetRequestCount(clientId string) int {
	if gLimiter == nil {
		return 0
	}
	return gLimiter.GetRequestCount(clientId)
}

func AllowRequestWithStrategy(algorithm string, req Request) bool {
	if _, ok := gStrategyLimiters[algorithm]; !ok {
		gStrategyLimiters[algorithm] = CreateLimiter(algorithm, 100, 60)
	}
	return gStrategyLimiters[algorithm].AllowRequest(req)
}
