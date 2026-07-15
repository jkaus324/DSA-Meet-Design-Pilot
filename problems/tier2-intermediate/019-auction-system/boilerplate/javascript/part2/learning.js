// Data class (given — do not modify).
class AuctionOp {
  constructor(kind, s1 = "", s2 = "", s3 = "", i1 = 0, i2 = 0, i3 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function auction_simulate(ops) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { AuctionOp, auction_simulate };
