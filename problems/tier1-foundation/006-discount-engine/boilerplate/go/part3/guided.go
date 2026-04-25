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

// ─── Existing Strategies ──────────────────────────────────────────────────────

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

type StackedDiscount struct {
	Discounts []Discount
}

func (d *StackedDiscount) Apply(cart []CartItem) float64 {
	return 0 // TODO: implement
}

// ─── NEW: Eligibility Rules ───────────────────────────────────────────────────
// HINT: When a rule is not met, the discount is SKIPPED (original total returned).
// When eligibleCategory is non-empty, only items in that category are discounted.

type DiscountEngine struct {
	discount Discount
}

func NewDiscountEngine(d Discount) *DiscountEngine {
	return &DiscountEngine{discount: d}
}

func (e *DiscountEngine) ComputeTotal(cart []CartItem) float64 {
	return 0 // TODO: implement
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
	// TODO: Check eligibility rules:
	// 1. If raw total < minCartValue, return raw total (skip discount)
	// 2. If requireFirstTimeUser && !user.IsFirstTimeUser, return raw total
	// 3. If eligibleCategory is non-empty, only discount items in that category
	return 0
}
