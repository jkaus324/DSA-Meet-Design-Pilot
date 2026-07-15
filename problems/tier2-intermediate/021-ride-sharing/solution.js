"use strict";

class RideOp {
  constructor(kind, s1 = "", s2 = "", s3 = "", s4 = "", i1 = 0, i2 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.s4 = s4;
    this.i1 = i1;
    this.i2 = i2;
  }
}

class User {
  constructor(id, name) {
    this.id = id;
    this.name = name;
    this.rides_offered = 0;
    this.rides_taken = 0;
  }
}

class Vehicle {
  constructor(id, owner_id, model, reg_number) {
    this.id = id;
    this.owner_id = owner_id;
    this.model = model;
    this.reg_number = reg_number;
  }
}

class Ride {
  constructor(id, driver_id, vehicle_id, origin, destination, total_seats) {
    this.id = id;
    this.driver_id = driver_id;
    this.vehicle_id = vehicle_id;
    this.origin = origin;
    this.destination = destination;
    this.total_seats = total_seats;
    this.available_seats = total_seats;
    this.active = true;
  }
}

class MostVacantStrategy {
  select(candidates, seats_needed, preference) {
    let best = null;
    for (const ride of candidates) {
      if (ride.available_seats >= seats_needed) {
        if (best === null || ride.available_seats > best.available_seats) {
          best = ride;
        }
      }
    }
    return best;
  }
}

class PreferredVehicleStrategy {
  constructor(vehicle_store) {
    this.vehicle_store = vehicle_store;
  }
  select(candidates, seats_needed, preference) {
    for (const ride of candidates) {
      if (ride.available_seats >= seats_needed) {
        const v = this.vehicle_store.has(ride.vehicle_id)
          ? this.vehicle_store.get(ride.vehicle_id)
          : undefined;
        if (v !== undefined && v.model === preference) {
          return ride;
        }
      }
    }
    return null;
  }
}

class RideService {
  constructor() {
    this.users = new Map();
    this.vehicles = new Map();
    this.rides = new Map();
    this.active_vehicles = new Map(); // regNumber -> rideId
    this.ride_counter = 0;
  }

  add_user(name) {
    if (this.users.has(name)) return;
    this.users.set(name, new User(name, name));
  }

  add_vehicle(user_name, model, reg_number) {
    if (!this.users.has(user_name)) return;
    this.vehicles.set(reg_number, new Vehicle(reg_number, user_name, model, reg_number));
  }

  offer_ride(user_name, origin, dest, seats, vehicle_reg_number) {
    if (!this.users.has(user_name)) return "";
    if (!this.vehicles.has(vehicle_reg_number)) return "";
    if (this.vehicles.get(vehicle_reg_number).owner_id !== user_name) return "";
    if (this.active_vehicles.has(vehicle_reg_number)) return "";
    this.ride_counter += 1;
    const ride_id = "RIDE-" + String(this.ride_counter);
    this.rides.set(ride_id, new Ride(ride_id, user_name, vehicle_reg_number, origin, dest, seats));
    this.active_vehicles.set(vehicle_reg_number, ride_id);
    this.users.get(user_name).rides_offered += 1;
    return ride_id;
  }

  select_ride(passenger_name, origin, dest, seats, strategy, preference = "") {
    if (!this.users.has(passenger_name)) return "";
    const candidates = [];
    for (const ride of this.rides.values()) {
      if (ride.active && ride.origin === origin && ride.destination === dest
          && ride.available_seats >= seats && ride.driver_id !== passenger_name) {
        candidates.push(ride);
      }
    }
    const selected = strategy.select(candidates, seats, preference);
    if (selected !== null && selected !== undefined) {
      selected.available_seats -= seats;
      this.users.get(passenger_name).rides_taken += 1;
      return selected.id;
    }
    return "";
  }

