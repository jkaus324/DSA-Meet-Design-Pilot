#include <iostream>
#include <vector>
#include <string>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a Notification System that:
//   1. Lets users subscribe to notification channels (email, SMS, push)
//   2. Sends notifications to all subscribed users when an event occurs
//   3. Adding a new channel type requires NO changes to existing classes
//
// Think about:
//   - How do you model the relationship between an event source and its listeners?
//   - What if the same user is subscribed to multiple channels?
//   - How would you add WhatsApp notifications with zero changes to existing code?
//
// Entry point (must exist for tests):
//   void notify(const string& event, const vector<User>& users);
//
// ─────────────────────────────────────────────────────────────────────────────

