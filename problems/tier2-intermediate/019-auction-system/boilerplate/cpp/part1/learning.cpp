#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Models (given — do not modify) ────────────────────────────────────

enum class UserType { BUYER, SELLER };
enum class AuctionStatus { OPEN, CLOSED, NO_SALE };

struct User {
    int userId;
    string name;
    UserType type;
};

struct Bid {
    int bidderId;
    double amount;
};

struct Auction {
    int auctionId;
    int sellerId;
    string item;
    double basePrice;
    AuctionStatus status;
    vector<Bid> bids;
};

// ─── Auction System ─────────────────────────────────────────────────────────

class AuctionSystem {
    int nextUserId;
    int nextAuctionId;
    unordered_map<int, User> users;
    unordered_map<int, Auction> auctions;

public:
    AuctionSystem() : nextUserId(1), nextAuctionId(1) {}

    int registerUser(string name, string type) {
        // TODO: Parse type ("BUYER" or "SELLER") into UserType enum
        // TODO: Create User{nextUserId, name, parsedType}
        // TODO: Store in users map, increment nextUserId, return the ID
        return -1;
    }

    int createAuction(int sellerId, string item, double basePrice) {
        // TODO: Check users.find(sellerId) != users.end()
        // TODO: Check users[sellerId].type == UserType::SELLER
        // TODO: Create Auction{nextAuctionId, sellerId, item, basePrice, OPEN, {}}
        // TODO: Store in auctions map, increment nextAuctionId, return the ID
        return -1;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        // TODO: Validate auctionId exists in auctions map
        // TODO: Validate buyerId exists and is a BUYER
        // TODO: Validate auction status is OPEN
        // TODO: Validate buyerId != auction.sellerId
        // TODO: Find currentHighest = max of basePrice and all existing bids
        // TODO: If amount <= currentHighest, return false
        // TODO: Push Bid{buyerId, amount} into auction.bids
        // TODO: Return true
        return false;
    }

    double getWinningBid(int auctionId) {
        // TODO: If auction has no bids, return -1
        // TODO: Find and return the maximum bid amount
        return -1;
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Auction System — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
