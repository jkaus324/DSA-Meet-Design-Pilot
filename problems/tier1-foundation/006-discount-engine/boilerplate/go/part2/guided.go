package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type CartItem struct {
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ─── Discount Interface ───────────────────────────────────────────────────────

type Discount interface {
	Apply(cart []CartItem) float64
}

// ─── Existing Strategies ──────────────────────────────────────────────────────
// TODO: Copy your Part 1 strategies here (or extend them)

type PercentageDiscount struct {
	Percentage float64
}

func (d *PercentageDiscount) Apply(cart []CartItem) float64 {
	return 0 // TODO: implement
}

type FlatDiscount struct {
	Amount float64
}

func (d *FlatDiscount) Apply(cart []CartItem) float64 {
	return 0 // TODO: implement
}

type BuyXGetYDiscount struct {
	BuyCount  int
	FreeCount int
}

func (d *BuyXGetYDiscount) Apply(cart []CartItem) float64 {
	return 0 // TODO: implement
}

// ─── NEW: Stacked Discount (Decorator) ───────────────────────────────────────
// HINT: A StackedDiscount holds a slice of other discounts.
// It applies the first discount to the raw total, then applies the second
// discount to the result, and so on. This is the Decorator pattern.

type StackedDiscount struct {
	Discounts []Discount
}

func (d *StackedDiscount) Apply(cart []CartItem) float64 {
	// TODO: Compute raw total from cart.
	// Then apply each discount sequentially to the running total.
	// HINT: Create a temporary single-item cart with the running total.
	return 0
}

// ─── Engine ───────────────────────────────────────────────────────────────────

type DiscountEngine struct {
	discount Discount
}

func NewDiscountEngine(d Discount) *DiscountEngine {
	return &DiscountEngine{discount: d}
}

func (e *DiscountEngine) ComputeTotal(cart []CartItem) float64 {
	return 0 // TODO: delegate to e.discount.Apply()
}

// ─── Test Entry Points ────────────────────────────────────────────────────────

func ApplyPercentageDiscount(cart []CartItem, percentage float64) float64 {
	d := &PercentageDiscount{Percentage: percentage}
	return NewDiscountEngine(d).ComputeTotal(cart)
}

func ApplyFlatDiscount(cart []CartItem, amount float64) float64 {
	d := &FlatDiscount{Amount: amount}
	return NewDiscountEngine(d).ComputeTotal(cart)
}

func ApplyBogo(cart []CartItem, buyCount int, freeCount int) float64 {
	d := &BuyXGetYDiscount{BuyCount: buyCount, FreeCount: freeCount}
	return NewDiscountEngine(d).ComputeTotal(cart)
}

func ApplyStackedDiscounts(cart []CartItem, discounts []Discount) float64 {
	d := &StackedDiscount{Discounts: discounts}
	return NewDiscountEngine(d).ComputeTotal(cart)
}
