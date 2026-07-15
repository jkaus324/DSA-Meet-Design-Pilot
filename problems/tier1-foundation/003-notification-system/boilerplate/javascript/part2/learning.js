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

function notify_event(event, userIds, subscribedChannels) {
  // TODO: implement this
  return null;
}

function notify_priority(event, priority, userIds, subscribedChannels, minPriority) {
  // TODO: implement this
  return null;
}

function notify_priority_level(priority) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { reset_service, notify_event, notify_priority, notify_priority_level };
