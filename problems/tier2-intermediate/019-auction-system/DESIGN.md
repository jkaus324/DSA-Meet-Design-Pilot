# Design Walkthrough — Online Auction System

> This file is the answer guide. Only read after you've attempted the problem.

---

## The Core Design Decision

The auction system has two independent axes of variation:
1. **Auction strategy** — how bids are accepted and winners determined (Ascending, SealedBid, BuyNow)
2. **Auction state** — lifecycle transitions (Open -> Closed / NoSale)

Trying to handle both with if-else chains creates an explosion of conditions. Instead:

```
AuctionSystem (Singleton)
    ├── UserRegistry: HashMap<userId, User>
    ├── AuctionRegistry: HashMap<auctionId, Auction>
    └── Auction
            ├── AuctionStrategy (interface)
            │       ├── AscendingStrategy
            │       ├── SealedBidStrategy
            │       └── BuyNowStrategy
            └── AuctionStatus (state tracking)
                    ├── OPEN
                    ├── CLOSED
                    └── NO_SALE
```

---

## Part 1: Core Auction System

The simplest correct implementation uses HashMaps for storage and direct comparisons for bid validation.

```cpp
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

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

class AuctionSystem {
    int nextUserId = 1;
    int nextAuctionId = 1;
    unordered_map<int, User> users;
    unordered_map<int, Auction> auctions;

public:
    int registerUser(string name, string type) {
        UserType ut = (type == "SELLER") ? UserType::SELLER : UserType::BUYER;
        int id = nextUserId++;
        users[id] = {id, name, ut};
        return id;
    }

    int createAuction(int sellerId, string item, double basePrice) {
        if (users.find(sellerId) == users.end()) return -1;
        if (users[sellerId].type != UserType::SELLER) return -1;
        int id = nextAuctionId++;
        auctions[id] = {id, sellerId, item, basePrice, AuctionStatus::OPEN, {}};
        return id;
    }

    bool placeBid(int auctionId, int buyerId, double amount) {
        if (auctions.find(auctionId) == auctions.end()) return false;
        if (users.find(buyerId) == users.end()) return false;
        if (users[buyerId].type != UserType::BUYER) return false;
        auto& auction = auctions[auctionId];
        if (auction.status != AuctionStatus::OPEN) return false;
        if (auction.sellerId == buyerId) return false;

        double currentHighest = auction.basePrice;
        for (auto& b : auction.bids) {
            currentHighest = max(currentHighest, b.amount);
        }
        if (amount <= currentHighest) return false;

        auction.bids.push_back({buyerId, amount});
        return true;
    }

    double getWinningBid(int auctionId) {
        if (auctions.find(auctionId) == auctions.end()) return -1;
        auto& auction = auctions[auctionId];
        if (auction.bids.empty()) return -1;
        double maxBid = 0;
        for (auto& b : auction.bids) {
            maxBid = max(maxBid, b.amount);
        }
        return maxBid;
    }
};
```

---

## Part 2: Auction Closing and State Management

Adding close functionality requires tracking state transitions.

```cpp
bool closeAuction(int auctionId) {
    if (auctions.find(auctionId) == auctions.end()) return false;
    auto& auction = auctions[auctionId];
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
```

**State transitions:**
```
OPEN ──(close with bids)──> CLOSED
OPEN ──(close no bids)───> NO_SALE
CLOSED ──(close)──────────> false (already closed)
NO_SALE ─(close)──────────> false (already closed)
```

---

## Part 3: Strategy Pattern for Auction Types

This is where the design shines. Each auction type encapsulates its own rules.

