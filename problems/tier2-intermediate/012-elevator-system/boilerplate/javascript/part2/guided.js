// Data class (given).
class ElevOp {
  constructor(kind, s1 = "", i1 = 0, i2 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.i1 = i1;
    this.i2 = i2;
  }
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function elevator_simulate(ops) {
  // TODO: write your solution
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { ElevOp, elevator_simulate };
