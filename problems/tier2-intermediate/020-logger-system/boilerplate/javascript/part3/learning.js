// Data class (given — do not modify).
class LogOp {
  constructor(kind, s1 = "", s2 = "", i1 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.i1 = i1;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function logger_simulate(ops) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { LogOp, logger_simulate };
