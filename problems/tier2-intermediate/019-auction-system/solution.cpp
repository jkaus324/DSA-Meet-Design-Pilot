#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
using namespace std;

// ─── Data Models ────────────────────────────────────────────────────────────

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
        double currentHighest = auction.basePrice;
        for (const auto& bid : auction.bids) {
            if (bid.amount > currentHighest) currentHighest = bid.amount;
        }
        return amount > currentHighest;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        if (auction.bids.empty()) return -1;
        double maxBid = auction.bids[0].amount;
        for (const auto& bid : auction.bids) {
            if (bid.amount > maxBid) maxBid = bid.amount;
        }
        return maxBid;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        return false;
    }
};

class SealedBidStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        return amount > auction.basePrice;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        if (auction.status == AuctionStatus::OPEN) return -1;
        if (auction.bids.empty()) return -1;
        double maxBid = auction.bids[0].amount;
        for (const auto& bid : auction.bids) {
            if (bid.amount > maxBid) maxBid = bid.amount;
        }
        return maxBid;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        return false;
    }
};

class BuyNowStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        return amount >= auction.basePrice * 1.5;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        if (auction.bids.empty()) return -1;
        return auction.bids.back().amount;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        return true;
    }
};

// ─── Strategy Factory ───────────────────────────────────────────────────────

AuctionStrategy* createStrategy(const string& type) {
    if (type == "SEALED") return new SealedBidStrategy();
    if (type == "BUYNOW") return new BuyNowStrategy();
    return new AscendingStrategy();
}

// ─── Auction System ─────────────────────────────────────────────────────────

class AuctionSystem {
    int nextUserId;
    int nextAuctionId;
    unordered_map<int, User> users;
    unordered_map<int, Auction> auctions;
    unordered_map<int, AuctionStrategy*> strategies;

public:
    AuctionSystem() : nextUserId(1), nextAuctionId(1) {}

    ~AuctionSystem() {
        for (auto& [id, s] : strategies) delete s;
    }

    int registerUser(string name, string type) {
        UserType ut = (type == "SELLER") ? UserType::SELLER : UserType::BUYER;
        int id = nextUserId++;
        users[id] = User{id, name, ut};
        return id;
    }

    int createAuction(int sellerId, string item, double basePrice, string strategyType = "ASCENDING") {
        if (users.find(sellerId) == users.end()) return -1;
        if (users[sellerId].type != UserType::SELLER) return -1;
        int id = nextAuctionId++;
        auctions[id] = Auction{id, sellerId, item, basePrice, AuctionStatus::OPEN, {}};
        strategies[id] = createStrategy(strategyType);
        return id;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        if (auctions.find(auctionId) == auctions.end()) return false;
        if (users.find(buyerId) == users.end()) return false;
        if (users[buyerId].type != UserType::BUYER) return false;
        Auction& auction = auctions[auctionId];
        if (auction.status != AuctionStatus::OPEN) return false;
        if (buyerId == auction.sellerId) return false;

        if (!strategies[auctionId]->acceptBid(auction, amount)) return false;

        auction.bids.push_back(Bid{buyerId, amount});

        if (strategies[auctionId]->shouldAutoClose(auction, amount)) {
            auction.status = AuctionStatus::CLOSED;
        }
        return true;
    }

    double getWinningBid(int auctionId) {
        if (auctions.find(auctionId) == auctions.end()) return -1;
        return strategies[auctionId]->getVisibleWinningBid(auctions[auctionId]);
    }

    bool closeAuction(int auctionId) {
        if (auctions.find(auctionId) == auctions.end()) return false;
        Auction& auction = auctions[auctionId];
        if (auction.status != AuctionStatus::OPEN) return false;
        if (auction.bids.empty()) {
            auction.status = AuctionStatus::NO_SALE;
        } else {
            auction.status = AuctionStatus::CLOSED;
        }
        return true;
    }

    string getAuctionStatus(int auctionId) {
        if (auctions.find(auctionId) == auctions.end()) return "UNKNOWN";
        switch (auctions[auctionId].status) {
            case AuctionStatus::OPEN: return "OPEN";
            case AuctionStatus::CLOSED: return "CLOSED";
            case AuctionStatus::NO_SALE: return "NO_SALE";
        }
        return "UNKNOWN";
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Auction System — all parts implemented." << endl;
    return 0;
}
#endif
