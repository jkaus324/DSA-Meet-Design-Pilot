// Data class (given — do not modify).
class LruOp {
  constructor(kind, i1 = 0, i2 = 0, i3 = 0, i4 = 0) {
    this.kind = kind;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
    this.i4 = i4;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function lru_simulate(ops) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { LruOp, lru_simulate };
