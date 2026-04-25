package main

import "fmt"

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

	// Test 1: rank_by_rewards — highest cashback first
	test("test_rewards_ranking", func() {
		methods := []PaymentMethod{
			{"UPI", 0.01, 0.0, 1000},
			{"Credit Card A", 0.05, 5.0, 500},
			{"Credit Card B", 0.10, 8.0, 300},
		}
		ranked := RankByRewards(methods)
		if len(ranked) != 3 {
			panic("expected 3 results")
		}
		if ranked[0].Name != "Credit Card B" {
			panic("expected Credit Card B first (10% cashback)")
		}
		if ranked[1].Name != "Credit Card A" {
			panic("expected Credit Card A second (5% cashback)")
		}
		if ranked[2].Name != "UPI" {
			panic("expected UPI third (1% cashback)")
		}
	})

	// Test 2: rank_by_low_fee — lowest fee first
	test("test_low_fee_ranking", func() {
		methods := []PaymentMethod{
			{"Debit Card", 0.0, 2.0, 800},
			{"Credit Card A", 0.05, 5.0, 500},
			{"Credit Card B", 0.10, 8.0, 300},
			{"UPI", 0.01, 0.0, 1000},
		}
		ranked := RankByLowFee(methods)
		if len(ranked) != 4 {
			panic("expected 4 results")
		}
		if ranked[0].Name != "UPI" {
			panic("expected UPI first (0 fee)")
		}
		if ranked[1].Name != "Debit Card" {
			panic("expected Debit Card second (2.0 fee)")
		}
		if ranked[2].Name != "Credit Card A" {
			panic("expected Credit Card A third (5.0 fee)")
		}
		if ranked[3].Name != "Credit Card B" {
			panic("expected Credit Card B fourth (8.0 fee)")
		}
	})

	// Test 3: rank_by_trust — highest usage count first
	test("test_trust_ranking", func() {
		methods := []PaymentMethod{
			{"Credit Card A", 0.05, 5.0, 500},
			{"UPI", 0.01, 0.0, 1000},
			{"Debit Card", 0.0, 2.0, 800},
		}
		ranked := RankByTrust(methods)
		if len(ranked) != 3 {
			panic("expected 3 results")
		}
		if ranked[0].Name != "UPI" {
			panic("expected UPI first (1000 uses)")
		}
		if ranked[1].Name != "Debit Card" {
			panic("expected Debit Card second (800 uses)")
		}
		if ranked[2].Name != "Credit Card A" {
			panic("expected Credit Card A third (500 uses)")
		}
	})

	// Test 4: empty input returns empty
	test("test_empty_input", func() {
		if len(RankByRewards([]PaymentMethod{})) != 0 {
			panic("expected empty result for RankByRewards")
		}
		if len(RankByLowFee([]PaymentMethod{})) != 0 {
			panic("expected empty result for RankByLowFee")
		}
		if len(RankByTrust([]PaymentMethod{})) != 0 {
			panic("expected empty result for RankByTrust")
		}
	})

	// Test 5: single item returns itself
	test("test_single_item", func() {
		single := []PaymentMethod{{"UPI", 0.01, 0.0, 100}}
		result := RankByRewards(single)
		if len(result) != 1 {
			panic("expected 1 result")
		}
		if result[0].Name != "UPI" {
			panic("expected UPI")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
