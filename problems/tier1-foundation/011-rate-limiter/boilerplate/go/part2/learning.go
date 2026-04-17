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

// ─── Fixed-Window (complete from Part 1) ──────────────────────────────────────

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
// TODO: Implement AllowRequest() and GetRequestCount()

type SlidingWindowLimiter struct {
	maxRequests       int
	windowSizeSeconds int64
	requestQueues     map[string][]int64
}

func NewSlidingWindowLimiter(maxReq int, windowSize int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		maxRequests:       maxReq,
		windowSizeSeconds: int64(windowSize),
		requestQueues:     make(map[string][]int64),
	}
}

func (l *SlidingWindowLimiter) AllowRequest(req Request) bool {
	// TODO: Remove expired timestamps from the front of the queue
	// TODO: If queue length >= maxRequests, return false
	// TODO: Append current timestamp and return true
	return false
}

func (l *SlidingWindowLimiter) GetRequestCount(clientId string) int {
	// TODO: Return current queue length for this client
	return 0
}

// ─── Token-Bucket Limiter ─────────────────────────────────────────────────────
// TODO: Implement AllowRequest() and GetRequestCount()

type TokenBucketLimiter struct {
	maxTokens    int
	refillRate   float64
	tokens       map[string]float64
	lastRefill   map[string]int64
}

func NewTokenBucketLimiter(maxTokens int, windowSize int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		maxTokens:  maxTokens,
		refillRate: float64(maxTokens) / float64(windowSize),
		tokens:     make(map[string]float64),
		lastRefill: make(map[string]int64),
	}
}

func (l *TokenBucketLimiter) AllowRequest(req Request) bool {
	// TODO: Initialize tokens for new clients
	// TODO: Refill tokens based on elapsed time
	// TODO: If tokens < 1, return false
	// TODO: Consume 1 token and return true
	return false
}

func (l *TokenBucketLimiter) GetRequestCount(clientId string) int {
	// TODO: Return tokens consumed (maxTokens - remaining)
	return 0
}

// ─── Factory ──────────────────────────────────────────────────────────────────

func CreateLimiter(algorithm string, maxRequests int, windowSize int) RateLimiter {
	switch algorithm {
	case "fixed-window":
		return NewFixedWindowLimiter(maxRequests, windowSize)
	case "sliding-window":
		return NewSlidingWindowLimiter(maxRequests, windowSize)
	case "token-bucket":
		return NewTokenBucketLimiter(maxRequests, windowSize)
	}
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
	if gStrategyLimiters[algorithm] == nil {
		return false
	}
	return gStrategyLimiters[algorithm].AllowRequest(req)
}
