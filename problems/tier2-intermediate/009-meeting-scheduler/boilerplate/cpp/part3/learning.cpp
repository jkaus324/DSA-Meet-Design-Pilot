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

struct Attendee {
    string id;
    string name;
    string email;
};

// ─── Observer Interface ─────────────────────────────────────────────────────

class MeetingObserver {
public:
    virtual void onMeetingBooked(const Meeting& meeting) = 0;
    virtual void onMeetingCancelled(const Meeting& meeting) = 0;
    virtual void onMeetingRescheduled(const Meeting& oldMeeting,
                                      const Meeting& newMeeting) = 0;
    virtual ~MeetingObserver() = default;
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
    unordered_map<string, Meeting> meetingsById;
    unordered_map<string, vector<MeetingObserver*>> observers;
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
        meetingsById[meeting.id] = meeting;
        // TODO: notify observers for this meeting
        // HINT: if (observers.count(meeting.id)) { for (auto* obs : observers[meeting.id]) obs->onMeetingBooked(meeting); }
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
        string roomId = strategy->selectRoom(allRooms, *this,
                                             startTime, endTime, attendeeCount);
        if (roomId.empty()) return "";
        Meeting m{meetingId, title, startTime, endTime, roomId};
        if (bookMeeting(m)) return roomId;
        return "";
    }

    void subscribeAttendee(const string& meetingId, MeetingObserver* obs) {
        observers[meetingId].push_back(obs);
    }

    bool cancelMeeting(const string& meetingId) {
        auto it = meetingsById.find(meetingId);
        if (it == meetingsById.end()) return false;
        Meeting meeting = it->second;
        // TODO: remove meeting from schedule[meeting.roomId]
        // TODO: erase from meetingsById
        // TODO: notify observers with onMeetingCancelled
        return false;
    }

    bool rescheduleMeeting(const string& meetingId, int newStart, int newEnd) {
        auto it = meetingsById.find(meetingId);
        if (it == meetingsById.end()) return false;
        Meeting oldMeeting = it->second;
        // TODO: temporarily remove old meeting
        // TODO: check if new time is available
        // TODO: if available, add new meeting and notify observers
        // TODO: if not available, restore old meeting and return false
        return false;
    }
};

// ─── Concrete Strategies ────────────────────────────────────────────────────

class FirstAvailable : public AllocationStrategy {
public:
    string selectRoom(const vector<Room>& rooms,
                      const MeetingScheduler& scheduler,
                      int startTime, int endTime,
                      int attendeeCount) override {
        for (auto& room : rooms) {
            if (room.capacity >= attendeeCount &&
                scheduler.isAvailable(room.id, startTime, endTime))
                return room.id;
        }
        return "";
    }
};

class BestFit : public AllocationStrategy {
public:
    string selectRoom(const vector<Room>& rooms,
                      const MeetingScheduler& scheduler,
                      int startTime, int endTime,
                      int attendeeCount) override {
        string bestId;
        int bestCap = INT_MAX;
        for (auto& room : rooms) {
            if (room.capacity >= attendeeCount &&
                scheduler.isAvailable(room.id, startTime, endTime)) {
                if (room.capacity < bestCap) {
                    bestCap = room.capacity;
                    bestId = room.id;
                }
            }
        }
        return bestId;
    }
};

class PriorityBased : public AllocationStrategy {
public:
    string selectRoom(const vector<Room>& rooms,
                      const MeetingScheduler& scheduler,
                      int startTime, int endTime,
                      int attendeeCount) override {
        string bestAV, bestNonAV;
        int capAV = INT_MAX, capNonAV = INT_MAX;
        for (auto& room : rooms) {
            if (room.capacity >= attendeeCount &&
                scheduler.isAvailable(room.id, startTime, endTime)) {
                if (room.hasAV && room.capacity < capAV) {
                    capAV = room.capacity;
                    bestAV = room.id;
                } else if (!room.hasAV && room.capacity < capNonAV) {
                    capNonAV = room.capacity;
                    bestNonAV = room.id;
                }
            }
        }
        return bestAV.empty() ? bestNonAV : bestAV;
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

void subscribe_attendee(const string& meetingId, MeetingObserver* observer) {
    scheduler.subscribeAttendee(meetingId, observer);
}

bool cancel_meeting(const string& meetingId) {
    return scheduler.cancelMeeting(meetingId);
}

bool reschedule_meeting(const string& meetingId, int newStart, int newEnd) {
    return scheduler.rescheduleMeeting(meetingId, newStart, newEnd);
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Observer notifications — implement the TODOs above." << endl;
    return 0;
}
#endif
