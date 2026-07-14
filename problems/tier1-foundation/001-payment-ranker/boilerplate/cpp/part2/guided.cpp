#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct PaymentMethod {
    string name;
    double cashbackRate;
    double transactionFee;
    int usageCount;
    bool easyRefundEligible;
    PaymentMethod(const string& name_, double cashbackRate_, double transactionFee_, int usageCount_, bool easyRefundEligible_ = false)
      : name(name_), cashbackRate(cashbackRate_), transactionFee(transactionFee_), usageCount(usageCount_), easyRefundEligible(easyRefundEligible_) {}
};

// Forward declaration so signatures compile; design and implement your own.
class RankingStrategy;

// HINT: introduce an abstraction so new ranking rules don't change existing code.
// HINT: keep the comparator small — one rule per class.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods) {
    // TODO: write your solution
    return methods;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods) {
    // TODO: write your solution
    return methods;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods) {
    // TODO: write your solution
    return methods;
}

// HINT: think about how to compose multiple criteria into a single decision.
vector<PaymentMethod> rank_composite(vector<PaymentMethod> methods, vector<RankingStrategy*> criteria) {
    // TODO: write your solution
    return methods;
}
