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

class NotificationObserver {
public:
    virtual void send(const string& userId, const string& message) = 0;
    virtual string channelName() const = 0;
    virtual ~NotificationObserver() = default;
};

// ─── Concrete Observers ───────────────────────────────────────────────────────

class EmailNotifier : public NotificationObserver {
public:
    string channelName() const override { return "email"; }
    void send(const string& userId, const string& message) override {
        // TODO: Print a formatted email notification
        // e.g., "[EMAIL] To: <userId> — <message>"
    }
};

class SMSNotifier : public NotificationObserver {
public:
    string channelName() const override { return "sms"; }
    void send(const string& userId, const string& message) override {
        // TODO: Print a formatted SMS notification
        // e.g., "[SMS] To: <userId> — <message>"
    }
};

class PushNotifier : public NotificationObserver {
public:
    string channelName() const override { return "push"; }
    void send(const string& userId, const string& message) override {
        // TODO: Print a formatted push notification
        // e.g., "[PUSH] To: <userId> — <message>"
    }
};

// ─── Notification Manager ────────────────────────────────────────────────────

class NotificationManager {
private:
    vector<NotificationObserver*> observers;
public:
    void subscribe(NotificationObserver* obs) {
        // TODO: Add observer to the list
    }

    void unsubscribe(const string& channel) {
        // TODO: Remove observer matching channelName()
    }

    void notify(const string& event, const vector<User>& users) {
        // TODO: For each user, for each of their subscribedChannels,
        //       find the matching observer and call send()
    }
};

// ─── Test Entry Point ────────────────────────────────────────────────────────

void notify(const string& event, const vector<User>& users) {
    NotificationManager mgr;
    mgr.subscribe(new EmailNotifier());
    mgr.subscribe(new SMSNotifier());
    mgr.subscribe(new PushNotifier());
    mgr.notify(event, users);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Notification System — implement the TODO methods above." << endl;
    return 0;
}
#endif