  end_ride(ride_id) {
    const ride = this.rides.has(ride_id) ? this.rides.get(ride_id) : null;
    if (ride === null) return;
    if (!ride.active) return;
    ride.active = false;
    if (this.active_vehicles.has(ride.vehicle_id)) {
      this.active_vehicles.delete(ride.vehicle_id);
    }
  }

  has_user(name) { return this.users.has(name); }
  has_vehicle(reg) { return this.vehicles.has(reg); }
  has_ride(rid) { return this.rides.has(rid); }
  get_user(name) { return this.users.get(name); }
  get_ride(rid) { return this.rides.get(rid); }
  get_vehicles() { return this.vehicles; }
}

function ride_simulate(ops) {
  const out = [];
  let svc = new RideService();
  let ride_slots = new Array(32).fill("");
  for (const op of ops) {
    const k = op.kind;
    if (k === "new") {
      svc = new RideService();
      ride_slots = new Array(32).fill("");
      out.push("ok");
    } else if (k === "add_user") {
      svc.add_user(op.s1);
      out.push("ok");
    } else if (k === "add_veh") {
      svc.add_vehicle(op.s1, op.s2, op.s3);
      out.push("ok");
    } else if (k === "offer") {
      const rid = svc.offer_ride(op.s1, op.s2, op.s3, op.i1, op.s4);
      if (0 <= op.i2 && op.i2 < ride_slots.length) ride_slots[op.i2] = rid;
      out.push(rid);
    } else if (k === "ride_active") {
      const rid = ride_slots[op.i2];
      out.push(svc.has_ride(rid) && svc.get_ride(rid).active ? "yes" : "no");
    } else if (k === "ride_origin") {
      const rid = ride_slots[op.i2];
      out.push(svc.has_ride(rid) ? svc.get_ride(rid).origin : "");
    } else if (k === "ride_dest") {
      const rid = ride_slots[op.i2];
      out.push(svc.has_ride(rid) ? svc.get_ride(rid).destination : "");
    } else if (k === "ride_total") {
      const rid = ride_slots[op.i2];
      out.push(svc.has_ride(rid) ? String(svc.get_ride(rid).total_seats) : "-1");
    } else if (k === "ride_avail") {
      const rid = ride_slots[op.i2];
      out.push(svc.has_ride(rid) ? String(svc.get_ride(rid).available_seats) : "-1");
    } else if (k === "ride_driver") {
      const rid = ride_slots[op.i2];
      out.push(svc.has_ride(rid) ? svc.get_ride(rid).driver_id : "");
    } else if (k === "select_mv") {
      const rid = svc.select_ride(op.s1, op.s2, op.s3, op.i1, new MostVacantStrategy());
      if (0 <= op.i2 && op.i2 < ride_slots.length) ride_slots[op.i2] = rid;
      out.push(rid);
    } else if (k === "select_pv") {
      const rid = svc.select_ride(op.s1, op.s2, op.s3, op.i1,
        new PreferredVehicleStrategy(svc.get_vehicles()), op.s4);
      if (0 <= op.i2 && op.i2 < ride_slots.length) ride_slots[op.i2] = rid;
      out.push(rid);
    } else if (k === "end") {
      svc.end_ride(ride_slots[op.i2]);
      out.push("ok");
    } else if (k === "end_id") {
      svc.end_ride(op.s1);
      out.push("ok");
    } else if (k === "user_offered") {
      out.push(svc.has_user(op.s1) ? String(svc.get_user(op.s1).rides_offered) : "0");
    } else if (k === "user_taken") {
      out.push(svc.has_user(op.s1) ? String(svc.get_user(op.s1).rides_taken) : "0");
    } else if (k === "has_user") {
      out.push(svc.has_user(op.s1) ? "yes" : "no");
    } else if (k === "has_vehicle") {
      out.push(svc.has_vehicle(op.s1) ? "yes" : "no");
    } else if (k === "has_ride") {
      out.push(svc.has_ride(ride_slots[op.i2]) ? "yes" : "no");
    } else {
      out.push("unknown:" + k);
    }
  }
  return out;
}

module.exports = { RideOp, ride_simulate };
