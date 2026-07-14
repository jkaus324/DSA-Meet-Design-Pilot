#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given — do not modify).
struct CartItem {
    string name;
    double price;
    int quantity;
    string category;
    CartItem(const string& name_, double price_, int quantity_, const string& category_ = "")
      : name(name_), price(price_), quantity(quantity_), category(category_) {}
};

double apply_percentage_discount(vector<CartItem> cart, double percentage) {
    // TODO: implement this
    return {};
}

double apply_flat_discount(vector<CartItem> cart, double amount) {
    // TODO: implement this
    return {};
}

double apply_bogo(vector<CartItem> cart, int buyCount, int freeCount) {
    // TODO: implement this
    return {};
}
