# Problem 006 — Billing & Discount Engine

**Tier:** 1 (Foundation) | **Pattern:** Strategy + Decorator | **DSA:** HashMap
**Companies:** Flipkart, Amazon, Meesho | **Time:** 45 minutes

---

## Problem Statement

You're building the billing engine for an e-commerce platform. The system must apply **discount strategies** to a shopping cart. Different discount types exist — percentage off, flat amount off, buy-X-get-Y — and the business needs to add new discount types without modifying the billing engine itself.

**Your task:** Design and implement a `DiscountEngine` that can apply any discount strategy to a cart, stack multiple discounts using the Decorator pattern, and enforce eligibility rules.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What varies here? The *discount algorithm* varies. The *billing engine* stays the same.
2. If you used `if-else` inside `applyDiscount()`, what happens when a 4th discount type is added? You modify existing code — violating Open/Closed Principle.
3. How does the Strategy pattern solve Part 1? Each discount type becomes a separate class implementing a common interface.
4. How does the Decorator pattern solve Part 2? Each decorator wraps the previous discount, adding its own logic on top.

**The key insight:** A discount strategy computes how much to subtract. A decorator *wraps* another discount and augments the result. This lets you compose arbitrary discount chains at runtime.

---

## Data Structures

```cpp
struct CartItem {
    std::string name;     // "Laptop", "Phone Case", "Headphones"
    double price;         // unit price in rupees
    int quantity;         // number of units
    std::string category; // "electronics", "accessories", "clothing"
};
```

---

## Base Requirement — Apply single discount type

Implement a discount engine that accepts any discount strategy and applies it to a list of cart items to compute the final bill amount.

You must support these three discount types:

| Discount Type | Rule |
|---------------|------|
| Percentage | Reduce total by a given percentage (e.g., 10% off) |
| Flat | Subtract a fixed amount from total (e.g., Rs. 200 off) |
| Buy-X-Get-Y | For every X items of a product, Y items are free (e.g., buy 2 get 1 free) |

**Design goal:** Adding a 4th discount type must require **zero changes** to the billing engine itself.

**Entry points (tests will call these):**
```cpp
double apply_percentage_discount(std::vector<CartItem> cart, double percentage);
double apply_flat_discount(std::vector<CartItem> cart, double amount);
double apply_bogo(std::vector<CartItem> cart, int buyCount, int freeCount);
```

**What to implement:**
```cpp
class Discount {
public:
    virtual double apply(const std::vector<CartItem>& cart) = 0;
    virtual ~Discount() = default;
};

class PercentageDiscount : public Discount { ... };
class FlatDiscount       : public Discount { ... };
class BuyXGetYDiscount   : public Discount { ... };

class DiscountEngine {
public:
    DiscountEngine(Discount* discount);
    void setDiscount(Discount* discount);
    double computeTotal(const std::vector<CartItem>& cart);
};
```

---

## Extension 1 — Stack multiple discounts

The product team now wants **stacked discounts**: apply a coupon discount, then a seasonal discount on top of the result, then a membership discount on top of that.

> Example: Cart total Rs. 1000 → 10% coupon → Rs. 900 → Rs. 50 flat seasonal → Rs. 850 → 5% membership → Rs. 807.50

**Design challenge:** How do you chain discount computations **without modifying** `PercentageDiscount`, `FlatDiscount`, or `DiscountEngine`?

**New entry point:**
```cpp
double apply_stacked_discounts(std::vector<CartItem> cart,
                               std::vector<Discount*> discounts);
```

The function accepts an ordered list of discounts. Each discount applies to the **result** of the previous one (Decorator pattern). The first discount applies to the raw cart total.

**Hint:** A `StackedDiscount` wraps a base discount and applies additional discounts sequentially on the resulting amount.

---

## Extension 2 — Discount eligibility rules

The compliance team has added eligibility rules. Not every discount can be applied to every cart:

| Rule | Condition |
|------|-----------|
| Minimum cart value | Discount only applies if raw cart total >= threshold |
| First-time user | Discount only applies if user is a first-time buyer |
| Category-specific | Discount only applies to items in a specific category |

**New entry point:**
```cpp
struct UserContext {
    bool isFirstTimeUser;
};

double apply_with_eligibility(std::vector<CartItem> cart,
                              Discount* discount,
                              double minCartValue,
                              bool requireFirstTimeUser,
                              UserContext user,
                              std::string eligibleCategory);
```

When a rule is not met, the discount is **skipped** (original total is returned for that rule). When `eligibleCategory` is non-empty, only items in that category contribute to the discountable subtotal; other items are billed at full price.

**Design challenge:** Are eligibility rules a new discount type, a decorator, or a separate concern? Can your existing design handle this without modification?

---

## Running Tests

```bash
./run-tests.sh 006-discount-engine cpp
```
