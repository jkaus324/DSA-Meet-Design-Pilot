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
		paidItems := groups*d.BuyCount + min(remainder, d.BuyCount)
		total += item.Price * float64(paidItems)
	}
	return total
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ─── StackedDiscount (Decorator) ─────────────────────────────────────────────
// Chains multiple discounts: applies each to the result of the previous.

type StackedDiscount struct {
	Discounts []Discount
}

func (d *StackedDiscount) Apply(cart []CartItem) float64 {
	current := 0.0
	for _, item := range cart {
		current += item.Price * float64(item.Quantity)
	}

	for _, disc := range d.Discounts {
		// TODO: Apply each discount to the running total.
		// HINT: Create a temporary cart with a single item representing
		// the current subtotal, then call disc.Apply() on it.
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

func (e *DiscountEngine) SetDiscount(d Discount) {
	e.discount = d
}

func (e *DiscountEngine) ComputeTotal(cart []CartItem) float64 {
	// TODO: Use e.discount.Apply() to compute the final total
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
