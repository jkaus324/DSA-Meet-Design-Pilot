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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a VendingMachine that:
//   1. Has distinct states: Idle, ItemSelected, PaymentPending, Dispensing
//   2. Handles transitions between states based on user actions
//   3. Adding a new state requires NO changes to existing state logic
//
// Think about:
//   - What happens if the user tries to pay before selecting an item?
//   - How do you prevent invalid state transitions?
//   - What pattern lets each state handle its own logic independently?
//
// Entry points:
//   void selectItem(const string& itemName);
//   void insertMoney(double amount);
//   void dispense();
//   void cancel();
//   string getState();
//
// ─────────────────────────────────────────────────────────────────────────────

