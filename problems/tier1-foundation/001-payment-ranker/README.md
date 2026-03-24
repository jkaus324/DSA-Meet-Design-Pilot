# Problem 001 — Payment Method Ranker

**Tier:** 1 (Foundation) | **Pattern:** Strategy + Comparator | **DSA:** Sorting
**Companies:** Amazon, Flipkart | **Time:** 45 minutes

---

## Problem Statement

You're building the checkout page for an e-commerce platform. When a user is about to pay, the system must **rank their available payment methods** from most recommended to least recommended.

Different users have different ranking criteria:
- A **rewards maximizer** wants the card with the highest cashback rate first
- A **low-fee seeker** wants the option with the lowest transaction fee first
- A **trust-based** ranker wants the most commonly used payment method first (by usage count)

**Your task:** Design and implement a `PaymentRanker` class that can rank a list of payment methods using any of these strategies — and allows new strategies to be added without modifying the ranker itself.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What varies here? The *ranking algorithm* varies. The *ranker itself* stays the same.
2. If you used `if-else` inside `rank()`, what happens when a 4th strategy is added? You modify existing code — violating Open/Closed Principle.
3. How does the Strategy pattern solve this? Each ranking algorithm becomes a separate class implementing a common interface.

**The key insight:** A comparator *is* a strategy. When you pass a lambda or a comparator object to `std::sort`, you are already using the Strategy pattern implicitly.

---

## Data Structures

```cpp
struct PaymentMethod {
    std::string name;       // "HDFC Credit Card", "UPI", "Amazon Pay"
    double cashbackRate;    // 0.05 = 5% cashback
    double transactionFee;  // in rupees
    int usageCount;         // how many times user has used this method
};
```

---

## Part 1

**Base requirement — Rank by a single criterion**

Implement a `PaymentRanker` that accepts any ranking strategy at construction time and uses it to sort a list of payment methods.

You must support these three strategies:

| Strategy | Rule |
|----------|------|
| Rewards Maximizer | Highest `cashbackRate` first |
| Low-Fee Seeker | Lowest `transactionFee` first |
| Trust-Based Ranker | Highest `usageCount` first |

**Design goal:** Adding a 4th strategy must require **zero changes** to `PaymentRanker` itself.

**Entry points (tests will call these):**
```cpp
vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods);
vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods);
vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods);
```

**What to implement:**
```cpp
class RankingStrategy {
public:
    virtual bool compare(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~RankingStrategy() = default;
};

class RewardsMaximizer : public RankingStrategy { ... };
class LowFeeSeeker    : public RankingStrategy { ... };
class TrustBased      : public RankingStrategy { ... };

class PaymentRanker {
public:
    PaymentRanker(RankingStrategy* strategy);
    void setStrategy(RankingStrategy* strategy);
    vector<PaymentMethod> rank(vector<PaymentMethod> methods);
};
```

---

## Part 2

**Extension 1 — Composite ranking**

The product team now wants **composite ranking**: rank by cashback first, and use transaction fee as a tiebreaker when two methods have equal cashback.

> Example: Card A (10% cashback, Rs. 8 fee) vs Card B (10% cashback, Rs. 3 fee) → Card B wins because tiebreaker is lower fee.

**Design challenge:** How do you chain ranking criteria **without modifying** `RewardsMaximizer`, `LowFeeSeeker`, or `PaymentRanker`?

**New entry point:**
```cpp
vector<PaymentMethod> rank_composite(vector<PaymentMethod> methods,
                                     vector<RankingStrategy*> criteria);
```

The function accepts an ordered list of criteria. The first criterion is the primary sort key; subsequent criteria break ties.

**Hint:** A `CompositeStrategy` holds a list of strategies. It tries the first; if tied, tries the second; and so on.

---

## Part 3

**Extension 2 — Easy-refund eligibility**

The compliance team has added a new field to `PaymentMethod`:

```cpp
struct PaymentMethod {
    // ... existing fields ...
    bool easyRefundEligible;  // NEW
};
```

Some payment methods don't support easy refunds. When the user has enabled "prefer easy refund" in settings, those methods should rank lower — regardless of cashback or fees.

**New entry point:**
```cpp
vector<PaymentMethod> rank_with_refund_filter(vector<PaymentMethod> methods,
                                              bool preferEasyRefund);
```

When `preferEasyRefund = true`: methods with `easyRefundEligible = true` always rank above those without it. Among methods with the same refund eligibility, rank by cashback as tiebreaker.

When `preferEasyRefund = false`: refund eligibility is ignored; rank only by cashback.

**Design challenge:** Is this a new strategy, a filter, or both? Can your existing `CompositeStrategy` handle it?

---

## Running Tests

```bash
./run-tests.sh 001-payment-ranker cpp
```