```cpp
class AuctionStrategy {
public:
    virtual bool acceptBid(Auction& auction, double amount) = 0;
    virtual double getVisibleWinningBid(const Auction& auction) = 0;
    virtual bool shouldAutoClose(const Auction& auction, double amount) = 0;
    virtual ~AuctionStrategy() = default;
};

class AscendingStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        double currentHighest = auction.basePrice;
        for (auto& b : auction.bids) currentHighest = max(currentHighest, b.amount);
        return amount > currentHighest;
    }

    double getVisibleWinningBid(const Auction& auction) override {
        if (auction.bids.empty()) return -1;
        double maxBid = 0;
        for (auto& b : auction.bids) maxBid = max(maxBid, b.amount);
        return maxBid;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        return false;  // ascending never auto-closes
    }
};

class SealedBidStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        return amount > auction.basePrice;  // any bid above base price
    }

    double getVisibleWinningBid(const Auction& auction) override {
        if (auction.status == AuctionStatus::OPEN) return -1;  // sealed!
        if (auction.bids.empty()) return -1;
        double maxBid = 0;
        for (auto& b : auction.bids) maxBid = max(maxBid, b.amount);
        return maxBid;
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        return false;
    }
};

class BuyNowStrategy : public AuctionStrategy {
public:
    bool acceptBid(Auction& auction, double amount) override {
        return amount >= auction.basePrice * 1.5;  // must meet buy-now price
    }

    double getVisibleWinningBid(const Auction& auction) override {
        if (auction.bids.empty()) return -1;
        return auction.bids.back().amount;  // only one bid matters
    }

    bool shouldAutoClose(const Auction& auction, double amount) override {
        return true;  // auto-close on successful bid
    }
};
```

The `AuctionSystem.placeBid` now delegates to the strategy:

```cpp
bool placeBid(int auctionId, int buyerId, double amount) {
    // ... validation checks ...
    auto& auction = auctions[auctionId];
    auto* strategy = strategies[auctionId];

    if (!strategy->acceptBid(auction, amount)) return false;

    auction.bids.push_back({buyerId, amount});

    if (strategy->shouldAutoClose(auction, amount)) {
        auction.status = AuctionStatus::CLOSED;
    }
    return true;
}
```

---

## Pattern Insights

### Strategy Pattern
Each auction type has different rules for:
- **Bid acceptance** — ascending requires exceeding current max; sealed accepts any bid above base; buy-now requires meeting the premium
- **Winner visibility** — ascending shows the current leader; sealed hides everything until close
- **Auto-close behavior** — only buy-now closes automatically

Encapsulating these in strategy objects means the `AuctionSystem` never uses `if (type == ...)`.

### Factory Pattern
Creating the right strategy based on a string type name is a Factory responsibility:

```cpp
AuctionStrategy* createStrategy(string type) {
    if (type == "ASCENDING") return new AscendingStrategy();
    if (type == "SEALED") return new SealedBidStrategy();
    if (type == "BUYNOW") return new BuyNowStrategy();
    return new AscendingStrategy();  // default
}
```

### State Pattern Connection
The auction status (OPEN, CLOSED, NO_SALE) controls which operations are valid. A full State pattern implementation would create state objects, but for this scale, an enum with transition guards is sufficient.

### Observer Pattern Connection
In production, you'd notify:
- Bidders when outbid (ascending)
- Sellers when a new bid arrives
- Winners when the auction closes

---

## What Interviewers Look For

1. **Strategy identification** — Did you recognize that auction types have different bidding rules?
2. **Clean separation** — Is bid validation logic in the strategy, not in the main system?
3. **State management** — Did you handle all state transitions correctly (open->closed, open->no_sale)?
4. **Edge cases** — Seller bidding on own auction, bids on closed auctions, bids below current highest
5. **Extensibility** — Can you add a new auction type without modifying `AuctionSystem`?

---

## Common Interview Follow-ups

- *"What if you need to support time-based auto-close?"* — Add a timer or check current time against auction end time on each operation
- *"How would you notify outbid users?"* — Observer pattern: register bidders as observers, notify when a new bid exceeds theirs
- *"What about concurrent bids?"* — Mutex per auction, or optimistic locking with version numbers
- *"How would you persist this?"* — Serialize auctions to a database; reconstruct strategy from stored type on load
