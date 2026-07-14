// Elevator system — SCAN ordering + dispatch strategies (JavaScript).

class _SortedSet {
  constructor() {
    this._s = [];
    this._set = new Set();
  }

  add(x) {
    if (!this._set.has(x)) {
      this._set.add(x);
      // insert sorted
      let lo = 0;
      let hi = this._s.length;
      while (lo < hi) {
        const mid = Math.floor((lo + hi) / 2);
        if (this._s[mid] < x) lo = mid + 1;
        else hi = mid;
      }
      this._s.splice(lo, 0, x);
    }
  }

  remove(x) {
    if (this._set.has(x)) {
      this._set.delete(x);
      const i = this._s.indexOf(x);
      if (i !== -1) this._s.splice(i, 1);
    }
  }

  has(x) { return this._set.has(x); }

  get length() { return this._s.length; }

  empty() { return this._s.length === 0; }

  first() { return this._s[0]; }

  last() { return this._s[this._s.length - 1]; }
}

class ElevOp {
  constructor(kind, s1 = '', i1 = 0, i2 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.i1 = i1;
    this.i2 = i2;
  }
}

// State constants
const IDLE = 'IDLE';
const MOVING_UP = 'MOVING_UP';
const MOVING_DOWN = 'MOVING_DOWN';
const DOOR_OPEN = 'DOOR_OPEN';

// Direction constants
const DIR_UP = 'UP';
const DIR_DOWN = 'DOWN';
const DIR_NONE = 'NONE';

class Elevator {
  constructor(id = 0) {
    this.id = id;
    this.currentFloor = 0;
    this.state = IDLE;
    this.currentDirection = DIR_NONE;
    this.upRequests = new _SortedSet();
    this.downRequests = new _SortedSet();
  }

  getId() { return this.id; }

  getCurrentFloor() { return this.currentFloor; }

  getState() { return this.state; }

  getCurrentDirection() { return this.currentDirection; }

  getPendingCount() { return this.upRequests.length + this.downRequests.length; }

  addRequest(floor, direction) {
    if (floor === this.currentFloor && this.state === IDLE) {
      this.state = DOOR_OPEN;
      return;
    }
    if (direction === DIR_UP) {
      this.upRequests.add(floor);
    } else if (direction === DIR_DOWN) {
      this.downRequests.add(floor);
    } else {
      if (floor > this.currentFloor) this.upRequests.add(floor);
      else this.downRequests.add(floor);
    }
    if (this.state === IDLE) {
      if (!this.upRequests.empty() && (
        this.downRequests.empty()
        || Math.abs(this.upRequests.first() - this.currentFloor)
          <= Math.abs(this.downRequests.last() - this.currentFloor)
      )) {
        this.currentDirection = DIR_UP;
        this.state = MOVING_UP;
      } else {
        this.currentDirection = DIR_DOWN;
        this.state = MOVING_DOWN;
      }
    }
  }

  step() {
    if (this.state === IDLE) return;
    if (this.state === MOVING_UP) {
      this.currentFloor += 1;
      if (this.upRequests.has(this.currentFloor)) {
        this.upRequests.remove(this.currentFloor);
        this.state = DOOR_OPEN;
      }
      return;
    }
    if (this.state === MOVING_DOWN) {
      this.currentFloor -= 1;
      if (this.downRequests.has(this.currentFloor)) {
        this.downRequests.remove(this.currentFloor);
        this.state = DOOR_OPEN;
      }
      return;
    }
    if (this.state === DOOR_OPEN) {
      if (this.currentDirection === DIR_UP) {
        if (!this.upRequests.empty()) {
          this.state = MOVING_UP;
        } else if (!this.downRequests.empty()) {
          this.currentDirection = DIR_DOWN;
          this.state = MOVING_DOWN;
        } else {
          this.currentDirection = DIR_NONE;
          this.state = IDLE;
        }
      } else if (this.currentDirection === DIR_DOWN) {
        if (!this.downRequests.empty()) {
          this.state = MOVING_DOWN;
        } else if (!this.upRequests.empty()) {
          this.currentDirection = DIR_UP;
          this.state = MOVING_UP;
        } else {
          this.currentDirection = DIR_NONE;
          this.state = IDLE;
        }
      } else {
        if (!this.upRequests.empty()) {
          this.currentDirection = DIR_UP;
          this.state = MOVING_UP;
        } else if (!this.downRequests.empty()) {
          this.currentDirection = DIR_DOWN;
          this.state = MOVING_DOWN;
        } else {
          this.state = IDLE;
        }
      }
    }
  }
}

