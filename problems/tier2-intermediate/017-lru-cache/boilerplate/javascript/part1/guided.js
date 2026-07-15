// Data class (given).
class LruOp {
  constructor(kind, i1 = 0, i2 = 0, i3 = 0, i4 = 0) {
    this.kind = kind;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
    this.i4 = i4;
  }
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function lru_simulate(ops) {
  // TODO: write your solution
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { LruOp, lru_simulate };
