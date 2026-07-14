#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given — do not modify).
struct PaymentMethod {
    string name;
    double cashbackRate;
    double transactionFee;
    int usageCount;
    bool easyRefundEligible;
    PaymentMethod(const string& name_, double cashbackRate_, double transactionFee_, int usageCount_, bool easyRefundEligible_ = false)
      : name(name_), cashbackRate(cashbackRate_), transactionFee(transactionFee_), usageCount(usageCount_), easyRefundEligible(easyRefundEligible_) {}
};

// RankingStrategy — strategy interface (given). Implement compare() in each concrete type.
class RankingStrategy {
public:
    virtual bool compare(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~RankingStrategy() = default;
};

// Concrete strategies — fill in compare() bodies.
class RewardsMaximizer : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: implement this
        return false;
    }
};

class LowFeeSeeker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: implement this
        return false;
    }
};

class TrustBasedRanker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: implement this
        return false;
    }
};

vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods) {
    // TODO: implement this
    return methods;
}

vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods) {
    // TODO: implement this
    return methods;
}

vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods) {
    // TODO: implement this
    return methods;
}

vector<PaymentMethod> rank_composite(vector<PaymentMethod> methods, vector<RankingStrategy*> criteria) {
    // TODO: implement this
    return methods;
}
