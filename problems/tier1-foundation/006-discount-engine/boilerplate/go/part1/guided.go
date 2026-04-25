package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type CartItem struct {
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ─── Discount Interface ───────────────────────────────────────────────────────
// HINT: This interface lets you swap discount logic at runtime.
// What method signature would let you compute a discounted total?

type Discount interface {
	// HINT: what method computes the discounted total for a cart?
	// Apply(cart []CartItem) float64
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────
// TODO: Implement a strategy for each discount type:
//   - Percentage discount (reduce total by X%)
//   - Flat discount (subtract fixed amount, minimum 0)
//   - Buy-X-Get-Y (for every X+Y items, Y are free)

// ─── Engine ───────────────────────────────────────────────────────────────────
// TODO: Implement a DiscountEngine struct that:
//   - Accepts any Discount strategy
//   - Has a ComputeTotal() method that returns the discounted amount
//   - Does NOT know about specific discount types

// type DiscountEngine struct {
//     discount Discount
// }
// func NewDiscountEngine(d Discount) *DiscountEngine { ... }
// func (e *DiscountEngine) ComputeTotal(cart []CartItem) float64 { ... }

// ─── Test Entry Points (must exist for tests) ─────────────────────────────────

func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64 {
	return 0
}

func ApplyFlatDiscount(cart []CartItem, amount float64) float64 {
	return 0
}

func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64 {
	return 0
}
