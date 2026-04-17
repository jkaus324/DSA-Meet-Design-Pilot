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
	total := 0.0
	for _, item := range cart {
		total += item.Price * float64(item.Quantity)
	}
	return total * (1.0 - d.Percentage/100.0)
}

type FlatDiscount struct {
	Amount float64
}

func (d *FlatDiscount) Apply(cart []CartItem) float64 {
	total := 0.0
	for _, item := range cart {
		total += item.Price * float64(item.Quantity)
	}
	if total-d.Amount < 0 {
		return 0
	}
	return total - d.Amount
}

type BuyXGetYDiscount struct {
	BuyCount  int
	FreeCount int
}

func (d *BuyXGetYDiscount) Apply(cart []CartItem) float64 {
	total := 0.0
	for _, item := range cart {
		groupSize := d.BuyCount + d.FreeCount
		groups := item.Quantity / groupSize
		remainder := item.Quantity % groupSize
		paidItems := groups*d.BuyCount + minInt(remainder, d.BuyCount)
		total += item.Price * float64(paidItems)
	}
	return total
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
	current := 0.0
	for _, item := range cart {
		current += item.Price * float64(item.Quantity)
	}
	for _, disc := range d.Discounts {
		temp := []CartItem{{Name: "subtotal", Price: current, Quantity: 1, Category: ""}}
		current = disc.Apply(temp)
	}
	return current
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
