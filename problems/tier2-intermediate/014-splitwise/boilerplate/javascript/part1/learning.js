// Data class (given — do not modify).
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

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function splitwise_simulate(ops) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { SplitOp, splitwise_simulate };
