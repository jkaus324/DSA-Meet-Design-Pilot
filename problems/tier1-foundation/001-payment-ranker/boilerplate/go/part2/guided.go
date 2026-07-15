package main

// Data class (given).
type PaymentMethod struct {
	name string
	cashbackRate float64
	transactionFee float64
	usageCount int
	easyRefundEligible bool
}

// RankingStrategy — implement this interface with your own strategy types.
type RankingStrategy interface {
	// TODO: define the method(s) your strategies share.
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func rank_by_rewards(methods []PaymentMethod) []PaymentMethod {
	// TODO: write your solution
	return methods
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func rank_by_low_fee(methods []PaymentMethod) []PaymentMethod {
	// TODO: write your solution
	return methods
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func rank_by_trust(methods []PaymentMethod) []PaymentMethod {
	// TODO: write your solution
	return methods
}

// HINT: think about how to compose multiple criteria into a single decision.
func rank_composite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod {
	// TODO: write your solution
	return methods
}
