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
// HINT: This interface lets you swap ranking logic at runtime.
// What method signature would let you compare two PaymentMethods?

class /* YourInterfaceName */ {
public:
    virtual bool /* yourMethodName */(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~/* YourInterfaceName */() = default;
};

// ─── Concrete Strategies ─────────────────────────────────────────────────────
// TODO: Implement a strategy for each ranking criterion:
//   - Rewards maximizer (highest cashback first)
//   - Low-fee seeker (lowest transaction fee first)
//   - Trust-based ranker (highest usage count first)


// ─── Ranker ──────────────────────────────────────────────────────────────────
// TODO: Implement a PaymentRanker class that:
//   - Accepts any strategy (via constructor or setter)
//   - Has a rank() method that returns sorted payment methods
//   - Does NOT know about specific ranking criteria

// class PaymentRanker {
// public:
//     PaymentRanker(/* what goes here? */);
//     vector<PaymentMethod> rank(vector<PaymentMethod> methods);
// };


// ─── Test Entry Points (must exist for tests to compile) ─────────────────────
// Your solution must provide these functions:
//
//   vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods);
//
// How you implement them internally is up to you.
// ─────────────────────────────────────────────────────────────────────────────

