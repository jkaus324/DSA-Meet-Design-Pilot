// Order Management — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testCreateOrder() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 2}}, 500.0);
            boolean pass = !id.isEmpty()
                && get_order_state(id) == OrderState.Created;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCreateOrder");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCreateOrder (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFullLifecycle() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            boolean pass = confirm_order(id) == true
                && get_order_state(id) == OrderState.Confirmed
                && ship_order(id) == true
                && get_order_state(id) == OrderState.Shipped
                && deliver_order(id) == true
                && get_order_state(id) == OrderState.Delivered;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFullLifecycle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFullLifecycle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidSkipToShipped() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            boolean pass = ship_order(id) == false
                && get_order_state(id) == OrderState.Created); // unchanged;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidSkipToShipped");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidSkipToShipped (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidSkipToDelivered() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            boolean pass = deliver_order(id) == false
                && get_order_state(id) == OrderState.Created;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidSkipToDelivered");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidSkipToDelivered (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidBackwardTransition() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            confirm_order(id);
            ship_order(id);
            boolean pass = confirm_order(id) == false); // can't go back
                && get_order_state(id) == OrderState.Shipped;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidBackwardTransition");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidBackwardTransition (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleOrdersIndependent() {
        try {
            reset_manager();
            var id1 = create_order({{"PROD-1", 1}}, 100.0);
            var id2 = create_order({{"PROD-2", 1}}, 200.0);
            confirm_order(id1);
            boolean pass = get_order_state(id1) == OrderState.Confirmed
                && get_order_state(id2) == OrderState.Created;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleOrdersIndependent");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleOrdersIndependent (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNonexistentOrder() {
        try {
            reset_manager();
            boolean pass = confirm_order("NONEXISTENT") == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNonexistentOrder");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNonexistentOrder (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testCreateOrder()) passed++;
        total++; if (testFullLifecycle()) passed++;
        total++; if (testInvalidSkipToShipped()) passed++;
        total++; if (testInvalidSkipToDelivered()) passed++;
        total++; if (testInvalidBackwardTransition()) passed++;
        total++; if (testMultipleOrdersIndependent()) passed++;
        total++; if (testNonexistentOrder()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