class NearestFirst {
  selectElevator(elevators, requestFloor, requestDirection) {
    let bestIdx = 0;
    let bestScore = Infinity;
    const PENALTY = 10000;
    for (let i = 0; i < elevators.length; i++) {
      const e = elevators[i];
      const dist = Math.abs(e.getCurrentFloor() - requestFloor);
      const st = e.getState();
      const d = e.getCurrentDirection();
      let score;
      if (st === IDLE || d === DIR_NONE) {
        score = dist;
      } else if (d === DIR_UP && requestDirection === DIR_UP && e.getCurrentFloor() <= requestFloor) {
        score = dist;
      } else if (d === DIR_DOWN && requestDirection === DIR_DOWN && e.getCurrentFloor() >= requestFloor) {
        score = dist;
      } else {
        score = dist + PENALTY;
      }
      if (score < bestScore) {
        bestScore = score;
        bestIdx = i;
      }
    }
    return bestIdx;
  }
}

class LeastLoaded {
  selectElevator(elevators, requestFloor, requestDirection) {
    let bestIdx = 0;
    let bestCount = Infinity;
    for (let i = 0; i < elevators.length; i++) {
      const cnt = elevators[i].getPendingCount();
      if (cnt < bestCount) {
        bestCount = cnt;
        bestIdx = i;
      }
    }
    return bestIdx;
  }
}

class ElevatorSystem {
  constructor() {
    this.elevators = [];
    this.strategy = null;
  }

  addElevator(id) { this.elevators.push(new Elevator(id)); }

  setDispatchStrategy(s) { this.strategy = s; }

  getElevator(index) {
    if (index < 0 || index >= this.elevators.length) return null;
    return this.elevators[index];
  }

  getElevatorCount() { return this.elevators.length; }

  addRequest(floor, direction) {
    if (this.elevators.length === 0) return;
    let idx = 0;
    if (this.strategy) {
      idx = this.strategy.selectElevator(this.elevators, floor, direction);
    }
    this.elevators[idx].addRequest(floor, direction);
  }

  step() {
    for (const e of this.elevators) e.step();
  }
}

function _dir_from(s) {
  if (s === 'up') return DIR_UP;
  if (s === 'down') return DIR_DOWN;
  return DIR_NONE;
}

function elevator_simulate(ops) {
  const out = [];
  let single = null;
  let sys = null;
  const nf = new NearestFirst();
  const ll = new LeastLoaded();
  for (const op of ops) {
    const k = op.kind;
    if (k === 'new_elev') {
      single = new Elevator();
      sys = null;
      out.push('ok');
    } else if (k === 'new_sys') {
      sys = new ElevatorSystem();
      single = null;
      out.push('ok');
    } else if (k === 'add_elev') {
      sys.addElevator(op.i1);
      out.push('ok');
    } else if (k === 'set_strategy') {
      if (op.s1 === 'nearest') sys.setDispatchStrategy(nf);
      else if (op.s1 === 'least_loaded') sys.setDispatchStrategy(ll);
      out.push('ok');
    } else if (k === 'req') {
      single.addRequest(op.i1, _dir_from(op.s1));
      out.push('ok');
    } else if (k === 'sys_req') {
      sys.addRequest(op.i1, _dir_from(op.s1));
      out.push('ok');
    } else if (k === 'elev_req') {
      sys.getElevator(op.i1).addRequest(op.i2, _dir_from(op.s1));
      out.push('ok');
    } else if (k === 'step') {
      single.step();
      out.push('ok');
    } else if (k === 'sys_step') {
      sys.step();
      out.push('ok');
    } else if (k === 'elev_step') {
      sys.getElevator(op.i1).step();
      out.push('ok');
    } else if (k === 'floor') {
      out.push(String(single.getCurrentFloor()));
    } else if (k === 'elev_floor') {
      out.push(String(sys.getElevator(op.i1).getCurrentFloor()));
    } else if (k === 'state') {
      out.push(single.getState());
    } else if (k === 'elev_state') {
      out.push(sys.getElevator(op.i1).getState());
    } else if (k === 'elev_pending') {
      out.push(String(sys.getElevator(op.i1).getPendingCount()));
    } else if (k === 'count') {
      out.push(String(sys.getElevatorCount()));
    } else if (k === 'elev_null') {
      out.push(sys.getElevator(op.i1) === null ? 'yes' : 'no');
    } else {
      out.push('unknown:' + k);
    }
  }
  return out;
}

module.exports = {
  ElevOp, Elevator, NearestFirst, LeastLoaded, ElevatorSystem,
  elevator_simulate,
};
