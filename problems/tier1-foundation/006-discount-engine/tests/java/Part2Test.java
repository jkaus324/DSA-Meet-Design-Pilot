// Discount Engine — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testStackPercentageThenFlat() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Laptop", 10000.0, 1, "electronics"),
            };
            // Total = 10000 → 10% off → 9000 → Rs.500 flat off → 8500
            PercentageDiscount pct = new PercentageDiscount(10.0);
            FlatDiscount flat = new FlatDiscount(500.0);
            double result = apply_stacked_discounts(cart, { pct,  flat});
            boolean pass = approx2(result, 8500.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testStackPercentageThenFlat");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStackPercentageThenFlat (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStackFlatThenPercentage() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Laptop", 10000.0, 1, "electronics"),
            };
            // Total = 10000 → Rs.500 flat off → 9500 → 10% off → 8550
            FlatDiscount flat = new FlatDiscount(500.0);
            PercentageDiscount pct = new PercentageDiscount(10.0);
            double result = apply_stacked_discounts(cart, { flat,  pct});
            boolean pass = approx2(result, 8550.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testStackFlatThenPercentage");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStackFlatThenPercentage (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStackThreeDiscounts() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Phone", 20000.0, 1, "electronics"),
            };
            // Total = 20000 → 10% coupon → 18000 → Rs.1000 seasonal → 17000 → 5% membership → 16150
            PercentageDiscount coupon = new PercentageDiscount(10.0);
            FlatDiscount seasonal = new FlatDiscount(1000.0);
            PercentageDiscount membership = new PercentageDiscount(5.0);
            double result = apply_stacked_discounts(cart, { coupon,  seasonal,  membership});
            boolean pass = approx2(result, 16150.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testStackThreeDiscounts");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStackThreeDiscounts (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSingleDiscountStack() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Book", 500.0, 2, "books"),
            };
            PercentageDiscount pct = new PercentageDiscount(20.0);
            double stacked = apply_stacked_discounts(cart, { pct});
            double direct = apply_percentage_discount(cart, 20.0);
            boolean pass = approx2(stacked, direct);
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingleDiscountStack");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingleDiscountStack (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStackReducesToZero() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Sticker", 100.0, 1, "accessories"),
            };
            FlatDiscount flat1 = new FlatDiscount(60.0);
            FlatDiscount flat2 = new FlatDiscount(60.0);
            double result = apply_stacked_discounts(cart, { flat1,  flat2});
            boolean pass = approx2(result, 0.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testStackReducesToZero");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStackReducesToZero (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testStackPercentageThenFlat()) passed++;
        total++; if (testStackFlatThenPercentage()) passed++;
        total++; if (testStackThreeDiscounts()) passed++;
        total++; if (testSingleDiscountStack()) passed++;
        total++; if (testStackReducesToZero()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
