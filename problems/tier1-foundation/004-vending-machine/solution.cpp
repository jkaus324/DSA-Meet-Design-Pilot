#include <iostream>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

struct Item {
    string name;
    double price;
    int    quantity;
};

// ─── Forward Declarations ───────────────────────────────────────────────────

class VendingMachine;

// ─── State Interface ────────────────────────────────────────────────────────

class VMState {
public:
    virtual void selectItem(VendingMachine& vm, const string& item) = 0;
    virtual void insertMoney(VendingMachine& vm, double amount) = 0;
    virtual void dispense(VendingMachine& vm) = 0;
    virtual void cancel(VendingMachine& vm) = 0;
    virtual string name() const = 0;
    virtual ~VMState() = default;
};

// ─── Concrete State Declarations ────────────────────────────────────────────

class IdleState : public VMState {
    VendingMachine* m;
public:
    IdleState(VendingMachine* machine) : m(machine) {}
    void selectItem(VendingMachine& vm, const string& item) override;
    void insertMoney(VendingMachine& vm, double amount) override {
        cout << "[Error] Select an item first." << endl;
    }
    void dispense(VendingMachine& vm) override {
        cout << "[Error] No item selected." << endl;
    }
    void cancel(VendingMachine& vm) override {
        cout << "[Info] Nothing to cancel." << endl;
    }
    string name() const override { return "Idle"; }
};

class PaymentPendingState : public VMState {
    VendingMachine* m;
public:
    PaymentPendingState(VendingMachine* machine) : m(machine) {}
    void selectItem(VendingMachine& vm, const string& item) override {
        cout << "[Info] Item already selected. Cancel first." << endl;
    }
    void insertMoney(VendingMachine& vm, double amount) override;
    void dispense(VendingMachine& vm) override {
        cout << "[Error] Insert payment first." << endl;
    }
    void cancel(VendingMachine& vm) override;
    string name() const override { return "PaymentPending"; }
};

class DispensingState : public VMState {
    VendingMachine* m;
public:
    DispensingState(VendingMachine* machine) : m(machine) {}
    void selectItem(VendingMachine& vm, const string& item) override {
        cout << "[Error] Dispensing in progress." << endl;
    }
    void insertMoney(VendingMachine& vm, double amount) override {
        cout << "[Error] Dispensing in progress." << endl;
    }
    void dispense(VendingMachine& vm) override;
    void cancel(VendingMachine& vm) override {
        cout << "[Error] Cannot cancel during dispensing." << endl;
    }
    string name() const override { return "Dispensing"; }
};

class MaintenanceState : public VMState {
    VendingMachine* m;
public:
    MaintenanceState(VendingMachine* machine) : m(machine) {}
    void selectItem(VendingMachine& vm, const string& item) override {
        cout << "[Info] Machine in maintenance." << endl;
    }
    void insertMoney(VendingMachine& vm, double amount) override {
        cout << "[Info] Machine in maintenance." << endl;
    }
    void dispense(VendingMachine& vm) override {
        cout << "[Info] Machine in maintenance." << endl;
    }
    void cancel(VendingMachine& vm) override {
        cout << "[Info] Machine in maintenance." << endl;
    }
    string name() const override { return "Maintenance"; }
};

// ─── Vending Machine Context ────────────────────────────────────────────────

class VendingMachine {
public:
    VMState* currentState;
    string   selectedItem;
    double   insertedMoney = 0;
    unordered_map<string, Item> inventory;
    string   operatorPin = "1234";

    IdleState            idle{this};
    PaymentPendingState  paymentPending{this};
    DispensingState      dispensing{this};
    MaintenanceState     maintenance{this};

    VendingMachine() : currentState(&idle) {
        inventory["Cola"]  = {"Cola",  25.0, 5};
        inventory["Chips"] = {"Chips", 15.0, 3};
    }

    void setState(VMState* s)               { currentState = s; }
    void selectItem(const string& item)     { currentState->selectItem(*this, item); }
    void insertMoney(double amt)            { currentState->insertMoney(*this, amt); }
    void dispense()                         { currentState->dispense(*this); }
    void cancel()                           { currentState->cancel(*this); }
    string getState() const                 { return currentState->name(); }

