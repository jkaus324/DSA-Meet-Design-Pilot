# Problem 013 — Parking Lot System

**Tier:** 2 (Intermediate) | **Pattern:** Factory + Strategy + Singleton | **DSA:** HashMap + Queue
**Companies:** Salesforce | **Time:** 45 minutes

---

## Problem Statement

You're designing a multi-floor parking lot for a commercial building. The system must handle different vehicle types, allocate appropriate spots, track parking duration for billing, and support multiple pricing strategies.

**Your task:** Design and implement a `ParkingLot` that parks vehicles in compatible spots, calculates fees based on pluggable pricing strategies, and manages entry/exit gates.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. A motorcycle takes a small spot, a car takes a medium spot, and a truck takes a large spot. Should the parking logic know about every vehicle type? Or should vehicle and spot creation be delegated to a factory?
2. When the business wants to switch from flat-rate to hourly pricing, should you modify the parking lot class? Or should the pricing algorithm be a separate, swappable component?
3. Should there be multiple parking lot instances, or should the system guarantee exactly one? What happens if two threads create separate instances with different state?

**The key insight:** The **Factory** pattern handles object creation (vehicles and spots) without coupling the lot to specific types. The **Strategy** pattern decouples pricing from the lot. The **Singleton** pattern (optional) ensures a single lot instance with consistent state.

---

## Data Structures

```cpp
enum class VehicleType {
    MOTORCYCLE,
    CAR,
    TRUCK
};

enum class SpotSize {
    SMALL,      // fits motorcycle
    MEDIUM,     // fits car
    LARGE       // fits truck
};

struct Vehicle {
    string licensePlate;
    VehicleType type;
};

struct ParkingSpot {
    string spotId;
    int floor;
    SpotSize size;
    bool isOccupied;
    string vehicleLicensePlate;  // empty if unoccupied
};

struct Ticket {
    string ticketId;
    string licensePlate;
    string spotId;
    int floor;
    long entryTime;  // epoch seconds (use time(nullptr) or pass explicitly)
};
```

---

## Part 1

**Base requirement — Multi-floor parking with vehicle-spot matching**

Implement a `ParkingLot` with multiple floors, each containing spots of various sizes. Vehicles must be parked in the nearest compatible spot (smallest floor number first, then smallest spot ID). Unparking returns a ticket with duration.

**Compatibility rule:**
| Vehicle Type | Minimum Spot Size |
|---|---|
| MOTORCYCLE | SMALL |
| CAR | MEDIUM |
| TRUCK | LARGE |

A vehicle can park in a spot of its minimum size or larger (e.g., a motorcycle can park in a MEDIUM spot if no SMALL spots are available).

**Entry points (tests will call these):**
```cpp
Ticket* parkVehicle(const Vehicle& vehicle, long entryTime);
double unparkVehicle(const string& ticketId, long exitTime);
int getAvailableSpots(SpotSize size);
int getAvailableSpotsByFloor(int floor, SpotSize size);
```

**What to implement:**
```cpp
// Factory for creating spots
class SpotFactory {
public:
    static ParkingSpot createSpot(const string& spotId, int floor, SpotSize size);
};

class ParkingLot {
    vector<vector<ParkingSpot>> floors;  // floors[i] = spots on floor i
    unordered_map<string, Ticket> activeTickets;  // ticketId -> ticket
    int nextTicketId;
public:
    ParkingLot(int numFloors);
    void addSpot(int floor, SpotSize size);
    Ticket* parkVehicle(const Vehicle& vehicle, long entryTime);
    double unparkVehicle(const string& ticketId, long exitTime);
    int getAvailableSpots(SpotSize size);
    int getAvailableSpotsByFloor(int floor, SpotSize size);
};
```

**Design goal:** The factory encapsulates spot creation. Vehicle-to-spot matching uses a clear compatibility function. Spot allocation scans floors in order for the nearest available compatible spot.

---

## Part 2

**Extension — Pluggable pricing strategies and gate management**

The parking lot owner wants different pricing models and entry/exit gate tracking.

| Strategy | Rule |
|----------|------|
| FlatRate | Fixed fee regardless of duration (e.g., $10 per visit) |
| Hourly | Fee per hour (rounded up). E.g., $5/hour, 2.5 hours = $15 |
| Tiered | First hour: base rate. 1-3 hours: mid rate/hour. 3+ hours: high rate/hour |

**Gate tracking:** Each gate has an ID and a type (ENTRY or EXIT). The system tracks which gate a vehicle entered/exited through.

**New entry points:**
```cpp
void setPricingStrategy(PricingStrategy* strategy);
void addGate(const string& gateId, GateType type);
Ticket* parkVehicle(const Vehicle& vehicle, long entryTime, const string& gateId);
double unparkVehicle(const string& ticketId, long exitTime, const string& gateId);
vector<string> getGates(GateType type);
```

**What to implement:**
```cpp
enum class GateType {
    ENTRY,
    EXIT
};

struct Gate {
    string gateId;
    GateType type;
};

class PricingStrategy {
public:
    virtual double calculateFee(long durationSeconds) = 0;
    virtual ~PricingStrategy() = default;
};

class FlatRate : public PricingStrategy { ... };   // fixed fee
class Hourly : public PricingStrategy { ... };     // rate * ceil(hours)
class Tiered : public PricingStrategy { ... };     // tiered rates by duration bracket
```

**Hint:** The pricing strategy receives duration in seconds and returns the fee. The parking lot delegates fee calculation entirely to the strategy. Gate tracking extends the Ticket struct with entry/exit gate IDs.

---

## Running Tests

```bash
./run-tests.sh 013-parking-lot cpp
```
