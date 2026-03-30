// Part 3 Tests — Multiple Auction Strategies
// Tests Ascending, SealedBid, and BuyNow auction types

#include "solution.cpp"
#include <cassert>
#include <iostream>
#include <cmath>
using namespace std;

int part3_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Ascending strategy (default) — same behavior as Part 1/2
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer1 = sys.registerUser("Bob", "BUYER");
        int buyer2 = sys.registerUser("Charlie", "BUYER");
        int aId = sys.createAuction(seller, "Laptop", 500.0, "ASCENDING");
        assert(sys.placeBid(aId, buyer1, 600.0) == true);
        assert(sys.placeBid(aId, buyer2, 550.0) == false);  // must exceed 600
        assert(sys.placeBid(aId, buyer2, 700.0) == true);
        assert(sys.getWinningBid(aId) == 700.0);
        cout << "PASS test_ascending_strategy" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_ascending_strategy" << endl;
        failed++;
    }

    // Test 2: SealedBid — any bid above base accepted, winner hidden while open
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer1 = sys.registerUser("Bob", "BUYER");
        int buyer2 = sys.registerUser("Charlie", "BUYER");
        int aId = sys.createAuction(seller, "Art", 100.0, "SEALED");

        // Any bid above base price is accepted
        assert(sys.placeBid(aId, buyer1, 500.0) == true);
        assert(sys.placeBid(aId, buyer2, 200.0) == true);  // lower than 500 but still valid

        // While open, winning bid is hidden
        assert(sys.getWinningBid(aId) == -1);

        // After close, winner is revealed
        sys.closeAuction(aId);
        assert(sys.getWinningBid(aId) == 500.0);  // highest bid wins
        cout << "PASS test_sealed_bid_strategy" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_sealed_bid_strategy" << endl;
        failed++;
    }

    // Test 3: SealedBid — bid at or below base price is rejected
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Vase", 200.0, "SEALED");
        assert(sys.placeBid(aId, buyer, 100.0) == false);  // below base
        assert(sys.placeBid(aId, buyer, 200.0) == false);  // equal to base
        assert(sys.placeBid(aId, buyer, 201.0) == true);   // above base
        cout << "PASS test_sealed_rejects_low_bids" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_sealed_rejects_low_bids" << endl;
        failed++;
    }

    // Test 4: BuyNow — instant purchase at premium price
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Guitar", 100.0, "BUYNOW");
        // Buy-now price = 100 * 1.5 = 150

        assert(sys.placeBid(aId, buyer, 120.0) == false);  // below buy-now price
        assert(sys.placeBid(aId, buyer, 150.0) == true);   // meets buy-now price

        // Should auto-close
        assert(sys.getAuctionStatus(aId) == "CLOSED");
        assert(sys.getWinningBid(aId) == 150.0);
        cout << "PASS test_buynow_strategy" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_buynow_strategy" << endl;
        failed++;
    }

    // Test 5: BuyNow — no bids after auto-close
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer1 = sys.registerUser("Bob", "BUYER");
        int buyer2 = sys.registerUser("Charlie", "BUYER");
        int aId = sys.createAuction(seller, "Drum", 200.0, "BUYNOW");
        // Buy-now price = 200 * 1.5 = 300

        assert(sys.placeBid(aId, buyer1, 300.0) == true);   // auto-closes
        assert(sys.placeBid(aId, buyer2, 400.0) == false);  // auction already closed
        cout << "PASS test_buynow_no_bids_after_close" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_buynow_no_bids_after_close" << endl;
        failed++;
    }

    // Test 6: Default strategy is ASCENDING
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Mouse", 50.0);  // no strategy specified
        assert(sys.placeBid(aId, buyer, 60.0) == true);
        assert(sys.getWinningBid(aId) == 60.0);  // visible (ascending behavior)
        cout << "PASS test_default_ascending" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_default_ascending" << endl;
        failed++;
    }

    // Test 7: Mixed strategies in same system
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int a1 = sys.createAuction(seller, "Item1", 100.0, "ASCENDING");
        int a2 = sys.createAuction(seller, "Item2", 100.0, "SEALED");
        int a3 = sys.createAuction(seller, "Item3", 100.0, "BUYNOW");

        sys.placeBid(a1, buyer, 200.0);
        sys.placeBid(a2, buyer, 200.0);
        sys.placeBid(a3, buyer, 150.0);

        assert(sys.getWinningBid(a1) == 200.0);   // ascending: visible
        assert(sys.getWinningBid(a2) == -1);       // sealed: hidden while open
        assert(sys.getAuctionStatus(a3) == "CLOSED");  // buynow: auto-closed
        assert(sys.getWinningBid(a3) == 150.0);   // buynow: visible after close
        cout << "PASS test_mixed_strategies" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_mixed_strategies" << endl;
        failed++;
    }

    cout << "PART3_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
