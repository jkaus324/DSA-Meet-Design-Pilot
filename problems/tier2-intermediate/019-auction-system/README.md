# Problem 019 — Online Auction System

**Tier:** 2 (Intermediate) | **Patterns:** Strategy, Observer, Factory, State | **DSA:** HashMap, PriorityQueue, Sorting
**Companies:** Flipkart | **Time:** 90 minutes

---

## Problem Statement

You are building an online auction platform. Sellers list items for auction with a base price. Buyers place bids. The system tracks the highest bid and determines winners when auctions close.

**Your task:** Design and implement an auction system that supports multiple auction strategies (ascending, sealed-bid, buy-now), proper state management (open/closed), and correct winner determination.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What varies here? The **auction strategy** (how bids are accepted and winners determined) varies. The auction lifecycle is the same.
2. What states can an auction be in? Open, Closed, NoSale. What transitions are valid?
3. If you hardcode strategy logic with `if-else` inside the bid method, what happens when the product team adds a "Dutch auction" strategy?

**Naive approach:** One monolithic `Auction` class with `if (strategy == "ascending") { ... } else if (strategy == "sealed") { ... }` inside every method. Every new strategy requires modifying every method.

**Pattern approach:** Use the **Strategy pattern** to encapsulate each auction type's bidding rules. Use the **State pattern** (or state tracking) to manage auction lifecycle. Use **Factory** to create auctions. Use **Observer** to notify interested parties when events occur.

**The key insight:** The auction *type* determines bid acceptance rules and winner declaration logic. Encapsulate these rules in strategy objects — the Auction class delegates to its strategy without knowing the specifics.

---

## Data Structures

```cpp
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
    int timestamp;
};

struct Auction {
    int auctionId;
    int sellerId;
    string item;
    double basePrice;
    AuctionStatus status;
    vector<Bid> bids;
};
```

---

## Part 1

**Base requirement — Core auction system**

Implement an `AuctionSystem` that supports user registration, auction creation, bidding, and querying the current winner.

| Operation | Description |
|-----------|-------------|
| `registerUser(name, type)` | Register a new user as BUYER or SELLER. Return the assigned userId. |
| `createAuction(sellerId, item, basePrice)` | Seller creates an auction. Return the assigned auctionId. Only SELLER users can create auctions. |
| `placeBid(auctionId, buyerId, amount)` | Buyer places a bid. Must exceed the current highest bid (or base price if no bids). Only BUYER users can bid. Auction must be OPEN. Return `true` if bid accepted, `false` otherwise. |
| `getWinningBid(auctionId)` | Return the current highest bid amount. Return `-1` if no bids exist. |

**Constraints:**
- User IDs and auction IDs are auto-assigned starting from 1
- A seller cannot bid on their own auction
- Bids must strictly exceed the current highest (no equal bids)
- Base price is always > 0

**Entry points (tests will call these):**
```cpp
AuctionSystem();
int registerUser(string name, string type);  // type: "BUYER" or "SELLER"
int createAuction(int sellerId, string item, double basePrice);
bool placeBid(int auctionId, int buyerId, double amount);
double getWinningBid(int auctionId);
```

---

## Part 2

**Extension 1 — Auction closing and state management**

Auctions can now be closed. Once closed, no further bids are accepted.

| Operation | Description |
|-----------|-------------|
| `closeAuction(auctionId)` | Close the auction. If there are bids, the highest bidder wins. If no bids, auction status becomes `NO_SALE`. Only open auctions can be closed. Return `true` if successfully closed, `false` if already closed. |
| `getAuctionStatus(auctionId)` | Return the auction's current status as a string: `"OPEN"`, `"CLOSED"`, or `"NO_SALE"`. |

**Rules:**
- Closing an already-closed or no-sale auction returns `false`
- After closing, `placeBid` on that auction must return `false`
- After closing with bids, `getWinningBid` returns the final winning amount
- After closing with no bids, `getWinningBid` returns `-1`

**Entry points (tests will call these):**
```cpp
bool closeAuction(int auctionId);
string getAuctionStatus(int auctionId);
```

---

## Part 3

**Extension 2 — Multiple auction strategies**

The platform now supports three types of auctions:

| Strategy | Behavior |
|----------|----------|
| **Ascending** (default) | Standard auction. Bids must exceed the current highest. Winner is the highest bidder at close. |
| **SealedBid** | Bids are hidden from other buyers. Any bid above base price is accepted (even if lower than previous bids). Winner is the highest bidder, revealed only at close. `getWinningBid` returns `-1` while the auction is open (bids are sealed). |
| **BuyNow** | No competitive bidding. The first buyer to bid at or above `basePrice * 1.5` (the buy-now premium) instantly wins. The auction auto-closes. Bids below the buy-now price are rejected. |

**Modified entry point:**
```cpp
int createAuction(int sellerId, string item, double basePrice, string strategyType);
// strategyType: "ASCENDING", "SEALED", "BUYNOW"
```

The original `createAuction(sellerId, item, basePrice)` should default to `"ASCENDING"`.

**Design goal:** Adding a 4th auction strategy (e.g., Dutch auction) must require **zero changes** to the `AuctionSystem` class. Each strategy encapsulates its own bid validation and winner determination logic.

---

## Running Tests

```bash
./run-tests.sh 019-auction-system cpp
```
