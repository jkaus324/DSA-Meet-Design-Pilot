#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct Notification {
    string userId;
    string message;
    string channel;
};

struct User {
    string id;
    string email;
    string phone;
    vector<string> subscribedChannels;
};

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The product team wants notification PRIORITIES and FILTERING:
//   - Each event now has a priority: "critical", "info", "promotional"
//   - Users can set minimum priority per channel
//     (e.g., only receive SMS for "critical" events, not "promotional")
//   - The system must respect these per-user, per-channel preferences
//
// Think about:
//   - Where does priority filtering belong in your Observer design?
//   - Is filtering a responsibility of the subject, the observer, or a decorator?
//   - How do you store per-user preferences without coupling User to channels?
//
// Entry points:
//   void notify(const string& event, const string& priority,
//               const vector<User>& users,
//               const unordered_map<string, string>& userMinPriority);
//               // userMinPriority: userId -> min priority for their channel
//
// ─────────────────────────────────────────────────────────────────────────────


