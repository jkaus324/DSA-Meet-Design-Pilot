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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Extend the auction system with closing and state management:
//
//   closeAuction(auctionId):
//     - Close an open auction. If it has bids, status -> CLOSED.
//     - If no bids, status -> NO_SALE.
//     - Return true on success, false if already closed.
//
//   getAuctionStatus(auctionId):
//     - Return "OPEN", "CLOSED", or "NO_SALE"
//
// Think about:
//   - What state transitions are valid?
//   - What happens when someone bids on a closed auction?
//   - What does getWinningBid return after a no-sale close?
//
// Entry points (in addition to Part 1):
//   bool closeAuction(int auctionId)
//   string getAuctionStatus(int auctionId)
//
// ─────────────────────────────────────────────────────────────────────────────


