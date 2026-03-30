#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <unordered_map>
using namespace std;

// ─── Data Structure ──────────────────────────────────────────────────────────

struct CartItem {
    string name;       // "Laptop", "Phone Case", "Headphones"
    double price;      // unit price in rupees
    int    quantity;   // number of units
    string category;   // "electronics", "accessories", "clothing"
};

// ─── Strategy Interface ───────────────────────────────────────────────────────

class Discount {
public:
    virtual double apply(const vector<CartItem>& cart) = 0;
    virtual ~Discount() = default;
};

// ─── TODO: Implement Concrete Strategies ─────────────────────────────────────

class PercentageDiscount : public Discount {
    double pct;
public:
    PercentageDiscount(double percentage) : pct(percentage) {}
    double apply(const vector<CartItem>& cart) override {
        double total = 0;
        for (auto& item : cart) total += item.price * item.quantity;
        return total * (1.0 - pct / 100.0);
    }
};

class FlatDiscount : public Discount {
    double amount;
public:
    FlatDiscount(double amt) : amount(amt) {}
    double apply(const vector<CartItem>& cart) override {
        double total = 0;
        for (auto& item : cart) total += item.price * item.quantity;
        return max(0.0, total - amount);
    }
};

class BuyXGetYDiscount : public Discount {
    int buyCount, freeCount;
public:
    BuyXGetYDiscount(int buy, int free) : buyCount(buy), freeCount(free) {}
    double apply(const vector<CartItem>& cart) override {
        double total = 0;
        int groupSize = buyCount + freeCount;
        for (auto& item : cart) {
            int groups = item.quantity / groupSize;
            int remainder = item.quantity % groupSize;
            int paidItems = groups * buyCount + min(remainder, buyCount);
            total += paidItems * item.price;
        }
        return total;
    }
};

// ─── TODO: Implement DiscountEngine ──────────────────────────────────────────

class DiscountEngine {
    Discount* discount;
public:
    DiscountEngine(Discount* d) : discount(d) {}
    void setDiscount(Discount* d) { discount = d; }

    double computeTotal(const vector<CartItem>& cart) {
        return discount->apply(cart);
    }
};

// ─── TODO: Implement StackedDiscount (Decorator) ─────────────────────────────

class StackedDiscount : public Discount {
    vector<Discount*> discounts;
public:
    StackedDiscount(vector<Discount*> d) : discounts(d) {}

    double apply(const vector<CartItem>& cart) override {
        double current = 0;
        for (auto& item : cart) current += item.price * item.quantity;
        for (auto* d : discounts) {
            vector<CartItem> temp = {{"subtotal", current, 1, ""}};
            current = d->apply(temp);
        }
        return current;
    }
};

// ─── TODO: Implement Eligibility Rules ───────────────────────────────────────

struct UserContext {
    bool isFirstTimeUser;
};

// ─── Test Entry Points ───────────────────────────────────────────────────────

double apply_percentage_discount(vector<CartItem> cart, double percentage) {
    PercentageDiscount d(percentage);
    return DiscountEngine(&d).computeTotal(cart);
}

double apply_flat_discount(vector<CartItem> cart, double amount) {
    FlatDiscount d(amount);
    return DiscountEngine(&d).computeTotal(cart);
}

double apply_bogo(vector<CartItem> cart, int buyCount, int freeCount) {
    BuyXGetYDiscount d(buyCount, freeCount);
    return DiscountEngine(&d).computeTotal(cart);
}

double apply_stacked_discounts(vector<CartItem> cart, vector<Discount*> discounts) {
    StackedDiscount d(discounts);
    return DiscountEngine(&d).computeTotal(cart);
}

double apply_with_eligibility(vector<CartItem> cart,
                              Discount* discount,
                              double minCartValue,
                              bool requireFirstTimeUser,
                              UserContext user,
                              string eligibleCategory) {
    double rawTotal = 0;
    for (auto& item : cart) rawTotal += item.price * item.quantity;

    // Rule 1: minimum cart value check
    if (rawTotal < minCartValue) return rawTotal;

    // Rule 2: first-time user check
    if (requireFirstTimeUser && !user.isFirstTimeUser) return rawTotal;

    // Rule 3: category-specific discount
    if (!eligibleCategory.empty()) {
        vector<CartItem> eligible;
        double nonEligibleTotal = 0;
        for (auto& item : cart) {
            if (item.category == eligibleCategory) {
                eligible.push_back(item);
            } else {
                nonEligibleTotal += item.price * item.quantity;
            }
        }
        return discount->apply(eligible) + nonEligibleTotal;
    }

    return discount->apply(cart);
}

// ─── Main (test your implementation) ─────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    vector<CartItem> cart = {
        {"Laptop",     50000.0, 1, "electronics"},
        {"Phone Case",   500.0, 3, "accessories"},
        {"Headphones",  2000.0, 2, "electronics"},
    };

    cout << "10% discount: " << apply_percentage_discount(cart, 10.0) << endl;
    cout << "Rs.200 flat:  " << apply_flat_discount(cart, 200.0) << endl;
    cout << "Buy 2 Get 1:  " << apply_bogo(cart, 2, 1) << endl;

    return 0;
}
#endif
