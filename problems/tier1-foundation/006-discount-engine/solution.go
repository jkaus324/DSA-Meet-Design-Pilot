// Discount Engine — Strategy + Decorator reference solution (Go).
package main

type CartItem struct {
	name     string
	price    float64
	quantity int
	category string
}

func cartTotal(cart []CartItem) float64 {
	total := 0.0
	for _, i := range cart {
		total += i.price * float64(i.quantity)
	}
	return total
}

func percentageDiscount(cart []CartItem, pct float64) float64 {
	return cartTotal(cart) * (1.0 - pct/100.0)
}

func apply_percentage_discount(cart []CartItem, percentage float64) float64 {
	return percentageDiscount(cart, percentage)
}

func apply_flat_discount(cart []CartItem, amount float64) float64 {
	total := cartTotal(cart)
	if total-amount < 0.0 {
		return 0.0
	}
	return total - amount
}

func apply_bogo(cart []CartItem, buyCount, freeCount int) float64 {
	group := buyCount + freeCount
	total := 0.0
	for _, it := range cart {
		groups := it.quantity / group
		remainder := it.quantity % group
		paid := groups*buyCount + min(remainder, buyCount)
		total += float64(paid) * it.price
	}
	return total
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func apply_percentage_with_eligibility(cart []CartItem, percentage, minCartValue float64,
	requireFirstTimeUser, isFirstTimeUser bool, eligibleCategory string) float64 {
	raw := cartTotal(cart)
	if raw < minCartValue {
		return raw
	}
	if requireFirstTimeUser && !isFirstTimeUser {
		return raw
	}
	if eligibleCategory != "" {
		var eligible []CartItem
		nonEligible := 0.0
		for _, i := range cart {
			if i.category == eligibleCategory {
				eligible = append(eligible, i)
			} else {
				nonEligible += i.price * float64(i.quantity)
			}
		}
		return percentageDiscount(eligible, percentage) + nonEligible
	}
	return percentageDiscount(cart, percentage)
}

func reset_service() {}
