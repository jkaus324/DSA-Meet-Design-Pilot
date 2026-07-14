// Payment ranker — Strategy pattern reference solution (JavaScript).

class PaymentMethod {
  constructor(name, cashbackRate, transactionFee, usageCount, easyRefundEligible = false) {
    this.name = name;
    this.cashbackRate = cashbackRate;
    this.transactionFee = transactionFee;
    this.usageCount = usageCount;
    this.easyRefundEligible = easyRefundEligible;
  }
}

class RankingStrategy {
  // Return true iff `a` should rank strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

class RewardsMaximizer extends RankingStrategy {
  compare(a, b) { return a.cashbackRate > b.cashbackRate; }
}

class LowFeeSeeker extends RankingStrategy {
  compare(a, b) { return a.transactionFee < b.transactionFee; }
}

class TrustBasedRanker extends RankingStrategy {
  compare(a, b) { return a.usageCount > b.usageCount; }
}

class CompositeStrategy extends RankingStrategy {
  constructor(criteria) { super(); this.criteria = criteria; }
  compare(a, b) {
    for (const s of this.criteria) {
      if (s.compare(a, b)) return true;
      if (s.compare(b, a)) return false;
    }
    return false;
  }
}

class EasyRefundStrategy extends RankingStrategy {
  constructor(prefer) { super(); this.prefer = prefer; }
  compare(a, b) {
    if (!this.prefer) return false;
    return a.easyRefundEligible && !b.easyRefundEligible;
  }
}

class PaymentRanker {
  constructor(strategy) { this.strategy = strategy; }
  setStrategy(strategy) { this.strategy = strategy; }
  rank(methods) {
    // Translate strict-less comparator into a stable sort via insertion.
    const result = [];
    for (const m of methods) {
      let inserted = false;
      for (let i = 0; i < result.length; i++) {
        if (this.strategy.compare(m, result[i])) {
          result.splice(i, 0, m);
          inserted = true;
          break;
        }
      }
      if (!inserted) result.push(m);
    }
    return result;
  }
}

function rank_by_rewards(methods) {
  return new PaymentRanker(new RewardsMaximizer()).rank(methods);
}

function rank_by_low_fee(methods) {
  return new PaymentRanker(new LowFeeSeeker()).rank(methods);
}

function rank_by_trust(methods) {
  return new PaymentRanker(new TrustBasedRanker()).rank(methods);
}

function rank_composite(methods, criteria) {
  return new PaymentRanker(new CompositeStrategy(criteria)).rank(methods);
}

function rank_with_refund_filter(methods, preferEasyRefund) {
  const composite = new CompositeStrategy([
    new EasyRefundStrategy(preferEasyRefund),
    new RewardsMaximizer(),
  ]);
  return new PaymentRanker(composite).rank(methods);
}

module.exports = {
  PaymentMethod, RankingStrategy, RewardsMaximizer, LowFeeSeeker, TrustBasedRanker,
  CompositeStrategy, EasyRefundStrategy, PaymentRanker,
  rank_by_rewards, rank_by_low_fee, rank_by_trust, rank_composite, rank_with_refund_filter,
};
