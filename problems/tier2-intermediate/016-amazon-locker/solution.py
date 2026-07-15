"""Amazon Locker — allocation/deposit/retrieval with expiry + notifications."""

from collections import deque


class LockerOp:
    def __init__(self, kind, s1="", s2="", i1=0, i2=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.i1 = i1
        self.i2 = i2


# Locker sizes — keep ordered SMALL < MEDIUM < LARGE
SMALL = "SMALL"
MEDIUM = "MEDIUM"
LARGE = "LARGE"


class Locker:
    def __init__(self, locker_id, size):
        self.locker_id = locker_id
        self.size = size
        self.occupied = False


class DepositRecord:
    def __init__(self, locker_id, package_id, pickup_code, deposit_time):
        self.locker_id = locker_id
        self.package_id = package_id
        self.pickup_code = pickup_code
        self.deposit_time = deposit_time


class SmallestFitStrategy:
    def allocate(self, package_size, available):
        if package_size == SMALL:
            try_order = [SMALL, MEDIUM, LARGE]
        elif package_size == MEDIUM:
            try_order = [MEDIUM, LARGE]
        else:
            try_order = [LARGE]
        for sz in try_order:
            q = available.get(sz)
            if q:
                return q.popleft()
        return ""


class CapturingChannel:
    def __init__(self, log):
        self.log = log

    def notify(self, package_id, message):
        self.log.append(package_id + ": " + message)


class LockerSystem:
    def __init__(self):
        self.lockers = {}
        self.available = {SMALL: deque(), MEDIUM: deque(), LARGE: deque()}
        self.active_deposits = {}
        self.strategy = SmallestFitStrategy()
        self.channels = []
        self.code_counter = 0
        self.expiry_hours = 0

    def _generate_code(self):
        self.code_counter += 1
        return "CODE-" + str(self.code_counter)

    def _notify_all(self, package_id, message):
        for ch in self.channels:
            ch.notify(package_id, message)

    def _free_locker(self, locker_id):
        if locker_id in self.lockers:
            locker = self.lockers[locker_id]
            locker.occupied = False
            self.available[locker.size].append(locker_id)

    def add_locker(self, locker_id, size):
        self.lockers[locker_id] = Locker(locker_id, size)
        self.available[size].append(locker_id)

    def deposit_package(self, package_id, size, deposit_time=0):
        locker_id = self.strategy.allocate(size, self.available)
        if not locker_id:
            return ""
        self.lockers[locker_id].occupied = True
        code = self._generate_code()
        self.active_deposits[code] = DepositRecord(locker_id, package_id, code, deposit_time)
        self._notify_all(package_id, "Package " + package_id + " deposited. Code: " + code)
        return code

    def retrieve_package(self, code):
        if code not in self.active_deposits:
            return False
        rec = self.active_deposits[code]
        self._free_locker(rec.locker_id)
        del self.active_deposits[code]
        return True

    def set_code_expiry(self, hours):
        self.expiry_hours = hours

    def check_expired(self, current_time):
        expired = []
        if self.expiry_hours <= 0:
            return expired
        for code in list(self.active_deposits.keys()):
            rec = self.active_deposits[code]
            if current_time - rec.deposit_time > self.expiry_hours * 3600:
                self._free_locker(rec.locker_id)
                expired.append(rec.package_id)
                self._notify_all(rec.package_id, "Package " + rec.package_id + " expired. Locker freed.")
                del self.active_deposits[code]
        return expired

    def add_notification_channel(self, channel):
        self.channels.append(channel)


def _lsize_from(s):
    if s == "S":
        return SMALL
    if s == "M":
        return MEDIUM
    return LARGE


def locker_simulate(ops):
    out = []
    sys = LockerSystem()
    codes = [""] * 32
    chan_log = []
    chan = None
    last_expired = []
    for op in ops:
        k = op.kind
        if k == "new":
            sys = LockerSystem()
            codes = [""] * 32
            chan_log = []
            chan = None
            last_expired = []
            out.append("ok")
        elif k == "add_locker":
            sys.add_locker(op.s1, _lsize_from(op.s2))
            out.append("ok")
        elif k == "deposit":
            code = sys.deposit_package(op.s1, _lsize_from(op.s2), op.i1)
            if 0 <= op.i2 < len(codes):
                codes[op.i2] = code
            out.append(code)
        elif k == "code_at":
            out.append(codes[op.i2])
        elif k == "retrieve":
            out.append("ok" if sys.retrieve_package(codes[op.i2]) else "fail")
        elif k == "retrieve_id":
            out.append("ok" if sys.retrieve_package(op.s1) else "fail")
        elif k == "set_expiry":
            sys.set_code_expiry(op.i1)
            out.append("ok")
        elif k == "check_expired":
            last_expired = sys.check_expired(op.i1)
            out.append(str(len(last_expired)))
        elif k == "expired_at":
            if 0 <= op.i2 < len(last_expired):
                out.append(last_expired[op.i2])
            else:
                out.append("")
        elif k == "add_chan":
            chan = CapturingChannel(chan_log)
            sys.add_notification_channel(chan)
            out.append("ok")
        elif k == "chan_log_size":
            out.append(str(len(chan_log)))
        elif k == "chan_log_contains":
            found = any(op.s1 in entry for entry in chan_log)
            out.append("yes" if found else "no")
        else:
            out.append("unknown:" + k)
    return out
