package main

import "fmt"

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

	// Test 1: preferEasyRefund=true puts eligible methods first
	test("test_easy_refund_preferred", func() {
		methods := []PaymentMethod{
			{"Card A", 0.10, 5.0, 300, false}, // high cashback, no easy refund
			{"Card B", 0.02, 2.0, 500, true},  // low cashback, easy refund
			{"Card C", 0.05, 3.0, 400, true},  // medium cashback, easy refund
		}
		ranked := RankWithRefundFilter(methods, true)
		if len(ranked) != 3 {
			panic("expected 3 results")
		}
		// Both B and C have easy refund, so they come first (in cashback order)
		if !ranked[0].EasyRefundEligible {
			panic("expected easy-refund eligible method first")
		}
		if !ranked[1].EasyRefundEligible {
			panic("expected easy-refund eligible method second")
		}
		if ranked[2].Name != "Card A" {
			panic("expected Card A last (no easy refund)")
		}
	})

	// Test 2: preferEasyRefund=false should not reorder by refund eligibility
	test("test_refund_filter_disabled", func() {
		methods := []PaymentMethod{
			{"Card A", 0.10, 5.0, 300, false},
			{"Card B", 0.02, 2.0, 500, true},
		}
		ranked := RankWithRefundFilter(methods, false)
		// Without refund preference, Card A should still win (higher cashback)
		if ranked[0].Name != "Card A" {
			panic("expected Card A first (higher cashback, refund filter disabled)")
		}
	})

	// Test 3: all methods have easy refund — order falls back to other criteria
	test("test_all_refund_eligible_tiebreak", func() {
		methods := []PaymentMethod{
			{"Card A", 0.10, 5.0, 300, true},
			{"Card B", 0.02, 2.0, 500, true},
		}
		ranked := RankWithRefundFilter(methods, true)
		// All have easy refund, so tiebreak by cashback
		if ranked[0].Name != "Card A" {
			panic("expected Card A first (higher cashback)")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
