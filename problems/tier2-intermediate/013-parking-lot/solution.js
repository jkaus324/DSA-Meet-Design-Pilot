'use strict';

/* Parking Lot — multi-floor + spot matching + pricing strategies. */

class ParkOp {
  constructor(kind, s1 = '', s2 = '', s3 = '', i1 = 0, i2 = 0, i3 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
  }
}

// Constants
const M_TYPE = 'MOTORCYCLE';
const C_TYPE = 'CAR';
const T_TYPE = 'TRUCK';

const SMALL = 0;
const MEDIUM = 1;
const LARGE = 2;

class Vehicle {
  constructor(licensePlate, type) {
    this.licensePlate = licensePlate;
    this.type = type;
  }
}

class ParkingSpot {
  constructor(spotId, floor, size) {
    this.spotId = spotId;
    this.floor = floor;
    this.size = size;
    this.isOccupied = false;
    this.vehicleLicensePlate = '';
  }
}

class Ticket {
  constructor() {
    this.ticketId = '';
    this.licensePlate = '';
    this.spotId = '';
    this.floor = 0;
    this.entryTime = 0;
    this.entryGateId = '';
    this.exitGateId = '';
  }
}

class Gate {
  constructor(gateId, type) {
    this.gateId = gateId;
    this.type = type;
  }
}

function _min_spot_size(vtype) {
  if (vtype === M_TYPE) return SMALL;
  if (vtype === C_TYPE) return MEDIUM;
  return LARGE;
}

function _is_compatible(spotSize, minRequired) {
  return spotSize >= minRequired;
}

class FlatRate {
  constructor(fee) {
    this.fee = fee;
  }
  calculateFee(durationSeconds) {
    return this.fee;
  }
}

class Hourly {
  constructor(rate) {
    this.rate = rate;
  }
  calculateFee(durationSeconds) {
    const hours = Math.ceil(durationSeconds / 3600.0);
    return this.rate * hours;
  }
}

class Tiered {
  constructor(base, mid, high) {
    this.base = base;
    this.mid = mid;
    this.high = high;
  }
  calculateFee(durationSeconds) {
    const hours = Math.ceil(durationSeconds / 3600.0);
    if (hours <= 1) return this.base;
    if (hours <= 3) return this.base + this.mid * (hours - 1);
    return this.base + this.mid * 2 + this.high * (hours - 3);
  }
}

class ParkingLot {
  constructor(numFloors) {
    this.floors = [];
    for (let i = 0; i < numFloors; i++) this.floors.push([]);
    this.activeTickets = new Map();
    this.gates = [];
    this.strategy = null;
    this.nextTicketId = 1;
  }

  addSpot(floor, size) {
    if (floor < 0 || floor >= this.floors.length) return;
    const spotId = 'F' + String(floor) + 'S' + String(this.floors[floor].length);
    this.floors[floor].push(new ParkingSpot(spotId, floor, size));
  }

  setPricingStrategy(s) {
    this.strategy = s;
  }

  addGate(gateId, type) {
    this.gates.push(new Gate(gateId, type));
  }

  getGates(type) {
    return this.gates.filter(g => g.type === type).map(g => g.gateId);
  }

  parkVehicle(vehicle, entryTime, gateId = '') {
    const minSize = _min_spot_size(vehicle.type);
    for (let f = 0; f < this.floors.length; f++) {
      const spots = this.floors[f];
      for (const spot of spots) {
        if (!spot.isOccupied && _is_compatible(spot.size, minSize)) {
          spot.isOccupied = true;
          spot.vehicleLicensePlate = vehicle.licensePlate;
          const tid = 'T' + String(this.nextTicketId);
          this.nextTicketId += 1;
          const t = new Ticket();
          t.ticketId = tid;
          t.licensePlate = vehicle.licensePlate;
          t.spotId = spot.spotId;
          t.floor = f;
          t.entryTime = entryTime;
          t.entryGateId = gateId;
          t.exitGateId = '';
          this.activeTickets.set(tid, t);
          return t;
        }
      }
    }
    return null;
  }

  unparkVehicle(ticketId, exitTime, gateId = '') {
    if (!this.activeTickets.has(ticketId)) return -1.0;
    const ticket = this.activeTickets.get(ticketId);
    ticket.exitGateId = gateId;
    for (const floorSpots of this.floors) {
      let broke = false;
      for (const spot of floorSpots) {
        if (spot.spotId === ticket.spotId && spot.isOccupied) {
          spot.isOccupied = false;
          spot.vehicleLicensePlate = '';
          broke = true;
          break;
        }
      }
      if (broke) {
        // Python only breaks the inner loop; keep outer loop continuing.
      }
    }
    const duration = exitTime - ticket.entryTime;
    let fee;
    if (this.strategy) {
      fee = this.strategy.calculateFee(duration);
    } else {
      fee = duration;
    }
    this.activeTickets.delete(ticketId);
    return fee;
  }

