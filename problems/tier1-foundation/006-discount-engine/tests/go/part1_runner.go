package main

import (
	"fmt"
	"math"
)

func approx(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func part1Tests() int {
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

	// Test 1: percentage discount — 10% off
	test("test_percentage_discount", func() {
		cart := []CartItem{
			{Name: "Laptop", Price: 50000.0, Quantity: 1, Category: "electronics"},
			{Name: "Phone Case", Price: 500.0, Quantity: 2, Category: "accessories"},
		}
		// Total = 50000 + 1000 = 51000, 10% off = 45900
		result := ApplyPercentageDiscount(cart, 10.0)
		if !approx(result, 45900.0) {
			panic(fmt.Sprintf("expected 45900.0, got %f", result))
		}
	})

	// Test 2: flat discount — Rs. 200 off
	test("test_flat_discount", func() {
		cart := []CartItem{
			{Name: "Headphones", Price: 2000.0, Quantity: 1, Category: "electronics"},
			{Name: "Cable", Price: 300.0, Quantity: 3, Category: "accessories"},
		}
		// Total = 2000 + 900 = 2900, flat 200 off = 2700
		result := ApplyFlatDiscount(cart, 200.0)
		if !approx(result, 2700.0) {
			panic(fmt.Sprintf("expected 2700.0, got %f", result))
		}
	})

	// Test 3: flat discount exceeds total — should return 0
	test("test_flat_discount_exceeds_total", func() {
		cart := []CartItem{
			{Name: "Sticker", Price: 50.0, Quantity: 1, Category: "accessories"},
		}
		result := ApplyFlatDiscount(cart, 200.0)
		if !approx(result, 0.0) {
			panic(fmt.Sprintf("expected 0.0, got %f", result))
		}
	})

	// Test 4: buy 2 get 1 free
	test("test_bogo_exact_groups", func() {
		cart := []CartItem{
			{Name: "T-Shirt", Price: 500.0, Quantity: 6, Category: "clothing"},
		}
		// 6 / (2+1) = 2 groups, each group pays for 2 → 4 paid items
		result := ApplyBogo(cart, 2, 1)
		if !approx(result, 2000.0) {
			panic(fmt.Sprintf("expected 2000.0, got %f", result))
		}
	})

	// Test 5: buy 2 get 1 with remainder
	test("test_bogo_with_remainder", func() {
		cart := []CartItem{
			{Name: "Socks", Price: 200.0, Quantity: 5, Category: "clothing"},
		}
		// 5 / 3 = 1 group (pay 2), remainder 2 (pay 2) → 4 paid items
		result := ApplyBogo(cart, 2, 1)
		if !approx(result, 800.0) {
			panic(fmt.Sprintf("expected 800.0, got %f", result))
		}
	})

	// Test 6: empty cart
	test("test_empty_cart", func() {
		var empty []CartItem
		if !approx(ApplyPercentageDiscount(empty, 10.0), 0.0) {
			panic("expected 0.0 for empty cart percentage discount")
		}
		if !approx(ApplyFlatDiscount(empty, 100.0), 0.0) {
			panic("expected 0.0 for empty cart flat discount")
		}
		if !approx(ApplyBogo(empty, 2, 1), 0.0) {
			panic("expected 0.0 for empty cart bogo")
		}
	})

	// Test 7: single item, no discount effect (0%)
	test("test_zero_percentage", func() {
		cart := []CartItem{{Name: "Book", Price: 300.0, Quantity: 1, Category: "books"}}
		result := ApplyPercentageDiscount(cart, 0.0)
		if !approx(result, 300.0) {
			panic(fmt.Sprintf("expected 300.0, got %f", result))
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
