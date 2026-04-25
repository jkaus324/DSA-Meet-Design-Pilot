package main

import "fmt"

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

	// Test 1: composite ranking — cashback first, then fee as tiebreaker
	test("test_composite_cashback_then_fee", func() {
		methods := []PaymentMethod{
			{"Card A", 0.10, 8.0, 300}, // 10% cashback, high fee
			{"Card B", 0.10, 3.0, 400}, // 10% cashback, low fee
			{"Card C", 0.05, 1.0, 200}, // 5% cashback
		}
		ranked := RankComposite(methods, []RankingStrategy{
			RewardsMaximizer{},
			LowFeeSeeker{},
		})
		if len(ranked) != 3 {
			panic("expected 3 results")
		}
		if ranked[0].Name != "Card B" {
			panic("expected Card B first (tied cashback, lower fee wins)")
		}
		if ranked[1].Name != "Card A" {
			panic("expected Card A second (tied cashback, higher fee loses)")
		}
		if ranked[2].Name != "Card C" {
			panic("expected Card C third (lower cashback)")
		}
	})

	// Test 2: composite ranking — trust first, then cashback
	test("test_composite_trust_then_cashback", func() {
		methods := []PaymentMethod{
			{"UPI", 0.01, 0.0, 1000},
			{"Card A", 0.10, 5.0, 200},
			{"Card B", 0.05, 3.0, 1000}, // tied trust with UPI
		}
		ranked := RankComposite(methods, []RankingStrategy{
			TrustBasedRanker{},
			RewardsMaximizer{},
		})
		if len(ranked) != 3 {
			panic("expected 3 results")
		}
		// UPI and Card B both have 1000 uses — tiebreak by cashback
		// Card B has 5% cashback > UPI's 1%
		if ranked[0].Name != "Card B" {
			panic("expected Card B first")
		}
		if ranked[1].Name != "UPI" {
			panic("expected UPI second")
		}
		if ranked[2].Name != "Card A" {
			panic("expected Card A third")
		}
	})

	// Test 3: single criterion composite behaves like that criterion alone
	test("test_single_criterion_composite", func() {
		methods := []PaymentMethod{
			{"Card X", 0.02, 5.0, 100},
			{"Card Y", 0.08, 3.0, 200},
		}
		composite := RankComposite(methods, []RankingStrategy{RewardsMaximizer{}})
		direct := RankByRewards(methods)
		if composite[0].Name != direct[0].Name {
			panic("expected same first element")
		}
		if composite[1].Name != direct[1].Name {
			panic("expected same second element")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
