// Payment ranker — Strategy pattern reference solution (Go).
//
// Convention (must match harness/go/codegen.py):
//   • package main (compiled alongside the generated runner.go)
//   • data-class structs have fields in spec-declared order, named exactly as
//     in spec.yaml (lowercase/unexported is fine — same package)
//   • free functions use the spec's snake_case names
package main

type PaymentMethod struct {
	name               string
	cashbackRate       float64
	transactionFee     float64
	usageCount         int
	easyRefundEligible bool
}

// RankingStrategy reports whether a should rank strictly before b.
type RankingStrategy interface {
	compare(a, b PaymentMethod) bool
}

type RewardsMaximizer struct{}

func (RewardsMaximizer) compare(a, b PaymentMethod) bool { return a.cashbackRate > b.cashbackRate }

type LowFeeSeeker struct{}

func (LowFeeSeeker) compare(a, b PaymentMethod) bool { return a.transactionFee < b.transactionFee }

type TrustBasedRanker struct{}

func (TrustBasedRanker) compare(a, b PaymentMethod) bool { return a.usageCount > b.usageCount }

type CompositeStrategy struct {
	criteria []RankingStrategy
}

func (c CompositeStrategy) compare(a, b PaymentMethod) bool {
	for _, s := range c.criteria {
		if s.compare(a, b) {
			return true
		}
		if s.compare(b, a) {
			return false
		}
	}
	return false
}

type EasyRefundStrategy struct {
	prefer bool
}

func (e EasyRefundStrategy) compare(a, b PaymentMethod) bool {
	if !e.prefer {
		return false
	}
	return a.easyRefundEligible && !b.easyRefundEligible
}

func rankWith(strategy RankingStrategy, methods []PaymentMethod) []PaymentMethod {
	// Stable insertion sort driven by the strict-less comparator.
	result := []PaymentMethod{}
	for _, m := range methods {
		inserted := false
		for i, existing := range result {
			if strategy.compare(m, existing) {
				result = append(result[:i], append([]PaymentMethod{m}, result[i:]...)...)
				inserted = true
				break
			}
		}
		if !inserted {
			result = append(result, m)
		}
	}
	return result
}

func rank_by_rewards(methods []PaymentMethod) []PaymentMethod {
	return rankWith(RewardsMaximizer{}, methods)
}

func rank_by_low_fee(methods []PaymentMethod) []PaymentMethod {
	return rankWith(LowFeeSeeker{}, methods)
}

func rank_by_trust(methods []PaymentMethod) []PaymentMethod {
	return rankWith(TrustBasedRanker{}, methods)
}

func rank_composite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod {
	return rankWith(CompositeStrategy{criteria: criteria}, methods)
}

func rank_with_refund_filter(methods []PaymentMethod, preferEasyRefund bool) []PaymentMethod {
	composite := CompositeStrategy{criteria: []RankingStrategy{
		EasyRefundStrategy{prefer: preferEasyRefund},
		RewardsMaximizer{},
	}}
	return rankWith(composite, methods)
}
