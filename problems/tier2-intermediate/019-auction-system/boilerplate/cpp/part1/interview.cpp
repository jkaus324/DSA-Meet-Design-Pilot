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
// Design and implement an AuctionSystem that:
//   1. Registers users as BUYER or SELLER (registerUser)
//   2. Lets sellers create auctions with a base price (createAuction)
//   3. Lets buyers place bids that must exceed the current highest (placeBid)
//   4. Returns the current highest bid for an auction (getWinningBid)
//
// Think about:
//   - How do you store users and auctions for fast lookup?
//   - How do you validate that only buyers bid and only sellers create?
//   - How do you track the current highest bid efficiently?
//   - What happens if a seller tries to bid on their own auction?
//
// Entry points (must exist for tests):
//   AuctionSystem()
//   int registerUser(string name, string type)
//   int createAuction(int sellerId, string item, double basePrice)
//   bool placeBid(int auctionId, int buyerId, double amount)
//   double getWinningBid(int auctionId)
//
// ─────────────────────────────────────────────────────────────────────────────


