// Payment Ranker — Solution (Java)
// All classes in default package; free functions live as static methods on `Solution`.

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;

// ─── Data Structure ──────────────────────────────────────────────────────────

class PaymentMethod {
    public String name;
    public double cashbackRate;
    public double transactionFee;
    public int usageCount;
    public boolean easyRefundEligible;

    // 4-arg ctor (Parts 1 & 2): defaults easyRefundEligible to false.
    public PaymentMethod(String name, double cashbackRate, double transactionFee, int usageCount) {
        this(name, cashbackRate, transactionFee, usageCount, false);
    }

    // 5-arg ctor (Part 3).
    public PaymentMethod(String name, double cashbackRate, double transactionFee,
                         int usageCount, boolean easyRefundEligible) {
        this.name = name;
        this.cashbackRate = cashbackRate;
        this.transactionFee = transactionFee;
        this.usageCount = usageCount;
        this.easyRefundEligible = easyRefundEligible;
    }
}

// ─── Strategy Interface ──────────────────────────────────────────────────────

interface RankingStrategy {
    /** Returns true if `a` should come before `b`. */
    boolean compare(PaymentMethod a, PaymentMethod b);
}

// ─── Part 1: Concrete Strategies ─────────────────────────────────────────────

class RewardsMaximizer implements RankingStrategy {
    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        return a.cashbackRate > b.cashbackRate;
    }
}

class LowFeeSeeker implements RankingStrategy {
    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        return a.transactionFee < b.transactionFee;
    }
}

class TrustBasedRanker implements RankingStrategy {
    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        return a.usageCount > b.usageCount;
    }
}

// ─── Part 2: Composite Strategy ──────────────────────────────────────────────

class CompositeStrategy implements RankingStrategy {
    private final List<RankingStrategy> criteria;

    public CompositeStrategy(List<RankingStrategy> criteria) {
        this.criteria = new ArrayList<>(criteria);
    }

    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        for (RankingStrategy s : criteria) {
            if (s.compare(a, b)) return true;
            if (s.compare(b, a)) return false;
        }
        return false;
    }
}

// ─── Part 3: Easy-Refund Strategy ────────────────────────────────────────────

class EasyRefundStrategy implements RankingStrategy {
    private final boolean prefer;

    public EasyRefundStrategy(boolean preferRefund) {
        this.prefer = preferRefund;
    }

    @Override
    public boolean compare(PaymentMethod a, PaymentMethod b) {
        if (!prefer) return false;
        return a.easyRefundEligible && !b.easyRefundEligible;
    }
}

// ─── PaymentRanker (context class) ───────────────────────────────────────────

class PaymentRanker {
    private RankingStrategy strategy;

    public PaymentRanker(RankingStrategy s) {
        this.strategy = s;
    }

    public void setStrategy(RankingStrategy s) {
        this.strategy = s;
    }

    public List<PaymentMethod> rank(List<PaymentMethod> methods) {
        List<PaymentMethod> out = new ArrayList<>(methods);
        // Convert "a before b" relation into a stable, transitive Comparator.
        out.sort((a, b) -> {
            if (strategy.compare(a, b)) return -1;
            if (strategy.compare(b, a)) return 1;
            return 0;
        });
        return out;
    }
}

// ─── Free-function wrappers (used by tests) ──────────────────────────────────

public class Solution {
    public static List<PaymentMethod> rank_by_rewards(List<PaymentMethod> methods) {
        return new PaymentRanker(new RewardsMaximizer()).rank(methods);
    }

    public static List<PaymentMethod> rank_by_low_fee(List<PaymentMethod> methods) {
        return new PaymentRanker(new LowFeeSeeker()).rank(methods);
    }

    public static List<PaymentMethod> rank_by_trust(List<PaymentMethod> methods) {
        return new PaymentRanker(new TrustBasedRanker()).rank(methods);
    }

    public static List<PaymentMethod> rank_composite(
            List<PaymentMethod> methods, List<RankingStrategy> criteria) {
        return new PaymentRanker(new CompositeStrategy(criteria)).rank(methods);
    }

    public static List<PaymentMethod> rank_with_refund_filter(
            List<PaymentMethod> methods, boolean preferEasyRefund) {
        List<RankingStrategy> chain = new ArrayList<>();
        chain.add(new EasyRefundStrategy(preferEasyRefund));
        chain.add(new RewardsMaximizer());
        return new PaymentRanker(new CompositeStrategy(chain)).rank(methods);
    }
}
