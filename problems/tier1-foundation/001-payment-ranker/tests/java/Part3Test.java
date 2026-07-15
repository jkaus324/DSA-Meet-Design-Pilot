// Payment Ranker — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testEasyRefundPreferred() {
        try {
            List<PaymentMethod> methods = {
            {"Card A", 0.10, 5.0, 300, false},  // high cashback, no easy refund
            {"Card B", 0.02, 2.0, 500, true},   // low cashback, easy refund
            {"Card C", 0.05, 3.0, 400, true},   // medium cashback, easy refund
            };
            var ranked = rank_with_refund_filter(methods, true);
            // Both B and C have easy refund, so they come first (in cashback order)
            boolean pass = ranked.size() == 3
                && ranked[0].easyRefundEligible == true
                && ranked[1].easyRefundEligible == true
                && ranked[2].name == "Card A"); // no easy refund, goes last;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEasyRefundPreferred");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEasyRefundPreferred (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRefundFilterDisabled() {
        try {
            List<PaymentMethod> methods = {
            {"Card A", 0.10, 5.0, 300, false},
            {"Card B", 0.02, 2.0, 500, true},
            };
            var ranked = rank_with_refund_filter(methods, false);
            // Without refund preference, Card A should still win (higher cashback)
            boolean pass = ranked[0].name == "Card A";
            System.out.println((pass ? "PASS" : "FAIL") + ": testRefundFilterDisabled");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRefundFilterDisabled (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAllRefundEligibleTiebreak() {
        try {
            List<PaymentMethod> methods = {
            {"Card A", 0.10, 5.0, 300, true},
            {"Card B", 0.02, 2.0, 500, true},
            };
            var ranked = rank_with_refund_filter(methods, true);
            // All have easy refund, so tiebreak by cashback
            boolean pass = ranked[0].name == "Card A"); // higher cashback;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAllRefundEligibleTiebreak");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAllRefundEligibleTiebreak (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testEasyRefundPreferred()) passed++;
        total++; if (testRefundFilterDisabled()) passed++;
        total++; if (testAllRefundEligibleTiebreak()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
