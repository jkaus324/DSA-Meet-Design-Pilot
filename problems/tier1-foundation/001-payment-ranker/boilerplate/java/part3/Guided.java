import java.util.*;

// Data class (given).
class PaymentMethod {
    public String name;
    public double cashbackRate;
    public double transactionFee;
    public int usageCount;
    public boolean easyRefundEligible;

    public PaymentMethod(String name, double cashbackRate, double transactionFee, int usageCount, boolean easyRefundEligible) {
        this.name = name;
        this.cashbackRate = cashbackRate;
        this.transactionFee = transactionFee;
        this.usageCount = usageCount;
        this.easyRefundEligible = easyRefundEligible;
    }

    public PaymentMethod(String name, double cashbackRate, double transactionFee, int usageCount) {
        this(name, cashbackRate, transactionFee, usageCount, false);
    }
}

// Marker interface so signatures compile; you supply the methods.
interface RankingStrategy {}

// HINT: introduce an abstraction so new ranking rules don't change existing code.
public class Solution {
    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static List<PaymentMethod> rank_by_rewards(List<PaymentMethod> methods) {
        // TODO: write your solution
        return methods;
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static List<PaymentMethod> rank_by_low_fee(List<PaymentMethod> methods) {
        // TODO: write your solution
        return methods;
    }

    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static List<PaymentMethod> rank_by_trust(List<PaymentMethod> methods) {
        // TODO: write your solution
        return methods;
    }

    // HINT: think about how to compose multiple criteria into a single decision.
    public static List<PaymentMethod> rank_composite(List<PaymentMethod> methods, List<RankingStrategy> criteria) {
        // TODO: write your solution
        return methods;
    }

    // HINT: a boolean flag changes ranking — handle it as a separate piece you can chain.
    public static List<PaymentMethod> rank_with_refund_filter(List<PaymentMethod> methods, boolean preferEasyRefund) {
        // TODO: write your solution
        return methods;
    }

}
