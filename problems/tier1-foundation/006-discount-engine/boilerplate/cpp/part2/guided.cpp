#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct CartItem {
    string name;
    double price;
    int    quantity;
    string category;
};

// ─── Discount Interface ─────────────────────────────────────────────────────

class Discount {
public:
    virtual double apply(const vector<CartItem>& cart) = 0;
    virtual ~Discount() = default;
};

// ─── Existing Strategies ─────────────────────────────────────────────────────
// TODO: Copy your Part 1 strategies here (or extend them)

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

// ─── NEW: Stacked Discount (Decorator) ──────────────────────────────────────
// HINT: A StackedDiscount holds a list of other discounts.
// It applies the first discount to the raw total, then applies the second
// discount to the result, and so on. This is the Decorator pattern.

class StackedDiscount : public Discount {
private:
    vector<Discount*> discounts;
public:
    StackedDiscount(vector<Discount*> d) : discounts(d) {}

    double apply(const vector<CartItem>& cart) override {
        // TODO: Compute raw total from cart.
        // Then apply each discount sequentially to the running total.
        // HINT: Create a temporary single-item cart with the running total.
        return 0;
    }
};

// ─── Engine ──────────────────────────────────────────────────────────────────

class DiscountEngine {
private:
    Discount* discount;
public:
    DiscountEngine(Discount* d) : discount(d) {}
    void setDiscount(Discount* d) { discount = d; }
    double computeTotal(const vector<CartItem>& cart) {
        // TODO: Sort using discount
        return 0;
    }
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

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Stacked discounts — implement the TODOs above." << endl;
    return 0;
}
#endif
