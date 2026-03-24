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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a PaymentRanker that:
//   1. Ranks payment methods by the criteria described in the problem
//   2. Allows new ranking strategies to be added WITHOUT modifying
//      the ranker itself
//
// Think about:
//   - What abstraction lets you swap ranking logic at runtime?
//   - How would you add a 4th ranking criterion with zero changes
//     to existing code?
//   - What happens to your code when Extension 1 (cashback) is added?
//
// Entry point (must exist for tests):
//   vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods);
//
// ─────────────────────────────────────────────────────────────────────────────


