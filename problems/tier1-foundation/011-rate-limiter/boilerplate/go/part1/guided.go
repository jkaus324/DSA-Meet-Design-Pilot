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
// HINT: This interface lets you swap rate-limiting algorithms at runtime.
// What methods would a rate limiter need?

type RateLimiter interface {
	// HINT: what methods does a rate limiter expose?
	// AllowRequest(req Request) bool
	// GetRequestCount(clientId string) int
}

// ─── Concrete Strategy ────────────────────────────────────────────────────────
// TODO: Implement a fixed-window rate limiter:
//   - Divide time into windows of windowSizeSeconds
//   - Track request count per client per window (use map)
//   - If count >= maxRequests, reject
//   - When a new window starts, reset the count
//
// HINT: windowStart = (timestamp / windowSize) * windowSize

// ─── Global Entry Points ──────────────────────────────────────────────────────
// TODO: Implement these functions using your rate limiter:

func InitLimiter(maxRequests int, windowSize int) {}

func AllowRequest(req Request) bool {
	return false
}

func GetRequestCount(clientId string) int {
	return 0
}
