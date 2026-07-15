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

// TODO: design and implement your solution.
// Required free functions:
//   double apply_percentage_discount(vector<CartItem> cart, double percentage);
//   double apply_flat_discount(vector<CartItem> cart, double amount);
//   double apply_bogo(vector<CartItem> cart, int buyCount, int freeCount);
//   double apply_percentage_with_eligibility(vector<CartItem> cart, double percentage, double minCartValue, bool requireFirstTimeUser, bool isFirstTimeUser, string eligibleCategory);

double apply_percentage_discount(vector<CartItem> cart, double percentage) {
    // TODO: write your solution
    return {};
}

double apply_flat_discount(vector<CartItem> cart, double amount) {
    // TODO: write your solution
    return {};
}

double apply_bogo(vector<CartItem> cart, int buyCount, int freeCount) {
    // TODO: write your solution
    return {};
}

double apply_percentage_with_eligibility(vector<CartItem> cart, double percentage, double minCartValue, bool requireFirstTimeUser, bool isFirstTimeUser, string eligibleCategory) {
    // TODO: write your solution
    return {};
}
