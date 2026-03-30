// Part 1 Tests — Core Auction System
// Tests user registration, auction creation, bidding, and winning bid queries

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Register users and create an auction
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        assert(seller == 1);
        assert(buyer == 2);
        int auctionId = sys.createAuction(seller, "Laptop", 500.0);
        assert(auctionId == 1);
        cout << "PASS test_register_and_create" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_register_and_create" << endl;
        failed++;
    }

    // Test 2: Place a valid bid
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Phone", 100.0);
        bool result = sys.placeBid(aId, buyer, 150.0);
        assert(result == true);
        assert(sys.getWinningBid(aId) == 150.0);
        cout << "PASS test_place_valid_bid" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_place_valid_bid" << endl;
        failed++;
    }

    // Test 3: Bid must exceed current highest
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer1 = sys.registerUser("Bob", "BUYER");
        int buyer2 = sys.registerUser("Charlie", "BUYER");
        int aId = sys.createAuction(seller, "Watch", 100.0);
        assert(sys.placeBid(aId, buyer1, 200.0) == true);
        assert(sys.placeBid(aId, buyer2, 150.0) == false);  // below current highest
        assert(sys.placeBid(aId, buyer2, 200.0) == false);  // equal, not exceeding
        assert(sys.placeBid(aId, buyer2, 250.0) == true);   // exceeds
        assert(sys.getWinningBid(aId) == 250.0);
        cout << "PASS test_bid_must_exceed" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_bid_must_exceed" << endl;
        failed++;
    }

    // Test 4: Bid must exceed base price when no prior bids
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(seller, "Book", 50.0);
        assert(sys.placeBid(aId, buyer, 30.0) == false);  // below base price
        assert(sys.placeBid(aId, buyer, 50.0) == false);  // equal to base price
        assert(sys.placeBid(aId, buyer, 51.0) == true);   // above base price
        cout << "PASS test_bid_exceeds_base_price" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_bid_exceeds_base_price" << endl;
        failed++;
    }

    // Test 5: Only buyers can bid
    try {
        AuctionSystem sys;
        int seller1 = sys.registerUser("Alice", "SELLER");
        int seller2 = sys.registerUser("Bob", "SELLER");
        int aId = sys.createAuction(seller1, "Tablet", 200.0);
        assert(sys.placeBid(aId, seller2, 300.0) == false);  // seller cannot bid
        cout << "PASS test_only_buyers_bid" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_only_buyers_bid" << endl;
        failed++;
    }

    // Test 6: Only sellers can create auctions
    try {
        AuctionSystem sys;
        int buyer = sys.registerUser("Bob", "BUYER");
        int aId = sys.createAuction(buyer, "Camera", 300.0);
        assert(aId == -1);  // buyer cannot create auction
        cout << "PASS test_only_sellers_create" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_only_sellers_create" << endl;
        failed++;
    }

    // Test 7: getWinningBid returns -1 when no bids
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int aId = sys.createAuction(seller, "Keyboard", 75.0);
        assert(sys.getWinningBid(aId) == -1);
        cout << "PASS test_no_bids_returns_negative" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_no_bids_returns_negative" << endl;
        failed++;
    }

    // Test 8: Multiple auctions are independent
    try {
        AuctionSystem sys;
        int seller = sys.registerUser("Alice", "SELLER");
        int buyer = sys.registerUser("Bob", "BUYER");
        int a1 = sys.createAuction(seller, "Item1", 100.0);
        int a2 = sys.createAuction(seller, "Item2", 200.0);
        sys.placeBid(a1, buyer, 150.0);
        sys.placeBid(a2, buyer, 300.0);
        assert(sys.getWinningBid(a1) == 150.0);
        assert(sys.getWinningBid(a2) == 300.0);
        cout << "PASS test_independent_auctions" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_independent_auctions" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
