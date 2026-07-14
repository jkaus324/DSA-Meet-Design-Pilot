# Problem 019 — Online Auction System

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer + State | **DSA:** Priority Queue + HashMap
**Companies:** Flipkart, Amazon, eBay | **Time:** 90 minutes

---

## Problem Statement

You are building an online auction platform. Sellers list items for auction with a base price; buyers place bids. The system tracks the highest bid, manages auction lifecycle (OPEN → CLOSED / NO_SALE), and supports three auction strategies with different bidding rules. Adding a new auction type must require zero changes to the `AuctionSystem` class.

**Constraints:**
- User IDs and auction IDs are auto-assigned starting from 1
- Only SELLER users can create auctions; only BUYER users can bid
- A seller cannot bid on their own auction
- Bids must be strictly greater than the current highest (for ascending auctions)
- Base price is always > 0

---

## Base Requirement — Core Auction System

Implement an `AuctionSystem` supporting user registration, auction creation, bidding, and querying the current winning bid.

**Example:**
```
sys.registerUser("Alice", "SELLER")  →  1
sys.registerUser("Bob",   "BUYER")   →  2
sys.registerUser("Carol", "BUYER")   →  3

sys.createAuction(sellerId=1, item="Laptop", basePrice=500.0)  →  1

sys.placeBid(auctionId=1, buyerId=2, amount=600.0)   →  true
sys.placeBid(auctionId=1, buyerId=3, amount=550.0)   →  false  // must exceed 600
sys.placeBid(auctionId=1, buyerId=3, amount=700.0)   →  true

sys.getWinningBid(1)  →  700.0
```

**Public methods:**
- `AuctionSystem()`
- `int registerUser(string name, string type)`
- `int createAuction(int sellerId, string item, double basePrice)`
- `bool placeBid(int auctionId, int buyerId, double amount)`
- `double getWinningBid(int auctionId)`

---

## Extension 1 — Auction Closing and State Management

Auctions can now be closed. No bids are accepted after closing. If the auction had bids, the highest bidder wins; if no bids were placed, the status becomes NO_SALE.

**Example:**
```
sys.closeAuction(1)              →  true    // closed with winner Carol at $700
sys.getAuctionStatus(1)          →  "CLOSED"
sys.placeBid(1, 2, 800.0)        →  false   // auction is closed
sys.closeAuction(1)              →  false   // already closed

sys.createAuction(1, "Mouse", 20.0)  →  2
sys.closeAuction(2)              →  true
sys.getAuctionStatus(2)          →  "NO_SALE"
sys.getWinningBid(2)             →  -1.0
```

**Public methods:**
- `bool closeAuction(int auctionId)`
- `string getAuctionStatus(int auctionId)`

---

## Extension 2 — Multiple Auction Strategies

Support three auction types at creation time. Each type encapsulates its own bid validation and winner-determination logic.

| Strategy | Behavior |
|---|---|
| ASCENDING (default) | Bids must exceed current highest. Winner is highest bidder at close. |
| SEALED | All bids above base price are accepted (bidders cannot see others' bids). `getWinningBid` returns -1 while open. Winner (highest bidder) revealed only at close. |
| BUYNOW | No competitive bidding. First buyer to bid >= `basePrice * 1.5` instantly wins and auto-closes the auction. Bids below that threshold are rejected. |

**Example:**
```
// BUYNOW auction
sys.createAuction(1, "Phone", 1000.0, "BUYNOW")  →  3
// Buy-now threshold = 1500.0
sys.placeBid(3, 2, 1200.0)  →  false  // below threshold
sys.placeBid(3, 2, 1500.0)  →  true   // accepted and auction auto-closes
sys.getAuctionStatus(3)     →  "CLOSED"
sys.getWinningBid(3)        →  1500.0

// SEALED auction
sys.createAuction(1, "Bag", 50.0, "SEALED")  →  4
sys.placeBid(4, 2, 80.0)   →  true
sys.placeBid(4, 3, 60.0)   →  true   // sealed: both accepted even though 60 < 80
sys.getWinningBid(4)       →  -1.0   // bids hidden while open
sys.closeAuction(4)
sys.getWinningBid(4)       →  80.0   // winner revealed at close
```

**Modified method:**
- `int createAuction(int sellerId, string item, double basePrice, string strategyType = "ASCENDING")`

---

## Running Tests

```bash
./run-tests.sh 019-auction-system cpp
```
