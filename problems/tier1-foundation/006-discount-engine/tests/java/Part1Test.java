// Discount Engine — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testPercentageDiscount() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Laptop",     50000.0, 1, "electronics"),
            Arrays.asList("Phone Case",   500.0, 2, "accessories"),
            };
            // Total = 50000 + 1000 = 51000, 10% off = 45900
            double result = apply_percentage_discount(cart, 10.0);
            boolean pass = approx(result, 45900.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testPercentageDiscount");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPercentageDiscount (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFlatDiscount() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Headphones", 2000.0, 1, "electronics"),
            Arrays.asList("Cable",       300.0, 3, "accessories"),
            };
            // Total = 2000 + 900 = 2900, flat 200 off = 2700
            double result = apply_flat_discount(cart, 200.0);
            boolean pass = approx(result, 2700.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testFlatDiscount");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFlatDiscount (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFlatDiscountExceedsTotal() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Sticker", 50.0, 1, "accessories"),
            };
            double result = apply_flat_discount(cart, 200.0);
            boolean pass = approx(result, 0.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testFlatDiscountExceedsTotal");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFlatDiscountExceedsTotal (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBogoExactGroups() {
        try {
            List<CartItem> cart = {
            Arrays.asList("T-Shirt", 500.0, 6, "clothing"),  // 6 shirts: pay for 4
            };
            // 6 / (2+1) = 2 groups, each group pays for 2 → 4 paid items
            double result = apply_bogo(cart, 2, 1);
            boolean pass = approx(result, 2000.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testBogoExactGroups");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBogoExactGroups (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBogoWithRemainder() {
        try {
            List<CartItem> cart = {
            Arrays.asList("Socks", 200.0, 5, "clothing"),  // 5 socks: 1 group of 3 (pay 2) + 2 remainder (pay 2)
            };
            // 5 / 3 = 1 group (pay 2), remainder 2 (pay 2) → 4 paid items
            double result = apply_bogo(cart, 2, 1);
            boolean pass = approx(result, 800.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testBogoWithRemainder");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBogoWithRemainder (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyCart() {
        try {
            List<CartItem> empty;
            boolean pass = approx(apply_percentage_discount(empty, 10.0), 0.0)
                && approx(apply_flat_discount(empty, 100.0), 0.0)
                && approx(apply_bogo(empty, 2, 1), 0.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyCart");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyCart (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testZeroPercentage() {
        try {
            List<CartItem> cart = {Arrays.asList("Book", 300.0, 1, "books")};
            double result = apply_percentage_discount(cart, 0.0);
            boolean pass = approx(result, 300.0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testZeroPercentage");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testZeroPercentage (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testPercentageDiscount()) passed++;
        total++; if (testFlatDiscount()) passed++;
        total++; if (testFlatDiscountExceedsTotal()) passed++;
        total++; if (testBogoExactGroups()) passed++;
        total++; if (testBogoWithRemainder()) passed++;
        total++; if (testEmptyCart()) passed++;
        total++; if (testZeroPercentage()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
