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

// ─── NEW in Extension 2 ───────────────────────────────────────────────────────
//
// Different user tiers have different rate limits:
//   FREE       = 10 requests per minute
//   PRO        = 100 requests per minute
//   ENTERPRISE = 1000 requests per minute
//
// Think about:
//   - How does the Factory pattern adapt to handle tier-based creation?
//   - Can you combine tier limits with per-endpoint algorithm selection?
//   - What changes if a new tier is added (e.g., STARTUP = 50 req/min)?
//
// Entry points (must exist for tests — all previous entry points plus):
//   func AllowRequestForTier(tier UserTier, req Request) bool
//
// ─────────────────────────────────────────────────────────────────────────────

type RateLimiter interface {
	AllowRequest(req Request) bool
	GetRequestCount(clientId string) int
}

func InitLimiter(maxRequests int, windowSize int)                                  {}
func AllowRequest(req Request) bool                                                { return false }
func GetRequestCount(clientId string) int                                          { return 0 }
func CreateLimiter(algorithm string, maxRequests int, windowSize int) RateLimiter  { return nil }
func AllowRequestWithStrategy(algorithm string, req Request) bool                  { return false }
func AllowRequestForTier(tier UserTier, req Request) bool                          { return false }
