#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

struct PaymentMethod {
    string name;
    double cashbackRate;
    double transactionFee;
    int    usageCount;
    bool   easyRefundEligible;  // NEW in Part 3
};

class RankingStrategy {
public:
    virtual bool compare(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~RankingStrategy() = default;
};

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

class CompositeStrategy : public RankingStrategy {
private:
    vector<RankingStrategy*> criteria;
public:
    CompositeStrategy(vector<RankingStrategy*> c) : criteria(c) {}
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        for (auto* s : criteria) {
            if (s->compare(a, b)) return true;
            if (s->compare(b, a)) return false;
        }
        return false;
    }
};

// ─── EasyRefundStrategy ──────────────────────────────────────────────────────

class EasyRefundStrategy : public RankingStrategy {
private:
    bool prefer;
public:
    EasyRefundStrategy(bool preferRefund) : prefer(preferRefund) {}
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        if (!prefer) return false;
        // a wins if it has easy refund and b doesn't
        return a.easyRefundEligible && !b.easyRefundEligible;
    }
};

class PaymentRanker {
private:
    RankingStrategy* strategy;
public:
    PaymentRanker(RankingStrategy* s) : strategy(s) {}
    void setStrategy(RankingStrategy* s) { strategy = s; }
    vector<PaymentMethod> rank(vector<PaymentMethod> methods) {
        sort(methods.begin(), methods.end(), [this](const PaymentMethod& a, const PaymentMethod& b) {
            return strategy->compare(a, b);
        });
        return methods;
    }
};

vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods) {
    RewardsMaximizer s; return PaymentRanker(&s).rank(methods);
}
vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods) {
    LowFeeSeeker s; return PaymentRanker(&s).rank(methods);
}
vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods) {
    TrustBasedRanker s; return PaymentRanker(&s).rank(methods);
}
vector<PaymentMethod> rank_composite(vector<PaymentMethod> methods, vector<RankingStrategy*> criteria) {
    CompositeStrategy s(criteria); return PaymentRanker(&s).rank(methods);
}
vector<PaymentMethod> rank_with_refund_filter(vector<PaymentMethod> methods, bool preferEasyRefund) {
    EasyRefundStrategy refundStrat(preferEasyRefund);
    RewardsMaximizer rewardsStrat;
    // Refund filter takes priority over rewards
    CompositeStrategy s({&refundStrat, &rewardsStrat});
    return PaymentRanker(&s).rank(methods);
}

int main() {
    cout << "Part 3: Easy-refund filter — full scaffolding provided." << endl;
    return 0;
}
