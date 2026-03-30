#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

struct CartItem {
    string name;
    double price;
    int    quantity;
    string category;
};

class Discount {
public:
    virtual double apply(const vector<CartItem>& cart) = 0;
    virtual ~Discount() = default;
};

class PercentageDiscount : public Discount {
    double pct;
public:
    PercentageDiscount(double percentage) : pct(percentage) {}
    double apply(const vector<CartItem>& cart) override {
        return 0; // TODO: implement
    }
};

class FlatDiscount : public Discount {
    double amount;
public:
    FlatDiscount(double amt) : amount(amt) {}
    double apply(const vector<CartItem>& cart) override {
        return 0; // TODO: implement
    }
};

class BuyXGetYDiscount : public Discount {
    int buyCount, freeCount;
public:
    BuyXGetYDiscount(int buy, int free) : buyCount(buy), freeCount(free) {}
    double apply(const vector<CartItem>& cart) override {
        return 0; // TODO: implement
    }
};

class StackedDiscount : public Discount {
    vector<Discount*> discounts;
public:
    StackedDiscount(vector<Discount*> d) : discounts(d) {}
    double apply(const vector<CartItem>& cart) override {
        return 0; // TODO: implement
    }
};

// ─── NEW: Eligibility Rules ─────────────────────────────────────────────────
// HINT: When a rule is not met, the discount is SKIPPED (original total returned).
// When eligibleCategory is non-empty, only items in that category are discounted.

struct UserContext {
    bool isFirstTimeUser;
};

class DiscountEngine {
    Discount* discount;
public:
    DiscountEngine(Discount* d) : discount(d) {}
    void setDiscount(Discount* d) { discount = d; }
    double computeTotal(const vector<CartItem>& cart) {
        return 0; // TODO: implement
    }
};

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
    // TODO: Check eligibility rules:
    // 1. If raw total < minCartValue, return raw total (skip discount)
    // 2. If requireFirstTimeUser && !user.isFirstTimeUser, return raw total
    // 3. If eligibleCategory is non-empty, only discount items in that category
    return 0;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Eligibility rules — implement the TODOs above." << endl;
    return 0;
}
#endif
