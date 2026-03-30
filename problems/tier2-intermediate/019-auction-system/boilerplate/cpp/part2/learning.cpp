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

// ─── Auction System with Closing ────────────────────────────────────────────

class AuctionSystem {
    int nextUserId;
    int nextAuctionId;
    unordered_map<int, User> users;
    unordered_map<int, Auction> auctions;

public:
    AuctionSystem() : nextUserId(1), nextAuctionId(1) {}

    int registerUser(string name, string type) {
        // TODO: Same as Part 1 — parse type, create User, store, return ID
        return -1;
    }

    int createAuction(int sellerId, string item, double basePrice) {
        // TODO: Same as Part 1 — validate seller, create Auction, store, return ID
        return -1;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        // TODO: Same as Part 1 — all validations + bid must exceed current highest
        return false;
    }

    double getWinningBid(int auctionId) {
        // TODO: Same as Part 1 — return highest bid amount or -1
        return -1;
    }

    bool closeAuction(int auctionId) {
        // TODO: Check auctions.find(auctionId) != auctions.end()
        // TODO: Check auction.status == AuctionStatus::OPEN (only open can close)
        // TODO: If auction.bids.empty() -> set status to NO_SALE
        // TODO: Else -> set status to CLOSED
        // TODO: Return true on success
        return false;
    }

    string getAuctionStatus(int auctionId) {
        // TODO: Switch on auctions[auctionId].status
        // TODO: Return "OPEN", "CLOSED", or "NO_SALE"
        return "UNKNOWN";
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Auction System (Part 2) — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
