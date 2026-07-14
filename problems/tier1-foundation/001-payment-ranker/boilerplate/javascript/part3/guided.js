// Data class (given).
class PaymentMethod {
  constructor(name, cashbackRate, transactionFee, usageCount, easyRefundEligible = false) {
    this.name = name;
    this.cashbackRate = cashbackRate;
    this.transactionFee = transactionFee;
    this.usageCount = usageCount;
    this.easyRefundEligible = easyRefundEligible;
  }
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function rank_by_rewards(methods) {
  // TODO: write your solution
  return methods;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function rank_by_low_fee(methods) {
  // TODO: write your solution
  return methods;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function rank_by_trust(methods) {
  // TODO: write your solution
  return methods;
}

// HINT: think about how to compose multiple criteria into a single decision.
function rank_composite(methods, criteria) {
  // TODO: write your solution
  return methods;
}

// HINT: a boolean flag changes ranking — handle it as a separate piece you can chain.
function rank_with_refund_filter(methods, preferEasyRefund) {
  // TODO: write your solution
  return methods;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { PaymentMethod, rank_by_rewards, rank_by_low_fee, rank_by_trust, rank_composite, rank_with_refund_filter };
