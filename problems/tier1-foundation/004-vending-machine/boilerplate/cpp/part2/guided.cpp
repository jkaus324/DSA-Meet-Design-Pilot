#include <iostream>
#include <string>
#include <unordered_map>
using namespace std;

struct Item { string name; double price; int quantity; };

// ─── State Interface ──────────────────────────────────────────────────────────

class VendingMachineState {
public:
    virtual void selectItem(const string& item) = 0;
    virtual void insertMoney(double amount) = 0;
    virtual void dispense() = 0;
    virtual void cancel() = 0;
    virtual string getName() = 0;
    virtual ~VendingMachineState() = default;
};

class VendingMachine; // forward declare

// ─── Copy your Part 1 states here ────────────────────────────────────────────

class IdleState : public VendingMachineState {
    VendingMachine* machine;
public:
    IdleState(VendingMachine* m) : machine(m) {}
    void selectItem(const string& item) override { /* TODO */ }
    void insertMoney(double) override { cout << "Select an item first." << endl; }
    void dispense() override { cout << "Select an item first." << endl; }
    void cancel() override { cout << "Nothing to cancel." << endl; }
    string getName() override { return "Idle"; }
};

// TODO: Add ItemSelected, PaymentPending, Dispensing states from Part 1

// ─── NEW: MaintenanceState ────────────────────────────────────────────────────
// HINT: All user-facing operations should print "Machine in maintenance" and return.
// Only operator operations (restock, exit) are allowed.

class MaintenanceState : public VendingMachineState {
    VendingMachine* machine;
public:
    MaintenanceState(VendingMachine* m) : machine(m) {}
    void selectItem(const string&) override { cout << "Machine in maintenance mode." << endl; }
    void insertMoney(double) override { cout << "Machine in maintenance mode." << endl; }
    void dispense() override { cout << "Machine in maintenance mode." << endl; }
    void cancel() override { cout << "Machine in maintenance mode." << endl; }
    string getName() override { return "Maintenance"; }
};

// ─── VendingMachine ───────────────────────────────────────────────────────────

class VendingMachine {
    VendingMachineState* currentState;
    unordered_map<string, Item> inventory;
    double insertedMoney = 0;
    string selectedItem;
    string operatorPin = "1234";
public:
    VendingMachine();
    void setState(VendingMachineState* s) { currentState = s; }
    void selectItem(const string& item) { currentState->selectItem(item); }
    void insertMoney(double amt) { currentState->insertMoney(amt); }
    void dispense() { currentState->dispense(); }
    void cancel() { currentState->cancel(); }
    string getState() { return currentState->getName(); }

    // Operator operations
    void enterMaintenance(const string& pin) {
        // TODO: if pin == operatorPin, switch to MaintenanceState
    }
    void exitMaintenance(const string& pin) {
        // TODO: if pin == operatorPin and in maintenance, switch to IdleState
    }
    void restock(const string& itemName, int qty) {
        // TODO: only works in MaintenanceState
    }
};

int main() {
    cout << "Part 2: Maintenance mode — implement TODOs above." << endl;
    return 0;
}
