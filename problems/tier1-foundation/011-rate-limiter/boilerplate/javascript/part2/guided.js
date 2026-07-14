// Data class (given).

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function reset_service() {
  // TODO: write your solution
  return null;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function init_limiter(maxRequests, windowSize) {
  // TODO: write your solution
  return null;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function allow_request_simple(clientId, timestamp, endpoint) {
  // TODO: write your solution
  return null;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function get_request_count(clientId) {
  // TODO: write your solution
  return null;
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
function allow_request_with_strategy_simple(algorithm, clientId, timestamp, endpoint) {
  // TODO: write your solution
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { reset_service, init_limiter, allow_request_simple, get_request_count, allow_request_with_strategy_simple };
