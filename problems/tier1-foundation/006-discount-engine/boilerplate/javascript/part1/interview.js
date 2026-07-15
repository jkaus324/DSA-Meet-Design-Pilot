// Data class (given).
class CartItem {
  constructor(name, price, quantity, category = "") {
    this.name = name;
    this.price = price;
    this.quantity = quantity;
    this.category = category;
  }
}

// TODO: design and implement your solution.
// Required functions:
//   function apply_percentage_discount(cart, percentage)
//   function apply_flat_discount(cart, amount)
//   function apply_bogo(cart, buyCount, freeCount)

function apply_percentage_discount(cart, percentage) {
  // TODO: write your solution
  return null;
}

function apply_flat_discount(cart, amount) {
  // TODO: write your solution
  return null;
}

function apply_bogo(cart, buyCount, freeCount) {
  // TODO: write your solution
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { CartItem, apply_percentage_discount, apply_flat_discount, apply_bogo };
