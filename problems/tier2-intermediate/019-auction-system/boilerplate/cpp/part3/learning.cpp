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

// ─── Strategy Interface ─────────────────────────────────────────────────────

class AuctionStrategy {
public:
    virtual bool acceptBid(Auction& auction, double amount) = 0;
    virtual double getVisibleWinningBid(const Auction& auction) = 0;
    virtual bool shouldAutoClose(const Auction& auction, double amount) = 0;
    virtual ~AuctionStrategy() = default;
};

// ─── Concrete Strategies ────────────────────────────────────────────────────

class AscendingStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        // TODO: Find currentHighest = max(basePrice, all bid amounts)
        // TODO: Return amount > currentHighest
        return false;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        // TODO: Return the maximum bid amount, or -1 if no bids
        return -1;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        // TODO: Ascending never auto-closes
        return false;
    }
};

class SealedBidStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        // TODO: Return amount > auction.basePrice (any bid above base is valid)
        return false;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        // TODO: If auction.status == OPEN, return -1 (sealed — bids are hidden)
        // TODO: If closed, return the maximum bid amount, or -1 if no bids
        return -1;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        // TODO: Sealed bids never auto-close
        return false;
    }
};

class BuyNowStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        // TODO: Return amount >= auction.basePrice * 1.5
        return false;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        // TODO: Return the bid amount if bids exist, -1 otherwise
        return -1;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        // TODO: BuyNow always auto-closes on successful bid
        return true;
    }
};

// ─── Strategy Factory ───────────────────────────────────────────────────────

AuctionStrategy* createStrategy(const string& type) {
    // TODO: Return new AscendingStrategy() for "ASCENDING"
    // TODO: Return new SealedBidStrategy() for "SEALED"
    // TODO: Return new BuyNowStrategy() for "BUYNOW"
    // TODO: Default to AscendingStrategy
    return nullptr;
}

// ─── Auction System ─────────────────────────────────────────────────────────

class AuctionSystem {
    int nextUserId;
    int nextAuctionId;
    unordered_map<int, User> users;
    unordered_map<int, Auction> auctions;
    unordered_map<int, AuctionStrategy*> strategies;  // auctionId -> strategy

public:
    AuctionSystem() : nextUserId(1), nextAuctionId(1) {}

    ~AuctionSystem() {
        for (auto& [id, s] : strategies) delete s;
    }

    int registerUser(string name, string type) {
        // TODO: Same as before — parse type, create User, store, return ID
        return -1;
    }

    int createAuction(int sellerId, string item, double basePrice, string strategyType = "ASCENDING") {
        // TODO: Validate seller exists and is a SELLER
        // TODO: Create Auction with auto-assigned ID
        // TODO: Use createStrategy(strategyType) to create and store the strategy
        // TODO: Return the auction ID
        return -1;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        // TODO: Validate auctionId, buyerId, auction is OPEN, buyer is not seller
        // TODO: Delegate to strategies[auctionId]->acceptBid(auction, amount)
        // TODO: If accepted, push Bid{buyerId, amount}
        // TODO: If strategies[auctionId]->shouldAutoClose(...), set status to CLOSED
        // TODO: Return true if accepted
        return false;
    }

    double getWinningBid(int auctionId) {
        // TODO: Delegate to strategies[auctionId]->getVisibleWinningBid(auction)
        return -1;
    }

    bool closeAuction(int auctionId) {
        // TODO: Same as Part 2 — validate, transition state
        return false;
    }

    string getAuctionStatus(int auctionId) {
        // TODO: Same as Part 2 — switch on status enum
        return "UNKNOWN";
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Auction System (Part 3) — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
