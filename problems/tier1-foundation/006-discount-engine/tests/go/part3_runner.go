package main

import (
	"fmt"
	"math"
)

func approx3(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func part3Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	// Test 1: minimum cart value met — discount applies
	test("test_min_cart_value_met", func() {
		cart := []CartItem{
			{Name: "Laptop", Price: 50000.0, Quantity: 1, Category: "electronics"},
		}
		pct := &PercentageDiscount{Percentage: 10.0}
		user := UserContext{IsFirstTimeUser: false}
		result := ApplyWithEligibility(cart, pct, 1000.0, false, user, "")
		// Total 50000 >= 1000, discount applies: 45000
		if !approx3(result, 45000.0) {
			panic(fmt.Sprintf("expected 45000.0, got %f", result))
		}
	})

	// Test 2: minimum cart value NOT met — discount skipped
	test("test_min_cart_value_not_met", func() {
		cart := []CartItem{
			{Name: "Sticker", Price: 50.0, Quantity: 1, Category: "accessories"},
		}
		pct := &PercentageDiscount{Percentage: 10.0}
		user := UserContext{IsFirstTimeUser: false}
		result := ApplyWithEligibility(cart, pct, 1000.0, false, user, "")
		// Total 50 < 1000, discount skipped: 50
		if !approx3(result, 50.0) {
			panic(fmt.Sprintf("expected 50.0, got %f", result))
		}
	})

	// Test 3: first-time user required and user IS first-time
	test("test_first_time_user_eligible", func() {
		cart := []CartItem{
			{Name: "Phone", Price: 20000.0, Quantity: 1, Category: "electronics"},
		}
		flat := &FlatDiscount{Amount: 2000.0}
		user := UserContext{IsFirstTimeUser: true}
		result := ApplyWithEligibility(cart, flat, 0.0, true, user, "")
		// First-time user, discount applies: 20000 - 2000 = 18000
		if !approx3(result, 18000.0) {
			panic(fmt.Sprintf("expected 18000.0, got %f", result))
		}
	})

	// Test 4: first-time user required but user is NOT first-time
	test("test_first_time_user_not_eligible", func() {
		cart := []CartItem{
			{Name: "Phone", Price: 20000.0, Quantity: 1, Category: "electronics"},
		}
		flat := &FlatDiscount{Amount: 2000.0}
		user := UserContext{IsFirstTimeUser: false}
		result := ApplyWithEligibility(cart, flat, 0.0, true, user, "")
		// Not first-time, discount skipped: 20000
		if !approx3(result, 20000.0) {
			panic(fmt.Sprintf("expected 20000.0, got %f", result))
		}
	})

	// Test 5: category-specific discount — only electronics discounted
	test("test_category_specific_discount", func() {
		cart := []CartItem{
			{Name: "Laptop", Price: 50000.0, Quantity: 1, Category: "electronics"},
			{Name: "Phone Case", Price: 500.0, Quantity: 2, Category: "accessories"},
		}
		pct := &PercentageDiscount{Percentage: 10.0}
		user := UserContext{IsFirstTimeUser: false}
		result := ApplyWithEligibility(cart, pct, 0.0, false, user, "electronics")
		// Electronics: 50000 * 0.9 = 45000, Accessories: 1000 full price → 46000
		if !approx3(result, 46000.0) {
			panic(fmt.Sprintf("expected 46000.0, got %f", result))
		}
	})

	// Test 6: all rules combined — min cart met, first-time, category filter
	test("test_all_rules_combined", func() {
		cart := []CartItem{
			{Name: "Laptop", Price: 50000.0, Quantity: 1, Category: "electronics"},
			{Name: "T-Shirt", Price: 1000.0, Quantity: 3, Category: "clothing"},
		}
		pct := &PercentageDiscount{Percentage: 20.0}
		user := UserContext{IsFirstTimeUser: true}
		result := ApplyWithEligibility(cart, pct, 5000.0, true, user, "electronics")
		// Total = 53000 >= 5000, first-time user OK
		// Electronics: 50000 * 0.8 = 40000, Clothing: 3000 full price → 43000
		if !approx3(result, 43000.0) {
			panic(fmt.Sprintf("expected 43000.0, got %f", result))
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
