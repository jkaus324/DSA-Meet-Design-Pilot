#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

enum class OrderState { Created, Confirmed, Shipped, Delivered, Cancelled };

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
    unordered_map<string, int> inventory;
    int nextId = 1;

    bool transition(const string& orderId, OrderState expected, OrderState next) {
        auto it = orders.find(orderId);
        if (it == orders.end()) return false;
        if (it->second.state != expected) return false;
        it->second.state = next;
        return true;
    }

public:
    void setInventory(const string& productId, int qty) {
        inventory[productId] = qty;
    }

    int getInventory(const string& productId) {
        auto it = inventory.find(productId);
        return it != inventory.end() ? it->second : 0;
    }

    string createOrder(vector<OrderItem> items, double totalAmount) {
        string id = "ORD-" + to_string(nextId++);
        // TODO: Decrement inventory for each item
        // TODO: Store the order in the map with state Created
        orders[id] = {id, items, totalAmount, OrderState::Created};
        return id;
    }

    bool confirmOrder(const string& orderId) {
        return transition(orderId, OrderState::Created, OrderState::Confirmed);
    }

    bool shipOrder(const string& orderId) {
        return transition(orderId, OrderState::Confirmed, OrderState::Shipped);
    }

    bool deliverOrder(const string& orderId) {
        return transition(orderId, OrderState::Shipped, OrderState::Delivered);
    }

    bool cancelOrder(const string& orderId) {
        auto it = orders.find(orderId);
        if (it == orders.end()) return false;
        // TODO: Check if state is Created or Confirmed (otherwise return false)
        // TODO: Iterate through order items and restore inventory
        // TODO: Set state to Cancelled
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

bool cancel_order(const string& orderId) {
    return manager.cancelOrder(orderId);
}

OrderState get_order_state(const string& orderId) {
    return manager.getOrderState(orderId);
}

void set_inventory(const string& productId, int qty) {
    manager.setInventory(productId, qty);
}

int get_inventory(const string& productId) {
    return manager.getInventory(productId);
}

void reset_manager() {
    manager = OrderManager();
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Cancellation + Inventory — implement the TODOs above." << endl;
    return 0;
}
#endif
