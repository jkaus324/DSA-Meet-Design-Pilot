package main

import (
	"fmt"
	"math"
)

func approx2(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func part2Tests() int {
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

	// Test 1: stack percentage then flat
	test("test_stack_percentage_then_flat", func() {
		cart := []CartItem{
			{Name: "Laptop", Price: 10000.0, Quantity: 1, Category: "electronics"},
		}
		// Total = 10000 → 10% off → 9000 → Rs.500 flat off → 8500
		pct := &PercentageDiscount{Percentage: 10.0}
		flat := &FlatDiscount{Amount: 500.0}
		result := ApplyStackedDiscounts(cart, []Discount{pct, flat})
		if !approx2(result, 8500.0) {
			panic(fmt.Sprintf("expected 8500.0, got %f", result))
		}
	})

	// Test 2: stack flat then percentage
	test("test_stack_flat_then_percentage", func() {
		cart := []CartItem{
			{Name: "Laptop", Price: 10000.0, Quantity: 1, Category: "electronics"},
		}
		// Total = 10000 → Rs.500 flat off → 9500 → 10% off → 8550
		flat := &FlatDiscount{Amount: 500.0}
		pct := &PercentageDiscount{Percentage: 10.0}
		result := ApplyStackedDiscounts(cart, []Discount{flat, pct})
		if !approx2(result, 8550.0) {
			panic(fmt.Sprintf("expected 8550.0, got %f", result))
		}
	})

	// Test 3: stack three discounts
	test("test_stack_three_discounts", func() {
		cart := []CartItem{
			{Name: "Phone", Price: 20000.0, Quantity: 1, Category: "electronics"},
		}
		// Total = 20000 → 10% coupon → 18000 → Rs.1000 seasonal → 17000 → 5% membership → 16150
		coupon := &PercentageDiscount{Percentage: 10.0}
		seasonal := &FlatDiscount{Amount: 1000.0}
		membership := &PercentageDiscount{Percentage: 5.0}
		result := ApplyStackedDiscounts(cart, []Discount{coupon, seasonal, membership})
		if !approx2(result, 16150.0) {
			panic(fmt.Sprintf("expected 16150.0, got %f", result))
		}
	})

	// Test 4: single discount in stack behaves like direct application
	test("test_single_discount_stack", func() {
		cart := []CartItem{
			{Name: "Book", Price: 500.0, Quantity: 2, Category: "books"},
		}
		pct := &PercentageDiscount{Percentage: 20.0}
		stacked := ApplyStackedDiscounts(cart, []Discount{pct})
		direct := ApplyPercentageDiscount(cart, 20.0)
		if !approx2(stacked, direct) {
			panic(fmt.Sprintf("expected stacked %f == direct %f", stacked, direct))
		}
	})

	// Test 5: stacked discounts that reduce to zero
	test("test_stack_reduces_to_zero", func() {
		cart := []CartItem{
			{Name: "Sticker", Price: 100.0, Quantity: 1, Category: "accessories"},
		}
		flat1 := &FlatDiscount{Amount: 60.0}
		flat2 := &FlatDiscount{Amount: 60.0}
		result := ApplyStackedDiscounts(cart, []Discount{flat1, flat2})
		if !approx2(result, 0.0) {
			panic(fmt.Sprintf("expected 0.0, got %f", result))
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
