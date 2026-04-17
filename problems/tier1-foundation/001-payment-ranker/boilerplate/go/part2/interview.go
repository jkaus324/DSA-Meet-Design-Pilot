package main

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type PaymentMethod struct {
	Name           string
	CashbackRate   float64 // e.g. 0.05 = 5%
	TransactionFee float64 // in rupees
	UsageCount     int
}

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The product team now wants COMPOSITE ranking:
// rank by cashback first, then use transaction fee as tiebreaker.
//
// Think about:
//   - How do you chain ranking criteria without modifying existing strategies?
//   - What if the product team adds a 4th criterion tomorrow?
//   - Is your Part 1 design extensible enough to handle this?
//
// Entry points (must exist for tests):
//   func RankByRewards(methods []PaymentMethod) []PaymentMethod
//   func RankByLowFee(methods []PaymentMethod) []PaymentMethod
//   func RankByTrust(methods []PaymentMethod) []PaymentMethod
//   func RankComposite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod
//
// ─────────────────────────────────────────────────────────────────────────────
