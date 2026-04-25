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

// ─── Existing Strategies ─────────────────────────────────────────────────────
// TODO: Copy your Part 1 strategies here (or extend them)

type RewardsMaximizer struct{}

func (s RewardsMaximizer) Compare(a, b PaymentMethod) bool {
	return false // TODO: implement
}

type LowFeeSeeker struct{}

func (s LowFeeSeeker) Compare(a, b PaymentMethod) bool {
	return false // TODO: implement
}

type TrustBasedRanker struct{}

func (s TrustBasedRanker) Compare(a, b PaymentMethod) bool {
	return false // TODO: implement
}

// ─── NEW: CompositeStrategy ───────────────────────────────────────────────────
// HINT: A CompositeStrategy holds a list of other strategies.
// It tries the first strategy; if tied, falls back to the second, then third...
// This is the Composite pattern applied to a comparator.

type CompositeStrategy struct {
	criteria []RankingStrategy
}

func NewCompositeStrategy(criteria []RankingStrategy) *CompositeStrategy {
	return &CompositeStrategy{criteria: criteria}
}

func (s *CompositeStrategy) Compare(a, b PaymentMethod) bool {
	// TODO: Iterate through criteria.
	// If criteria[i] says a > b, return true.
	// If criteria[i] says b > a, return false.
	// If tied, move to criteria[i+1].
	return false
}

// ─── Ranker ──────────────────────────────────────────────────────────────────

type PaymentRanker struct {
	strategy RankingStrategy
}

func NewPaymentRanker(s RankingStrategy) *PaymentRanker {
	return &PaymentRanker{strategy: s}
}

func (r *PaymentRanker) Rank(methods []PaymentMethod) []PaymentMethod {
	// TODO: Sort using strategy
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

func RankComposite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod {
	return NewPaymentRanker(NewCompositeStrategy(criteria)).Rank(methods)
}
