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

class VendingMachine;

// ─── State Interface ─────────────────────────────────────────────────────────

class VMState {
public:
    virtual void selectItem(VendingMachine& vm, const string& item) = 0;
    virtual void insertMoney(VendingMachine& vm, double amount) = 0;
    virtual void dispense(VendingMachine& vm) = 0;
    virtual void cancel(VendingMachine& vm) = 0;
    virtual string name() const = 0;
    virtual ~VMState() = default;
};

// ─── Vending Machine Context ─────────────────────────────────────────────────

class VendingMachine {
public:
    VMState* state;
    string   selectedItem;
    double   insertedAmount = 0;
    unordered_map<string, Item> inventory;

    VendingMachine() {}

    void setState(VMState* s) { state = s; }
    void selectItem(const string& item)  { state->selectItem(*this, item); }
    void insertMoney(double amount)      { state->insertMoney(*this, amount); }
    void dispense()                      { state->dispense(*this); }
    void cancel()                        { state->cancel(*this); }
    string getState() const              { return state->name(); }
};

// ─── Concrete States ─────────────────────────────────────────────────────────

class IdleState : public VMState {
public:
    string name() const override { return "Idle"; }
    void selectItem(VendingMachine& vm, const string& item) override {
        // TODO: Check if item exists in inventory and has quantity > 0
        //       If yes: set vm.selectedItem, transition to SelectedState
        //       If no:  print "Item not available"
    }
    void insertMoney(VendingMachine& vm, double amount) override {
        cout << "[Error] Select an item first." << endl;
    }
    void dispense(VendingMachine& vm) override {
        cout << "[Error] No item selected." << endl;
    }
    void cancel(VendingMachine& vm) override {
        cout << "[Info] Nothing to cancel." << endl;
    }
};

class SelectedState : public VMState {
public:
    string name() const override { return "ItemSelected"; }
    void selectItem(VendingMachine& vm, const string& item) override {
        cout << "[Info] Item already selected. Cancel first." << endl;
    }
    void insertMoney(VendingMachine& vm, double amount) override {
        // TODO: Add amount to vm.insertedAmount
        //       If insertedAmount >= item price: transition to PaidState
        //       Else: print how much more is needed
    }
    void dispense(VendingMachine& vm) override {
        cout << "[Error] Insert payment first." << endl;
    }
    void cancel(VendingMachine& vm) override {
        // TODO: Reset selectedItem and insertedAmount, go back to IdleState
    }
};

class PaidState : public VMState {
public:
    string name() const override { return "PaymentReceived"; }
    void selectItem(VendingMachine& vm, const string& item) override {
        cout << "[Error] Payment already made. Dispense or cancel." << endl;
    }
    void insertMoney(VendingMachine& vm, double amount) override {
        cout << "[Error] Payment already received." << endl;
    }
    void dispense(VendingMachine& vm) override {
        // TODO: Dispense the item (decrement quantity, print confirmation)
        //       Return change if overpaid
        //       Transition to IdleState
    }
    void cancel(VendingMachine& vm) override {
        // TODO: Refund insertedAmount, reset state, go to IdleState
    }
};

int main() {
    cout << "Vending Machine — implement the TODO methods in each state." << endl;
    return 0;
}
