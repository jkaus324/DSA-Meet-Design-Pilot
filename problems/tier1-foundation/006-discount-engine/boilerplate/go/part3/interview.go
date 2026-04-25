package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type CartItem struct {
	Name     string
	Price    float64
	Quantity int
	Category string
}

type UserContext struct {
	IsFirstTimeUser bool
}

// ─── NEW in Extension 2 ───────────────────────────────────────────────────────
//
// The compliance team wants to add eligibility rules for discounts:
// - Minimum cart value threshold
// - First-time user only
// - Category-specific discounts (only items in a certain category)
//
// Think about:
//   - Are eligibility rules a new discount type, a decorator, or a filter?
//   - How does your existing design accommodate this without modification?
//   - What if the eligibility rules themselves need to be composable?
//
// Entry points (must exist for tests):
//   func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64
//   func ApplyFlatDiscount(cart []CartItem, amount float64) float64
//   func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64
//   func ApplyStackedDiscounts(cart []CartItem, discounts []Discount) float64
//   func ApplyWithEligibility(cart []CartItem, discount Discount,
//       minCartValue float64, requireFirstTimeUser bool,
//       user UserContext, eligibleCategory string) float64
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

func ApplyWithEligibility(cart []CartItem, discount Discount,
	minCartValue float64, requireFirstTimeUser bool,
	user UserContext, eligibleCategory string) float64 {
	return 0
}
