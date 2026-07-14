// Discount Engine — Solution (Java, Strategy)
import java.util.*;

class CartItem {
    public String name;
    public double price;
    public int quantity;
    public String category;

    public CartItem(String name, double price, int quantity) {
        this(name, price, quantity, "");
    }

    public CartItem(String name, double price, int quantity, String category) {
        this.name = name; this.price = price; this.quantity = quantity; this.category = category;
    }
}

interface Discount {
    double apply(List<CartItem> cart);
}

class PercentageDiscount implements Discount {
    private final double pct;
    public PercentageDiscount(double pct) { this.pct = pct; }
    public double apply(List<CartItem> cart) {
        double total = 0.0;
        for (CartItem i : cart) total += i.price * i.quantity;
        return total * (1.0 - pct / 100.0);
    }
}

class FlatDiscount implements Discount {
    private final double amount;
    public FlatDiscount(double amount) { this.amount = amount; }
    public double apply(List<CartItem> cart) {
        double total = 0.0;
        for (CartItem i : cart) total += i.price * i.quantity;
        return Math.max(0.0, total - amount);
    }
}

class BuyXGetYDiscount implements Discount {
    private final int buy, free;
    public BuyXGetYDiscount(int buy, int free) { this.buy = buy; this.free = free; }
    public double apply(List<CartItem> cart) {
        int group = buy + free;
        double total = 0.0;
        for (CartItem it : cart) {
            int groups = it.quantity / group;
            int remainder = it.quantity % group;
            int paid = groups * buy + Math.min(remainder, buy);
            total += paid * it.price;
        }
        return total;
    }
}

public class Solution {
    public static void reset_service() {}

    public static double apply_percentage_discount(List<CartItem> cart, double percentage) {
        return new PercentageDiscount(percentage).apply(cart);
    }

    public static double apply_flat_discount(List<CartItem> cart, double amount) {
        return new FlatDiscount(amount).apply(cart);
    }

    public static double apply_bogo(List<CartItem> cart, int buyCount, int freeCount) {
        return new BuyXGetYDiscount(buyCount, freeCount).apply(cart);
    }

    public static double apply_percentage_with_eligibility(List<CartItem> cart, double percentage,
                                                           double minCartValue, boolean requireFirstTimeUser,
                                                           boolean isFirstTimeUser, String eligibleCategory) {
        double raw = 0.0;
        for (CartItem i : cart) raw += i.price * i.quantity;
        if (raw < minCartValue) return raw;
        if (requireFirstTimeUser && !isFirstTimeUser) return raw;
        if (eligibleCategory != null && !eligibleCategory.isEmpty()) {
            List<CartItem> eligible = new ArrayList<>();
            double nonEligible = 0.0;
            for (CartItem i : cart) {
                if (i.category != null && i.category.equals(eligibleCategory)) eligible.add(i);
                else nonEligible += i.price * i.quantity;
            }
            return new PercentageDiscount(percentage).apply(eligible) + nonEligible;
        }
        return new PercentageDiscount(percentage).apply(cart);
    }
}
