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

// ─── Fixed-Window Rate Limiter ────────────────────────────────────────────────
// TODO: Implement the AllowRequest() and GetRequestCount() methods

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

func (l *FixedWindowLimiter) getWindowStart(timestamp int64) int64 {
	return (timestamp / l.windowSizeSeconds) * l.windowSizeSeconds
}

func (l *FixedWindowLimiter) AllowRequest(req Request) bool {
	// TODO: Check if we're in a new window (reset count if so)
	// TODO: If count >= maxRequests, return false
	// TODO: Increment count and return true
	return false
}

func (l *FixedWindowLimiter) GetRequestCount(clientId string) int {
	// TODO: Return the current request count for this client
	return 0
}

// ─── Global Entry Points ──────────────────────────────────────────────────────

var gLimiter RateLimiter

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
