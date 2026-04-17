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

	// Test 1: Register users and create an auction
	test("test_register_and_create", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		if seller != 1 {
			panic("expected seller ID 1")
		}
		if buyer != 2 {
			panic("expected buyer ID 2")
		}
		auctionId := sys.CreateAuction(seller, "Laptop", 500.0)
		if auctionId != 1 {
			panic("expected auction ID 1")
		}
	})

	// Test 2: Place a valid bid
	test("test_place_valid_bid", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Phone", 100.0)
		result := sys.PlaceBid(aId, buyer, 150.0)
		if !result {
			panic("expected bid to succeed")
		}
		if sys.GetWinningBid(aId) != 150.0 {
			panic("expected winning bid 150.0")
		}
	})

	// Test 3: Bid must exceed current highest
	test("test_bid_must_exceed", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer1 := sys.RegisterUser("Bob", "BUYER")
		buyer2 := sys.RegisterUser("Charlie", "BUYER")
		aId := sys.CreateAuction(seller, "Watch", 100.0)
		if !sys.PlaceBid(aId, buyer1, 200.0) {
			panic("expected first bid to succeed")
		}
		if sys.PlaceBid(aId, buyer2, 150.0) {
			panic("expected bid below current highest to fail")
		}
		if sys.PlaceBid(aId, buyer2, 200.0) {
			panic("expected equal bid to fail")
		}
		if !sys.PlaceBid(aId, buyer2, 250.0) {
			panic("expected bid exceeding highest to succeed")
		}
		if sys.GetWinningBid(aId) != 250.0 {
			panic("expected winning bid 250.0")
		}
	})

	// Test 4: Bid must exceed base price when no prior bids
	test("test_bid_exceeds_base_price", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Book", 50.0)
		if sys.PlaceBid(aId, buyer, 30.0) {
			panic("expected bid below base price to fail")
		}
		if sys.PlaceBid(aId, buyer, 50.0) {
			panic("expected bid equal to base price to fail")
		}
		if !sys.PlaceBid(aId, buyer, 51.0) {
			panic("expected bid above base price to succeed")
		}
	})

	// Test 5: Only buyers can bid
	test("test_only_buyers_bid", func() {
		sys := NewAuctionSystem()
		seller1 := sys.RegisterUser("Alice", "SELLER")
		seller2 := sys.RegisterUser("Bob", "SELLER")
		aId := sys.CreateAuction(seller1, "Tablet", 200.0)
		if sys.PlaceBid(aId, seller2, 300.0) {
			panic("expected seller bidding to fail")
		}
	})

	// Test 6: Only sellers can create auctions
	test("test_only_sellers_create", func() {
		sys := NewAuctionSystem()
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(buyer, "Camera", 300.0)
		if aId != -1 {
			panic("expected buyer creating auction to return -1")
		}
	})

	// Test 7: GetWinningBid returns -1 when no bids
	test("test_no_bids_returns_negative", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		aId := sys.CreateAuction(seller, "Keyboard", 75.0)
		if sys.GetWinningBid(aId) != -1 {
			panic("expected -1 when no bids")
		}
	})

	// Test 8: Multiple auctions are independent
	test("test_independent_auctions", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		a1 := sys.CreateAuction(seller, "Item1", 100.0)
		a2 := sys.CreateAuction(seller, "Item2", 200.0)
		sys.PlaceBid(a1, buyer, 150.0)
		sys.PlaceBid(a2, buyer, 300.0)
		if sys.GetWinningBid(a1) != 150.0 {
			panic("expected 150.0 for auction 1")
		}
		if sys.GetWinningBid(a2) != 300.0 {
			panic("expected 300.0 for auction 2")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
