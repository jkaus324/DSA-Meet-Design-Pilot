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

// ─── NEW: EasyRefundStrategy ─────────────────────────────────────────────────
// HINT: When preferEasyRefund=true, methods with easyRefundEligible=true
// should always rank above those without, regardless of other criteria.

class EasyRefundStrategy : public RankingStrategy {
private:
    bool prefer;
public:
    EasyRefundStrategy(bool preferRefund) : prefer(preferRefund) {}
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: If prefer=true, a wins if a.easyRefundEligible && !b.easyRefundEligible
        return false;
    }
};

class PaymentRanker {
private:
    RankingStrategy* strategy;
public:
    PaymentRanker(RankingStrategy* s) : strategy(s) {}
    void setStrategy(RankingStrategy* s) { strategy = s; }
    vector<PaymentMethod> rank(vector<PaymentMethod> methods) {
        // TODO: implement
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
    // TODO: Use EasyRefundStrategy as the first criterion in a CompositeStrategy
    return methods;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Easy-refund filter — implement the TODOs above." << endl;
    return 0;
}
#endif
