#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <climits>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct Room {
    string id;
    string name;
    int capacity;
    bool hasAV;
};

struct Meeting {
    string id;
    string title;
    int startTime;
    int endTime;
    string roomId;
};

// ─── Forward declaration ────────────────────────────────────────────────────

class MeetingScheduler;

// ─── Strategy Interface ─────────────────────────────────────────────────────

class AllocationStrategy {
public:
    virtual string selectRoom(const vector<Room>& rooms,
                              const MeetingScheduler& scheduler,
                              int startTime, int endTime,
                              int attendeeCount) = 0;
    virtual ~AllocationStrategy() = default;
};

// ─── Scheduler ──────────────────────────────────────────────────────────────

class MeetingScheduler {
    unordered_map<string, Room> rooms;
    unordered_map<string, vector<Meeting>> schedule;
    AllocationStrategy* strategy;

public:
    MeetingScheduler() : strategy(nullptr) {}

    void setStrategy(AllocationStrategy* s) { strategy = s; }

    void addRoom(const Room& room) {
        rooms[room.id] = room;
    }

    vector<Room> getAllRooms() const {
        vector<Room> result;
        for (auto& [id, room] : rooms) result.push_back(room);
        sort(result.begin(), result.end(),
             [](const Room& a, const Room& b) { return a.id < b.id; });
        return result;
    }

    bool isAvailable(const string& roomId, int startTime, int endTime) const {
        auto it = schedule.find(roomId);
        if (it == schedule.end()) return true;
        for (auto& m : it->second) {
            if (startTime < m.endTime && m.startTime < endTime) return false;
        }
        return true;
    }

    bool bookMeeting(const Meeting& meeting) {
        if (rooms.find(meeting.roomId) == rooms.end()) return false;
        if (!isAvailable(meeting.roomId, meeting.startTime, meeting.endTime))
            return false;
        schedule[meeting.roomId].push_back(meeting);
        return true;
    }

    vector<Meeting> getRoomSchedule(const string& roomId) const {
        auto it = schedule.find(roomId);
        if (it == schedule.end()) return {};
        auto result = it->second;
        sort(result.begin(), result.end(),
             [](const Meeting& a, const Meeting& b) {
                 return a.startTime < b.startTime;
             });
        return result;
    }

    string bookWithStrategy(const string& meetingId, const string& title,
                            int startTime, int endTime, int attendeeCount) {
        if (!strategy) return "";
        auto allRooms = getAllRooms();
        // TODO: use strategy->selectRoom() to pick a room, then bookMeeting()
        // HINT: string roomId = strategy->selectRoom(allRooms, *this, startTime, endTime, attendeeCount);
        return "";
    }
};

// ─── Concrete Strategies ────────────────────────────────────────────────────

class FirstAvailable : public AllocationStrategy {
public:
    string selectRoom(const vector<Room>& rooms,
                      const MeetingScheduler& scheduler,
                      int startTime, int endTime,
                      int attendeeCount) override {
        // TODO: iterate rooms, return first that fits capacity and is available
        return "";
    }
};

class BestFit : public AllocationStrategy {
public:
    string selectRoom(const vector<Room>& rooms,
                      const MeetingScheduler& scheduler,
                      int startTime, int endTime,
                      int attendeeCount) override {
        // TODO: find room with smallest capacity >= attendeeCount that is available
        // HINT: track bestId and bestCapacity, update when you find a smaller fit
        return "";
    }
};

class PriorityBased : public AllocationStrategy {
public:
    string selectRoom(const vector<Room>& rooms,
                      const MeetingScheduler& scheduler,
                      int startTime, int endTime,
                      int attendeeCount) override {
        // TODO: prefer AV rooms; among AV rooms, pick smallest that fits
        // If no AV rooms available, fall back to smallest non-AV room
        return "";
    }
};

// ─── Test Entry Points ──────────────────────────────────────────────────────

MeetingScheduler scheduler;

bool book_meeting(const Meeting& meeting) {
    return scheduler.bookMeeting(meeting);
}

vector<Meeting> get_room_schedule(const string& roomId) {
    return scheduler.getRoomSchedule(roomId);
}

bool is_available(const string& roomId, int startTime, int endTime) {
    return scheduler.isAvailable(roomId, startTime, endTime);
}

string book_with_strategy(const string& meetingId, const string& title,
                          int startTime, int endTime, int attendeeCount) {
    return scheduler.bookWithStrategy(meetingId, title, startTime, endTime, attendeeCount);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Allocation strategies — implement the TODOs above." << endl;
    return 0;
}
#endif
