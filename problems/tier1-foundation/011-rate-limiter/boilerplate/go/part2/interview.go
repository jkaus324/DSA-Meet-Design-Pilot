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

// ─── NEW in Extension 1 ───────────────────────────────────────────────────────
//
// The platform team wants different endpoints to use different algorithms:
//   - fixed-window: simple counter per time window
//   - sliding-window: rolling window using a queue of timestamps
//   - token-bucket: tokens replenish at a fixed rate, allows bursts
//
// Think about:
//   - How do you create the right limiter without the caller knowing the algo?
//   - What pattern encapsulates object creation decisions?
//   - How would you add a 4th algorithm tomorrow?
//
// Entry points (must exist for tests):
//   func InitLimiter(maxRequests int, windowSize int)
//   func AllowRequest(req Request) bool
//   func GetRequestCount(clientId string) int
//   func CreateLimiter(algorithm string, maxRequests int, windowSize int) RateLimiter
//   func AllowRequestWithStrategy(algorithm string, req Request) bool
//
// ─────────────────────────────────────────────────────────────────────────────

type RateLimiter interface {
	AllowRequest(req Request) bool
	GetRequestCount(clientId string) int
}

func InitLimiter(maxRequests int, windowSize int)                            {}
func AllowRequest(req Request) bool                                          { return false }
func GetRequestCount(clientId string) int                                    { return 0 }
func CreateLimiter(algorithm string, maxRequests int, windowSize int) RateLimiter { return nil }
func AllowRequestWithStrategy(algorithm string, req Request) bool            { return false }
