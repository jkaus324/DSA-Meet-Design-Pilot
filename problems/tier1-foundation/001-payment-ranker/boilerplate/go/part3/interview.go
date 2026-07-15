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

// TODO: design and implement your solution.
// Required free functions:
//   func rank_by_rewards(methods []PaymentMethod) []PaymentMethod
//   func rank_by_low_fee(methods []PaymentMethod) []PaymentMethod
//   func rank_by_trust(methods []PaymentMethod) []PaymentMethod
//   func rank_composite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod
//   func rank_with_refund_filter(methods []PaymentMethod, preferEasyRefund bool) []PaymentMethod

func rank_by_rewards(methods []PaymentMethod) []PaymentMethod {
	// TODO: write your solution
	return methods
}

func rank_by_low_fee(methods []PaymentMethod) []PaymentMethod {
	// TODO: write your solution
	return methods
}

func rank_by_trust(methods []PaymentMethod) []PaymentMethod {
	// TODO: write your solution
	return methods
}

func rank_composite(methods []PaymentMethod, criteria []RankingStrategy) []PaymentMethod {
	// TODO: write your solution
	return methods
}

func rank_with_refund_filter(methods []PaymentMethod, preferEasyRefund bool) []PaymentMethod {
	// TODO: write your solution
	return methods
}
