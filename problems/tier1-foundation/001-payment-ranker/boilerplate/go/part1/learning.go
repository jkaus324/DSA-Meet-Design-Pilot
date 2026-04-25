package main

import "sort"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type PaymentMethod struct {
	Name           string
	CashbackRate   float64 // e.g. 0.05 = 5%
	TransactionFee float64 // in rupees
	UsageCount     int
}

// ─── Strategy Interface ──────────────────────────────────────────────────────

type RankingStrategy interface {
	Compare(a, b PaymentMethod) bool
}

// ─── Concrete Strategies ─────────────────────────────────────────────────────
// TODO: Implement the Compare() method for each strategy

type RewardsMaximizer struct{}

func (s RewardsMaximizer) Compare(a, b PaymentMethod) bool {
	// TODO: return true if 'a' should rank higher than 'b'
	// Higher cashback rate = better ranking
	return false
}

type LowFeeSeeker struct{}

func (s LowFeeSeeker) Compare(a, b PaymentMethod) bool {
	// TODO: return true if 'a' should rank higher than 'b'
	// Lower transaction fee = better ranking
	return false
}

type TrustBasedRanker struct{}

func (s TrustBasedRanker) Compare(a, b PaymentMethod) bool {
	// TODO: return true if 'a' should rank higher than 'b'
	// Higher usage count = better ranking
	return false
}

// ─── Ranker ──────────────────────────────────────────────────────────────────

type PaymentRanker struct {
	strategy RankingStrategy
}

func NewPaymentRanker(s RankingStrategy) *PaymentRanker {
	return &PaymentRanker{strategy: s}
}

func (r *PaymentRanker) SetStrategy(s RankingStrategy) {
	r.strategy = s
}

func (r *PaymentRanker) Rank(methods []PaymentMethod) []PaymentMethod {
	// TODO: Sort methods using the current strategy's Compare()
	// Return the sorted slice
	result := make([]PaymentMethod, len(methods))
	copy(result, methods)
	sort.Slice(result, func(i, j int) bool {
		return r.strategy.Compare(result[i], result[j])
	})
	return result
}

// ─── Test Entry Points ───────────────────────────────────────────────────────

func RankByRewards(methods []PaymentMethod) []PaymentMethod {
	return NewPaymentRanker(RewardsMaximizer{}).Rank(methods)
}

func RankByLowFee(methods []PaymentMethod) []PaymentMethod {
	return NewPaymentRanker(LowFeeSeeker{}).Rank(methods)
}

func RankByTrust(methods []PaymentMethod) []PaymentMethod {
	return NewPaymentRanker(TrustBasedRanker{}).Rank(methods)
}
