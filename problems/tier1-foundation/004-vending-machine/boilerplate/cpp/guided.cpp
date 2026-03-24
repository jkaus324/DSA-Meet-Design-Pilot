#include <iostream>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct Item {
    string name;
    double price;
    int    quantity;
};

// ─── State Interface ─────────────────────────────────────────────────────────
// HINT: Each state handles its own version of each user action.
// If the action is invalid in this state, it prints an error message.

class VendingMachine; // forward declaration

class /* YourStateName */ {
public:
    virtual void selectItem(VendingMachine& vm, const string& itemName) = 0;
    virtual void insertMoney(VendingMachine& vm, double amount) = 0;
    virtual void dispense(VendingMachine& vm) = 0;
    virtual void cancel(VendingMachine& vm) = 0;
    virtual string name() const = 0;
    virtual ~/* YourStateName */() = default;
};

// ─── Concrete States ─────────────────────────────────────────────────────────
// TODO: Implement each state:
//   - IdleState       — waiting for item selection
//   - SelectedState   — item chosen, waiting for payment
//   - PaidState       — payment received, ready to dispense
//   - DispensingState — currently dispensing item


// ─── Vending Machine Context ─────────────────────────────────────────────────
// TODO: Implement the VendingMachine class that:
//   - Holds the current state
//   - Delegates all actions to the current state
//   - Has setState() to switch between states

// class VendingMachine {
// public:
//     void setState(/* YourStateName* */ state);
//     void selectItem(const string& itemName);
//     void insertMoney(double amount);
//     void dispense();
//     void cancel();
//     string getState();
// };

// ─────────────────────────────────────────────────────────────────────────────

