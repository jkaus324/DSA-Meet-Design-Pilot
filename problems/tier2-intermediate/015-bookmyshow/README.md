# Problem 015 — BookMyShow Ticket Booking System

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer + State | **DSA:** HashMap + 2D Matrix
**Companies:** BookMyShow, Flipkart, Amazon, DoorDash | **Time:** 60 minutes

---

## Problem Statement

You are building a movie ticket booking system. Theaters in multiple cities host shows on 2D seating grids. Users search for movies in their city, view available seats, and book them. The system must prevent double-booking atomically and support temporary seat locking with TTL-based auto-release for a checkout flow.

**Constraints:**
- Seats identified by (row, col) pairs, 0-indexed
- Booking is all-or-nothing: if any requested seat is unavailable, no seats are booked
- Seat lock TTL is configurable; time is passed explicitly (tests control the clock)
- An expired lock means the seat is AVAILABLE, regardless of stored status

---

## Base Requirement — Theater, Show, and Seat Booking

Implement a `BookingSystem` that registers theaters and shows, returns available seats for a show, and books seats. Seats are modeled as a 2D grid with states: AVAILABLE or BOOKED.

**Example:**
```
addTheater("T1", name="PVR", city="Bangalore")
addShow("S1", theaterId="T1", movie="Dune", time="18:00", rows=5, cols=10)

searchMovies("Bangalore")       →  ["Dune"]
getAvailableSeats("S1").size()  →  50

bookSeats("B1", showId="S1", seats=[(0,0),(0,1)], userId="u1")  →  true
bookSeats("B2", showId="S1", seats=[(0,0),(0,2)], userId="u2")  →  false  // (0,0) already booked
getAvailableSeats("S1").size()  →  48
```

**Public methods:**
- `void addTheater(const string& theaterId, const string& name, const string& city)`
- `void addShow(const string& showId, const string& theaterId, const string& movie, const string& time, int rows, int cols)`
- `vector<string> searchMovies(const string& city)`
- `vector<pair<int,int>> getAvailableSeats(const string& showId)`
- `bool bookSeats(const string& bookingId, const string& showId, const vector<pair<int,int>>& seats, const string& userId)`

---

## Extension 1 — Seat Locking with TTL

Add a seat-locking step to the checkout flow. When a user selects seats, lock them temporarily. Locked seats cannot be booked by others. If the user does not confirm within the TTL, the lock auto-expires and seats return to AVAILABLE.

**Lock lifecycle:**
1. `lockSeats(...)` → seats become LOCKED; returns a `lockId`
2. `confirmBooking(lockId, currentTime)` → seats become BOOKED (only if lock has not expired)
3. TTL expires → seats treated as AVAILABLE on next availability check
4. `releaseLock(lockId, currentTime)` → seats immediately return to AVAILABLE

**Example:**
```
lockId = lockSeats("S1", seats=[(1,0),(1,1)], userId="u3", ttlMinutes=5, currentTime=1000)
// At currentTime=1100: lock still active (1100 < 1000 + 300)
confirmBooking(lockId, currentTime=1100)  →  true

lockId2 = lockSeats("S1", seats=[(2,0)], userId="u4", ttlMinutes=5, currentTime=2000)
confirmBooking(lockId2, currentTime=2500)  →  false  // expired: 2500 >= 2000 + 300
getAvailableSeats("S1") includes (2,0)  →  true
```

**Public methods:**
- `string lockSeats(const string& showId, const vector<pair<int,int>>& seats, const string& userId, int ttlMinutes, long currentTime)`
- `bool confirmBooking(const string& lockId, long currentTime)`
- `bool releaseLock(const string& lockId, long currentTime)`

---

## Running Tests

```bash
./run-tests.sh 015-bookmyshow cpp
```
