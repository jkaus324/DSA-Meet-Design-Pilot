#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

enum class OrderState { Created, Confirmed, Shipped, Delivered };

struct OrderItem {
    string productId;
    int quantity;
};

struct Order {
    string id;
    vector<OrderItem> items;
    double totalAmount;
    OrderState state;
};

// ─── OrderManager ───────────────────────────────────────────────────────────

class OrderManager {
    unordered_map<string, Order> orders;
    int nextId = 1;

    bool transition(const string& orderId, OrderState expected, OrderState next) {
        auto it = orders.find(orderId);
        if (it == orders.end()) return false;
        // TODO: Check if current state matches 'expected'
        // If yes, update to 'next' and return true
        // If no, return false (invalid transition)
        return false;
    }

public:
    string createOrder(vector<OrderItem> items, double totalAmount) {
        string id = "ORD-" + to_string(nextId++);
        // TODO: Create an Order with state Created and store it in the map
        return id;
    }

    bool confirmOrder(const string& orderId) {
        // TODO: transition from Created to Confirmed
        return false;
    }

    bool shipOrder(const string& orderId) {
        // TODO: transition from Confirmed to Shipped
        return false;
    }

    bool deliverOrder(const string& orderId) {
        // TODO: transition from Shipped to Delivered
        return false;
    }

    OrderState getOrderState(const string& orderId) {
        return orders.at(orderId).state;
    }
};

// ─── Test Entry Points ──────────────────────────────────────────────────────

OrderManager manager;

string create_order(vector<OrderItem> items, double totalAmount) {
    return manager.createOrder(items, totalAmount);
}

bool confirm_order(const string& orderId) {
    return manager.confirmOrder(orderId);
}

bool ship_order(const string& orderId) {
    return manager.shipOrder(orderId);
}

bool deliver_order(const string& orderId) {
    return manager.deliverOrder(orderId);
}

OrderState get_order_state(const string& orderId) {
    return manager.getOrderState(orderId);
}

void reset_manager() {
    manager = OrderManager();
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Order Management — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
