#include <iostream>
#include <memory>
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

// ─── Ops simulator (used by spec-based tests) ──────────────────────────────
//
// Drives one AuctionSystem through a sequence of operations. Created users
// and auctions are stored in slots i1/i2 so subsequent ops can refer to them.
//
// Op fields:
//   "new"                                                                  -> "ok"
//   "register"     s1=name s2="SELLER"|"BUYER"  i1=user_slot              -> userId as string
//   "create"       i1=seller_slot s2=item       i3=basePrice s3=strategy(""|"ASCENDING"|"SEALED"|"BUYNOW") i2=auction_slot
//                                                                           -> auctionId as string
//   "bid"          i1=auction_slot i2=buyer_slot i3=amount                  -> "ok"/"fail"
//   "close"        i1=auction_slot                                           -> "ok"/"fail"
//   "winning"      i1=auction_slot                                           -> int as string ("-1" if none / hidden)
//   "status"       i1=auction_slot                                           -> "OPEN"|"CLOSED"|"NO_SALE"|"UNKNOWN"
//   "user_id_eq"   i1=user_slot i2=expected                                  -> "yes"/"no"

#include <unordered_map>

struct AuctionOp {
    string kind;
    string s1;
    string s2;
    string s3;
    int    i1;
    int    i2;
    int    i3;
};

vector<string> auction_simulate(vector<AuctionOp> ops) {
    vector<string> out;
    unique_ptr<AuctionSystem> sys(new AuctionSystem());
    unordered_map<int,int> userSlot, auctionSlot;
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "new") {
            sys.reset(new AuctionSystem());
            userSlot.clear();
            auctionSlot.clear();
            out.push_back("ok");
        } else if (k == "register") {
            int id = sys->registerUser(op.s1, op.s2);
            userSlot[op.i1] = id;
            out.push_back(to_string(id));
        } else if (k == "create") {
            int sid = userSlot.count(op.i1) ? userSlot[op.i1] : op.i1;
            string strat = op.s3.empty() ? "ASCENDING" : op.s3;
            int aid = sys->createAuction(sid, op.s2, (double)op.i3, strat);
            auctionSlot[op.i2] = aid;
            out.push_back(to_string(aid));
        } else if (k == "bid") {
            int aid = auctionSlot.count(op.i1) ? auctionSlot[op.i1] : op.i1;
            int bid = userSlot.count(op.i2) ? userSlot[op.i2] : op.i2;
            out.push_back(sys->placeBid(aid, bid, (double)op.i3) ? "ok" : "fail");
        } else if (k == "close") {
            int aid = auctionSlot.count(op.i1) ? auctionSlot[op.i1] : op.i1;
            out.push_back(sys->closeAuction(aid) ? "ok" : "fail");
        } else if (k == "winning") {
            int aid = auctionSlot.count(op.i1) ? auctionSlot[op.i1] : op.i1;
            double w = sys->getWinningBid(aid);
            // Use integer formatting if integral, else .2f
            if (w < 0) out.push_back("-1");
            else {
                char buf[32];
                if (w == (long long)w) snprintf(buf, sizeof(buf), "%lld", (long long)w);
                else                   snprintf(buf, sizeof(buf), "%.2f", w);
                out.push_back(buf);
            }
        } else if (k == "status") {
            int aid = auctionSlot.count(op.i1) ? auctionSlot[op.i1] : op.i1;
            out.push_back(sys->getAuctionStatus(aid));
        } else if (k == "user_id_eq") {
            int uid = userSlot.count(op.i1) ? userSlot[op.i1] : op.i1;
            out.push_back(uid == op.i2 ? "yes" : "no");
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Auction System — all parts implemented." << endl;
    return 0;
}
#endif
