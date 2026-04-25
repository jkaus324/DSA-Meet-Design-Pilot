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

// ─── Concrete Strategies ─────────────────────────────────────────────────────

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

// ─── CompositeStrategy ───────────────────────────────────────────────────────

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

// ─── EasyRefundStrategy ──────────────────────────────────────────────────────

type EasyRefundStrategy struct {
	prefer bool
}

func NewEasyRefundStrategy(preferRefund bool) *EasyRefundStrategy {
	return &EasyRefundStrategy{prefer: preferRefund}
}

func (s *EasyRefundStrategy) Compare(a, b PaymentMethod) bool {
	if !s.prefer {
		return false
	}
	// TODO: a wins if it has easy refund and b doesn't
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
	refundStrat := NewEasyRefundStrategy(preferEasyRefund)
	rewardsStrat := RewardsMaximizer{}
	// Refund filter takes priority over rewards
	composite := NewCompositeStrategy([]RankingStrategy{refundStrat, rewardsStrat})
	return NewPaymentRanker(composite).Rank(methods)
}
