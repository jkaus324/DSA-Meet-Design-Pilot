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

// TODO: design and implement your solution.
// Required functions:
//   function rank_by_rewards(methods)
//   function rank_by_low_fee(methods)
//   function rank_by_trust(methods)
//   function rank_composite(methods, criteria)
//   function rank_with_refund_filter(methods, preferEasyRefund)

function rank_by_rewards(methods) {
  // TODO: write your solution
  return methods;
}

function rank_by_low_fee(methods) {
  // TODO: write your solution
  return methods;
}

function rank_by_trust(methods) {
  // TODO: write your solution
  return methods;
}

function rank_composite(methods, criteria) {
  // TODO: write your solution
  return methods;
}

function rank_with_refund_filter(methods, preferEasyRefund) {
  // TODO: write your solution
  return methods;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { PaymentMethod, rank_by_rewards, rank_by_low_fee, rank_by_trust, rank_composite, rank_with_refund_filter };
