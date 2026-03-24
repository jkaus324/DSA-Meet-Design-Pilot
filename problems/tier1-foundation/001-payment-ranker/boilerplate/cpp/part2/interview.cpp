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

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The product team now wants COMPOSITE ranking:
// rank by cashback first, then use transaction fee as tiebreaker.
//
// Think about:
//   - How do you chain ranking criteria without modifying existing strategies?
//   - What if the product team adds a 4th criterion tomorrow?
//   - Is your Part 1 design extensible enough to handle this?
//
// Entry points (must exist for tests):
//   vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_composite(vector<PaymentMethod> methods,
//       vector<???> criteria);  // you decide the type
//
// ─────────────────────────────────────────────────────────────────────────────


