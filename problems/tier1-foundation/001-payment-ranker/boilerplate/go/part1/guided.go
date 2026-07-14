package main

// Data class (given).
type PaymentMethod struct {
	name string
	cashbackRate float64
	transactionFee float64
	usageCount int
	easyRefundEligible bool
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
