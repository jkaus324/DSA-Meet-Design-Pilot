// Data class (given — do not modify).
class ElevOp {
  constructor(kind, s1 = "", i1 = 0, i2 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.i1 = i1;
    this.i2 = i2;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function elevator_simulate(ops) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { ElevOp, elevator_simulate };
