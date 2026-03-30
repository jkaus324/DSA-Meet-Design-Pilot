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
// Extend the auction system with multiple auction strategies:
//
//   ASCENDING (default):
//     - Bids must exceed current highest. Winner = highest bidder at close.
//
//   SEALED:
//     - Any bid above base price is accepted (even below previous bids).
//     - getWinningBid returns -1 while OPEN (bids are hidden).
//     - Winner revealed only after close.
//
//   BUYNOW:
//     - First bid >= basePrice * 1.5 instantly wins and auto-closes.
//     - Bids below the buy-now threshold are rejected.
//
// Think about:
//   - What abstraction lets each auction type have its own bid rules?
//   - How do you create the right strategy from a string type name?
//   - How do you add a 4th strategy without modifying the system class?
//
// Modified entry point:
//   int createAuction(int sellerId, string item, double basePrice, string strategyType)
//
// ─────────────────────────────────────────────────────────────────────────────


