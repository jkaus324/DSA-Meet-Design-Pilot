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

// ─── Concrete Strategies ─────────────────────────────────────────────────────

class RewardsMaximizer : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return a.cashbackRate > b.cashbackRate; // higher cashback = better
    }
};

class LowFeeSeeker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return a.transactionFee < b.transactionFee; // lower fee = better
    }
};

class TrustBasedRanker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return a.usageCount > b.usageCount; // higher usage = more trusted
    }
};

// ─── CompositeStrategy ───────────────────────────────────────────────────────
// Chains multiple strategies: tries first, falls back on tie.

class CompositeStrategy : public RankingStrategy {
private:
    vector<RankingStrategy*> criteria;
public:
    CompositeStrategy(vector<RankingStrategy*> c) : criteria(c) {}

    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        for (auto* s : criteria) {
            if (s->compare(a, b)) return true;   // a beats b on this criterion
            if (s->compare(b, a)) return false;  // b beats a on this criterion
            // tied — try next criterion
        }
        return false; // fully tied
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
        // TODO: Sort 'methods' using strategy->compare()
        // HINT: std::sort with a lambda that calls strategy->compare()
        sort(methods.begin(), methods.end(), [this](const PaymentMethod& a, const PaymentMethod& b) {
            return strategy->compare(a, b);
        });
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
    cout << "Part 2: Composite ranking — all scaffolding provided, implement rank() if not done." << endl;
    return 0;
}
#endif
