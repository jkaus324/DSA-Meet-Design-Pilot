// Order Management — Solution (Java, State + Observer)
import java.util.*;

class OrderItem {
    public String productId;
    public int quantity;
    public OrderItem(String productId, int quantity) {
        this.productId = productId; this.quantity = quantity;
    }
}

class Order {
    public String id;
    public List<OrderItem> items;
    public double totalAmount;
    public String state;
    public Order(String id, List<OrderItem> items, double totalAmount) {
        this.id = id; this.items = items; this.totalAmount = totalAmount;
        this.state = "Created";
    }
}

interface OrderObserver {
    void onStateChange(String orderId, String fromState, String toState);
}

class OrderManager {
    public Map<String, Order> orders = new LinkedHashMap<>();
    public Map<String, Integer> inventory = new LinkedHashMap<>();
    public Map<String, List<String[]>> history = new LinkedHashMap<>();
    public List<OrderObserver> observers = new ArrayList<>();
    public int nextId = 1;

    private boolean transition(String orderId, String expected, String next) {
        Order o = orders.get(orderId);
        if (o == null) return false;
        if (!o.state.equals(expected)) return false;
        String old = o.state;
        o.state = next;
        history.computeIfAbsent(orderId, k -> new ArrayList<>()).add(new String[]{old, next});
        for (OrderObserver obs : observers) obs.onStateChange(orderId, old, next);
        return true;
    }

    public void setInventory(String productId, int qty) { inventory.put(productId, qty); }

    public int getInventory(String productId) { return inventory.getOrDefault(productId, 0); }

    public String createOrder(List<OrderItem> items, double totalAmount) {
        String oid = "ORD-" + nextId;
        nextId++;
        for (OrderItem it : items) {
            inventory.put(it.productId, inventory.getOrDefault(it.productId, 0) - it.quantity);
        }
        orders.put(oid, new Order(oid, new ArrayList<>(items), totalAmount));
        List<String[]> h = new ArrayList<>();
        h.add(new String[]{"Created", "Created"});
        history.put(oid, h);
        return oid;
    }

    public boolean confirmOrder(String oid) { return transition(oid, "Created", "Confirmed"); }
    public boolean shipOrder(String oid) { return transition(oid, "Confirmed", "Shipped"); }
    public boolean deliverOrder(String oid) { return transition(oid, "Shipped", "Delivered"); }

    public boolean cancelOrder(String oid) {
        Order o = orders.get(oid);
        if (o == null) return false;
        if (!o.state.equals("Created") && !o.state.equals("Confirmed")) return false;
        for (OrderItem it : o.items) {
            inventory.put(it.productId, inventory.getOrDefault(it.productId, 0) + it.quantity);
        }
        String old = o.state;
        o.state = "Cancelled";
        history.computeIfAbsent(oid, k -> new ArrayList<>()).add(new String[]{old, "Cancelled"});
        for (OrderObserver obs : observers) obs.onStateChange(oid, old, "Cancelled");
        return true;
    }

    public String getOrderState(String oid) {
        Order o = orders.get(oid); return o == null ? "" : o.state;
    }

    public List<String[]> getOrderHistory(String oid) {
        return history.getOrDefault(oid, new ArrayList<>());
    }

    public void addObserver(OrderObserver obs) { observers.add(obs); }
}

class CountingObs implements OrderObserver {
    public int count = 0;
    public String lastOrderId = "", lastFrom = "", lastTo = "";
    public void onStateChange(String oid, String fromState, String toState) {
        count++;
        lastOrderId = oid; lastFrom = fromState; lastTo = toState;
    }
}

public class Solution {
    private static OrderManager mgr = new OrderManager();
    private static CountingObs obs = null;

    public static void reset_service() {
        mgr = new OrderManager();
        obs = null;
    }

    public static void set_inventory(String productId, int qty) {
        mgr.setInventory(productId, qty);
    }

    public static int get_inventory(String productId) {
        return mgr.getInventory(productId);
    }

    public static String create_order_simple(String productId, int quantity, double totalAmount) {
        List<OrderItem> items = new ArrayList<>();
        items.add(new OrderItem(productId, quantity));
        return mgr.createOrder(items, totalAmount);
    }

    public static boolean confirm_order(String oid) { return mgr.confirmOrder(oid); }
    public static boolean ship_order(String oid) { return mgr.shipOrder(oid); }
    public static boolean deliver_order(String oid) { return mgr.deliverOrder(oid); }
    public static boolean cancel_order(String oid) { return mgr.cancelOrder(oid); }
    public static String get_order_state_str(String oid) { return mgr.getOrderState(oid); }
    public static int get_history_size(String oid) { return mgr.getOrderHistory(oid).size(); }

    public static void om_attach_observer() {
        obs = new CountingObs();
        mgr.addObserver(obs);
    }

    public static int om_observer_count() {
        return obs == null ? 0 : obs.count;
    }

    public static String om_observer_last_to() {
        return obs == null ? "" : obs.lastTo;
    }
}
