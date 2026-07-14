// Discount Engine — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testMinCartValueMet() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Laptop", 50000.0, 1, "electronics"),
            };
            PercentageDiscount pct = new PercentageDiscount(10.0);
            UserContext user{false};
            double result = apply_with_eligibility(cart,  pct, 1000.0, false, user, "");
            // Total 50000 >= 1000, discount applies: 45000
            boolean pass = approx3(result, 45000.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testMinCartValueMet");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMinCartValueMet (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMinCartValueNotMet() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Sticker", 50.0, 1, "accessories"),
            };
            PercentageDiscount pct = new PercentageDiscount(10.0);
            UserContext user{false};
            double result = apply_with_eligibility(cart,  pct, 1000.0, false, user, "");
            // Total 50 < 1000, discount skipped: 50
            boolean pass = approx3(result, 50.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testMinCartValueNotMet");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMinCartValueNotMet (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFirstTimeUserEligible() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Phone", 20000.0, 1, "electronics"),
            };
            FlatDiscount flat = new FlatDiscount(2000.0);
            UserContext user{true};
            double result = apply_with_eligibility(cart,  flat, 0.0, true, user, "");
            // First-time user, discount applies: 20000 - 2000 = 18000
            boolean pass = approx3(result, 18000.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testFirstTimeUserEligible");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFirstTimeUserEligible (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFirstTimeUserNotEligible() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Phone", 20000.0, 1, "electronics"),
            };
            FlatDiscount flat = new FlatDiscount(2000.0);
            UserContext user{false};
            double result = apply_with_eligibility(cart,  flat, 0.0, true, user, "");
            // Not first-time, discount skipped: 20000
            boolean pass = approx3(result, 20000.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testFirstTimeUserNotEligible");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFirstTimeUserNotEligible (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCategorySpecificDiscount() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Laptop",     50000.0, 1, "electronics"),
            Arrays.asList("Phone Case",   500.0, 2, "accessories"),
            };
            PercentageDiscount pct = new PercentageDiscount(10.0);
            UserContext user{false};
            double result = apply_with_eligibility(cart,  pct, 0.0, false, user, "electronics");
            // Electronics: 50000 * 0.9 = 45000, Accessories: 1000 full price → 46000
            boolean pass = approx3(result, 46000.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testCategorySpecificDiscount");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCategorySpecificDiscount (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAllRulesCombined() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Laptop",     50000.0, 1, "electronics"),
            Arrays.asList("T-Shirt",     1000.0, 3, "clothing"),
            };
            PercentageDiscount pct = new PercentageDiscount(20.0);
            UserContext user{true};
            double result = apply_with_eligibility(cart,  pct, 5000.0, true, user, "electronics");
            // Total = 53000 >= 5000, first-time user OK
            // Electronics: 50000 * 0.8 = 40000, Clothing: 3000 full price → 43000
            boolean pass = approx3(result, 43000.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testAllRulesCombined");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAllRulesCombined (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testMinCartValueMet()) passed++;
        total++; if (testMinCartValueNotMet()) passed++;
        total++; if (testFirstTimeUserEligible()) passed++;
        total++; if (testFirstTimeUserNotEligible()) passed++;
        total++; if (testCategorySpecificDiscount()) passed++;
        total++; if (testAllRulesCombined()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
