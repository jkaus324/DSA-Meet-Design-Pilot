// Vending Machine — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testInitialStateIdle() {
        try {
            boolean pass = getState() == "Idle";
            System.out.println((pass ? "PASS" : "FAIL") + ": testInitialStateIdle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInitialStateIdle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSelectTransitionsState() {
        try {
            reset(); // reset machine to clean state
            selectItem("Cola");
            String s = getState();
            boolean pass = s == "PaymentPending" || s == "ItemSelected";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSelectTransitionsState");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSelectTransitionsState (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFullPurchaseCycle() {
        try {
            reset();
            selectItem("Cola");
            insertMoney(25.0);  // assume Cola costs 25
            dispense();
            boolean pass = getState() == "Idle";
            System.out.println((pass ? "PASS" : "FAIL") + ": testFullPurchaseCycle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFullPurchaseCycle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelReturnsIdle() {
        try {
            reset();
            selectItem("Cola");
            insertMoney(10.0);
            cancel();
            boolean pass = getState() == "Idle";
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelReturnsIdle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelReturnsIdle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPayBeforeSelectSafe() {
        try {
            reset();
            insertMoney(50.0);  // should be ignored or print warning
            boolean pass = getState() == "Idle"); // should stay idle;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPayBeforeSelectSafe");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPayBeforeSelectSafe (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testInitialStateIdle()) passed++;
        total++; if (testSelectTransitionsState()) passed++;
        total++; if (testFullPurchaseCycle()) passed++;
        total++; if (testCancelReturnsIdle()) passed++;
        total++; if (testPayBeforeSelectSafe()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
