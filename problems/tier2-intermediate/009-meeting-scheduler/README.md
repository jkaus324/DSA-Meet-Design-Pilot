# Problem 009 — Meeting Room Scheduler

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer | **DSA:** Interval Checking + HashMap + Priority Queue
**Companies:** Flipkart, Razorpay, Groww, Microsoft | **Time:** 60 minutes

---

## Problem Statement

You are building a meeting room booking system for a corporate office. The system manages rooms with varying capacities and AV equipment. Users book rooms by specifying a time window; the system detects conflicts and rejects double-bookings. The office manager can switch between room allocation policies without rewriting the scheduler.

**Constraints:**
- Up to 10^4 bookings per room
- Time is represented as minutes since midnight (0–1439)
- Two meetings overlap if `start1 < end2 && start2 < end1`
- Rooms have a fixed capacity and optional AV equipment flag

---

## Base Requirement — Room Booking with Conflict Detection

Implement a `MeetingScheduler` that manages rooms and meetings. When a booking is requested, verify the room is free for the entire requested window. Reject the booking if any existing meeting on that room overlaps.

**Example:**
```
addRoom("R1", capacity=10, hasAV=true)
bookMeeting({id="M1", roomId="R1", start=540, end=600})  →  true   // 9:00–10:00 AM booked
bookMeeting({id="M2", roomId="R1", start=570, end=630})  →  false  // overlaps M1
isAvailable("R1", 600, 660)                               →  true   // 10:00–11:00 AM is free
getRoomSchedule("R1")                                     →  [M1]
```

**Public methods:**
- `void addRoom(const Room& room)`
- `bool bookMeeting(const Meeting& meeting)`
- `vector<Meeting> getRoomSchedule(const string& roomId)`
- `bool isAvailable(const string& roomId, int startTime, int endTime)`

---

## Extension 1 — Multiple Allocation Strategies

Add automatic room selection. Instead of the user specifying a room, the scheduler picks one based on a pluggable strategy. Adding a new strategy must require zero changes to the scheduler class.

| Strategy | Rule |
|---|---|
| FirstAvailable | First free room by room ID lexicographic order |
| BestFit | Smallest capacity room (by attendee count) that is free |
| PriorityBased | Prefer AV-equipped rooms; among those, smallest capacity that fits |

**Example:**
```
// 3 rooms: R1 (cap=20, AV), R2 (cap=8, no AV), R3 (cap=8, AV)
// Strategy: BestFit, attendeeCount=6
bookWithStrategy("M3", "Design Review", start=480, end=540, attendees=6)
→  "R2"   // smallest room that fits 6 people and is free
```

**Public method:**
- `string bookWithStrategy(const string& meetingId, const string& title, int startTime, int endTime, int attendeeCount)`

---

## Extension 2 — Attendee Notifications

When a meeting is booked, cancelled, or rescheduled, all subscribed observers are notified. The scheduler must not know about email, SMS, or Slack — it only fires events through the observer interface.

**Example:**
```
subscribeAttendee("M1", &emailNotifier)
cancelMeeting("M1")           // emailNotifier.onMeetingCancelled(M1) is called
reschedule("M1", 660, 720)    // emailNotifier.onMeetingRescheduled(old, new) is called
```

**Public methods:**
- `void subscribeAttendee(const string& meetingId, MeetingObserver* observer)`
- `bool cancelMeeting(const string& meetingId)`
- `bool rescheduleMeeting(const string& meetingId, int newStart, int newEnd)`

---

## Running Tests

```bash
./run-tests.sh 009-meeting-scheduler cpp
```
