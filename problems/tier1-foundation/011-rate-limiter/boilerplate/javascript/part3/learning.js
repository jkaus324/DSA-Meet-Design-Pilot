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

function init_limiter(maxRequests, windowSize) {
  // TODO: implement this
  return null;
}

function allow_request_simple(clientId, timestamp, endpoint) {
  // TODO: implement this
  return null;
}

function get_request_count(clientId) {
  // TODO: implement this
  return null;
}

function allow_request_with_strategy_simple(algorithm, clientId, timestamp, endpoint) {
  // TODO: implement this
  return null;
}

function allow_request_for_tier_str(tier, clientId, timestamp, endpoint) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { reset_service, init_limiter, allow_request_simple, get_request_count, allow_request_with_strategy_simple, allow_request_for_tier_str };
