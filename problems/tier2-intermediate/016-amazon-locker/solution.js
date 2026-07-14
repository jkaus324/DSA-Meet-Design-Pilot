'use strict';

/* Amazon Locker — allocation/deposit/retrieval with expiry + notifications. */

class LockerOp {
  constructor(kind, s1 = '', s2 = '', i1 = 0, i2 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.i1 = i1;
    this.i2 = i2;
  }
}

// Locker sizes — keep ordered SMALL < MEDIUM < LARGE
const SMALL = 'SMALL';
const MEDIUM = 'MEDIUM';
const LARGE = 'LARGE';

class Locker {
  constructor(lockerId, size) {
    this.locker_id = lockerId;
    this.size = size;
    this.occupied = false;
  }
}

class DepositRecord {
  constructor(lockerId, packageId, pickupCode, depositTime) {
    this.locker_id = lockerId;
    this.package_id = packageId;
    this.pickup_code = pickupCode;
    this.deposit_time = depositTime;
  }
}

class SmallestFitStrategy {
  allocate(packageSize, available) {
    let tryOrder;
    if (packageSize === SMALL) tryOrder = [SMALL, MEDIUM, LARGE];
    else if (packageSize === MEDIUM) tryOrder = [MEDIUM, LARGE];
    else tryOrder = [LARGE];
    for (const sz of tryOrder) {
      const q = available.get(sz);
      if (q && q.length > 0) {
        return q.shift();
      }
    }
    return '';
  }
}

class CapturingChannel {
  constructor(log) {
    this.log = log;
  }
  notify(packageId, message) {
    this.log.push(packageId + ': ' + message);
  }
}

class LockerSystem {
  constructor() {
    this.lockers = new Map();
    this.available = new Map([
      [SMALL, []],
      [MEDIUM, []],
      [LARGE, []],
    ]);
    this.active_deposits = new Map();
    this.strategy = new SmallestFitStrategy();
    this.channels = [];
    this.code_counter = 0;
    this.expiry_hours = 0;
  }

  _generate_code() {
    this.code_counter += 1;
    return 'CODE-' + String(this.code_counter);
  }

  _notify_all(packageId, message) {
    for (const ch of this.channels) ch.notify(packageId, message);
  }

  _free_locker(lockerId) {
    if (this.lockers.has(lockerId)) {
      const locker = this.lockers.get(lockerId);
      locker.occupied = false;
      this.available.get(locker.size).push(lockerId);
    }
  }

  add_locker(lockerId, size) {
    this.lockers.set(lockerId, new Locker(lockerId, size));
    this.available.get(size).push(lockerId);
  }

  deposit_package(packageId, size, depositTime = 0) {
    const lockerId = this.strategy.allocate(size, this.available);
    if (!lockerId) return '';
    this.lockers.get(lockerId).occupied = true;
    const code = this._generate_code();
    this.active_deposits.set(code, new DepositRecord(lockerId, packageId, code, depositTime));
    this._notify_all(packageId, 'Package ' + packageId + ' deposited. Code: ' + code);
    return code;
  }

  retrieve_package(code) {
    if (!this.active_deposits.has(code)) return false;
    const rec = this.active_deposits.get(code);
    this._free_locker(rec.locker_id);
    this.active_deposits.delete(code);
    return true;
  }

  set_code_expiry(hours) {
    this.expiry_hours = hours;
  }

  check_expired(currentTime) {
    const expired = [];
    if (this.expiry_hours <= 0) return expired;
    for (const code of [...this.active_deposits.keys()]) {
      const rec = this.active_deposits.get(code);
      if (currentTime - rec.deposit_time > this.expiry_hours * 3600) {
        this._free_locker(rec.locker_id);
        expired.push(rec.package_id);
        this._notify_all(rec.package_id, 'Package ' + rec.package_id + ' expired. Locker freed.');
        this.active_deposits.delete(code);
      }
    }
    return expired;
  }

  add_notification_channel(channel) {
    this.channels.push(channel);
  }
}

function _lsize_from(s) {
  if (s === 'S') return SMALL;
  if (s === 'M') return MEDIUM;
  return LARGE;
}

function locker_simulate(ops) {
  const out = [];
  let sys = new LockerSystem();
  let codes = new Array(32).fill('');
  let chanLog = [];
  let chan = null;
  let lastExpired = [];
  for (const op of ops) {
    const k = op.kind;
    if (k === 'new') {
      sys = new LockerSystem();
      codes = new Array(32).fill('');
      chanLog = [];
      chan = null;
      lastExpired = [];
      out.push('ok');
    } else if (k === 'add_locker') {
      sys.add_locker(op.s1, _lsize_from(op.s2));
      out.push('ok');
    } else if (k === 'deposit') {
      const code = sys.deposit_package(op.s1, _lsize_from(op.s2), op.i1);
      if (op.i2 >= 0 && op.i2 < codes.length) codes[op.i2] = code;
      out.push(code);
    } else if (k === 'code_at') {
      out.push(codes[op.i2]);
    } else if (k === 'retrieve') {
      out.push(sys.retrieve_package(codes[op.i2]) ? 'ok' : 'fail');
    } else if (k === 'retrieve_id') {
      out.push(sys.retrieve_package(op.s1) ? 'ok' : 'fail');
    } else if (k === 'set_expiry') {
      sys.set_code_expiry(op.i1);
      out.push('ok');
    } else if (k === 'check_expired') {
      lastExpired = sys.check_expired(op.i1);
      out.push(String(lastExpired.length));
    } else if (k === 'expired_at') {
      if (op.i2 >= 0 && op.i2 < lastExpired.length) out.push(lastExpired[op.i2]);
      else out.push('');
    } else if (k === 'add_chan') {
      chan = new CapturingChannel(chanLog);
      sys.add_notification_channel(chan);
      out.push('ok');
    } else if (k === 'chan_log_size') {
      out.push(String(chanLog.length));
    } else if (k === 'chan_log_contains') {
      const found = chanLog.some(entry => entry.includes(op.s1));
      out.push(found ? 'yes' : 'no');
    } else {
      out.push('unknown:' + k);
    }
  }
  return out;
}

module.exports = { LockerOp, locker_simulate };
