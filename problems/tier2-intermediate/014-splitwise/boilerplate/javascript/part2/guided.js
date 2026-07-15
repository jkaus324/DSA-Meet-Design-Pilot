// Data class (given).
class SplitOp {
  constructor(kind, s1 = "", s2 = "", s3 = "", s4 = "", i1 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.s4 = s4;
    this.i1 = i1;
  }
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function splitwise_simulate(ops) {
  // TODO: write your solution
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { SplitOp, splitwise_simulate };
