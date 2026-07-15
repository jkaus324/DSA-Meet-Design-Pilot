// Data class (given — do not modify).
class CartItem {
  constructor(name, price, quantity, category = "") {
    this.name = name;
    this.price = price;
    this.quantity = quantity;
    this.category = category;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function apply_percentage_discount(cart, percentage) {
  // TODO: implement this
  return null;
}

function apply_flat_discount(cart, amount) {
  // TODO: implement this
  return null;
}

function apply_bogo(cart, buyCount, freeCount) {
  // TODO: implement this
  return null;
}

function apply_percentage_with_eligibility(cart, percentage, minCartValue, requireFirstTimeUser, isFirstTimeUser, eligibleCategory) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { CartItem, apply_percentage_discount, apply_flat_discount, apply_bogo, apply_percentage_with_eligibility };
