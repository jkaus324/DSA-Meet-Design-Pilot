// Vending Machine — Solution (Java, State pattern)
import java.util.*;

class VMItem {
    public String name;
    public double price;
    public int quantity;
    public VMItem(String name, double price, int quantity) {
        this.name = name; this.price = price; this.quantity = quantity;
    }
}

class VendingMachine {
    public static final String IDLE = "Idle";
    public static final String PAYMENT_PENDING = "PaymentPending";
    public static final String DISPENSING = "Dispensing";
    public static final String MAINTENANCE = "Maintenance";

    public String state = IDLE;
    public String selectedItem = "";
    public double insertedMoney = 0.0;
    public Map<String, VMItem> inventory = new LinkedHashMap<>();
    public String operatorPin = "1234";

    public VendingMachine() {
        stockDefaults();
    }

    private void stockDefaults() {
        inventory.put("Cola",  new VMItem("Cola",  25.0, 5));
        inventory.put("Chips", new VMItem("Chips", 15.0, 3));
    }

    public void reset() {
        state = IDLE;
        selectedItem = "";
        insertedMoney = 0.0;
        inventory.clear();
        stockDefaults();
    }

    public void selectItem(String item) {
        if (!state.equals(IDLE)) return;
        VMItem inv = inventory.get(item);
        if (inv != null && inv.quantity > 0) {
            selectedItem = item;
            state = PAYMENT_PENDING;
        }
    }

    public void insertMoney(double amount) {
        if (!state.equals(PAYMENT_PENDING)) return;
        insertedMoney += amount;
        if (insertedMoney >= inventory.get(selectedItem).price) {
            state = DISPENSING;
        }
    }

    public void dispense() {
        if (!state.equals(DISPENSING)) return;
        inventory.get(selectedItem).quantity -= 1;
        insertedMoney = 0.0;
        selectedItem = "";
        state = IDLE;
    }

    public void cancel() {
        if (!state.equals(PAYMENT_PENDING)) return;
        insertedMoney = 0.0;
        selectedItem = "";
        state = IDLE;
    }

    public void enterMaintenance(String pin) {
        if (pin.equals(operatorPin) && state.equals(IDLE)) {
            state = MAINTENANCE;
        }
    }

    public void exitMaintenance(String pin) {
        if (pin.equals(operatorPin) && state.equals(MAINTENANCE)) {
            state = IDLE;
        }
    }

    public void restock(String item, int qty) {
        if (!state.equals(MAINTENANCE)) return;
        VMItem inv = inventory.get(item);
        if (inv == null) {
            inv = new VMItem(item, 0.0, 0);
            inventory.put(item, inv);
        }
        inv.quantity += qty;
    }
}

public class Solution {
    private static VendingMachine vm = new VendingMachine();

    public static void reset_service() { vm = new VendingMachine(); }
    public static void reset() { vm.reset(); }
    public static String getState() { return vm.state; }
    public static void selectItem(String item) { vm.selectItem(item); }
    public static void insertMoney(double amount) { vm.insertMoney(amount); }
    public static void dispense() { vm.dispense(); }
    public static void cancel() { vm.cancel(); }
    public static void enterMaintenance(String pin) { vm.enterMaintenance(pin); }
    public static void exitMaintenance(String pin) { vm.exitMaintenance(pin); }
    public static void restock(String item, int qty) { vm.restock(item, qty); }

    public static int vm_get_quantity(String item) {
        VMItem inv = vm.inventory.get(item);
        return inv == null ? -1 : inv.quantity;
    }
}
