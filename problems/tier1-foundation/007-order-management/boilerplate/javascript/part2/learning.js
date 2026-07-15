// Data class (given — do not modify).

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function reset_service() {
  // TODO: implement this
  return null;
}

function set_inventory(productId, qty) {
  // TODO: implement this
  return null;
}

function get_inventory(productId) {
  // TODO: implement this
  return null;
}

function create_order_simple(productId, quantity, totalAmount) {
  // TODO: implement this
  return null;
}

function get_order_state_str(orderId) {
  // TODO: implement this
  return null;
}

function confirm_order(orderId) {
  // TODO: implement this
  return null;
}

function ship_order(orderId) {
  // TODO: implement this
  return null;
}

function deliver_order(orderId) {
  // TODO: implement this
  return null;
}

function cancel_order(orderId) {
  // TODO: implement this
  return null;
}

function get_history_size(orderId) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { reset_service, set_inventory, get_inventory, create_order_simple, get_order_state_str, confirm_order, ship_order, deliver_order, cancel_order, get_history_size };
