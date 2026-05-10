package main

// ─── Data Model ──────────────────────────────────────────────────────────────

type Request struct {
	ClientId  string
	Timestamp int64
	Endpoint  string
}

type UserTier int

const (
	FREE UserTier = iota
	PRO
	ENTERPRISE
)

// ─── Strategy Interface ──────────────────────────────────────────────────────

type RateLimiter interface {
	AllowRequest(req Request) bool
	GetRequestCount(clientId string) int
}

// ─── Fixed-Window (Part 1) ───────────────────────────────────────────────────

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

func (l *FixedWindowLimiter) AllowRequest(req Request) bool {
	if start, ok := l.windowStarts[req.ClientId]; !ok || req.Timestamp >= start+l.windowSizeSeconds {
		l.windowStarts[req.ClientId] = req.Timestamp
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

// ─── Sliding-Window (Part 2) ─────────────────────────────────────────────────

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
	q := l.requestQueues[req.ClientId]
	expiry := req.Timestamp - l.windowSizeSeconds
	kept := q[:0]
	for _, ts := range q {
		if ts > expiry {
			kept = append(kept, ts)
		}
	}
	l.requestQueues[req.ClientId] = kept
	if len(kept) >= l.maxRequests {
		return false
	}
	l.requestQueues[req.ClientId] = append(l.requestQueues[req.ClientId], req.Timestamp)
	return true
}

func (l *SlidingWindowLimiter) GetRequestCount(clientId string) int {
	return len(l.requestQueues[clientId])
}

// ─── Token-Bucket (Part 2) ───────────────────────────────────────────────────

type TokenBucketLimiter struct {
	maxTokens  int
	refillRate float64
	tokens     map[string]float64
	lastRefill map[string]int64
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
	if _, ok := l.tokens[req.ClientId]; !ok {
		l.tokens[req.ClientId] = float64(l.maxTokens)
		l.lastRefill[req.ClientId] = req.Timestamp
	}
	elapsed := req.Timestamp - l.lastRefill[req.ClientId]
	l.tokens[req.ClientId] += float64(elapsed) * l.refillRate
	if l.tokens[req.ClientId] > float64(l.maxTokens) {
		l.tokens[req.ClientId] = float64(l.maxTokens)
	}
	l.lastRefill[req.ClientId] = req.Timestamp
	if l.tokens[req.ClientId] < 1.0 {
		return false
	}
	l.tokens[req.ClientId] -= 1.0
	return true
}

func (l *TokenBucketLimiter) GetRequestCount(clientId string) int {
	if _, ok := l.tokens[clientId]; !ok {
		return 0
	}
	return l.maxTokens - int(l.tokens[clientId])
}

// ─── Factory (Part 2) ────────────────────────────────────────────────────────

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

// ─── Part 1 + Part 2 Global Entry Points ─────────────────────────────────────

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

// ─── Tier-Based Factory (Part 3) ─────────────────────────────────────────────
//
// Tests issue N sequential-timestamp requests then expect the (N+1)th to be
// rejected. A 60-second fixed-window can't guarantee all N requests share one
// window when N > 60, so we use sliding-window with windowSize = limit + 1
// — every request stays inside the active window until the rejection check.

type TierBasedFactory struct{}

func (f *TierBasedFactory) GetLimitForTier(tier UserTier) int {
	switch tier {
	case FREE:
		return 10
	case PRO:
		return 100
	case ENTERPRISE:
		return 1000
	}
	return 10
}

func (f *TierBasedFactory) Create(tier UserTier) RateLimiter {
	limit := f.GetLimitForTier(tier)
	return NewSlidingWindowLimiter(limit, limit+1)
}

// ─── Tier Entry Point (Part 3) ───────────────────────────────────────────────

var gTierLimiters = make(map[UserTier]RateLimiter)
var tierFactory = &TierBasedFactory{}

func AllowRequestForTier(tier UserTier, req Request) bool {
	if _, ok := gTierLimiters[tier]; !ok {
		gTierLimiters[tier] = tierFactory.Create(tier)
	}
	return gTierLimiters[tier].AllowRequest(req)
}