    void resetMachine() {
        insertedMoney = 0;
        selectedItem  = "";
        currentState  = &idle;
        inventory["Cola"]  = {"Cola",  25.0, 5};
        inventory["Chips"] = {"Chips", 15.0, 3};
    }

    void enterMaintenance(const string& pin) {
        if (pin == operatorPin) {
            setState(&maintenance);
            cout << "Entered maintenance mode." << endl;
        } else {
            cout << "Invalid PIN." << endl;
        }
    }

    void exitMaintenance(const string& pin) {
        if (pin == operatorPin && getState() == "Maintenance") {
            setState(&idle);
            cout << "Exited maintenance mode." << endl;
        } else {
            cout << "Invalid PIN or not in maintenance." << endl;
        }
    }

    void restock(const string& itemName, int qty) {
        if (getState() != "Maintenance") {
            cout << "Must be in maintenance to restock." << endl;
            return;
        }
        inventory[itemName].quantity += qty;
        cout << "Restocked " << itemName << " by " << qty << endl;
    }
};

// ─── State Method Implementations ───────────────────────────────────────────

void IdleState::selectItem(VendingMachine& vm, const string& item) {
    if (vm.inventory.count(item) && vm.inventory[item].quantity > 0) {
        vm.selectedItem = item;
        vm.setState(&vm.paymentPending);
        cout << "Selected: " << item << ". Insert Rs." << vm.inventory[item].price << endl;
    } else {
        cout << "Item unavailable." << endl;
    }
}

void PaymentPendingState::insertMoney(VendingMachine& vm, double amt) {
    vm.insertedMoney += amt;
    double price = vm.inventory[vm.selectedItem].price;
    if (vm.insertedMoney >= price) {
        vm.setState(&vm.dispensing);
        cout << "Payment accepted." << endl;
    } else {
        cout << "Need Rs." << (price - vm.insertedMoney) << " more." << endl;
    }
}

void PaymentPendingState::cancel(VendingMachine& vm) {
    cout << "Refunding Rs." << vm.insertedMoney << endl;
    vm.insertedMoney = 0;
    vm.selectedItem  = "";
    vm.setState(&vm.idle);
}

void DispensingState::dispense(VendingMachine& vm) {
    vm.inventory[vm.selectedItem].quantity--;
    double change = vm.insertedMoney - vm.inventory[vm.selectedItem].price;
    cout << "Dispensed: " << vm.selectedItem << endl;
    if (change > 0) {
        cout << "Change: Rs." << change << endl;
    }
    vm.insertedMoney = 0;
    vm.selectedItem  = "";
    vm.setState(&vm.idle);
}

// ─── Global Instance & Free Function Wrappers ───────────────────────────────

static VendingMachine g_vm;

string getState()                              { return g_vm.getState(); }
void   reset()                                 { g_vm.resetMachine(); }
void   selectItem(const string& item)          { g_vm.selectItem(item); }
void   insertMoney(double amount)              { g_vm.insertMoney(amount); }
void   dispense()                              { g_vm.dispense(); }
void   cancel()                                { g_vm.cancel(); }
void   enterMaintenance(const string& pin)     { g_vm.enterMaintenance(pin); }
void   exitMaintenance(const string& pin)      { g_vm.exitMaintenance(pin); }
void   restock(const string& item, int qty)    { g_vm.restock(item, qty); }

// ─── Main ───────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    cout << "=== Vending Machine — State Pattern Reference Solution ===" << endl;
    cout << "State: " << getState() << endl;

    cout << "\n--- Normal purchase ---" << endl;
    selectItem("Cola");
    insertMoney(25.0);
    dispense();
    cout << "State: " << getState() << endl;

    cout << "\n--- Cancel mid-payment ---" << endl;
    selectItem("Chips");
    insertMoney(5.0);
    cancel();
    cout << "State: " << getState() << endl;

    cout << "\n--- Maintenance mode ---" << endl;
    enterMaintenance("1234");
    cout << "State: " << getState() << endl;
    restock("Cola", 10);
    exitMaintenance("1234");
    cout << "State: " << getState() << endl;

    return 0;
}
#endif
