import java.util.*;

// Data class (given — do not modify).
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

interface RankingStrategy {
    boolean compare(PaymentMethod a, PaymentMethod b);
}

class RewardsMaximizer implements RankingStrategy {
    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        // TODO: implement this
        return false;
    }
}

class LowFeeSeeker implements RankingStrategy {
    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        // TODO: implement this
        return false;
    }
}

class TrustBasedRanker implements RankingStrategy {
    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        // TODO: implement this
        return false;
    }
}

public class Solution {
    public static List<PaymentMethod> rank_by_rewards(List<PaymentMethod> methods) {
        // TODO: implement this
        return methods;
    }

    public static List<PaymentMethod> rank_by_low_fee(List<PaymentMethod> methods) {
        // TODO: implement this
        return methods;
    }

    public static List<PaymentMethod> rank_by_trust(List<PaymentMethod> methods) {
        // TODO: implement this
        return methods;
    }

}
