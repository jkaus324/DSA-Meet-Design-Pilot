#include <iostream>
#include <string>
#include <unordered_map>
using namespace std;

struct Item {
    string name;
    double price;
    int    quantity;
};

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The vending machine now needs a MAINTENANCE mode:
//   - Operator can switch the machine into maintenance mode
//   - In maintenance mode: restock items, adjust prices, clear errors
//   - User-facing operations (select, pay, dispense) are blocked during maintenance
//   - Only the operator can exit maintenance mode
//
// Think about:
//   - Where does "maintenance" fit in your existing state diagram?
//   - Is it a state like Idle/PaymentPending, or a separate mode overlay?
//   - How do you prevent users from entering maintenance mode?
//
// Entry points (all from Part 1, plus):
//   void selectItem(const string& itemName);
//   void insertMoney(double amount);
//   void dispense();
//   void cancel();
//   string getState();
//   void enterMaintenance(const string& operatorPin);
//   void exitMaintenance(const string& operatorPin);
//   void restock(const string& itemName, int quantity);
//
// ─────────────────────────────────────────────────────────────────────────────


