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

// ─── Eligibility Rules ──────────────────────────────────────────────────────

struct UserContext {
    bool isFirstTimeUser;
};

class DiscountEngine {
    Discount* discount;
public:
    DiscountEngine(Discount* d) : discount(d) {}
    void setDiscount(Discount* d) { discount = d; }
    double computeTotal(const vector<CartItem>& cart) {
        return discount->apply(cart);
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
    double rawTotal = 0;
    for (auto& item : cart) rawTotal += item.price * item.quantity;

    // TODO: Rule 1 — if rawTotal < minCartValue, return rawTotal
    // TODO: Rule 2 — if requireFirstTimeUser && !user.isFirstTimeUser, return rawTotal
    // TODO: Rule 3 — if eligibleCategory is non-empty, split cart and discount only matching items

    return discount->apply(cart);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Eligibility rules — full scaffolding provided." << endl;
    return 0;
}
#endif
