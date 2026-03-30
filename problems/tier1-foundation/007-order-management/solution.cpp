#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <chrono>
using namespace std;

// ─── Data Structures ────────────────────────────────────────────────────────

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

// ─── TODO: Implement OrderManager ───────────────────────────────────────────

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
        history[orderId].push_back({from, next, now()});

        for (auto* obs : observers) {
            obs->onStateChange(orderId, from, next);
        }
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
        history[id].push_back({OrderState::Created, OrderState::Created, now()});
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
        history[orderId].push_back({current, OrderState::Cancelled, now()});

        for (auto* obs : observers) {
            obs->onStateChange(orderId, current, OrderState::Cancelled);
        }
        return true;
    }

    OrderState getOrderState(const string& orderId) {
        return orders.at(orderId).state;
    }

    vector<StateTransition> getOrderHistory(const string& orderId) {
        auto it = history.find(orderId);
        if (it == history.end()) return {};
        return it->second;
    }

    void addObserver(OrderObserver* obs) {
        observers.push_back(obs);
    }
};

// ─── Global Instance + Entry Points ─────────────────────────────────────────

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

bool cancel_order(const string& orderId) {
    return manager.cancelOrder(orderId);
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

// ─── Main (test your implementation) ────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    auto id = create_order({{"PROD-1", 2}}, 500.0);
    cout << "Created order: " << id << endl;

    confirm_order(id);
    cout << "Confirmed order" << endl;

    ship_order(id);
    cout << "Shipped order" << endl;

    deliver_order(id);
    cout << "Delivered order" << endl;

    return 0;
}
#endif
