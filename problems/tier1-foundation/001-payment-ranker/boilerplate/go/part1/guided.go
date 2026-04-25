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
// HINT: This interface lets you swap ranking logic at runtime.
// What method signature would let you compare two PaymentMethods?

type RankingStrategy interface {
	// HINT: return true if 'a' should rank higher than 'b'
	Compare(a, b PaymentMethod) bool
}

// ─── Concrete Strategies ─────────────────────────────────────────────────────
// TODO: Implement a strategy for each ranking criterion:
//   - RewardsMaximizer  (highest cashback first)
//   - LowFeeSeeker      (lowest transaction fee first)
//   - TrustBasedRanker  (highest usage count first)

// ─── Ranker ──────────────────────────────────────────────────────────────────
// TODO: Implement a PaymentRanker struct that:
//   - Accepts any RankingStrategy
//   - Has a Rank() method that returns sorted payment methods
//   - Does NOT know about specific ranking criteria

type PaymentRanker struct {
	// HINT: store a RankingStrategy here
}

// TODO: what parameter should NewPaymentRanker accept?
func NewPaymentRanker(s RankingStrategy) *PaymentRanker {
	return &PaymentRanker{}
}

func (r *PaymentRanker) Rank(methods []PaymentMethod) []PaymentMethod {
	// HINT: use sort.Slice with r.strategy.Compare()
	// sort.Slice(methods, func(i, j int) bool { return r.strategy.Compare(methods[i], methods[j]) })
	_ = sort.Search // keep sort import alive
	return methods
}

// ─── Test Entry Points (must exist for tests to compile) ─────────────────────
// Your solution must provide these functions:

func RankByRewards(methods []PaymentMethod) []PaymentMethod {
	return methods // TODO: use PaymentRanker with RewardsMaximizer
}

func RankByLowFee(methods []PaymentMethod) []PaymentMethod {
	return methods // TODO: use PaymentRanker with LowFeeSeeker
}

func RankByTrust(methods []PaymentMethod) []PaymentMethod {
	return methods // TODO: use PaymentRanker with TrustBasedRanker
}
