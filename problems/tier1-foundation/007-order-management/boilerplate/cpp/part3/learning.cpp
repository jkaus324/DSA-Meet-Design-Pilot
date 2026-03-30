#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <chrono>
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

struct StateTransition {
    OrderState fromState;
    OrderState toState;
    long long timestamp;
};

// ─── Observer Interface ─────────────────────────────────────────────────────

class OrderObserver {
public:
    virtual void onStateChange(const string& orderId,
                               OrderState from, OrderState to) = 0;
    virtual ~OrderObserver() = default;
};

// ─── OrderManager ───────────────────────────────────────────────────────────

class OrderManager {
    unordered_map<string, Order> orders;
    unordered_map<string, int> inventory;
    unordered_map<string, vector<StateTransition>> history;
    vector<OrderObserver*> observers;
    int nextId = 1;

    long long now() {
        return chrono::duration_cast<chrono::milliseconds>(
            chrono::system_clock::now().time_since_epoch()).count();
    }

    bool transition(const string& orderId, OrderState expected, OrderState next) {
        auto it = orders.find(orderId);
        if (it == orders.end()) return false;
        if (it->second.state != expected) return false;

        OrderState from = it->second.state;
        it->second.state = next;

        // TODO: Record transition in history with timestamp
        // TODO: Notify all observers
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
        for (auto& item : items) {
            inventory[item.productId] -= item.quantity;
        }
        orders[id] = {id, items, totalAmount, OrderState::Created};
        // TODO: Record initial history entry for creation
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

        OrderState current = it->second.state;
        if (current != OrderState::Created && current != OrderState::Confirmed)
            return false;

        for (auto& item : it->second.items) {
            inventory[item.productId] += item.quantity;
        }

        it->second.state = OrderState::Cancelled;
        // TODO: Record cancellation in history
        // TODO: Notify observers
        return true;
    }

    OrderState getOrderState(const string& orderId) {
        return orders.at(orderId).state;
    }

    vector<StateTransition> getOrderHistory(const string& orderId) {
        // TODO: Return history for this order (empty vector if not found)
        return {};
    }

    void addObserver(OrderObserver* obs) {
        observers.push_back(obs);
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

vector<StateTransition> get_order_history(const string& orderId) {
    return manager.getOrderHistory(orderId);
}

void add_observer(OrderObserver* obs) {
    manager.addObserver(obs);
}

void reset_manager() {
    manager = OrderManager();
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: History + Observer — implement the TODOs above." << endl;
    return 0;
}
#endif
