"""LRU Cache with TTL expiry and eviction listeners."""

from collections import OrderedDict


class LruOp:
    def __init__(self, kind, i1=0, i2=0, i3=0, i4=0):
        self.kind = kind
        self.i1 = i1
        self.i2 = i2
        self.i3 = i3
        self.i4 = i4


# Eviction reasons
CAPACITY = "CAPACITY"
TTL_EXPIRED = "TTL_EXPIRED"
EXPLICIT_DELETE = "EXPLICIT_DELETE"


class CapturingListener:
    def __init__(self):
        self.events = []  # list of (key, value, reason)

    def onEviction(self, key, value, reason):
        self.events.append((key, value, reason))


class LRUCache:
    def __init__(self, capacity):
        self.capacity = capacity
        # OrderedDict: key -> (value, expires_at)  (move_to_end -> most recent)
        # We treat the FIRST item as LRU, LAST as MRU.
        self.data = OrderedDict()
        self.listeners = []

    def _is_expired(self, expires_at, current_time):
        return expires_at > 0 and current_time >= expires_at

    def _notify(self, key, value, reason):
        for listener in self.listeners:
            listener.onEviction(key, value, reason)

    def _evict_key(self, key, reason):
        value, _ = self.data[key]
        del self.data[key]
        self._notify(key, value, reason)

    def _evict_lru(self):
        if not self.data:
            return
        key = next(iter(self.data))
        self._evict_key(key, CAPACITY)

    # Part 1
    def get(self, key):
        if key not in self.data:
            return -1
        self.data.move_to_end(key)
        return self.data[key][0]

    def put(self, key, value):
        if key in self.data:
            _, exp = self.data[key]
            self.data[key] = (value, exp)
            self.data.move_to_end(key)
            return
        if len(self.data) >= self.capacity:
            self._evict_lru()
        self.data[key] = (value, 0)

    # Part 2 (TTL)
    def get_t(self, key, current_time):
        if key not in self.data:
            return -1
        value, exp = self.data[key]
        if self._is_expired(exp, current_time):
            self._evict_key(key, TTL_EXPIRED)
            return -1
        self.data.move_to_end(key)
        return value

    def put_t(self, key, value, current_time, ttl=0):
        expires_at = current_time + ttl if ttl > 0 else 0
        if key in self.data:
            _, exp = self.data[key]
            if self._is_expired(exp, current_time):
                self._evict_key(key, TTL_EXPIRED)
            else:
                self.data[key] = (value, expires_at)
                self.data.move_to_end(key)
                return
        if len(self.data) >= self.capacity:
            self._evict_lru()
        self.data[key] = (value, expires_at)

    def delete_key(self, key):
        if key not in self.data:
            return False
        self._evict_key(key, EXPLICIT_DELETE)
        return True

    def size(self):
        return len(self.data)

    # Part 3
    def add_eviction_listener(self, listener):
        self.listeners.append(listener)

    def remove_eviction_listener(self, listener):
        self.listeners = [l for l in self.listeners if l is not listener]


def lru_simulate(ops):
    out = []
    cache = None
    listeners = []

    def ensure_listener(idx):
        while len(listeners) <= idx:
            listeners.append(CapturingListener())

    for op in ops:
        k = op.kind
        if k == "new":
            cache = LRUCache(op.i1)
            listeners = []
            out.append("ok")
        elif k == "put":
            cache.put(op.i1, op.i2)
            out.append("ok")
        elif k == "put_t":
            cache.put_t(op.i1, op.i2, op.i3, op.i4)
            out.append("ok")
        elif k == "get":
            out.append(str(cache.get(op.i1)))
        elif k == "get_t":
            out.append(str(cache.get_t(op.i1, op.i2)))
        elif k == "delete":
            out.append("ok" if cache.delete_key(op.i1) else "fail")
        elif k == "size":
            out.append(str(cache.size()))
        elif k == "add_listener":
            ensure_listener(op.i1)
            cache.add_eviction_listener(listeners[op.i1])
            out.append("ok")
        elif k == "remove_listener":
            if op.i1 < len(listeners):
                cache.remove_eviction_listener(listeners[op.i1])
            out.append("ok")
        elif k == "events_count":
            if op.i1 < len(listeners):
                out.append(str(len(listeners[op.i1].events)))
            else:
                out.append("0")
        elif k == "event_key":
            if op.i1 < len(listeners) and op.i2 < len(listeners[op.i1].events):
                out.append(str(listeners[op.i1].events[op.i2][0]))
            else:
                out.append("")
        elif k == "event_value":
            if op.i1 < len(listeners) and op.i2 < len(listeners[op.i1].events):
                out.append(str(listeners[op.i1].events[op.i2][1]))
            else:
                out.append("")
        elif k == "event_reason":
            if op.i1 < len(listeners) and op.i2 < len(listeners[op.i1].events):
                out.append(listeners[op.i1].events[op.i2][2])
            else:
                out.append("")
        else:
            out.append("unknown:" + k)
    return out
