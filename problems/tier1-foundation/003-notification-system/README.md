# Problem 003 — Notification System

**Tier:** 1 (Foundation) | **Pattern:** Observer | **DSA:** Queue
**Companies:** Flipkart, Swiggy | **Time:** 45 minutes

---

## Problem Statement

You're building the notification backbone for an e-commerce platform. When an **order event** occurs (order placed, order shipped, order delivered), the system must notify multiple channels simultaneously:
- Email service
- SMS service
- Push notification service
- In-app notification service

New notification channels should be addable **without changing the event publishing code**.

---

## Before You Code

> Read this section carefully.

**Ask yourself:**
1. What's the relationship here? One event → many notifications. Classic one-to-many.
2. What's the naive approach? Call `sendEmail()`, `sendSMS()`, `sendPush()` explicitly inside the order handler. What happens when Slack notifications are added?
3. How does Observer solve this? The order system publishes events. Notification handlers *subscribe* to events they care about. The publisher doesn't know who's listening.

**The key insight:** The order system should not import or know about email/SMS/push services. Decoupling is the goal.

---

## Data Structures

```cpp
struct Notification {
    string userId;
    string message;
    string channel;  // "email", "sms", "push"
};

struct User {
    string id;
    string email;
    string phone;
    vector<string> subscribedChannels;
};
```

---

## Part 1

**Base requirement — Multi-channel notifications**

Implement a notification system where:
- Users subscribe to one or more channels: `"email"`, `"sms"`, `"push"`
- When `notify(event, users)` is called, each user receives a notification on each of their subscribed channels
- Adding a new channel (e.g., WhatsApp) requires **zero changes** to existing channel classes

**Entry point (tests will call this):**
```cpp
void notify(const string& event, const vector<User>& users);
```

**What to implement:**
```cpp
// Observer interface
class NotificationObserver {
public:
    virtual void update(const string& event, const User& user) = 0;
    virtual string getChannel() = 0;
    virtual ~NotificationObserver() = default;
};

// Concrete observers
class EmailObserver  : public NotificationObserver { ... };
class SMSObserver    : public NotificationObserver { ... };
class PushObserver   : public NotificationObserver { ... };

// Subject (notification manager)
class NotificationManager {
public:
    void subscribe(NotificationObserver* obs);
    void notifyAll(const string& event, const vector<User>& users);
};
```

Each observer checks if the user has subscribed to its channel before sending.

---

## Part 2

**Extension 1 — Priority filtering**

The product team now wants **notification priorities**. Each event has a priority level: `"promotional"`, `"info"`, or `"critical"` (in increasing importance).

Users can set a **minimum priority per channel**. For example:
- SMS: only receive `"critical"` events
- Email: receive `"info"` and above
- Push: receive everything

**Updated entry point:**
```cpp
void notify(const string& event,
            const string& priority,
            const vector<User>& users,
            const unordered_map<string, string>& userMinPriority);
// userMinPriority: maps userId (or "*" for global) to minimum priority
```

**Design challenge:** Where does priority filtering belong — in the observer, the subject, or a wrapper? How do you add filtering without modifying `EmailObserver` or `SMSObserver`?

**Hint:** A `PriorityFilteredObserver` wraps any observer and skips the `update()` call if the event priority is below the user's minimum threshold. This is the Decorator pattern applied to an Observer.

---

## Running Tests

```bash
./run-tests.sh 003-notification-system cpp
```
