#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

struct Notification { string userId, message, channel; };
struct User { string id, email, phone; vector<string> subscribedChannels; };

// Priority order: critical > info > promotional
const vector<string> PRIORITY_ORDER = {"promotional", "info", "critical"};

int priorityLevel(const string& p) {
    for (int i = 0; i < (int)PRIORITY_ORDER.size(); i++)
        if (PRIORITY_ORDER[i] == p) return i;
    return 0;
}

// ─── Observer Interface ───────────────────────────────────────────────────────

class NotificationObserver {
public:
    virtual void update(const string& event, const string& priority,
                        const User& user) = 0;
    virtual string getChannel() = 0;
    virtual ~NotificationObserver() = default;
};

// ─── Existing Channels (copy from Part 1) ────────────────────────────────────

class EmailObserver : public NotificationObserver {
public:
    void update(const string& event, const string& priority, const User& user) override {
        // TODO: implement
    }
    string getChannel() override { return "email"; }
};

class SMSObserver : public NotificationObserver {
public:
    void update(const string& event, const string& priority, const User& user) override {
        // TODO: implement
    }
    string getChannel() override { return "sms"; }
};

class PushObserver : public NotificationObserver {
public:
    void update(const string& event, const string& priority, const User& user) override {
        // TODO: implement
    }
    string getChannel() override { return "push"; }
};

// ─── NEW: Filtered Observer Decorator ────────────────────────────────────────
// HINT: Wraps any observer and filters by priority.
// If event priority < user's minPriority for this channel, skip the notification.

class PriorityFilteredObserver : public NotificationObserver {
private:
    NotificationObserver* inner;
    string minPriority;
public:
    PriorityFilteredObserver(NotificationObserver* obs, string minP)
        : inner(obs), minPriority(minP) {}

    void update(const string& event, const string& priority, const User& user) override {
        // TODO: Only call inner->update() if priorityLevel(priority) >= priorityLevel(minPriority)
    }

    string getChannel() override { return inner->getChannel(); }
};

// ─── Notification Manager ─────────────────────────────────────────────────────

class NotificationManager {
    vector<NotificationObserver*> observers;
public:
    void subscribe(NotificationObserver* obs) { observers.push_back(obs); }
    void notifyAll(const string& event, const string& priority,
                   const vector<User>& users) {
        // TODO: for each user, for each observer matching user's channel, call update()
    }
};

void notify(const string& event, const string& priority,
            const vector<User>& users,
            const unordered_map<string, string>& userMinPriority) {
    // TODO: Build manager with PriorityFilteredObserver wrappers
}

int main() {
    cout << "Part 2: Priority filtering — implement TODOs above." << endl;
    return 0;
}
