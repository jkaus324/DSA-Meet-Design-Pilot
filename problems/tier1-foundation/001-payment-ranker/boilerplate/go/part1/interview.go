package main

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type PaymentMethod struct {
	Name            string
	CashbackRate    float64 // e.g. 0.05 = 5%
	TransactionFee  float64 // in rupees
	UsageCount      int
}

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a PaymentRanker that:
//   1. Ranks payment methods by the criteria described in the problem
//   2. Allows new ranking strategies to be added WITHOUT modifying
//      the ranker itself
//
// Think about:
//   - What abstraction lets you swap ranking logic at runtime?
//   - How would you add a 4th ranking criterion with zero changes
//     to existing code?
//   - What happens to your code when Extension 1 (cashback) is added?
//
// Entry points (must exist for tests):
//   func RankByRewards(methods []PaymentMethod) []PaymentMethod
//   func RankByLowFee(methods []PaymentMethod) []PaymentMethod
//   func RankByTrust(methods []PaymentMethod) []PaymentMethod
//
// ─────────────────────────────────────────────────────────────────────────────
