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
        for (auto& item : cart) {
            int groups = item.quantity / (buyCount + freeCount);
            int remainder = item.quantity % (buyCount + freeCount);
            int paidItems = groups * buyCount + min(remainder, buyCount);
            total += item.price * paidItems;
        }
        return total;
    }
};

// ─── StackedDiscount (Decorator) ────────────────────────────────────────────
// Chains multiple discounts: applies each to the result of the previous.

class StackedDiscount : public Discount {
private:
    vector<Discount*> discounts;
public:
    StackedDiscount(vector<Discount*> d) : discounts(d) {}

    double apply(const vector<CartItem>& cart) override {
        double current = 0;
        for (auto& item : cart) current += item.price * item.quantity;

        for (auto* d : discounts) {
            // TODO: Apply each discount to the running total.
            // HINT: Create a temporary cart with a single item representing
            // the current subtotal, then call d->apply() on it.
            vector<CartItem> temp = {{"subtotal", current, 1, ""}};
            current = d->apply(temp);
        }
        return current;
    }
};

// ─── Engine ──────────────────────────────────────────────────────────────────

class DiscountEngine {
private:
    Discount* discount;
public:
    DiscountEngine(Discount* d) : discount(d) {}
    void setDiscount(Discount* d) { discount = d; }
    vector<CartItem> cart;

    double computeTotal(const vector<CartItem>& cart) {
        // TODO: Use discount->apply() to compute the final total
        // HINT: Just delegate to the discount strategy
        return discount->apply(cart);
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
    cout << "Part 2: Stacked discounts — full scaffolding provided." << endl;
    return 0;
}
#endif
