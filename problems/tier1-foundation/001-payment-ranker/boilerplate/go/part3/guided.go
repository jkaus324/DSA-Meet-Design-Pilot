package main

import "sort"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type PaymentMethod struct {
	Name               string
	CashbackRate       float64 // e.g. 0.05 = 5%
	TransactionFee     float64 // in rupees
	UsageCount         int
	EasyRefundEligible bool // NEW in Part 3
}

// ─── Strategy Interface ──────────────────────────────────────────────────────

type RankingStrategy interface {
	Compare(a, b PaymentMethod) bool
}

// ─── Existing Strategies ─────────────────────────────────────────────────────

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

type CompositeStrategy struct {
	criteria []RankingStrategy
}

func NewCompositeStrategy(criteria []RankingStrategy) *CompositeStrategy {
	return &CompositeStrategy{criteria: criteria}
}

func (s *CompositeStrategy) Compare(a, b PaymentMethod) bool {
	for _, c := range s.criteria {
		if c.Compare(a, b) {
			return true
		}
		if c.Compare(b, a) {
			return false
		}
	}
	return false
}

// ─── NEW: EasyRefundStrategy ─────────────────────────────────────────────────
// HINT: When prefer=true, methods with EasyRefundEligible=true
// should always rank above those without, regardless of other criteria.

type EasyRefundStrategy struct {
	prefer bool
}

func NewEasyRefundStrategy(preferRefund bool) *EasyRefundStrategy {
	return &EasyRefundStrategy{prefer: preferRefund}
}

func (s *EasyRefundStrategy) Compare(a, b PaymentMethod) bool {
	// TODO: If prefer=true, a wins if a.EasyRefundEligible && !b.EasyRefundEligible
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
	// TODO: implement
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

func RankWithRefundFilter(methods []PaymentMethod, preferEasyRefund bool) []PaymentMethod {
	// TODO: Use EasyRefundStrategy as the first criterion in a CompositeStrategy
	return methods
}
