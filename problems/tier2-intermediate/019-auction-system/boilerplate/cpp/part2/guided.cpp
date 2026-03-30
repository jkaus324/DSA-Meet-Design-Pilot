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
// HINT: This extends Part 1. Include all Part 1 functionality.
// HINT: State transitions: OPEN -> CLOSED (has bids), OPEN -> NO_SALE (no bids)
// HINT: CLOSED and NO_SALE are terminal states — no further transitions.

class AuctionSystem {
    // HINT: Same data structures as Part 1

public:
    AuctionSystem() {
        // TODO: Initialize (same as Part 1)
    }

    int registerUser(string name, string type) {
        // TODO: Same as Part 1
        return -1;
    }

    int createAuction(int sellerId, string item, double basePrice) {
        // TODO: Same as Part 1
        return -1;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        // TODO: Same as Part 1
        // HINT: Already checks auction is OPEN, so closed auctions are handled
        return false;
    }

    double getWinningBid(int auctionId) {
        // TODO: Same as Part 1
        return -1;
    }

    bool closeAuction(int auctionId) {
        // TODO: Validate auction exists
        // HINT: Only OPEN auctions can be closed
        // HINT: If bids exist -> CLOSED; if no bids -> NO_SALE
        return false;
    }

    string getAuctionStatus(int auctionId) {
        // TODO: Return the status as a string
        // HINT: Map enum to string: OPEN, CLOSED, NO_SALE
        return "UNKNOWN";
    }
};
