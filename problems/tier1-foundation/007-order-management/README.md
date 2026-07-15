# Problem 007 — Order Management System

**Tier:** 1 (Foundation) | **Pattern:** State | **DSA:** HashMap
**Companies:** Meesho, PhonePe, Amazon | **Time:** 45 minutes

---

## Problem Statement

You're building the order management backend for an e-commerce platform. Every order goes through a defined lifecycle: it starts as **Created**, then moves to **Confirmed**, **Shipped**, and finally **Delivered**.

Not every transition is valid — you can't go from Shipped back to Created, and you can't deliver an order that hasn't been shipped yet. Your system must enforce these rules.

**Your task:** Design and implement an `OrderManager` that tracks orders in a `HashMap` and enforces a strict state machine for order lifecycle transitions.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What are the valid states? Created, Confirmed, Shipped, Delivered.
2. What transitions are allowed? Created->Confirmed, Confirmed->Shipped, Shipped->Delivered. No skipping, no going backward.
3. How do you store orders for O(1) lookup? A HashMap keyed by order ID.
4. What happens when someone tries an invalid transition? Return false / throw — the order state must not change.

**The key insight:** Each state defines which transitions are legal. The State pattern models this naturally — but even without formal State objects, you need a clear transition table.

---

## Data Structures

```cpp
enum class OrderState { Created, Confirmed, Shipped, Delivered };

struct OrderItem {
    std::string productId;
    int quantity;
};

struct Order {
    std::string id;
    std::vector<OrderItem> items;
    double totalAmount;
    OrderState state;
};
```

---

## Base Requirement — Order lifecycle with state transitions

Implement an `OrderManager` that creates orders and moves them through valid state transitions.

Valid transitions:
| From | To |
|------|------|
| Created | Confirmed |
| Confirmed | Shipped |
| Shipped | Delivered |

Any other transition (e.g., Shipped->Created, Created->Delivered) must be rejected.

**Entry points (tests will call these):**
```cpp
std::string create_order(std::vector<OrderItem> items, double totalAmount);
bool confirm_order(const std::string& orderId);
bool ship_order(const std::string& orderId);
bool deliver_order(const std::string& orderId);
OrderState get_order_state(const std::string& orderId);
```

**What to implement:**
```cpp
class OrderManager {
    std::unordered_map<std::string, Order> orders;
    int nextId = 1;
public:
    std::string createOrder(std::vector<OrderItem> items, double totalAmount);
    bool confirmOrder(const std::string& orderId);
    bool shipOrder(const std::string& orderId);
    bool deliverOrder(const std::string& orderId);
    OrderState getOrderState(const std::string& orderId);
};
```

---

## Extension 1 — Cancellation + inventory release

The product team now wants orders to be **cancellable** — but only from certain states.

Rules:
- An order can be cancelled from **Created** or **Confirmed** state.
- An order **cannot** be cancelled once it has been **Shipped** or **Delivered**.
- When an order is cancelled, its items must be **released back to inventory** (increment stock counts in a HashMap).

You must also maintain an inventory system:

```cpp
void set_inventory(const std::string& productId, int quantity);
int get_inventory(const std::string& productId);
bool cancel_order(const std::string& orderId);
```

When an order is created, inventory is decremented. When cancelled, inventory is restored.

**Design challenge:** How do you handle the inventory rollback cleanly? What if the order has multiple items?

---

## Extension 2 — Transition history tracking

The ops team wants to see the **full history** of every state transition an order has gone through, with timestamps.

```cpp
struct StateTransition {
    OrderState fromState;
    OrderState toState;
    long long timestamp;  // unix timestamp in milliseconds
};
```

**New entry points:**
```cpp
std::vector<StateTransition> get_order_history(const std::string& orderId);
```

Additionally, implement an **OrderObserver** interface that gets notified on every successful state transition:

```cpp
class OrderObserver {
public:
    virtual void onStateChange(const std::string& orderId,
                               OrderState from, OrderState to) = 0;
    virtual ~OrderObserver() = default;
};
```

The `OrderManager` should accept observers via `addObserver()` and notify all of them whenever a transition succeeds.

**Design challenge:** How do you decouple the notification logic from the state machine? What if you want to add logging, analytics, or alerts without modifying OrderManager?

---

## Running Tests

```bash
./run-tests.sh 007-order-management cpp
```
