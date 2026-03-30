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

// ─── Concrete Strategies ─────────────────────────────────────────────────────
// TODO: Implement the apply() method for each strategy

class PercentageDiscount : public Discount {
    double pct;
public:
    PercentageDiscount(double percentage) : pct(percentage) {}
    double apply(const vector<CartItem>& cart) override {
        // TODO: compute total, reduce by pct%
        // Higher percentage = more discount
        return 0;
    }
};

class FlatDiscount : public Discount {
    double amount;
public:
    FlatDiscount(double amt) : amount(amt) {}
    double apply(const vector<CartItem>& cart) override {
        // TODO: compute total, subtract flat amount (min 0)
        return 0;
    }
};

class BuyXGetYDiscount : public Discount {
    int buyCount, freeCount;
public:
    BuyXGetYDiscount(int buy, int free) : buyCount(buy), freeCount(free) {}
    double apply(const vector<CartItem>& cart) override {
        // TODO: for each item, calculate how many are paid
        // In groups of (buyCount + freeCount), only buyCount are paid
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
        // TODO: Use discount->apply() to compute the final total
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

#ifndef RUNNING_TESTS
int main() {
    cout << "Discount Engine — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
