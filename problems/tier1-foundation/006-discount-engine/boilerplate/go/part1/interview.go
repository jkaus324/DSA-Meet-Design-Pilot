package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type CartItem struct {
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement a DiscountEngine that:
//   1. Applies discount strategies to a shopping cart
//   2. Allows new discount types to be added WITHOUT modifying
//      the engine itself
//
// Think about:
//   - What abstraction lets you swap discount logic at runtime?
//   - How would you add a 4th discount type with zero changes
//     to existing code?
//
// Entry points (must exist for tests):
//   func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64
//   func ApplyFlatDiscount(cart []CartItem, amount float64) float64
//   func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64
//
// ─────────────────────────────────────────────────────────────────────────────

func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64 {
	return 0
}

func ApplyFlatDiscount(cart []CartItem, amount float64) float64 {
	return 0
}

func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64 {
	return 0
}
