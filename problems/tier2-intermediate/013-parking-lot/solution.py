"""Parking Lot — multi-floor + spot matching + pricing strategies."""

import math


class ParkOp:
    def __init__(self, kind, s1="", s2="", s3="", i1=0, i2=0, i3=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.i1 = i1
        self.i2 = i2
        self.i3 = i3


# Constants
M_TYPE = "MOTORCYCLE"
C_TYPE = "CAR"
T_TYPE = "TRUCK"

SMALL = 0
MEDIUM = 1
LARGE = 2


class Vehicle:
    def __init__(self, licensePlate, type):
        self.licensePlate = licensePlate
        self.type = type


class ParkingSpot:
    def __init__(self, spotId, floor, size):
        self.spotId = spotId
        self.floor = floor
        self.size = size
        self.isOccupied = False
        self.vehicleLicensePlate = ""


class Ticket:
    def __init__(self):
        self.ticketId = ""
        self.licensePlate = ""
        self.spotId = ""
        self.floor = 0
        self.entryTime = 0
        self.entryGateId = ""
        self.exitGateId = ""


class Gate:
    def __init__(self, gateId, type):
        self.gateId = gateId
        self.type = type


def _min_spot_size(vtype):
    if vtype == M_TYPE:
        return SMALL
    if vtype == C_TYPE:
        return MEDIUM
    return LARGE


def _is_compatible(spot_size, min_required):
    return spot_size >= min_required


class FlatRate:
    def __init__(self, fee):
        self.fee = fee

    def calculateFee(self, durationSeconds):
        return self.fee


class Hourly:
    def __init__(self, rate):
        self.rate = rate

    def calculateFee(self, durationSeconds):
        hours = math.ceil(durationSeconds / 3600.0)
        return self.rate * hours


class Tiered:
    def __init__(self, base, mid, high):
        self.base = base
        self.mid = mid
        self.high = high

    def calculateFee(self, durationSeconds):
        hours = math.ceil(durationSeconds / 3600.0)
        if hours <= 1:
            return self.base
        if hours <= 3:
            return self.base + self.mid * (hours - 1)
        return self.base + self.mid * 2 + self.high * (hours - 3)


class ParkingLot:
    def __init__(self, numFloors):
        self.floors = [[] for _ in range(numFloors)]
        self.activeTickets = {}
        self.gates = []
        self.strategy = None
        self.nextTicketId = 1

    def addSpot(self, floor, size):
        if floor < 0 or floor >= len(self.floors):
            return
        spotId = "F" + str(floor) + "S" + str(len(self.floors[floor]))
        self.floors[floor].append(ParkingSpot(spotId, floor, size))

    def setPricingStrategy(self, s):
        self.strategy = s

    def addGate(self, gateId, type):
        self.gates.append(Gate(gateId, type))

    def getGates(self, type):
        return [g.gateId for g in self.gates if g.type == type]

    def parkVehicle(self, vehicle, entryTime, gateId=""):
        minSize = _min_spot_size(vehicle.type)
        for f, spots in enumerate(self.floors):
            for spot in spots:
                if not spot.isOccupied and _is_compatible(spot.size, minSize):
                    spot.isOccupied = True
                    spot.vehicleLicensePlate = vehicle.licensePlate
                    tid = "T" + str(self.nextTicketId)
                    self.nextTicketId += 1
                    t = Ticket()
                    t.ticketId = tid
                    t.licensePlate = vehicle.licensePlate
                    t.spotId = spot.spotId
                    t.floor = f
                    t.entryTime = entryTime
                    t.entryGateId = gateId
                    t.exitGateId = ""
                    self.activeTickets[tid] = t
                    return t
        return None

    def unparkVehicle(self, ticketId, exitTime, gateId=""):
        if ticketId not in self.activeTickets:
            return -1.0
        ticket = self.activeTickets[ticketId]
        ticket.exitGateId = gateId
        for floorSpots in self.floors:
            for spot in floorSpots:
                if spot.spotId == ticket.spotId and spot.isOccupied:
                    spot.isOccupied = False
                    spot.vehicleLicensePlate = ""
                    break
        duration = exitTime - ticket.entryTime
        if self.strategy:
            fee = self.strategy.calculateFee(duration)
        else:
            fee = float(duration)
        del self.activeTickets[ticketId]
        return fee

    def getAvailableSpots(self, size):
        count = 0
        for floorSpots in self.floors:
            for spot in floorSpots:
                if not spot.isOccupied and spot.size == size:
                    count += 1
        return count

    def getAvailableSpotsByFloor(self, floor, size):
        if floor < 0 or floor >= len(self.floors):
            return 0
        count = 0
        for spot in self.floors[floor]:
            if not spot.isOccupied and spot.size == size:
                count += 1
        return count


def _size_from(s):
    if s == "S" or s == "small":
        return SMALL
    if s == "M" or s == "medium":
        return MEDIUM
    return LARGE


def _vtype_from(s):
    if s == "M" or s == "moto" or s == "motorcycle":
        return M_TYPE
    if s == "C" or s == "car":
        return C_TYPE
    return T_TYPE


def _gate_from(s):
    return "ENTRY" if s == "entry" else "EXIT"


def _fee_to_str(f):
    if f < 0:
        return "-1"
    return f"{f:.2f}"


def parking_simulate(ops):
    out = []
    lot = None
    tickets = [""] * 16
    snaps = [{"id": "", "floor": -1, "spotId": "", "entryGate": ""} for _ in range(16)]
    for op in ops:
        k = op.kind
        if k == "new":
            lot = ParkingLot(op.i1)
            tickets = [""] * 16
            snaps = [{"id": "", "floor": -1, "spotId": "", "entryGate": ""} for _ in range(16)]
            out.append("ok")
        elif k == "add_spot":
            lot.addSpot(op.i1, _size_from(op.s1))
            out.append("ok")
        elif k == "add_gate":
            lot.addGate(op.s1, _gate_from(op.s2))
            out.append("ok")
        elif k == "gates_count":
            out.append(str(len(lot.getGates(_gate_from(op.s1)))))
        elif k == "gate_at":
            g = lot.getGates(_gate_from(op.s1))
            out.append(g[op.i1] if 0 <= op.i1 < len(g) else "")
        elif k == "set_pricing":
            p = None
            if op.s1 == "flat":
                p = FlatRate(float(op.i1))
            elif op.s1 == "hourly":
                p = Hourly(float(op.i1))
            elif op.s1 == "tiered":
                p = Tiered(float(op.i1), float(op.i2), float(op.i3))
            if p:
                lot.setPricingStrategy(p)
            out.append("ok")
        elif k == "park":
            v = Vehicle(op.s1, _vtype_from(op.s2))
            t = lot.parkVehicle(v, op.i1, op.s3)
            if 0 <= op.i2 < len(tickets):
                if t:
                    tickets[op.i2] = t.ticketId
                    snaps[op.i2] = {"id": t.ticketId, "floor": t.floor, "spotId": t.spotId, "entryGate": t.entryGateId}
                else:
                    tickets[op.i2] = ""
                    snaps[op.i2] = {"id": "", "floor": -1, "spotId": "", "entryGate": ""}
            out.append(t.ticketId if t else "")
        elif k == "ticket_at":
            out.append(tickets[op.i1] if 0 <= op.i1 < len(tickets) else "")
        elif k == "ticket_floor":
            out.append(str(snaps[op.i1]["floor"]) if 0 <= op.i1 < len(snaps) else "-1")
        elif k == "ticket_spot_id":
            out.append(snaps[op.i1]["spotId"] if 0 <= op.i1 < len(snaps) else "")
        elif k == "ticket_entry":
            out.append(snaps[op.i1]["entryGate"] if 0 <= op.i1 < len(snaps) else "")
        elif k == "unpark":
            tid = tickets[op.i1] if 0 <= op.i1 < len(tickets) else ""
            fee = lot.unparkVehicle(tid, op.i2, op.s1)
            out.append(_fee_to_str(fee))
        elif k == "unpark_id":
            fee = lot.unparkVehicle(op.s1, op.i2, op.s2)
            out.append(_fee_to_str(fee))
        elif k == "available":
            out.append(str(lot.getAvailableSpots(_size_from(op.s1))))
        elif k == "available_floor":
            out.append(str(lot.getAvailableSpotsByFloor(op.i1, _size_from(op.s1))))
        else:
            out.append("unknown:" + k)
    return out
