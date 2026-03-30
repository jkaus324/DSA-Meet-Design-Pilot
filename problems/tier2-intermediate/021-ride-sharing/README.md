# Problem 021 — Ride-Sharing Application

**Tier:** 2 (Intermediate) | **Patterns:** Strategy, Factory | **DSA:** HashMap, Graph, BFS
**Companies:** Flipkart | **Time:** 90 minutes

---

## Problem Statement

Build a ride-sharing application that allows users to offer and select shared rides. The system must:

1. Onboard users and vehicles with in-memory storage
2. Allow users to offer rides with origin, destination, available seats, and vehicle
3. Allow passengers to search and select rides using pluggable selection strategies
4. Track ride statistics per user

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**The naive approach:** A single `RideManager` class with `if-else` chains for ride selection:

```cpp
Ride* selectRide(string origin, string dest, int seats, string strategy) {
    if (strategy == "preferred_vehicle") {
        // search by vehicle model match
    } else if (strategy == "most_vacant") {
        // search by most available seats
    }
}
```

**Why this breaks:** Adding a new selection strategy (e.g., "cheapest fare") requires modifying the `selectRide()` method. Every selection algorithm is tangled in one function.

**The pattern approach:**
- **Strategy**: Ride selection algorithms (PreferredVehicle, MostVacant) are encapsulated as separate classes implementing a common interface
- **Factory**: Vehicle creation can be delegated to a factory, but the key Strategy usage is in ride selection

**The DSA angle:** Users, vehicles, and rides are stored in HashMaps for O(1) lookup by ID. Ride search filters by origin and destination — this is essentially a graph adjacency lookup. BFS can model reachability if routes are extended.

---

## Data Structures

```cpp
struct User {
    string id;
    string name;
    int ridesOffered;
    int ridesTaken;
};

struct Vehicle {
    string id;
    string ownerId;
    string model;      // "Swift", "Activa", "XUV"
    string regNumber;
};

struct Ride {
    string id;
    string driverId;
    string vehicleId;
    string origin;
    string destination;
    int totalSeats;
    int availableSeats;
    bool active;
};
```

---

## Part 1

**Base requirement — User, vehicle, and ride onboarding**

Implement a `RideService` that supports user registration, vehicle registration, and ride offering.

**Rules:**
- Each user has a unique name (used as ID)
- Each vehicle belongs to exactly one user
- A vehicle can only have **one active ride** at a time — offering a second ride with the same vehicle while one is active should fail
- Offering a ride increments the driver's `ridesOffered` count
- `availableSeats` starts equal to `totalSeats`

**Entry points (tests will call these):**
```cpp
void RideService::addUser(const string& name);
void RideService::addVehicle(const string& userName, const string& model, const string& regNumber);
string RideService::offerRide(const string& userName, const string& origin,
                              const string& dest, int seats, const string& vehicleRegNumber);
// returns rideId or empty string on failure
```

**What to implement:**
```cpp
class RideService {
    unordered_map<string, User> users;
    unordered_map<string, Vehicle> vehicles;       // keyed by regNumber
    unordered_map<string, Ride> rides;             // keyed by rideId
    int rideCounter;

public:
    RideService();
    void addUser(const string& name);
    void addVehicle(const string& userName, const string& model, const string& regNumber);
    string offerRide(const string& userName, const string& origin,
                     const string& dest, int seats, const string& vehicleRegNumber);
};
```

---

## Part 2

**Extension 1 — Pluggable ride selection strategies**

Passengers can now search for rides between an origin and destination. The system must support multiple selection strategies:

| Strategy | Rule |
|----------|------|
| **MostVacant** | Select the ride with the **most available seats** |
| **PreferredVehicle** | Select the ride whose vehicle model matches the passenger's preference |

**Rules:**
- Only active rides matching origin and destination are considered
- The ride must have enough available seats for the request (1 or 2 seats)
- If no ride matches, return empty/null
- Selecting a ride decrements `availableSeats` and increments the passenger's `ridesTaken` count
- If `availableSeats` reaches 0, the ride is no longer selectable (but still active until ended)
- Adding a new strategy must require **zero changes** to `RideService`

**New entry points:**
```cpp
string RideService::selectRide(const string& passengerName, const string& origin,
                               const string& dest, int seats,
                               RideSelectionStrategy* strategy);
// returns rideId or empty string if no match
```

**Strategy interface:**
```cpp
class RideSelectionStrategy {
public:
    virtual Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                         const string& preference) = 0;
    virtual ~RideSelectionStrategy() = default;
};

class MostVacantStrategy : public RideSelectionStrategy { ... };
class PreferredVehicleStrategy : public RideSelectionStrategy { ... };
```

**Design challenge:** `PreferredVehicleStrategy` needs to know the vehicle model. Where does that information come from? The strategy receives candidate rides — each ride has a `vehicleId`. How do you resolve the vehicle model without coupling the strategy to the data store?

---

## Part 3

**Extension 2 — End rides and print statistics**

The system now needs to track ride lifecycle and provide per-user statistics.

**Rules:**
- `endRide(rideId)` marks a ride as inactive, freeing the vehicle for future rides
- `printRideStats()` outputs each user's offered and taken counts
- Only the driver (who offered the ride) can end it
- Ending an already-ended ride should be a no-op

**New entry points:**
```cpp
void RideService::endRide(const string& rideId);
vector<pair<string, pair<int, int>>> RideService::getRideStats() const;
// returns vector of {userName, {ridesOffered, ridesTaken}}
void RideService::printRideStats() const;
```

**Output format for printRideStats:**
```
User: Rohan — Offered: 2, Taken: 1
User: Deepa — Offered: 0, Taken: 3
```

**Design challenge:** Once a ride ends, the vehicle becomes available again. How do you track which vehicles are currently in use? A HashMap from vehicleId to active rideId provides O(1) lookup.

---

## Running Tests

```bash
./run-tests.sh 021-ride-sharing cpp
```
