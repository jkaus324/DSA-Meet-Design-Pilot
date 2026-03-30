# Design Walkthrough — Ride-Sharing Application

> This file is the answer guide. Only read after you've attempted the problem.

---

## The Core Design Decisions

```
RideService
    ├── unordered_map<string, User> users           // keyed by name
    ├── unordered_map<string, Vehicle> vehicles      // keyed by regNumber
    ├── unordered_map<string, Ride> rides            // keyed by rideId
    ├── unordered_map<string, string> activeVehicles // vehicleReg → rideId
    └── RideSelectionStrategy* (Strategy — swappable selection logic)
            ├── MostVacantStrategy
            └── PreferredVehicleStrategy
```

**Why HashMaps everywhere?** Every lookup in a ride-sharing system needs to be fast:
- User lookup by name: O(1)
- Vehicle lookup by registration: O(1)
- Ride lookup by ID: O(1)
- Active vehicle check: O(1)

Using vectors would make these O(n) searches, which would be unacceptable at scale.

**Why Strategy for selection?** The selection algorithm is the only thing that varies. The process of filtering candidates (by origin, destination, available seats) is the same for every strategy. Only the final selection criterion changes.

---

## Reference Implementation

### Part 1 — Onboarding

```cpp
class RideService {
    unordered_map<string, User> users;
    unordered_map<string, Vehicle> vehicles;
    unordered_map<string, Ride> rides;
    unordered_map<string, string> activeVehicles;  // regNumber → rideId
    int rideCounter = 0;

public:
    void addUser(const string& name) {
        if (users.find(name) != users.end()) return;  // already exists
        users[name] = {name, name, 0, 0};
    }

    void addVehicle(const string& userName, const string& model, const string& regNumber) {
        if (users.find(userName) == users.end()) return;
        vehicles[regNumber] = {regNumber, userName, model, regNumber};
    }

    string offerRide(const string& userName, const string& origin,
                     const string& dest, int seats, const string& vehicleRegNumber) {
        // Validate user and vehicle exist
        if (users.find(userName) == users.end()) return "";
        if (vehicles.find(vehicleRegNumber) == vehicles.end()) return "";
        // Vehicle must belong to user
        if (vehicles[vehicleRegNumber].ownerId != userName) return "";
        // Vehicle must not have an active ride
        if (activeVehicles.find(vehicleRegNumber) != activeVehicles.end()) return "";

        string rideId = "RIDE-" + to_string(++rideCounter);
        rides[rideId] = {rideId, userName, vehicleRegNumber, origin, dest, seats, seats, true};
        activeVehicles[vehicleRegNumber] = rideId;
        users[userName].ridesOffered++;
        return rideId;
    }
};
```

### Part 2 — Strategy Selection

```cpp
class RideSelectionStrategy {
public:
    virtual Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                         const string& preference) = 0;
    virtual ~RideSelectionStrategy() = default;
};

class MostVacantStrategy : public RideSelectionStrategy {
public:
    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        Ride* best = nullptr;
        for (auto* r : candidates) {
            if (r->availableSeats >= seatsNeeded) {
                if (!best || r->availableSeats > best->availableSeats) {
                    best = r;
                }
            }
        }
        return best;
    }
};

class PreferredVehicleStrategy : public RideSelectionStrategy {
    unordered_map<string, Vehicle>& vehicleStore;
public:
    PreferredVehicleStrategy(unordered_map<string, Vehicle>& vs) : vehicleStore(vs) {}

    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        for (auto* r : candidates) {
            if (r->availableSeats >= seatsNeeded) {
                if (vehicleStore[r->vehicleId].model == preference) {
                    return r;
                }
            }
        }
        return nullptr;
    }
};
```

The `selectRide()` method in `RideService`:

