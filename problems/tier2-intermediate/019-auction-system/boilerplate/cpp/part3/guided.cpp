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
// HINT: Each auction type needs different rules for:
//   1. Whether a bid should be accepted
//   2. What the visible winning bid is (sealed hides it while open)
//   3. Whether the auction should auto-close after a bid

class /* YourStrategyInterfaceName */ {
public:
    virtual bool /* acceptBid */(Auction& auction, double amount) = 0;
    virtual double /* getVisibleWinningBid */(const Auction& auction) = 0;
    virtual bool /* shouldAutoClose */(const Auction& auction, double amount) = 0;
    virtual ~/* YourStrategyInterfaceName */() = default;
};

// ─── Concrete Strategies ────────────────────────────────────────────────────
// TODO: Implement AscendingStrategy
//   HINT: acceptBid: amount must exceed current highest (or base price)
//   HINT: getVisibleWinningBid: return max bid amount
//   HINT: shouldAutoClose: always false

// TODO: Implement SealedBidStrategy
//   HINT: acceptBid: any amount above base price
//   HINT: getVisibleWinningBid: return -1 if OPEN (sealed), max bid if CLOSED
//   HINT: shouldAutoClose: always false

// TODO: Implement BuyNowStrategy
//   HINT: acceptBid: amount must be >= basePrice * 1.5
//   HINT: getVisibleWinningBid: return the bid amount (only one bid)
//   HINT: shouldAutoClose: always true (auto-close on successful bid)


// ─── Auction System ─────────────────────────────────────────────────────────
// HINT: Store a strategy per auction (map from auctionId to strategy pointer)
// HINT: Delegate bid validation to the strategy
// HINT: Use a factory function to create strategies from string names

class AuctionSystem {
    // HINT: Same data as Part 2, plus a map from auctionId to strategy

public:
    AuctionSystem() {
        // TODO: Initialize
    }

    int registerUser(string name, string type) {
        // TODO: Same as before
        return -1;
    }

    int createAuction(int sellerId, string item, double basePrice, string strategyType = "ASCENDING") {
        // TODO: Same validation as before
        // HINT: Create the appropriate strategy based on strategyType
        // HINT: Store the strategy mapped to this auctionId
        return -1;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        // TODO: Basic validations (exists, is buyer, auction is open, not seller)
        // HINT: Delegate bid acceptance to the auction's strategy
        // HINT: If strategy says auto-close, set status to CLOSED
        return false;
    }

    double getWinningBid(int auctionId) {
        // HINT: Delegate to the auction's strategy's getVisibleWinningBid
        return -1;
    }

    bool closeAuction(int auctionId) {
        // TODO: Same as Part 2
        return false;
    }

    string getAuctionStatus(int auctionId) {
        // TODO: Same as Part 2
        return "UNKNOWN";
    }
};
