#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

struct User {
    string id;
    string email;
    string phone;
    vector<string> subscribedChannels;
};

// ─── Priority Helpers ───────────────────────────────────────────────────────

const vector<string> PRIORITY_ORDER = {"promotional", "info", "critical"};

int priorityLevel(const string& p) {
    for (int i = 0; i < (int)PRIORITY_ORDER.size(); i++)
        if (PRIORITY_ORDER[i] == p) return i;
    return 0;
}

// ─── Observer Interface ─────────────────────────────────────────────────────

class NotificationObserver {
public:
    virtual void send(const string& userId, const string& message) = 0;
    virtual void update(const string& event, const string& priority,
                        const User& user) = 0;
    virtual string channelName() const = 0;
    virtual ~NotificationObserver() = default;
};

// ─── Concrete Observers ─────────────────────────────────────────────────────

class EmailNotifier : public NotificationObserver {
public:
    string channelName() const override { return "email"; }

    void send(const string& userId, const string& message) override {
        cout << "[EMAIL] To: " << userId << " — " << message << endl;
    }

    void update(const string& event, const string& priority,
                const User& user) override {
        cout << "[EMAIL] " << user.email << ": " << event
             << " [" << priority << "]" << endl;
    }
};

class SMSNotifier : public NotificationObserver {
public:
    string channelName() const override { return "sms"; }

    void send(const string& userId, const string& message) override {
        cout << "[SMS] To: " << userId << " — " << message << endl;
    }

    void update(const string& event, const string& priority,
                const User& user) override {
        cout << "[SMS] " << user.phone << ": " << event
             << " [" << priority << "]" << endl;
    }
};

class PushNotifier : public NotificationObserver {
public:
    string channelName() const override { return "push"; }

    void send(const string& userId, const string& message) override {
        cout << "[PUSH] To: " << userId << " — " << message << endl;
    }

    void update(const string& event, const string& priority,
                const User& user) override {
        cout << "[PUSH] " << user.id << ": " << event
             << " [" << priority << "]" << endl;
    }
};

// ─── Priority Filter Decorator (Part 2) ─────────────────────────────────────

class PriorityFilteredObserver : public NotificationObserver {
private:
    NotificationObserver* inner;
    string minPriority;
public:
    PriorityFilteredObserver(NotificationObserver* obs, const string& minP)
        : inner(obs), minPriority(minP) {}

    string channelName() const override { return inner->channelName(); }

    void send(const string& userId, const string& message) override {
        inner->send(userId, message);
    }

    void update(const string& event, const string& priority,
                const User& user) override {
        if (priorityLevel(priority) >= priorityLevel(minPriority)) {
            inner->update(event, priority, user);
        }
    }
};

// ─── Notification Manager ───────────────────────────────────────────────────

class NotificationManager {
private:
    vector<NotificationObserver*> observers;
public:
    void subscribe(NotificationObserver* obs) {
        observers.push_back(obs);
    }

    void unsubscribe(const string& channel) {
        observers.erase(
            remove_if(observers.begin(), observers.end(),
                [&](NotificationObserver* obs) {
                    return obs->channelName() == channel;
                }),
            observers.end());
    }

    // Part 1: simple notify — sends message to each user's subscribed channels
    void notify(const string& event, const vector<User>& users) {
        for (auto& user : users) {
            for (auto* obs : observers) {
                auto& ch = user.subscribedChannels;
                if (find(ch.begin(), ch.end(), obs->channelName()) != ch.end()) {
                    obs->send(user.id, event);
                }
            }
        }
    }

    // Part 2: priority-aware notify
    void notifyAll(const string& event, const string& priority,
                   const vector<User>& users) {
        for (auto& user : users) {
            for (auto* obs : observers) {
                auto& ch = user.subscribedChannels;
                if (find(ch.begin(), ch.end(), obs->channelName()) != ch.end()) {
                    obs->update(event, priority, user);
                }
            }
        }
    }
};

// ─── Free Function: Part 1 (2-arg) ─────────────────────────────────────────

void notify(const string& event, const vector<User>& users) {
    NotificationManager mgr;
    mgr.subscribe(new EmailNotifier());
    mgr.subscribe(new SMSNotifier());
    mgr.subscribe(new PushNotifier());
    mgr.notify(event, users);
}

// ─── Free Function: Part 2 (4-arg with priority filtering) ──────────────────

void notify(const string& event, const string& priority,
            const vector<User>& users,
            const unordered_map<string, string>& userMinPriority) {
    string minP = userMinPriority.count("*") ? userMinPriority.at("*") : "promotional";

    EmailNotifier email;
    SMSNotifier sms;
    PushNotifier push;

    PriorityFilteredObserver fe(&email, minP);
    PriorityFilteredObserver fs(&sms, minP);
    PriorityFilteredObserver fp(&push, minP);

    NotificationManager mgr;
    mgr.subscribe(&fe);
    mgr.subscribe(&fs);
    mgr.subscribe(&fp);
    mgr.notifyAll(event, priority, users);
}

// ─── Main ───────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    // Part 1 demo
    User u1 = {"user1", "user1@example.com", "+91-9000000001", {"email", "sms"}};
    User u2 = {"user2", "user2@example.com", "+91-9000000002", {"push"}};
    cout << "=== Part 1: Basic Notification ===" << endl;
    notify("Order shipped", {u1, u2});

    // Part 2 demo
    cout << "\n=== Part 2: Priority Filtering ===" << endl;
    unordered_map<string, string> prefs = {{"*", "info"}};
    notify("System update available", "info", {u1, u2}, prefs);

    cout << "\n=== Part 2: Promotional (filtered by info+ pref) ===" << endl;
    notify("50% off sale!", "promotional", {u1, u2}, prefs);

    return 0;
}
#endif
