#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

struct Notification { string userId, message, channel; };
struct User { string id, email, phone; vector<string> subscribedChannels; };

const vector<string> PRIORITY_ORDER = {"promotional", "info", "critical"};

int priorityLevel(const string& p) {
    for (int i = 0; i < (int)PRIORITY_ORDER.size(); i++)
        if (PRIORITY_ORDER[i] == p) return i;
    return 0;
}

class NotificationObserver {
public:
    virtual void update(const string& event, const string& priority,
                        const User& user) = 0;
    virtual string getChannel() = 0;
    virtual ~NotificationObserver() = default;
};

class EmailObserver : public NotificationObserver {
public:
    void update(const string& event, const string& priority, const User& user) override {
        cout << "[EMAIL] " << user.email << ": " << event << " [" << priority << "]" << endl;
    }
    string getChannel() override { return "email"; }
};

class SMSObserver : public NotificationObserver {
public:
    void update(const string& event, const string& priority, const User& user) override {
        cout << "[SMS] " << user.phone << ": " << event << " [" << priority << "]" << endl;
    }
    string getChannel() override { return "sms"; }
};

class PushObserver : public NotificationObserver {
public:
    void update(const string& event, const string& priority, const User& user) override {
        cout << "[PUSH] " << user.id << ": " << event << " [" << priority << "]" << endl;
    }
    string getChannel() override { return "push"; }
};

class PriorityFilteredObserver : public NotificationObserver {
private:
    NotificationObserver* inner;
    string minPriority;
public:
    PriorityFilteredObserver(NotificationObserver* obs, string minP)
        : inner(obs), minPriority(minP) {}

    void update(const string& event, const string& priority, const User& user) override {
        if (priorityLevel(priority) >= priorityLevel(minPriority)) {
            inner->update(event, priority, user);
        }
    }
    string getChannel() override { return inner->getChannel(); }
};

class NotificationManager {
    vector<NotificationObserver*> observers;
public:
    void subscribe(NotificationObserver* obs) { observers.push_back(obs); }
    void notifyAll(const string& event, const string& priority, const vector<User>& users) {
        for (auto& user : users) {
            for (auto* obs : observers) {
                auto& ch = user.subscribedChannels;
                if (find(ch.begin(), ch.end(), obs->getChannel()) != ch.end()) {
                    obs->update(event, priority, user);
                }
            }
        }
    }
};

void notify(const string& event, const string& priority,
            const vector<User>& users,
            const unordered_map<string, string>& userMinPriority) {
    // For simplicity, apply same minPriority to all users
    string minP = userMinPriority.count("*") ? userMinPriority.at("*") : "promotional";
    EmailObserver email; SMSObserver sms; PushObserver push;
    PriorityFilteredObserver fe(&email, minP), fs(&sms, minP), fp(&push, minP);
    NotificationManager mgr;
    mgr.subscribe(&fe); mgr.subscribe(&fs); mgr.subscribe(&fp);
    mgr.notifyAll(event, priority, users);
}

int main() {
    cout << "Part 2: Priority filtering — full scaffolding provided." << endl;
    return 0;
}
