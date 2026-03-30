#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <climits>
using namespace std;

// ─── Data Structures ────────────────────────────────────────────────────────

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

// ─── TODO: Implement MeetingScheduler ───────────────────────────────────────

class MeetingScheduler {
    unordered_map<string, Room> rooms;
    unordered_map<string, vector<Meeting>> schedule; // roomId -> meetings
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
            // TODO: check overlap — two intervals overlap if start1 < end2 && start2 < end1
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
        // Notify observers
        if (observers.count(meeting.id)) {
            for (auto* obs : observers[meeting.id]) {
                obs->onMeetingBooked(meeting);
            }
        }
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
        auto& roomMeetings = schedule[meeting.roomId];
        roomMeetings.erase(
            remove_if(roomMeetings.begin(), roomMeetings.end(),
                      [&](const Meeting& m) { return m.id == meetingId; }),
            roomMeetings.end());
        meetingsById.erase(it);
        // Notify observers
        if (observers.count(meetingId)) {
            for (auto* obs : observers[meetingId]) {
                obs->onMeetingCancelled(meeting);
            }
        }
        return true;
    }

    bool rescheduleMeeting(const string& meetingId, int newStart, int newEnd) {
        auto it = meetingsById.find(meetingId);
        if (it == meetingsById.end()) return false;
        Meeting oldMeeting = it->second;
        // Temporarily remove to check availability
        auto& roomMeetings = schedule[oldMeeting.roomId];
        roomMeetings.erase(
            remove_if(roomMeetings.begin(), roomMeetings.end(),
                      [&](const Meeting& m) { return m.id == meetingId; }),
            roomMeetings.end());
        if (!isAvailable(oldMeeting.roomId, newStart, newEnd)) {
            roomMeetings.push_back(oldMeeting);
            return false;
        }
        Meeting newMeeting = oldMeeting;
        newMeeting.startTime = newStart;
        newMeeting.endTime = newEnd;
        roomMeetings.push_back(newMeeting);
        meetingsById[meetingId] = newMeeting;
        // Notify observers
        if (observers.count(meetingId)) {
            for (auto* obs : observers[meetingId]) {
                obs->onMeetingRescheduled(oldMeeting, newMeeting);
            }
        }
        return true;
    }
};

// ─── TODO: Implement Concrete Strategies ────────────────────────────────────

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
        string bestId = "";
        int bestCapacity = INT_MAX;
        for (auto& room : rooms) {
            if (room.capacity >= attendeeCount &&
                scheduler.isAvailable(room.id, startTime, endTime)) {
                if (room.capacity < bestCapacity) {
                    bestCapacity = room.capacity;
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
        string bestAV = "", bestNonAV = "";
        int bestAVCap = INT_MAX, bestNonAVCap = INT_MAX;
        for (auto& room : rooms) {
            if (room.capacity >= attendeeCount &&
                scheduler.isAvailable(room.id, startTime, endTime)) {
                if (room.hasAV) {
                    if (room.capacity < bestAVCap) {
                        bestAVCap = room.capacity;
                        bestAV = room.id;
                    }
                } else {
                    if (room.capacity < bestNonAVCap) {
                        bestNonAVCap = room.capacity;
                        bestNonAV = room.id;
                    }
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

// ─── Main (test your implementation) ────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    scheduler.addRoom({"R1", "Small Room", 4, false});
    scheduler.addRoom({"R2", "Large Room", 20, true});

    Meeting m1{"M1", "Standup", 540, 570, "R1"};
    cout << "Book M1: " << (book_meeting(m1) ? "OK" : "FAIL") << endl;
    cout << "Double-book M1 slot: " << (is_available("R1", 550, 580) ? "available" : "conflict") << endl;

    return 0;
}
#endif
