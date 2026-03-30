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
// HINT: Use HashMaps for O(1) user and auction lookup.
// HINT: Auto-increment IDs for users and auctions.

class AuctionSystem {
    // HINT: You need counters for auto-assigning IDs
    // HINT: You need maps from ID to User and ID to Auction

public:
    AuctionSystem() {
        // TODO: Initialize ID counters
    }

    int registerUser(string name, string type) {
        // TODO: Create a new User with auto-assigned ID
        // HINT: Parse the type string into UserType enum
        return -1;
    }

    int createAuction(int sellerId, string item, double basePrice) {
        // TODO: Validate that sellerId exists and is a SELLER
        // TODO: Create auction with auto-assigned ID, status OPEN, empty bids
        // HINT: Return -1 if validation fails
        return -1;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        // TODO: Validate auction exists and is OPEN
        // TODO: Validate buyer exists and is a BUYER
        // TODO: Validate buyer is not the seller of this auction
        // HINT: Find the current highest bid (or use basePrice if no bids)
        // HINT: New bid must STRICTLY exceed the current highest
        return false;
    }

    double getWinningBid(int auctionId) {
        // TODO: Return the highest bid amount, or -1 if no bids
        // HINT: Iterate through bids to find the maximum
        return -1;
    }
};
