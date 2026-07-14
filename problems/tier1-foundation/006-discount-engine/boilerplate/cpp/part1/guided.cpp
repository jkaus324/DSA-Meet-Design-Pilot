#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct CartItem {
    string name;
    double price;
    int quantity;
    string category;
    CartItem(const string& name_, double price_, int quantity_, const string& category_ = "")
      : name(name_), price(price_), quantity(quantity_), category(category_) {}
};

// HINT: introduce an abstraction so new ranking rules don't change existing code.
// HINT: keep the comparator small — one rule per class.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
double apply_percentage_discount(vector<CartItem> cart, double percentage) {
    // TODO: write your solution
    return {};
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
double apply_flat_discount(vector<CartItem> cart, double amount) {
    // TODO: write your solution
    return {};
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
double apply_bogo(vector<CartItem> cart, int buyCount, int freeCount) {
    // TODO: write your solution
    return {};
}