```cpp
string selectRide(const string& passengerName, const string& origin,
                  const string& dest, int seats, RideSelectionStrategy* strategy,
                  const string& preference = "") {
    // 1. Filter rides by origin + destination + active
    vector<Ride*> candidates;
    for (auto& [id, ride] : rides) {
        if (ride.active && ride.origin == origin && ride.destination == dest
            && ride.availableSeats >= seats && ride.driverId != passengerName) {
            candidates.push_back(&ride);
        }
    }

    // 2. Delegate selection to strategy
    Ride* selected = strategy->select(candidates, seats, preference);
    if (!selected) return "";

    // 3. Update state
    selected->availableSeats -= seats;
    users[passengerName].ridesTaken++;
    return selected->id;
}
```

### Part 3 — End Rides and Stats

```cpp
void endRide(const string& rideId) {
    if (rides.find(rideId) == rides.end()) return;
    Ride& ride = rides[rideId];
    if (!ride.active) return;  // already ended

    ride.active = false;
    activeVehicles.erase(ride.vehicleId);
}

void printRideStats() const {
    for (const auto& [name, user] : users) {
        cout << "User: " << name
             << " — Offered: " << user.ridesOffered
             << ", Taken: " << user.ridesTaken << endl;
    }
}
```

---

## Key Structural Decisions

### Why track `activeVehicles` separately?

Checking "is this vehicle in an active ride?" by scanning all rides is O(n). A dedicated `activeVehicles` HashMap makes it O(1). When a ride ends, we erase the entry. This is a classic space-time tradeoff.

### Why does PreferredVehicleStrategy need access to the vehicle store?

The ride only stores `vehicleId`, not the vehicle model. The strategy needs to resolve the model to compare with the passenger's preference. Two approaches:

1. **Store model in Ride** (denormalization) — simple but redundant
2. **Pass vehicle store reference to strategy** — keeps data normalized but couples strategy to store

In an interview, option 2 is preferred because it demonstrates understanding of dependency injection and separation of concerns.

### Why can't a passenger select their own offered ride?

The `selectRide()` filter excludes rides where `driverId == passengerName`. This prevents a user from booking their own ride, which is a common edge case interviewers check.

---

## Pattern Interaction Diagram

```
User calls selectRide("Deepa", "Bangalore", "Mysore", 1, mostVacantStrategy)
    │
    ├─ Filter: find all active rides where origin="Bangalore", dest="Mysore"
    │   └─ Returns candidate rides: [RIDE-1 (3 seats), RIDE-3 (1 seat)]
    │
    ├─ Strategy: mostVacantStrategy->select(candidates, 1, "")
    │   └─ Returns RIDE-1 (has most available seats)
    │
    └─ Update:
        ├─ RIDE-1.availableSeats -= 1
        └─ Deepa.ridesTaken += 1
```

---

## Interview Tips

1. **Start with data models.** Define User, Vehicle, Ride structs first. Interviewers want to see you think about data before behavior.
2. **Use HashMaps explicitly.** Say "I'll use an unordered_map keyed by regNumber for O(1) vehicle lookup." This shows DSA awareness.
3. **Introduce Strategy only when Part 2 is revealed.** Don't pre-engineer. Show that your Part 1 design survives the extension.
4. **Handle edge cases vocally.** Mention: "What if the vehicle is already in an active ride?" before the interviewer asks.
5. **Track stats incrementally.** Don't recompute stats from ride history — increment counters when events happen.

---

## Common Interview Follow-ups

- *"How would you handle ride cancellation?"* — Add a `cancelRide(rideId)` that restores `availableSeats`, decrements counters, and marks ride inactive.
- *"How would you add fare calculation?"* — Another Strategy: `FareStrategy` with implementations like FlatFare, DistanceBasedFare, SurgeFare.
- *"How would you find rides with intermediate stops?"* — Model routes as a graph. Use BFS to find rides whose path includes both origin and destination in the correct order.
- *"How would you handle concurrent ride requests?"* — Optimistic locking on `availableSeats` or a mutex per ride. First-come-first-served with atomic decrement.
