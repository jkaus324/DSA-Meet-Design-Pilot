package main

// Data class (given).
type CartItem struct {
	name string
	price float64
	quantity int
	category string
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func apply_percentage_discount(cart []CartItem, percentage float64) float64 {
	// TODO: write your solution
	return 0.0
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func apply_flat_discount(cart []CartItem, amount float64) float64 {
	// TODO: write your solution
	return 0.0
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func apply_bogo(cart []CartItem, buyCount int, freeCount int) float64 {
	// TODO: write your solution
	return 0.0
}
