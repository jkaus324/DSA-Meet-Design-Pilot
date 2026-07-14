"""Ride-sharing service — onboarding, selection strategies, end + statistics."""


class RideOp:
    def __init__(self, kind, s1="", s2="", s3="", s4="", i1=0, i2=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.s4 = s4
        self.i1 = i1
        self.i2 = i2


class User:
    def __init__(self, id, name):
        self.id = id
        self.name = name
        self.rides_offered = 0
        self.rides_taken = 0


class Vehicle:
    def __init__(self, id, owner_id, model, reg_number):
        self.id = id
        self.owner_id = owner_id
        self.model = model
        self.reg_number = reg_number


class Ride:
    def __init__(self, id, driver_id, vehicle_id, origin, destination, total_seats):
        self.id = id
        self.driver_id = driver_id
        self.vehicle_id = vehicle_id
        self.origin = origin
        self.destination = destination
        self.total_seats = total_seats
        self.available_seats = total_seats
        self.active = True


class MostVacantStrategy:
    def select(self, candidates, seats_needed, preference):
        best = None
        for ride in candidates:
            if ride.available_seats >= seats_needed:
                if best is None or ride.available_seats > best.available_seats:
                    best = ride
        return best


class PreferredVehicleStrategy:
    def __init__(self, vehicle_store):
        self.vehicle_store = vehicle_store

    def select(self, candidates, seats_needed, preference):
        for ride in candidates:
            if ride.available_seats >= seats_needed:
                v = self.vehicle_store.get(ride.vehicle_id)
                if v is not None and v.model == preference:
                    return ride
        return None


class RideService:
    def __init__(self):
        self.users = {}
        self.vehicles = {}
        self.rides = {}
        self.active_vehicles = {}  # regNumber -> rideId
        self.ride_counter = 0

    def add_user(self, name):
        if name in self.users:
            return
        self.users[name] = User(name, name)

    def add_vehicle(self, user_name, model, reg_number):
        if user_name not in self.users:
            return
        self.vehicles[reg_number] = Vehicle(reg_number, user_name, model, reg_number)

    def offer_ride(self, user_name, origin, dest, seats, vehicle_reg_number):
        if user_name not in self.users:
            return ""
        if vehicle_reg_number not in self.vehicles:
            return ""
        if self.vehicles[vehicle_reg_number].owner_id != user_name:
            return ""
        if vehicle_reg_number in self.active_vehicles:
            return ""
        self.ride_counter += 1
        ride_id = "RIDE-" + str(self.ride_counter)
        self.rides[ride_id] = Ride(ride_id, user_name, vehicle_reg_number, origin, dest, seats)
        self.active_vehicles[vehicle_reg_number] = ride_id
        self.users[user_name].rides_offered += 1
        return ride_id

    def select_ride(self, passenger_name, origin, dest, seats, strategy, preference=""):
        if passenger_name not in self.users:
            return ""
        candidates = []
        for ride in self.rides.values():
            if (ride.active and ride.origin == origin and ride.destination == dest
                    and ride.available_seats >= seats and ride.driver_id != passenger_name):
                candidates.append(ride)
        selected = strategy.select(candidates, seats, preference)
        if selected is not None:
            selected.available_seats -= seats
            self.users[passenger_name].rides_taken += 1
            return selected.id
        return ""

    def end_ride(self, ride_id):
        ride = self.rides.get(ride_id)
        if ride is None:
            return
        if not ride.active:
            return
        ride.active = False
        if ride.vehicle_id in self.active_vehicles:
            del self.active_vehicles[ride.vehicle_id]

    def has_user(self, name):
        return name in self.users

    def has_vehicle(self, reg):
        return reg in self.vehicles

    def has_ride(self, rid):
        return rid in self.rides

    def get_user(self, name):
        return self.users[name]

    def get_ride(self, rid):
        return self.rides[rid]

    def get_vehicles(self):
        return self.vehicles


def ride_simulate(ops):
    out = []
    svc = RideService()
    ride_slots = [""] * 32
    for op in ops:
        k = op.kind
        if k == "new":
            svc = RideService()
            ride_slots = [""] * 32
            out.append("ok")
        elif k == "add_user":
            svc.add_user(op.s1)
            out.append("ok")
        elif k == "add_veh":
            svc.add_vehicle(op.s1, op.s2, op.s3)
            out.append("ok")
        elif k == "offer":
            rid = svc.offer_ride(op.s1, op.s2, op.s3, op.i1, op.s4)
            if 0 <= op.i2 < len(ride_slots):
                ride_slots[op.i2] = rid
            out.append(rid)
        elif k == "ride_active":
            rid = ride_slots[op.i2]
            out.append("yes" if svc.has_ride(rid) and svc.get_ride(rid).active else "no")
        elif k == "ride_origin":
            rid = ride_slots[op.i2]
            out.append(svc.get_ride(rid).origin if svc.has_ride(rid) else "")
        elif k == "ride_dest":
            rid = ride_slots[op.i2]
            out.append(svc.get_ride(rid).destination if svc.has_ride(rid) else "")
        elif k == "ride_total":
            rid = ride_slots[op.i2]
            out.append(str(svc.get_ride(rid).total_seats) if svc.has_ride(rid) else "-1")
        elif k == "ride_avail":
            rid = ride_slots[op.i2]
            out.append(str(svc.get_ride(rid).available_seats) if svc.has_ride(rid) else "-1")
        elif k == "ride_driver":
            rid = ride_slots[op.i2]
            out.append(svc.get_ride(rid).driver_id if svc.has_ride(rid) else "")
        elif k == "select_mv":
            rid = svc.select_ride(op.s1, op.s2, op.s3, op.i1, MostVacantStrategy())
            if 0 <= op.i2 < len(ride_slots):
                ride_slots[op.i2] = rid
            out.append(rid)
        elif k == "select_pv":
            rid = svc.select_ride(op.s1, op.s2, op.s3, op.i1,
                                  PreferredVehicleStrategy(svc.get_vehicles()), op.s4)
            if 0 <= op.i2 < len(ride_slots):
                ride_slots[op.i2] = rid
            out.append(rid)
        elif k == "end":
            svc.end_ride(ride_slots[op.i2])
            out.append("ok")
        elif k == "end_id":
            svc.end_ride(op.s1)
            out.append("ok")
        elif k == "user_offered":
            out.append(str(svc.get_user(op.s1).rides_offered) if svc.has_user(op.s1) else "0")
        elif k == "user_taken":
            out.append(str(svc.get_user(op.s1).rides_taken) if svc.has_user(op.s1) else "0")
        elif k == "has_user":
            out.append("yes" if svc.has_user(op.s1) else "no")
        elif k == "has_vehicle":
            out.append("yes" if svc.has_vehicle(op.s1) else "no")
        elif k == "has_ride":
            out.append("yes" if svc.has_ride(ride_slots[op.i2]) else "no")
        else:
            out.append("unknown:" + k)
    return out
