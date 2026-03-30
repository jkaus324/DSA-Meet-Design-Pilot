#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <queue>
#include <unordered_map>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

enum class LockerSize { SMALL, MEDIUM, LARGE };

struct Package {
    string packageId;
    LockerSize size;
};

struct Locker {
    string lockerId;
    LockerSize size;
    bool occupied;
};

// ─── Allocation Strategy ────────────────────────────────────────────────────

class LockerAllocationStrategy {
public:
    virtual string allocate(LockerSize packageSize,
                            map<LockerSize, queue<string>>& available) = 0;
    virtual ~LockerAllocationStrategy() = default;
};

class SmallestFitStrategy : public LockerAllocationStrategy {
public:
    string allocate(LockerSize packageSize,
                    map<LockerSize, queue<string>>& available) override {
        vector<LockerSize> tryOrder;
        if (packageSize == LockerSize::SMALL)
            tryOrder = {LockerSize::SMALL, LockerSize::MEDIUM, LockerSize::LARGE};
        else if (packageSize == LockerSize::MEDIUM)
            tryOrder = {LockerSize::MEDIUM, LockerSize::LARGE};
        else
            tryOrder = {LockerSize::LARGE};

        for (auto sz : tryOrder) {
            if (!available[sz].empty()) {
                string id = available[sz].front();
                available[sz].pop();
                return id;
            }
        }
        return "";
    }
};

// ─── Notification Channel ───────────────────────────────────────────────────

class NotificationChannel {
public:
    virtual void notify(const string& packageId, const string& message) = 0;
    virtual ~NotificationChannel() = default;
};

// ─── Deposit Record ─────────────────────────────────────────────────────────

struct DepositRecord {
    string lockerId;
    string packageId;
    string pickupCode;
    long   depositTime;
};

// ─── Locker System ──────────────────────────────────────────────────────────

class LockerSystem {
private:
    map<string, Locker> lockers;
    map<LockerSize, queue<string>> available;
    unordered_map<string, DepositRecord> activeDeposits;
    LockerAllocationStrategy* strategy;
    vector<NotificationChannel*> channels;
    int codeCounter = 0;
    int expiryHours = 0;

    string generateCode() {
        return "CODE-" + to_string(++codeCounter);
    }

    void notifyAll(const string& packageId, const string& message) {
        for (auto* ch : channels) {
            ch->notify(packageId, message);
        }
    }

    void freeLocker(const string& lockerId) {
        auto it = lockers.find(lockerId);
        if (it != lockers.end()) {
            it->second.occupied = false;
            available[it->second.size].push(lockerId);
        }
    }

public:
    LockerSystem() {
        strategy = new SmallestFitStrategy();
    }

    ~LockerSystem() {
        delete strategy;
    }

    void addLocker(const string& lockerId, LockerSize size) {
        lockers[lockerId] = {lockerId, size, false};
        available[size].push(lockerId);
    }

    string depositPackage(const string& packageId, LockerSize size,
                          long depositTime = 0) {
        string lockerId = strategy->allocate(size, available);
        if (lockerId.empty()) return "";

        lockers[lockerId].occupied = true;
        string code = generateCode();
        activeDeposits[code] = {lockerId, packageId, code, depositTime};
        notifyAll(packageId, "Package " + packageId + " deposited. Code: " + code);
        return code;
    }

    bool retrievePackage(const string& code) {
        auto it = activeDeposits.find(code);
        if (it == activeDeposits.end()) return false;

        freeLocker(it->second.lockerId);
        activeDeposits.erase(it);
        return true;
    }

    void setCodeExpiry(int hours) {
        expiryHours = hours;
    }

    vector<string> checkExpired(long currentTime) {
        vector<string> expired;
        if (expiryHours <= 0) return expired;

        for (auto it = activeDeposits.begin(); it != activeDeposits.end(); ) {
            if (currentTime - it->second.depositTime > (long)expiryHours * 3600) {
                freeLocker(it->second.lockerId);
                expired.push_back(it->second.packageId);
                notifyAll(it->second.packageId,
                          "Package " + it->second.packageId + " expired. Locker freed.");
                it = activeDeposits.erase(it);
            } else {
                ++it;
            }
        }
        return expired;
    }

    void addNotificationChannel(NotificationChannel* channel) {
        channels.push_back(channel);
    }
};

// ─── Global Entry Points ─────────────────────────────────────────────────────

static LockerSystem* g_system = nullptr;

void initLockerSystem() {
    delete g_system;
    g_system = new LockerSystem();
}

void addLocker(const string& lockerId, LockerSize size) {
    if (!g_system) initLockerSystem();
    g_system->addLocker(lockerId, size);
}

string depositPackage(const string& packageId, LockerSize size,
                      long depositTime = 0) {
    if (!g_system) initLockerSystem();
    return g_system->depositPackage(packageId, size, depositTime);
}

bool retrievePackage(const string& code) {
    if (!g_system) return false;
    return g_system->retrievePackage(code);
}

void setCodeExpiry(int hours) {
    if (!g_system) initLockerSystem();
    g_system->setCodeExpiry(hours);
}

vector<string> checkExpired(long currentTime) {
    if (!g_system) return {};
    return g_system->checkExpired(currentTime);
}

void addNotificationChannel(NotificationChannel* channel) {
    if (!g_system) initLockerSystem();
    g_system->addNotificationChannel(channel);
}

// ─── Main ────────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    initLockerSystem();

    addLocker("S1", LockerSize::SMALL);
    addLocker("S2", LockerSize::SMALL);
    addLocker("M1", LockerSize::MEDIUM);
    addLocker("L1", LockerSize::LARGE);

    string code1 = depositPackage("PKG001", LockerSize::SMALL);
    cout << "Deposited PKG001, code: " << code1 << endl;

    string code2 = depositPackage("PKG002", LockerSize::MEDIUM);
    cout << "Deposited PKG002, code: " << code2 << endl;

    bool ok = retrievePackage(code1);
    cout << "Retrieved PKG001: " << (ok ? "success" : "failed") << endl;

    return 0;
}
#endif
