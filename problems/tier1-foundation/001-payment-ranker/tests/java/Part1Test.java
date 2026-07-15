// Payment Ranker — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testRewardsRanking() {
        try {
            List<PaymentMethod> methods = {
            {"UPI",           0.01, 0.0, 1000},
            {"Credit Card A", 0.05, 5.0, 500},
            {"Credit Card B", 0.10, 8.0, 300},
            };
            var ranked = rank_by_rewards(methods);
            boolean pass = ranked.size() == 3
                && ranked[0].name == "Credit Card B"); // 10% cashback
                && ranked[1].name == "Credit Card A"); // 5% cashback
                && ranked[2].name == "UPI");           // 1% cashback;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRewardsRanking");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRewardsRanking (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLowFeeRanking() {
        try {
            List<PaymentMethod> methods = {
            {"Debit Card",    0.0,  2.0, 800},
            {"Credit Card A", 0.05, 5.0, 500},
            {"Credit Card B", 0.10, 8.0, 300},
            {"UPI",           0.01, 0.0, 1000},
            };
            var ranked = rank_by_low_fee(methods);
            boolean pass = ranked.size() == 4
                && ranked[0].name == "UPI");         // 0 fee
                && ranked[1].name == "Debit Card");  // 2.0 fee
                && ranked[2].name == "Credit Card A"); // 5.0 fee
                && ranked[3].name == "Credit Card B"); // 8.0 fee;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLowFeeRanking");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLowFeeRanking (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTrustRanking() {
        try {
            List<PaymentMethod> methods = {
            {"Credit Card A", 0.05, 5.0, 500},
            {"UPI",           0.01, 0.0, 1000},
            {"Debit Card",    0.0,  2.0, 800},
            };
            var ranked = rank_by_trust(methods);
            boolean pass = ranked.size() == 3
                && ranked[0].name == "UPI");          // 1000 uses
                && ranked[1].name == "Debit Card");   // 800 uses
                && ranked[2].name == "Credit Card A"); // 500 uses;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTrustRanking");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTrustRanking (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyInput() {
        try {
            List<PaymentMethod> empty;
            boolean pass = rank_by_rewards(empty).isEmpty()
                && rank_by_low_fee(empty).isEmpty()
                && rank_by_trust(empty).isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyInput");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyInput (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSingleItem() {
        try {
            List<PaymentMethod> single = {{"UPI", 0.01, 0.0, 100}};
            boolean pass = rank_by_rewards(single).size() == 1
                && rank_by_rewards(single)[0].name == "UPI";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingleItem");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingleItem (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testRewardsRanking()) passed++;
        total++; if (testLowFeeRanking()) passed++;
        total++; if (testTrustRanking()) passed++;
        total++; if (testEmptyInput()) passed++;
        total++; if (testSingleItem()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
