// Payment Ranker — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testCompositeCashbackThenFee() {
        try {
            List<PaymentMethod> methods = {
            {"Card A", 0.10, 8.0, 300},  // 10% cashback, high fee
            {"Card B", 0.10, 3.0, 400},  // 10% cashback, low fee
            {"Card C", 0.05, 1.0, 200},  // 5% cashback
            };
            RewardsMaximizer rewards = new RewardsMaximizer();
            LowFeeSeeker fee = new LowFeeSeeker();
            var ranked = rank_composite(methods, { rewards,  fee});
            boolean pass = ranked.size() == 3
                && ranked[0].name == "Card B"); // tied cashback, lower fee wins
                && ranked[1].name == "Card A"); // tied cashback, higher fee loses
                && ranked[2].name == "Card C"); // lower cashback;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCompositeCashbackThenFee");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCompositeCashbackThenFee (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCompositeTrustThenCashback() {
        try {
            List<PaymentMethod> methods = {
            {"UPI",    0.01, 0.0, 1000},
            {"Card A", 0.10, 5.0, 200},
            {"Card B", 0.05, 3.0, 1000},  // tied trust with UPI
            };
            TrustBasedRanker trust = new TrustBasedRanker();
            RewardsMaximizer rewards = new RewardsMaximizer();
            var ranked = rank_composite(methods, { trust,  rewards});
            // UPI and Card B both have 1000 uses — tiebreak by cashback
            // Card B has 5% cashback > UPI's 1%
            boolean pass = ranked.size() == 3
                && ranked[0].name == "Card B"
                && ranked[1].name == "UPI"
                && ranked[2].name == "Card A";
            System.out.println((pass ? "PASS" : "FAIL") + ": testCompositeTrustThenCashback");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCompositeTrustThenCashback (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSingleCriterionComposite() {
        try {
            List<PaymentMethod> methods = {
            {"Card X", 0.02, 5.0, 100},
            {"Card Y", 0.08, 3.0, 200},
            };
            RewardsMaximizer rewards = new RewardsMaximizer();
            var composite = rank_composite(methods, { rewards});
            var direct    = rank_by_rewards(methods);
            boolean pass = composite[0].name == direct[0].name
                && composite[1].name == direct[1].name;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingleCriterionComposite");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingleCriterionComposite (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testCompositeCashbackThenFee()) passed++;
        total++; if (testCompositeTrustThenCashback()) passed++;
        total++; if (testSingleCriterionComposite()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
