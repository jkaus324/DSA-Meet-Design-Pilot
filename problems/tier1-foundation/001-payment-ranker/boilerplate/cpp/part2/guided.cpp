#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct PaymentMethod {
    string name;
    double cashbackRate;    // e.g. 0.05 = 5%
    double transactionFee;  // in rupees
    int    usageCount;
};

// ─── Strategy Interface ──────────────────────────────────────────────────────

class RankingStrategy {
public:
    virtual bool compare(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~RankingStrategy() = default;
};

// ─── Existing Strategies ─────────────────────────────────────────────────────
// TODO: Copy your Part 1 strategies here (or extend them)

class RewardsMaximizer : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return false; // TODO: implement
    }
};

class LowFeeSeeker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return false; // TODO: implement
    }
};

class TrustBasedRanker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return false; // TODO: implement
    }
};

// ─── NEW: Composite Strategy ─────────────────────────────────────────────────
// HINT: A CompositeStrategy holds a list of other strategies.
// It tries the first strategy; if tied, falls back to the second, then third...
// This is the Composite pattern applied to a comparator.

class CompositeStrategy : public RankingStrategy {
private:
    vector<RankingStrategy*> criteria;
public:
    CompositeStrategy(vector<RankingStrategy*> c) : criteria(c) {}

    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: Iterate through criteria.
        // If criteria[i] says a > b, return true.
        // If criteria[i] says b > a, return false.
        // If tied, move to criteria[i+1].
        return false;
    }
};

// ─── Ranker ──────────────────────────────────────────────────────────────────

class PaymentRanker {
private:
    RankingStrategy* strategy;
public:
    PaymentRanker(RankingStrategy* s) : strategy(s) {}
    void setStrategy(RankingStrategy* s) { strategy = s; }
    vector<PaymentMethod> rank(vector<PaymentMethod> methods) {
        // TODO: Sort using strategy
        return methods;
    }
};

// ─── Test Entry Points ───────────────────────────────────────────────────────

vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods) {
    RewardsMaximizer s;
    return PaymentRanker(&s).rank(methods);
}

vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods) {
    LowFeeSeeker s;
    return PaymentRanker(&s).rank(methods);
}

vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods) {
    TrustBasedRanker s;
    return PaymentRanker(&s).rank(methods);
}

vector<PaymentMethod> rank_composite(vector<PaymentMethod> methods, vector<RankingStrategy*> criteria) {
    CompositeStrategy s(criteria);
    return PaymentRanker(&s).rank(methods);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Composite ranking — implement the TODOs above." << endl;
    return 0;
}
#endif
