# Problem 004 — Vending Machine

**Tier:** 1 (Foundation) | **Pattern:** State | **DSA:** HashMap
**Companies:** Amazon, Flipkart | **Time:** 45 minutes

---

## Problem Statement

Model a vending machine that transitions through distinct states based on user actions. The machine holds inventory, accepts money, dispenses items, and handles cancellations — all with correct state-based behavior.

**The key design question:** Where do the transitions live? If you use `if (state == IDLE) { ... }` inside every method, adding a new state means touching every method. The State pattern gives each state its own class, so new states are added without modifying existing ones.

---

## Before You Code

> The key question: where do the transitions live?

**Naive approach:** One big class with `if (state == IDLE) { ... } else if (state == HAS_COIN) { ... }` in every method. What happens when you add a "maintenance mode" state? You touch every method.

**State pattern approach:** Each state is its own class. `IdleState::insertMoney()` transitions to `PaymentPendingState`. `PaymentPendingState::insertMoney()` accumulates more money. The machine just delegates to whatever state is current.

---

## Inventory

```cpp
struct Item {
    string name;
    double price;
    int quantity;
};

// Machine holds inventory as: unordered_map<string, Item>
// O(1) lookup by item name
```

---

## Part 1

**Base requirement — State transitions**

Implement a vending machine with these states:

| State | Description |
|-------|-------------|
| `Idle` | Waiting for item selection |
| `PaymentPending` | Item selected, waiting for enough money |
| `Dispensing` | Money sufficient, dispense item |

**Transition rules:**
- `Idle` + `selectItem(name)` → `PaymentPending` (if item exists and has stock)
- `PaymentPending` + `insertMoney(amount)` → `Dispensing` (if total ≥ price)
- `PaymentPending` + `cancel()` → `Idle` (refund inserted money)
- `Dispensing` + `dispense()` → `Idle` (decrement stock, return change)
- Any invalid action (e.g., `dispense()` in `Idle`) prints a warning and does nothing

**Entry points (tests will call these):**
```cpp
void selectItem(const string& itemName);
void insertMoney(double amount);
void dispense();
void cancel();
string getState();
void reset();  // reset machine to initial Idle state for tests
```

**What to implement:**
```cpp
class VendingMachineState {
public:
    virtual void selectItem(const string& item) = 0;
    virtual void insertMoney(double amount) = 0;
    virtual void dispense() = 0;
    virtual void cancel() = 0;
    virtual string getName() = 0;
};

class IdleState          : public VendingMachineState { ... };
class PaymentPendingState: public VendingMachineState { ... };
class DispensingState    : public VendingMachineState { ... };

class VendingMachine { ... };  // delegates all actions to currentState
```

---

## Part 2

**Extension 1 — Maintenance mode**

The operations team needs a **Maintenance state** for restocking and price adjustments:

- Operator enters maintenance with a PIN: `enterMaintenance("1234")`
- In maintenance: user-facing operations (`selectItem`, `insertMoney`, `dispense`) are blocked
- Operator can restock: `restock("Cola", 10)` — only works in Maintenance state
- Operator exits maintenance with PIN: `exitMaintenance("1234")`
- Wrong PIN → stays in current state, prints error

**New entry points:**
```cpp
void enterMaintenance(const string& operatorPin);
void exitMaintenance(const string& operatorPin);
void restock(const string& itemName, int quantity);
```

**Design challenge:** Does `MaintenanceState` fit naturally into your existing State hierarchy? You should be able to add it without modifying `IdleState`, `PaymentPendingState`, or `DispensingState`.

---

## Running Tests

```bash
./run-tests.sh 004-vending-machine cpp
```
