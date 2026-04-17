package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type CartItem struct {
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ─── NEW in Extension 1 ───────────────────────────────────────────────────────
//
// The product team now wants STACKED discounts:
// apply a coupon, then a seasonal discount on top, then membership on top.
//
// Think about:
//   - How do you chain discounts without modifying existing strategies?
//   - What if the product team adds a 5th discount layer tomorrow?
//   - Is your Part 1 design extensible enough to handle this?
//
// Entry points (must exist for tests):
//   func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64
//   func ApplyFlatDiscount(cart []CartItem, amount float64) float64
//   func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64
//   func ApplyStackedDiscounts(cart []CartItem, discounts []Discount) float64
//
// ─────────────────────────────────────────────────────────────────────────────

type Discount interface {
	Apply(cart []CartItem) float64
}

func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64 {
	return 0
}

func ApplyFlatDiscount(cart []CartItem, amount float64) float64 {
	return 0
}

func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64 {
	return 0
}

func ApplyStackedDiscounts(cart []CartItem, discounts []Discount) float64 {
	return 0
}
