#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
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
    int startTime;  // minutes since midnight
    int endTime;
    string roomId;
};

// ─── Scheduler ──────────────────────────────────────────────────────────────

class MeetingScheduler {
    unordered_map<string, Room> rooms;
    unordered_map<string, vector<Meeting>> schedule; // roomId -> meetings

public:
    void addRoom(const Room& room) {
        rooms[room.id] = room;
    }

    bool isAvailable(const string& roomId, int startTime, int endTime) const {
        auto it = schedule.find(roomId);
        if (it == schedule.end()) return true;
        for (auto& m : it->second) {
            // TODO: return false if [startTime, endTime) overlaps [m.startTime, m.endTime)
            // HINT: overlap condition is startTime < m.endTime && m.startTime < endTime
        }
        return true;
    }

    bool bookMeeting(const Meeting& meeting) {
        if (rooms.find(meeting.roomId) == rooms.end()) return false;
        // TODO: check availability, then add meeting to schedule
        // HINT: if (!isAvailable(...)) return false; then push_back
        return false;
    }

    vector<Meeting> getRoomSchedule(const string& roomId) const {
        auto it = schedule.find(roomId);
        if (it == schedule.end()) return {};
        auto result = it->second;
        // TODO: sort result by startTime
        return result;
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

#ifndef RUNNING_TESTS
int main() {
    cout << "Meeting Scheduler — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
