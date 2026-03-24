#include <iostream>
#include <vector>
#include <string>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct User {
    string id;
    string email;
    string phone;
    vector<string> subscribedChannels;
};

// ─── Observer Interface ───────────────────────────────────────────────────────
// HINT: This interface represents something that "watches" for events.
// What method does it need to receive a notification?

class /* YourObserverName */ {
public:
    virtual void /* yourMethodName */(const string& userId, const string& message) = 0;
    virtual ~/* YourObserverName */() = default;
};

// ─── Concrete Observers ───────────────────────────────────────────────────────
// TODO: Implement one observer per notification channel:
//   - EmailNotifier
//   - SMSNotifier
//   - PushNotifier

// ─── Subject / Notification Manager ─────────────────────────────────────────
// TODO: Implement a NotificationManager that:
//   - Allows observers to subscribe/unsubscribe
//   - Notifies all subscribers when an event occurs
//   - Does NOT know which specific channels are used

// class NotificationManager {
// public:
//     void subscribe(/* what goes here? */);
//     void unsubscribe(/* what goes here? */);
//     void notify(const string& event, const vector<User>& users);
// };

// ─── Test Entry Points ───────────────────────────────────────────────────────
//   void notify(const string& event, const vector<User>& users);
// ─────────────────────────────────────────────────────────────────────────────

