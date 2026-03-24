#include <iostream>
#include <string>
#include <unordered_map>
using namespace std;

struct Item { string name; double price; int quantity; };

class VendingMachine;

class VendingMachineState {
public:
    virtual void selectItem(const string& item) = 0;
    virtual void insertMoney(double amount) = 0;
    virtual void dispense() = 0;
    virtual void cancel() = 0;
    virtual string getName() = 0;
    virtual ~VendingMachineState() = default;
};

class IdleState : public VendingMachineState {
    VendingMachine* m;
public:
    IdleState(VendingMachine* machine) : m(machine) {}
    void selectItem(const string& item) override;
    void insertMoney(double) override { cout << "Select item first." << endl; }
    void dispense() override { cout << "Select item first." << endl; }
    void cancel() override {}
    string getName() override { return "Idle"; }
};

class PaymentPendingState : public VendingMachineState {
    VendingMachine* m;
public:
    PaymentPendingState(VendingMachine* machine) : m(machine) {}
    void selectItem(const string&) override { cout << "Already selected." << endl; }
    void insertMoney(double amt) override;
    void dispense() override { cout << "Insert money first." << endl; }
    void cancel() override;
    string getName() override { return "PaymentPending"; }
};

class DispensingState : public VendingMachineState {
    VendingMachine* m;
public:
    DispensingState(VendingMachine* machine) : m(machine) {}
    void selectItem(const string&) override { cout << "Dispensing in progress." << endl; }
    void insertMoney(double) override { cout << "Dispensing in progress." << endl; }
    void dispense() override;
    void cancel() override { cout << "Cannot cancel during dispensing." << endl; }
    string getName() override { return "Dispensing"; }
};

class MaintenanceState : public VendingMachineState {
    VendingMachine* m;
public:
    MaintenanceState(VendingMachine* machine) : m(machine) {}
    void selectItem(const string&) override { cout << "Machine in maintenance." << endl; }
    void insertMoney(double) override { cout << "Machine in maintenance." << endl; }
    void dispense() override { cout << "Machine in maintenance." << endl; }
    void cancel() override { cout << "Machine in maintenance." << endl; }
    string getName() override { return "Maintenance"; }
};

class VendingMachine {
public:
    VendingMachineState* currentState;
    unordered_map<string, Item> inventory;
    double insertedMoney = 0;
    string selectedItem;
    string operatorPin = "1234";

    IdleState idle{this};
    PaymentPendingState paymentPending{this};
    DispensingState dispensing{this};
    MaintenanceState maintenance{this};

    VendingMachine() : currentState(&idle) {
        inventory["Cola"] = {"Cola", 25.0, 5};
        inventory["Chips"] = {"Chips", 15.0, 3};
    }

    void setState(VendingMachineState* s) { currentState = s; }
    void selectItem(const string& item) { currentState->selectItem(item); }
    void insertMoney(double amt) { currentState->insertMoney(amt); }
    void dispense() { currentState->dispense(); }
    void cancel() { currentState->cancel(); }
    string getState() { return currentState->getName(); }

    void enterMaintenance(const string& pin) {
        if (pin == operatorPin) { setState(&maintenance); cout << "Entered maintenance mode." << endl; }
        else cout << "Invalid PIN." << endl;
    }
    void exitMaintenance(const string& pin) {
        if (pin == operatorPin && getState() == "Maintenance") {
            setState(&idle); cout << "Exited maintenance mode." << endl;
        } else cout << "Invalid PIN or not in maintenance." << endl;
    }
    void restock(const string& itemName, int qty) {
        if (getState() != "Maintenance") { cout << "Must be in maintenance to restock." << endl; return; }
        inventory[itemName].quantity += qty;
        cout << "Restocked " << itemName << " by " << qty << endl;
    }
};

void IdleState::selectItem(const string& item) {
    if (m->inventory.count(item) && m->inventory[item].quantity > 0) {
        m->selectedItem = item;
        m->setState(&m->paymentPending);
        cout << "Selected: " << item << ". Insert Rs." << m->inventory[item].price << endl;
    } else cout << "Item unavailable." << endl;
}
void PaymentPendingState::insertMoney(double amt) {
    m->insertedMoney += amt;
    double price = m->inventory[m->selectedItem].price;
    if (m->insertedMoney >= price) { m->setState(&m->dispensing); cout << "Payment accepted." << endl; }
    else cout << "Need Rs." << (price - m->insertedMoney) << " more." << endl;
}
void PaymentPendingState::cancel() {
    cout << "Refunding Rs." << m->insertedMoney << endl;
    m->insertedMoney = 0; m->selectedItem = ""; m->setState(&m->idle);
}
void DispensingState::dispense() {
    m->inventory[m->selectedItem].quantity--;
    cout << "Dispensed: " << m->selectedItem << endl;
    m->insertedMoney = 0; m->selectedItem = ""; m->setState(&m->idle);
}

int main() {
    cout << "Part 2: Maintenance mode — full scaffolding provided." << endl;
    return 0;
}
