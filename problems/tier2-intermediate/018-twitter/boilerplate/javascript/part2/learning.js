// Data class (given — do not modify).
class TwitterOp {
  constructor(kind, i1 = 0, i2 = 0) {
    this.kind = kind;
    this.i1 = i1;
    this.i2 = i2;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function twitter_simulate(ops) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { TwitterOp, twitter_simulate };
