package main

// Data class (given — do not modify).
type PaymentMethod struct {
	name string
	cashbackRate float64
	transactionFee float64
	usageCount int
	easyRefundEligible bool
}

func rank_by_rewards(methods []PaymentMethod) []PaymentMethod {
	// TODO: implement this
	return methods
}

func rank_by_low_fee(methods []PaymentMethod) []PaymentMethod {
	// TODO: implement this
	return methods
}

func rank_by_trust(methods []PaymentMethod) []PaymentMethod {
	// TODO: implement this
	return methods
}
