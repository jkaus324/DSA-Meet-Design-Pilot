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
// TODO: Implement the compare() method for each strategy

class RewardsMaximizer : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: return true if 'a' should rank higher than 'b'
        // Higher cashback rate = better ranking
        return false;
    }
};

class LowFeeSeeker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: return true if 'a' should rank higher than 'b'
        // Lower transaction fee = better ranking
        return false;
    }
};

class TrustBasedRanker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: return true if 'a' should rank higher than 'b'
        // Higher usage count = better ranking
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
        // TODO: Sort methods using the current strategy's compare()
        // Return the sorted vector
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

int main() {
    cout << "Payment Ranker — implement the TODO methods above, then run tests." << endl;
    return 0;
}