  getAvailableSpots(size) {
    let count = 0;
    for (const floorSpots of this.floors) {
      for (const spot of floorSpots) {
        if (!spot.isOccupied && spot.size === size) count += 1;
      }
    }
    return count;
  }

  getAvailableSpotsByFloor(floor, size) {
    if (floor < 0 || floor >= this.floors.length) return 0;
    let count = 0;
    for (const spot of this.floors[floor]) {
      if (!spot.isOccupied && spot.size === size) count += 1;
    }
    return count;
  }
}

function _size_from(s) {
  if (s === 'S' || s === 'small') return SMALL;
  if (s === 'M' || s === 'medium') return MEDIUM;
  return LARGE;
}

function _vtype_from(s) {
  if (s === 'M' || s === 'moto' || s === 'motorcycle') return M_TYPE;
  if (s === 'C' || s === 'car') return C_TYPE;
  return T_TYPE;
}

function _gate_from(s) {
  return s === 'entry' ? 'ENTRY' : 'EXIT';
}

function _fee_to_str(f) {
  if (f < 0) return '-1';
  return f.toFixed(2);
}

function parking_simulate(ops) {
  const out = [];
  let lot = null;
  let tickets = new Array(16).fill('');
  let snaps = [];
  const freshSnaps = () => {
    const arr = [];
    for (let i = 0; i < 16; i++) arr.push({ id: '', floor: -1, spotId: '', entryGate: '' });
    return arr;
  };
  snaps = freshSnaps();
  for (const op of ops) {
    const k = op.kind;
    if (k === 'new') {
      lot = new ParkingLot(op.i1);
      tickets = new Array(16).fill('');
      snaps = freshSnaps();
      out.push('ok');
    } else if (k === 'add_spot') {
      lot.addSpot(op.i1, _size_from(op.s1));
      out.push('ok');
    } else if (k === 'add_gate') {
      lot.addGate(op.s1, _gate_from(op.s2));
      out.push('ok');
    } else if (k === 'gates_count') {
      out.push(String(lot.getGates(_gate_from(op.s1)).length));
    } else if (k === 'gate_at') {
      const g = lot.getGates(_gate_from(op.s1));
      out.push(op.i1 >= 0 && op.i1 < g.length ? g[op.i1] : '');
    } else if (k === 'set_pricing') {
      let p = null;
      if (op.s1 === 'flat') p = new FlatRate(Number(op.i1));
      else if (op.s1 === 'hourly') p = new Hourly(Number(op.i1));
      else if (op.s1 === 'tiered') p = new Tiered(Number(op.i1), Number(op.i2), Number(op.i3));
      if (p) lot.setPricingStrategy(p);
      out.push('ok');
    } else if (k === 'park') {
      const v = new Vehicle(op.s1, _vtype_from(op.s2));
      const t = lot.parkVehicle(v, op.i1, op.s3);
      if (op.i2 >= 0 && op.i2 < tickets.length) {
        if (t) {
          tickets[op.i2] = t.ticketId;
          snaps[op.i2] = { id: t.ticketId, floor: t.floor, spotId: t.spotId, entryGate: t.entryGateId };
        } else {
          tickets[op.i2] = '';
          snaps[op.i2] = { id: '', floor: -1, spotId: '', entryGate: '' };
        }
      }
      out.push(t ? t.ticketId : '');
    } else if (k === 'ticket_at') {
      out.push(op.i1 >= 0 && op.i1 < tickets.length ? tickets[op.i1] : '');
    } else if (k === 'ticket_floor') {
      out.push(op.i1 >= 0 && op.i1 < snaps.length ? String(snaps[op.i1].floor) : '-1');
    } else if (k === 'ticket_spot_id') {
      out.push(op.i1 >= 0 && op.i1 < snaps.length ? snaps[op.i1].spotId : '');
    } else if (k === 'ticket_entry') {
      out.push(op.i1 >= 0 && op.i1 < snaps.length ? snaps[op.i1].entryGate : '');
    } else if (k === 'unpark') {
      const tid = op.i1 >= 0 && op.i1 < tickets.length ? tickets[op.i1] : '';
      const fee = lot.unparkVehicle(tid, op.i2, op.s1);
      out.push(_fee_to_str(fee));
    } else if (k === 'unpark_id') {
      const fee = lot.unparkVehicle(op.s1, op.i2, op.s2);
      out.push(_fee_to_str(fee));
    } else if (k === 'available') {
      out.push(String(lot.getAvailableSpots(_size_from(op.s1))));
    } else if (k === 'available_floor') {
      out.push(String(lot.getAvailableSpotsByFloor(op.i1, _size_from(op.s1))));
    } else {
      out.push('unknown:' + k);
    }
  }
  return out;
}

module.exports = { ParkOp, parking_simulate };
