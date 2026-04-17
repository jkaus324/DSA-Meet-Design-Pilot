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

	// Test 1: Close auction with bids -> CLOSED
	test("test_close_with_bids", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Laptop", 500.0)
		sys.PlaceBid(aId, buyer, 600.0)
		if !sys.CloseAuction(aId) {
			panic("expected CloseAuction to return true")
		}
		if sys.GetAuctionStatus(aId) != "CLOSED" {
			panic("expected status CLOSED")
		}
		if sys.GetWinningBid(aId) != 600.0 {
			panic("expected winning bid 600.0")
		}
	})

	// Test 2: Close auction with no bids -> NO_SALE
	test("test_close_no_bids", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		aId := sys.CreateAuction(seller, "Phone", 300.0)
		if !sys.CloseAuction(aId) {
			panic("expected CloseAuction to return true")
		}
		if sys.GetAuctionStatus(aId) != "NO_SALE" {
			panic("expected status NO_SALE")
		}
		if sys.GetWinningBid(aId) != -1 {
			panic("expected -1 winning bid for no-sale")
		}
	})

	// Test 3: Cannot close an already closed auction
	test("test_double_close", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Watch", 100.0)
		sys.PlaceBid(aId, buyer, 150.0)
		if !sys.CloseAuction(aId) {
			panic("expected first close to succeed")
		}
		if sys.CloseAuction(aId) {
			panic("expected second close to fail")
		}
	})

	// Test 4: Cannot bid on a closed auction
	test("test_bid_on_closed", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer := sys.RegisterUser("Bob", "BUYER")
		aId := sys.CreateAuction(seller, "Tablet", 200.0)
		sys.CloseAuction(aId)
		if sys.PlaceBid(aId, buyer, 300.0) {
			panic("expected bid on closed auction to fail")
		}
	})

	// Test 5: New auction starts as OPEN
	test("test_initial_status_open", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		aId := sys.CreateAuction(seller, "Camera", 400.0)
		if sys.GetAuctionStatus(aId) != "OPEN" {
			panic("expected initial status OPEN")
		}
	})

	// Test 6: Cannot close a NO_SALE auction again
	test("test_close_nosale_again", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		aId := sys.CreateAuction(seller, "Book", 25.0)
		if !sys.CloseAuction(aId) {
			panic("expected first close to succeed (NO_SALE)")
		}
		if sys.CloseAuction(aId) {
			panic("expected second close to fail")
		}
		if sys.GetAuctionStatus(aId) != "NO_SALE" {
			panic("expected status to remain NO_SALE")
		}
	})

	// Test 7: Winning bid persists after close
	test("test_winning_bid_after_close", func() {
		sys := NewAuctionSystem()
		seller := sys.RegisterUser("Alice", "SELLER")
		buyer1 := sys.RegisterUser("Bob", "BUYER")
		buyer2 := sys.RegisterUser("Charlie", "BUYER")
		aId := sys.CreateAuction(seller, "Painting", 1000.0)
		sys.PlaceBid(aId, buyer1, 1500.0)
		sys.PlaceBid(aId, buyer2, 2000.0)
		sys.CloseAuction(aId)
		if sys.GetWinningBid(aId) != 2000.0 {
			panic("expected winning bid 2000.0 after close")
		}
		if sys.GetAuctionStatus(aId) != "CLOSED" {
			panic("expected status CLOSED")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
