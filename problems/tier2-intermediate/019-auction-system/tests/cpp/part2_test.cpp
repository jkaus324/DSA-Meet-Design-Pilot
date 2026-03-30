// Part 2 Tests — Auction Closing and State Management
// Tests closing auctions, state transitions, and edge cases

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Close auction with bids -> CLOSED
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Laptop", 500.0);
        sys.placeBid(aId, buyer, 600.0);
        assert(sys.closeAuction(aId) == true);
        assert(sys.getAuctionStatus(aId) == "CLOSED");
        assert(sys.getWinningBid(aId) == 600.0);
        cout << "PASS test_close_with_bids" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_close_with_bids" << endl;
        failed++;
    }

    // Test 2: Close auction with no bids -> NO_SALE
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int aId = sys.createAuction(seller, "Phone", 300.0);
        assert(sys.closeAuction(aId) == true);
        assert(sys.getAuctionStatus(aId) == "NO_SALE");
        assert(sys.getWinningBid(aId) == -1);
        cout << "PASS test_close_no_bids" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_close_no_bids" << endl;
        failed++;
    }

    // Test 3: Cannot close an already closed auction
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Watch", 100.0);
        sys.placeBid(aId, buyer, 150.0);
        assert(sys.closeAuction(aId) == true);
        assert(sys.closeAuction(aId) == false);  // already closed
        cout << "PASS test_double_close" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_double_close" << endl;
        failed++;
    }

    // Test 4: Cannot bid on a closed auction
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Tablet", 200.0);
        sys.closeAuction(aId);
        assert(sys.placeBid(aId, buyer, 300.0) == false);
        cout << "PASS test_bid_on_closed" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_bid_on_closed" << endl;
        failed++;
    }

    // Test 5: New auction starts as OPEN
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int aId = sys.createAuction(seller, "Camera", 400.0);
        assert(sys.getAuctionStatus(aId) == "OPEN");
        cout << "PASS test_initial_status_open" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_initial_status_open" << endl;
        failed++;
    }

    // Test 6: Cannot close a NO_SALE auction again
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int aId = sys.createAuction(seller, "Book", 25.0);
        assert(sys.closeAuction(aId) == true);   // NO_SALE
        assert(sys.closeAuction(aId) == false);   // already in terminal state
        assert(sys.getAuctionStatus(aId) == "NO_SALE");
        cout << "PASS test_close_nosale_again" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_close_nosale_again" << endl;
        failed++;
    }

    // Test 7: Winning bid persists after close
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer1 = sys.registerUser("Bob", "BUYER");
        int buyer2 = sys.registerUser("Charlie", "BUYER");
        int aId = sys.createAuction(seller, "Painting", 1000.0);
        sys.placeBid(aId, buyer1, 1500.0);
        sys.placeBid(aId, buyer2, 2000.0);
        sys.closeAuction(aId);
        assert(sys.getWinningBid(aId) == 2000.0);
        assert(sys.getAuctionStatus(aId) == "CLOSED");
        cout << "PASS test_winning_bid_after_close" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_winning_bid_after_close" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
