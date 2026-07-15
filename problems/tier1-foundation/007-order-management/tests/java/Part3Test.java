// Order Management — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class TestObserver implements OrderObserver {
    List<tuple<String, OrderState, OrderState>> notifications;
    void onStateChange(String orderId,
                       OrderState from, OrderState to)  {
        notifications.add({orderId, from, to});
    }
}

class Part3Test {
    static boolean testHistoryFullLifecycle() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            confirm_order(id);
            ship_order(id);
            deliver_order(id);
            var hist = get_order_history(id);
            // Should have: creation entry + 3 transitions = 4 entries
            boolean pass = hist.size() == 4
                && hist[1].fromState == OrderState.Created
                && hist[1].toState == OrderState.Confirmed
                && hist[2].fromState == OrderState.Confirmed
                && hist[2].toState == OrderState.Shipped
                && hist[3].fromState == OrderState.Shipped
                && hist[3].toState == OrderState.Delivered;
            System.out.println((pass ? "PASS" : "FAIL") + ": testHistoryFullLifecycle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testHistoryFullLifecycle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testHistoryTimestampsOrdered() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            confirm_order(id);
            ship_order(id);
            var hist = get_order_history(id);
            for (int i = 1; i < hist.size(); i++) {
            }
            boolean pass = hist.get(i).timestamp >= hist[i-1].timestamp;
            System.out.println((pass ? "PASS" : "FAIL") + ": testHistoryTimestampsOrdered");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testHistoryTimestampsOrdered (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFailedTransitionNoHistory() {
        try {
            reset_manager();
            var id = create_order({{"PROD-1", 1}}, 100.0);
            ship_order(id); // invalid — should fail
            var hist = get_order_history(id);
            boolean pass = hist.size() == 1); // only the creation entry;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFailedTransitionNoHistory");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFailedTransitionNoHistory (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testObserverNotified() {
        try {
            reset_manager();
            TestObserver obs = new TestObserver();
            add_observer( obs);
            var id = create_order({{"PROD-1", 1}}, 100.0);
            confirm_order(id);
            ship_order(id);
            boolean pass = obs.notifications.size() == 2
                && get<1>(obs.notifications[0]) == OrderState.Created
                && get<2>(obs.notifications[0]) == OrderState.Confirmed
                && get<1>(obs.notifications[1]) == OrderState.Confirmed
                && get<2>(obs.notifications[1]) == OrderState.Shipped;
            System.out.println((pass ? "PASS" : "FAIL") + ": testObserverNotified");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testObserverNotified (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testObserverNotNotifiedOnFailure() {
        try {
            reset_manager();
            TestObserver obs = new TestObserver();
            add_observer( obs);
            var id = create_order({{"PROD-1", 1}}, 100.0);
            ship_order(id); // invalid
            boolean pass = obs.notifications.size() == 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testObserverNotNotifiedOnFailure");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testObserverNotNotifiedOnFailure (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelInHistory() {
        try {
            reset_manager();
            set_inventory("PROD-1", 10);
            var id = create_order({{"PROD-1", 1}}, 100.0);
            cancel_order(id);
            var hist = get_order_history(id);
            boolean pass = hist.size() == 2); // creation + cancellation
                && hist[1].fromState == OrderState.Created
                && hist[1].toState == OrderState.Cancelled;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelInHistory");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelInHistory (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testHistoryFullLifecycle()) passed++;
        total++; if (testHistoryTimestampsOrdered()) passed++;
        total++; if (testFailedTransitionNoHistory()) passed++;
        total++; if (testObserverNotified()) passed++;
        total++; if (testObserverNotNotifiedOnFailure()) passed++;
        total++; if (testCancelInHistory()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
