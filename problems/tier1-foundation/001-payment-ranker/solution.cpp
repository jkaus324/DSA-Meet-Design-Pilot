#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Structure ──────────────────────────────────────────────────────────

struct PaymentMethod {
    string name;
    double cashbackRate;        // e.g. 0.05 = 5%
    double transactionFee;      // in rupees
    int    usageCount;
    bool   easyRefundEligible;  // Part 3: refund eligibility flag

    // Default easyRefundEligible to false so Part 1/2 tests
    // that construct with 4 args still compile.
    PaymentMethod(string n, double cb, double fee, int usage, bool refund = false)
        : name(move(n)), cashbackRate(cb), transactionFee(fee),
          usageCount(usage), easyRefundEligible(refund) {}
};

// ─── Strategy Interface ───────────────────────────────────────────────────────

class RankingStrategy {
public:
    virtual bool compare(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~RankingStrategy() = default;
};

// ─── Part 1: Concrete Strategies ─────────────────────────────────────────────

class RewardsMaximizer : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return a.cashbackRate > b.cashbackRate;
    }
};

class LowFeeSeeker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return a.transactionFee < b.transactionFee;
    }
};

class TrustBasedRanker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        return a.usageCount > b.usageCount;
    }
};

// ─── Part 2: Composite Strategy ──────────────────────────────────────────────

class CompositeStrategy : public RankingStrategy {
private:
    vector<RankingStrategy*> criteria;
public:
    CompositeStrategy(vector<RankingStrategy*> c) : criteria(move(c)) {}
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        for (auto* s : criteria) {
            if (s->compare(a, b)) return true;   // a wins on this criterion
            if (s->compare(b, a)) return false;   // b wins on this criterion
        }
        return false;  // complete tie across all criteria
    }
};

// ─── Part 3: Easy-Refund Strategy ────────────────────────────────────────────

class EasyRefundStrategy : public RankingStrategy {
private:
    bool prefer;
public:
    EasyRefundStrategy(bool preferRefund) : prefer(preferRefund) {}
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        if (!prefer) return false;
        // a wins if it has easy refund and b does not
        return a.easyRefundEligible && !b.easyRefundEligible;
    }
};

// ─── PaymentRanker (context class) ───────────────────────────────────────────

class PaymentRanker {
private:
    RankingStrategy* strategy;
public:
    PaymentRanker(RankingStrategy* s) : strategy(s) {}
    void setStrategy(RankingStrategy* s) { strategy = s; }
    vector<PaymentMethod> rank(vector<PaymentMethod> methods) {
        sort(methods.begin(), methods.end(),
             [this](const PaymentMethod& a, const PaymentMethod& b) {
                 return strategy->compare(a, b);
             });
        return methods;
    }
};

// ─── Free-function wrappers (used by tests) ──────────────────────────────────

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

vector<PaymentMethod> rank_with_refund_filter(vector<PaymentMethod> methods, bool preferEasyRefund) {
    EasyRefundStrategy refundStrat(preferEasyRefund);
    RewardsMaximizer rewardsStrat;
    // Refund filter takes priority; cashback is the tiebreaker
    CompositeStrategy s({&refundStrat, &rewardsStrat});
    return PaymentRanker(&s).rank(methods);
}

// ─── Main (guarded so test harness can include this file) ────────────────────

#ifndef RUNNING_TESTS
int main() {
    vector<PaymentMethod> methods = {
        {"HDFC Credit Card", 0.05, 2.0, 120},
        {"UPI",              0.00, 0.0, 450},
        {"Amazon Pay",       0.03, 0.0, 200},
        {"Debit Card",       0.01, 1.0,  80},
    };

    RewardsMaximizer rewardsStrategy;
    PaymentRanker ranker(&rewardsStrategy);

    auto ranked = ranker.rank(methods);
    cout << "Ranked by Rewards:\n";
    for (const auto& m : ranked) {
        cout << "  " << m.name << " (" << m.cashbackRate * 100 << "% cashback)\n";
    }

    cout << "\nRanked with Easy-Refund preference:\n";
    vector<PaymentMethod> refundMethods = {
        {"Card A", 0.10, 5.0, 300, false},
        {"Card B", 0.02, 2.0, 500, true},
        {"Card C", 0.05, 3.0, 400, true},
    };
    auto refundRanked = rank_with_refund_filter(refundMethods, true);
    for (const auto& m : refundRanked) {
        cout << "  " << m.name << " (refund=" << m.easyRefundEligible
             << ", cashback=" << m.cashbackRate * 100 << "%)\n";
    }

    return 0;
}
#endif
