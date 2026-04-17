package main

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type PaymentMethod struct {
	Name                string
	CashbackRate        float64 // e.g. 0.05 = 5%
	TransactionFee      float64 // in rupees
	UsageCount          int
	EasyRefundEligible  bool // NEW in Part 3
}

// ─── NEW in Extension 2 ──────────────────────────────────────────────────────
//
// The compliance team wants to add "easy-refund eligibility" as a filter.
// Some payment methods don't support easy refunds — these should be ranked
// lower regardless of cashback or fee, unless the user explicitly opts in.
//
// Think about:
//   - Is this a ranking criterion, a filter, or both?
//   - How does your existing CompositeStrategy handle a boolean filter?
//   - What if the "opt-in" flag is per-user, not per-session?
//
// Entry points (must exist for tests):
//   func RankByRewards(methods []PaymentMethod) []PaymentMethod
//   func RankByLowFee(methods []PaymentMethod) []PaymentMethod
//   func RankByTrust(methods []PaymentMethod) []PaymentMethod
//   func RankComposite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod
//   func RankWithRefundFilter(methods []PaymentMethod, preferEasyRefund bool) []PaymentMethod
//
// ─────────────────────────────────────────────────────────────────────────────
