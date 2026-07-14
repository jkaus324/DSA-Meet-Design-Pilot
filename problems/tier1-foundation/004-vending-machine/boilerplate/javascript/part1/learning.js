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

function reset() {
  // TODO: implement this
  return null;
}

function getState() {
  // TODO: implement this
  return null;
}

function selectItem(item) {
  // TODO: implement this
  return null;
}

function insertMoney(amount) {
  // TODO: implement this
  return null;
}

function dispense() {
  // TODO: implement this
  return null;
}

function cancel() {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { reset_service, reset, getState, selectItem, insertMoney, dispense, cancel };
