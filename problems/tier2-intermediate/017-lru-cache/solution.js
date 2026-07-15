"use strict";
// LRU Cache with TTL expiry and eviction listeners.

class LruOp {
  constructor(kind, i1 = 0, i2 = 0, i3 = 0, i4 = 0) {
    this.kind = kind;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
    this.i4 = i4;
  }
}

// Eviction reasons
const CAPACITY = "CAPACITY";
const TTL_EXPIRED = "TTL_EXPIRED";
const EXPLICIT_DELETE = "EXPLICIT_DELETE";

class CapturingListener {
  constructor() {
    this.events = []; // list of [key, value, reason]
  }
  onEviction(key, value, reason) {
    this.events.push([key, value, reason]);
  }
}

class LRUCache {
  constructor(capacity) {
    this.capacity = capacity;
    // Map: key -> [value, expires_at]. Insertion order: FIRST is LRU, LAST is MRU.
    this.data = new Map();
    this.listeners = [];
  }

  _isExpired(expiresAt, currentTime) {
    return expiresAt > 0 && currentTime >= expiresAt;
  }

  _notify(key, value, reason) {
    for (const listener of this.listeners) {
      listener.onEviction(key, value, reason);
    }
  }

  _evictKey(key, reason) {
    const [value] = this.data.get(key);
    this.data.delete(key);
    this._notify(key, value, reason);
  }

  _evictLru() {
    if (this.data.size === 0) return;
    const key = this.data.keys().next().value;
    this._evictKey(key, CAPACITY);
  }

  // Part 1
  get(key) {
    if (!this.data.has(key)) return -1;
    const entry = this.data.get(key);
    this.data.delete(key);
    this.data.set(key, entry);
    return entry[0];
  }

  put(key, value) {
    if (this.data.has(key)) {
      const [, exp] = this.data.get(key);
      this.data.delete(key);
      this.data.set(key, [value, exp]);
      return;
    }
    if (this.data.size >= this.capacity) {
      this._evictLru();
    }
    this.data.set(key, [value, 0]);
  }

  // Part 2 (TTL)
  get_t(key, currentTime) {
    if (!this.data.has(key)) return -1;
    const [value, exp] = this.data.get(key);
    if (this._isExpired(exp, currentTime)) {
      this._evictKey(key, TTL_EXPIRED);
      return -1;
    }
    const entry = this.data.get(key);
    this.data.delete(key);
    this.data.set(key, entry);
    return value;
  }

  put_t(key, value, currentTime, ttl = 0) {
    const expiresAt = ttl > 0 ? currentTime + ttl : 0;
    if (this.data.has(key)) {
      const [, exp] = this.data.get(key);
      if (this._isExpired(exp, currentTime)) {
        this._evictKey(key, TTL_EXPIRED);
      } else {
        this.data.delete(key);
        this.data.set(key, [value, expiresAt]);
        return;
      }
    }
    if (this.data.size >= this.capacity) {
      this._evictLru();
    }
    this.data.set(key, [value, expiresAt]);
  }

  delete_key(key) {
    if (!this.data.has(key)) return false;
    this._evictKey(key, EXPLICIT_DELETE);
    return true;
  }

  size() {
    return this.data.size;
  }

  // Part 3
  add_eviction_listener(listener) {
    this.listeners.push(listener);
  }

  remove_eviction_listener(listener) {
    this.listeners = this.listeners.filter((l) => l !== listener);
  }
}

function lru_simulate(ops) {
  const out = [];
  let cache = null;
  let listeners = [];

  const ensureListener = (idx) => {
    while (listeners.length <= idx) {
      listeners.push(new CapturingListener());
    }
  };

  for (const op of ops) {
    const k = op.kind;
    if (k === "new") {
      cache = new LRUCache(op.i1);
      listeners = [];
      out.push("ok");
    } else if (k === "put") {
      cache.put(op.i1, op.i2);
      out.push("ok");
    } else if (k === "put_t") {
      cache.put_t(op.i1, op.i2, op.i3, op.i4);
      out.push("ok");
    } else if (k === "get") {
      out.push(String(cache.get(op.i1)));
    } else if (k === "get_t") {
      out.push(String(cache.get_t(op.i1, op.i2)));
    } else if (k === "delete") {
      out.push(cache.delete_key(op.i1) ? "ok" : "fail");
    } else if (k === "size") {
      out.push(String(cache.size()));
    } else if (k === "add_listener") {
      ensureListener(op.i1);
      cache.add_eviction_listener(listeners[op.i1]);
      out.push("ok");
    } else if (k === "remove_listener") {
      if (op.i1 < listeners.length) {
        cache.remove_eviction_listener(listeners[op.i1]);
      }
      out.push("ok");
    } else if (k === "events_count") {
      if (op.i1 < listeners.length) {
        out.push(String(listeners[op.i1].events.length));
      } else {
        out.push("0");
      }
    } else if (k === "event_key") {
      if (op.i1 < listeners.length && op.i2 < listeners[op.i1].events.length) {
        out.push(String(listeners[op.i1].events[op.i2][0]));
      } else {
        out.push("");
      }
    } else if (k === "event_value") {
      if (op.i1 < listeners.length && op.i2 < listeners[op.i1].events.length) {
        out.push(String(listeners[op.i1].events[op.i2][1]));
      } else {
        out.push("");
      }
    } else if (k === "event_reason") {
      if (op.i1 < listeners.length && op.i2 < listeners[op.i1].events.length) {
        out.push(listeners[op.i1].events[op.i2][2]);
      } else {
        out.push("");
      }
    } else {
      out.push("unknown:" + k);
    }
  }
  return out;
}

module.exports = { LruOp, lru_simulate };
