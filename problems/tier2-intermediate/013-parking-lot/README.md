# Problem 013 — Parking Lot System

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Factory | **DSA:** HashMap + Heap
**Companies:** Amazon, Flipkart, Walmart, Salesforce | **Time:** 60 minutes

---

## Problem Statement

You are designing a multi-floor parking lot for a commercial building. The lot contains spots of three sizes (Small, Medium, Large) spread across multiple floors. Vehicles must be matched to the smallest compatible spot. On exit, the system calculates a fee using a pluggable pricing strategy. Entry and exit gates are tracked per transaction.

**Constraints:**
- Up to 10^3 spots across up to 10 floors
- Spot allocation uses nearest-first: lowest floor number, then spot ID order within that floor
- A vehicle may park in a larger spot if its exact size is full (e.g., motorcycle in MEDIUM)
- Fee calculation receives duration in seconds; strategies must handle fractional hours

---

## Base Requirement — Multi-Floor Parking with Vehicle-Spot Matching

Implement a `ParkingLot` that parks vehicles in the nearest compatible spot and issues a ticket on entry. On exit, `unparkVehicle` frees the spot and returns the fee (calculated by a default hourly strategy).

| Vehicle Type | Minimum Spot Size |
|---|---|
| MOTORCYCLE | SMALL |
| CAR | MEDIUM |
| TRUCK | LARGE |

**Example:**
```
// Floor 0: S1(SMALL), M1(MEDIUM). Floor 1: S2(SMALL)
parkVehicle({plate="KA01", type=CAR}, entryTime=0)
→  Ticket{ticketId="T1", spotId="M1", floor=0}

parkVehicle({plate="KA02", type=MOTORCYCLE}, entryTime=0)
→  Ticket{ticketId="T2", spotId="S1", floor=0}  // smallest compatible spot, floor 0 first

getAvailableSpots(SMALL)  →  1   // S2 on floor 1 still free
unparkVehicle("T1", exitTime=3600)  →  5.0  // 1 hour at default $5/hour rate
```

**Public methods:**
- `void addSpot(int floor, SpotSize size)`
- `Ticket* parkVehicle(const Vehicle& vehicle, long entryTime)`
- `double unparkVehicle(const string& ticketId, long exitTime)`
- `int getAvailableSpots(SpotSize size)`
- `int getAvailableSpotsByFloor(int floor, SpotSize size)`

---

## Extension 1 — Pluggable Pricing Strategies and Gate Management

Add swappable pricing strategies and gate tracking. Each entry/exit is associated with a specific gate. Adding a new pricing model must require zero changes to the `ParkingLot` class.

| Strategy | Rule |
|---|---|
| FlatRate | Fixed fee per visit, regardless of duration |
| Hourly | Rate per hour, duration rounded up to nearest hour |
| Tiered | First hour at base rate; hours 1–3 at mid rate; beyond 3 hours at high rate |

**Example:**
```
setPricingStrategy(new Hourly(rate=5.0))
addGate("G1", ENTRY), addGate("G2", EXIT)
parkVehicle({plate="KA03", type=CAR}, entryTime=0, gateId="G1")
→  Ticket{ticketId="T3", entryGate="G1", ...}

unparkVehicle("T3", exitTime=9000, gateId="G2")  // 2.5 hours
→  15.0   // Hourly: ceil(2.5) = 3 hours × $5 = $15
```

**Public methods:**
- `void setPricingStrategy(PricingStrategy* strategy)`
- `void addGate(const string& gateId, GateType type)`
- `Ticket* parkVehicle(const Vehicle& vehicle, long entryTime, const string& gateId)`
- `double unparkVehicle(const string& ticketId, long exitTime, const string& gateId)`
- `vector<string> getGates(GateType type)`

---

## Running Tests

```bash
./run-tests.sh 013-parking-lot cpp
```
