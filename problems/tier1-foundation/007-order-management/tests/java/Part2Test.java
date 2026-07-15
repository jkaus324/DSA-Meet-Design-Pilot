// Order Management — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testCancelFromCreated() {
        try {
            reset_manager();
            set_inventory("PROD-1", 10);
            var id = create_order({{"PROD-1", 3}}, 300.0);
            boolean pass = get_inventory("PROD-1") == 7); // decremented on create
                && cancel_order(id) == true
                && get_order_state(id) == OrderState.Cancelled
                && get_inventory("PROD-1") == 10); // restored on cancel;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelFromCreated");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelFromCreated (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelFromConfirmed() {
        try {
            reset_manager();
            set_inventory("PROD-1", 10);
            var id = create_order({{"PROD-1", 2}}, 200.0);
            confirm_order(id);
            boolean pass = cancel_order(id) == true
                && get_order_state(id) == OrderState.Cancelled
                && get_inventory("PROD-1") == 10;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelFromConfirmed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelFromConfirmed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelFromShippedFails() {
        try {
            reset_manager();
            set_inventory("PROD-1", 10);
            var id = create_order({{"PROD-1", 2}}, 200.0);
            confirm_order(id);
            ship_order(id);
            boolean pass = cancel_order(id) == false
                && get_order_state(id) == OrderState.Shipped
                && get_inventory("PROD-1") == 8); // not restored;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelFromShippedFails");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelFromShippedFails (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelFromDeliveredFails() {
        try {
            reset_manager();
            set_inventory("PROD-1", 10);
            var id = create_order({{"PROD-1", 1}}, 100.0);
            confirm_order(id);
            ship_order(id);
            deliver_order(id);
            boolean pass = cancel_order(id) == false
                && get_order_state(id) == OrderState.Delivered;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelFromDeliveredFails");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelFromDeliveredFails (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelMultiItemInventory() {
        try {
            reset_manager();
            set_inventory("PROD-A", 20);
            set_inventory("PROD-B", 15);
            var id = create_order({{"PROD-A", 5}, {"PROD-B", 3}}, 800.0);
            cancel_order(id);
            boolean pass = get_inventory("PROD-A") == 15
                && get_inventory("PROD-B") == 12
                && get_inventory("PROD-A") == 20
                && get_inventory("PROD-B") == 15;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelMultiItemInventory");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelMultiItemInventory (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDoubleCancelFails() {
        try {
            reset_manager();
            set_inventory("PROD-1", 10);
            var id = create_order({{"PROD-1", 2}}, 200.0);
            cancel_order(id);
            boolean pass = cancel_order(id) == false); // already cancelled;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDoubleCancelFails");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDoubleCancelFails (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testCancelFromCreated()) passed++;
        total++; if (testCancelFromConfirmed()) passed++;
        total++; if (testCancelFromShippedFails()) passed++;
        total++; if (testCancelFromDeliveredFails()) passed++;
        total++; if (testCancelMultiItemInventory()) passed++;
        total++; if (testDoubleCancelFails()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
