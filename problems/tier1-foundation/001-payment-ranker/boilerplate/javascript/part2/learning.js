// Data class (given — do not modify).
class PaymentMethod {
  constructor(name, cashbackRate, transactionFee, usageCount, easyRefundEligible = false) {
    this.name = name;
    this.cashbackRate = cashbackRate;
    this.transactionFee = transactionFee;
    this.usageCount = usageCount;
    this.easyRefundEligible = easyRefundEligible;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

// Concrete strategies — fill in compare() bodies.
class RewardsMaximizer extends Strategy {
  compare(a, b) {
    // TODO: implement this
    return false;
  }
}

class LowFeeSeeker extends Strategy {
  compare(a, b) {
    // TODO: implement this
    return false;
  }
}

class TrustBasedRanker extends Strategy {
  compare(a, b) {
    // TODO: implement this
    return false;
  }
}

function rank_by_rewards(methods) {
  // TODO: implement this
  return methods;
}

function rank_by_low_fee(methods) {
  // TODO: implement this
  return methods;
}

function rank_by_trust(methods) {
  // TODO: implement this
  return methods;
}

function rank_composite(methods, criteria) {
  // TODO: implement this
  return methods;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { PaymentMethod, Strategy, RewardsMaximizer, LowFeeSeeker, TrustBasedRanker, rank_by_rewards, rank_by_low_fee, rank_by_trust, rank_composite };
