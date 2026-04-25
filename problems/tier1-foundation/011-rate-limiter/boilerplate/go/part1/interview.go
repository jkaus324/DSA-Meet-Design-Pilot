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

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement a Rate Limiter that:
//   1. Tracks requests per client in a fixed time window
//   2. Rejects requests that exceed the limit
//   3. Allows new rate-limiting algorithms to be added WITHOUT modifying
//      existing code
//
// Think about:
//   - What abstraction lets you swap rate-limiting logic at runtime?
//   - How would you track per-client request counts efficiently?
//   - What happens when the time window rolls over?
//
// Entry points (must exist for tests):
//   func InitLimiter(maxRequests int, windowSize int)
//   func AllowRequest(req Request) bool
//   func GetRequestCount(clientId string) int
//
// ─────────────────────────────────────────────────────────────────────────────

func InitLimiter(maxRequests int, windowSize int) {}

func AllowRequest(req Request) bool {
	return false
}

func GetRequestCount(clientId string) int {
	return 0
}
