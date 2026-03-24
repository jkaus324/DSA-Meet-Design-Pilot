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
    bool   easyRefundEligible;  // NEW in Part 3
};

// ─── NEW in Extension 2 ──────────────────────────────────────────────────────
//
// The compliance team wants to add "easy-refund eligibility" as a filter.
// Some payment methods don't support easy refunds — these should be ranked
// lower regardless of cashback or fee, unless the user explicitly opts in.
//
// Think about:
//   - Is this a ranking criterion, a filter, or both?
//   - How does your existing Composite strategy handle a boolean filter?
//   - What if the "opt-in" flag is per-user, not per-session?
//
// Entry points (must exist for tests):
//   vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod>);
//   vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod>);
//   vector<PaymentMethod> rank_by_trust(vector<PaymentMethod>);
//   vector<PaymentMethod> rank_composite(vector<PaymentMethod>, vector<???> criteria);
//   vector<PaymentMethod> rank_with_refund_filter(vector<PaymentMethod>, bool preferEasyRefund);
//
// ─────────────────────────────────────────────────────────────────────────────


