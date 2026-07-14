# Problem 016 — Amazon Locker System

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + State | **DSA:** HashMap + Queue
**Companies:** Amazon | **Time:** 60 minutes

---

## Problem Statement

You are building an Amazon Locker package management system. Delivery agents deposit packages into lockers; customers retrieve them using a unique pickup code. The system allocates the smallest fitting locker, generates a pickup code on deposit, and frees the locker on successful retrieval. Pickup codes expire after a configurable time, and multiple notification channels are alerted on deposit and expiry.

**Constraints:**
- Three locker sizes: SMALL, MEDIUM, LARGE
- A smaller package may use a larger locker if its exact size is fully occupied
- Pickup codes are unique strings generated at deposit time
- Expiry check is triggered explicitly by calling `checkExpired(currentTime)`

---

## Base Requirement — Core Locker Allocation and Retrieval

Implement a `LockerSystem` that allocates the smallest fitting locker for a package and supports retrieval by pickup code. If no locker can fit the package, deposit fails.

**Allocation rule:** Try SMALL first, then MEDIUM, then LARGE. Use the smallest size that has an available locker.

| Package Size | Try sizes in order |
|---|---|
| SMALL | SMALL → MEDIUM → LARGE |
| MEDIUM | MEDIUM → LARGE |
| LARGE | LARGE only |

**Example:**
```
addLocker("L1", SMALL), addLocker("L2", SMALL), addLocker("L3", MEDIUM)

depositPackage("PKG1", SMALL)   →  "CODE-001"   // allocated L1
depositPackage("PKG2", SMALL)   →  "CODE-002"   // allocated L2
depositPackage("PKG3", SMALL)   →  "CODE-003"   // no SMALL left, uses L3 (MEDIUM)

retrievePackage("CODE-001")     →  true   // L1 freed
retrievePackage("CODE-001")     →  false  // code already used
depositPackage("PKG4", SMALL)   →  "CODE-004"   // L1 available again
```

**Public methods:**
- `void addLocker(const string& lockerId, LockerSize size)`
- `string depositPackage(const string& packageId, LockerSize size)`
- `bool retrievePackage(const string& code)`

---

## Extension 1 — Code Expiry and Notifications

Pickup codes now expire after a configurable number of hours. Expired packages are marked for return and their lockers are freed. Registered notification channels are alerted on deposit and on expiry.

**Expiry rules:**
- Each code records its creation timestamp
- `checkExpired(currentTime)` scans all active codes, expires any that exceeded the time limit, frees the locker, and returns the list of expired package IDs
- Notifications fire on deposit: `"Package <id> deposited. Code: <code>"`
- Notifications fire on expiry: `"Package <id> expired. Locker freed."`

**Example:**
```
setCodeExpiry(hours=24)
addNotificationChannel(&emailChannel)
depositPackage("PKG5", MEDIUM)  →  "CODE-005"
// emailChannel.notify("PKG5", "Package PKG5 deposited. Code: CODE-005") fires

checkExpired(currentTime = depositTime + 90000)  // 25 hours later
→  ["PKG5"]
// emailChannel.notify("PKG5", "Package PKG5 expired. Locker freed.") fires
```

**Public methods:**
- `void setCodeExpiry(int hours)`
- `vector<string> checkExpired(long currentTime)`
- `void addNotificationChannel(NotificationChannel* channel)`

---

## Running Tests

```bash
./run-tests.sh 016-amazon-locker cpp
```
