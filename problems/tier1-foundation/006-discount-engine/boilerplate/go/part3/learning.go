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

// ─── Discount Interface ───────────────────────────────────────────────────────

type Discount interface {
	Apply(cart []CartItem) float64
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────

type PercentageDiscount struct {
	Percentage float64
}

func (d *PercentageDiscount) Apply(cart []CartItem) float64 {
	// TODO: Sum all item.Price * item.Quantity, then apply d.Percentage discount
	return 0.0
}

type FlatDiscount struct {
	Amount float64
}

func (d *FlatDiscount) Apply(cart []CartItem) float64 {
	// TODO: Sum cart total, subtract d.Amount (floor at 0)
	return 0.0
}

type BuyXGetYDiscount struct {
	BuyCount  int
	FreeCount int
}

func (d *BuyXGetYDiscount) Apply(cart []CartItem) float64 {
	// TODO: For each item, compute paid quantity using BuyCount + FreeCount group logic
	// groupSize = BuyCount + FreeCount; paidItems = groups*BuyCount + min(remainder, BuyCount)
	return 0.0
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type StackedDiscount struct {
	Discounts []Discount
}

func (d *StackedDiscount) Apply(cart []CartItem) float64 {
	// TODO: Compute initial total from cart, then apply each discount in d.Discounts sequentially
	// HINT: Create a temp cart with a single item {Price: current, Quantity: 1} and call disc.Apply() on it
	return 0.0
}

// ─── Engine ───────────────────────────────────────────────────────────────────

type DiscountEngine struct {
	discount Discount
}

func NewDiscountEngine(d Discount) *DiscountEngine {
	return &DiscountEngine{discount: d}
}

func (e *DiscountEngine) ComputeTotal(cart []CartItem) float64 {
	return e.discount.Apply(cart)
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

func ApplyWithEligibility(cart []CartItem, discount Discount,
	minCartValue float64, requireFirstTimeUser bool,
	user UserContext, eligibleCategory string) float64 {
	rawTotal := 0.0
	for _, item := range cart {
		rawTotal += item.Price * float64(item.Quantity)
	}

	// TODO: Rule 1 — if rawTotal < minCartValue, return rawTotal
	// TODO: Rule 2 — if requireFirstTimeUser && !user.IsFirstTimeUser, return rawTotal
	// TODO: Rule 3 — if eligibleCategory is non-empty, split cart and discount only matching items

	return discount.Apply(cart)
}
