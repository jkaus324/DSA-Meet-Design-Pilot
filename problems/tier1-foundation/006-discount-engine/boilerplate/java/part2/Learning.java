import java.util.*;

// Data class (given — do not modify).
class CartItem {
    public String name;
    public double price;
    public int quantity;
    public String category;

    public CartItem(String name, double price, int quantity, String category) {
        this.name = name;
        this.price = price;
        this.quantity = quantity;
        this.category = category;
    }

    public CartItem(String name, double price, int quantity) {
        this(name, price, quantity, "");
    }
}

public class Solution {
    public static double apply_percentage_discount(List<CartItem> cart, double percentage) {
        // TODO: implement this
        return 0.0;
    }

    public static double apply_flat_discount(List<CartItem> cart, double amount) {
        // TODO: implement this
        return 0.0;
    }

    public static double apply_bogo(List<CartItem> cart, int buyCount, int freeCount) {
        // TODO: implement this
        return 0.0;
    }

}
