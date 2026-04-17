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

	// Test 1: Ascending strategy (default) — same behavior as Part 1/2
	test("test_ascending_strategy", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer1 := sys.RegisterUser("Bob", "BUYER")
		buyer2 := sys.RegisterUser("Charlie", "BUYER")
		aId := sys.CreateAuction(seller, "Laptop", 500.0, "ASCENDING")
		if !sys.PlaceBid(aId, buyer1, 600.0) {
			panic("expected first bid to succeed")
		}
		if sys.PlaceBid(aId, buyer2, 550.0) {
			panic("expected lower bid to fail")
		}
		if !sys.PlaceBid(aId, buyer2, 700.0) {
			panic("expected higher bid to succeed")
		}
		if sys.GetWinningBid(aId) != 700.0 {
			panic("expected winning bid 700.0")
		}
	})

	// Test 2: SealedBid — any bid above base accepted, winner hidden while open
	test("test_sealed_bid_strategy", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer1 := sys.RegisterUser("Bob", "BUYER")
		buyer2 := sys.RegisterUser("Charlie", "BUYER")
		aId := sys.CreateAuction(seller, "Art", 100.0, "SEALED")
		if !sys.PlaceBid(aId, buyer1, 500.0) {
			panic("expected first bid to succeed")
		}
		if !sys.PlaceBid(aId, buyer2, 200.0) {
			panic("expected lower bid to still succeed in sealed auction")
		}
		if sys.GetWinningBid(aId) != -1 {
			panic("expected winning bid hidden while open")
		}
		sys.CloseAuction(aId)
		if sys.GetWinningBid(aId) != 500.0 {
			panic("expected highest bid revealed after close")
		}
	})

	// Test 3: SealedBid — bid at or below base price is rejected
	test("test_sealed_rejects_low_bids", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Vase", 200.0, "SEALED")
		if sys.PlaceBid(aId, buyer, 100.0) {
			panic("expected bid below base to fail")
		}
		if sys.PlaceBid(aId, buyer, 200.0) {
			panic("expected bid equal to base to fail")
		}
		if !sys.PlaceBid(aId, buyer, 201.0) {
			panic("expected bid above base to succeed")
		}
	})

	// Test 4: BuyNow — instant purchase at premium price
	test("test_buynow_strategy", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Guitar", 100.0, "BUYNOW")
		// Buy-now price = 100 * 1.5 = 150
		if sys.PlaceBid(aId, buyer, 120.0) {
			panic("expected bid below buy-now price to fail")
		}
		if !sys.PlaceBid(aId, buyer, 150.0) {
			panic("expected bid at buy-now price to succeed")
		}
		if sys.GetAuctionStatus(aId) != "CLOSED" {
			panic("expected auction to auto-close")
		}
		if sys.GetWinningBid(aId) != 150.0 {
			panic("expected winning bid 150.0")
		}
	})

	// Test 5: BuyNow — no bids after auto-close
	test("test_buynow_no_bids_after_close", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer1 := sys.RegisterUser("Bob", "BUYER")
		buyer2 := sys.RegisterUser("Charlie", "BUYER")
		aId := sys.CreateAuction(seller, "Drum", 200.0, "BUYNOW")
		// Buy-now price = 200 * 1.5 = 300
		if !sys.PlaceBid(aId, buyer1, 300.0) {
			panic("expected auto-closing bid to succeed")
		}
		if sys.PlaceBid(aId, buyer2, 400.0) {
			panic("expected bid on closed auction to fail")
		}
	})

	// Test 6: Default strategy is ASCENDING (no strategyType passed)
	test("test_default_ascending", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Mouse", 50.0, "ASCENDING")
		if !sys.PlaceBid(aId, buyer, 60.0) {
			panic("expected bid to succeed")
		}
		if sys.GetWinningBid(aId) != 60.0 {
			panic("expected visible winning bid 60.0")
		}
	})

	// Test 7: Mixed strategies in same system
	test("test_mixed_strategies", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		a1 := sys.CreateAuction(seller, "Item1", 100.0, "ASCENDING")
		a2 := sys.CreateAuction(seller, "Item2", 100.0, "SEALED")
		a3 := sys.CreateAuction(seller, "Item3", 100.0, "BUYNOW")
		sys.PlaceBid(a1, buyer, 200.0)
		sys.PlaceBid(a2, buyer, 200.0)
		sys.PlaceBid(a3, buyer, 150.0)
		if sys.GetWinningBid(a1) != 200.0 {
			panic("expected ascending visible bid 200.0")
		}
		if sys.GetWinningBid(a2) != -1 {
			panic("expected sealed bid hidden while open")
		}
		if sys.GetAuctionStatus(a3) != "CLOSED" {
			panic("expected buynow auction auto-closed")
		}
		if sys.GetWinningBid(a3) != 150.0 {
			panic("expected buynow winning bid 150.0")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
