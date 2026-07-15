#include <iostream>
#include <memory>
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

// ─── Ops simulator (used by spec-based tests) ───────────────────────────────
//
// Each Op encodes one method call on a fresh-or-shared MeetingScheduler.
// `kind` selects the operation. To keep the schema simple, every op carries
// a few string fields (s1..s3) and a few int fields (i1..i3); each kind
// reads only the ones it needs. The simulator returns one string per op
// describing its outcome — checked against the expected sequence.
//
// kinds:
//   "reset"           — start a new MeetingScheduler -> "ok"
//   "add_room"        — s1=id, s2=name, i1=capacity, i2=hasAV(0/1) -> "ok"
//   "book"            — s1=mid, s2=roomId, i1=start, i2=end -> "ok"/"fail"
//   "is_available"    — s1=roomId, i1=start, i2=end -> "yes"/"no"
//   "sched_size"      — s1=roomId -> int-as-string
//   "sched_at"        — s1=roomId, i1=index -> meeting id at index (sorted by start)
//   "cancel"          — s1=mid -> "ok"/"fail"
//   "reschedule"      — s1=mid, i1=newStart, i2=newEnd -> "ok"/"fail"
//   "set_strategy"    — s1="first_available"|"best_fit"|"priority" -> "ok"
//   "book_strategy"   — s1=mid, s2=title, i1=start, i2=end, i3=attendees -> roomId or ""
//   "sub_obs"         — s1=mid, i1=observer index (auto-allocated) -> "ok"
//   "obs_booked"      — i1=observer index -> count
//   "obs_cancelled"   — i1=observer index -> count
//   "obs_rescheduled" — i1=observer index -> count
//   "obs_new_start"   — i1=observer index -> int
//   "obs_new_end"     — i1=observer index -> int

struct Op {
    string kind;
    string s1;
    string s2;
    string s3;
    int    i1;
    int    i2;
    int    i3;
};

class CountingObserver : public MeetingObserver {
public:
    int booked = 0;
    int cancelled = 0;
    int rescheduled = 0;
    int lastNewStart = 0;
    int lastNewEnd = 0;
    void onMeetingBooked(const Meeting&) override { booked++; }
    void onMeetingCancelled(const Meeting&) override { cancelled++; }
    void onMeetingRescheduled(const Meeting&, const Meeting& nm) override {
        rescheduled++;
        lastNewStart = nm.startTime;
        lastNewEnd = nm.endTime;
    }
};

vector<string> meeting_simulate(vector<Op> ops) {
    vector<string> out;
    MeetingScheduler s;
    FirstAvailable fa; BestFit bf; PriorityBased pb;
    vector<unique_ptr<CountingObserver>> obs;
    auto ensure_obs = [&](int idx) {
        while ((int)obs.size() <= idx) obs.push_back(unique_ptr<CountingObserver>(new CountingObserver()));
    };
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "reset") {
            s = MeetingScheduler();
            obs.clear();
            out.push_back("ok");
        } else if (k == "add_room") {
            s.addRoom({op.s1, op.s2, op.i1, op.i2 != 0});
            out.push_back("ok");
        } else if (k == "book") {
            Meeting m{op.s1, op.s1, op.i1, op.i2, op.s2};
            out.push_back(s.bookMeeting(m) ? "ok" : "fail");
        } else if (k == "is_available") {
            out.push_back(s.isAvailable(op.s1, op.i1, op.i2) ? "yes" : "no");
        } else if (k == "sched_size") {
            out.push_back(to_string((int)s.getRoomSchedule(op.s1).size()));
        } else if (k == "sched_at") {
            auto sched = s.getRoomSchedule(op.s1);
            out.push_back(op.i1 >= 0 && op.i1 < (int)sched.size() ? sched[op.i1].id : "");
        } else if (k == "cancel") {
            out.push_back(s.cancelMeeting(op.s1) ? "ok" : "fail");
        } else if (k == "reschedule") {
            out.push_back(s.rescheduleMeeting(op.s1, op.i1, op.i2) ? "ok" : "fail");
        } else if (k == "set_strategy") {
            if (op.s1 == "first_available") s.setStrategy(&fa);
            else if (op.s1 == "best_fit") s.setStrategy(&bf);
            else if (op.s1 == "priority") s.setStrategy(&pb);
            out.push_back("ok");
        } else if (k == "book_strategy") {
            string r = s.bookWithStrategy(op.s1, op.s2, op.i1, op.i2, op.i3);
            out.push_back(r);
        } else if (k == "sub_obs") {
            ensure_obs(op.i1);
            s.subscribeAttendee(op.s1, obs[op.i1].get());
            out.push_back("ok");
        } else if (k == "obs_booked") {
            out.push_back(op.i1 < (int)obs.size() ? to_string(obs[op.i1]->booked) : "0");
        } else if (k == "obs_cancelled") {
            out.push_back(op.i1 < (int)obs.size() ? to_string(obs[op.i1]->cancelled) : "0");
        } else if (k == "obs_rescheduled") {
            out.push_back(op.i1 < (int)obs.size() ? to_string(obs[op.i1]->rescheduled) : "0");
        } else if (k == "obs_new_start") {
            out.push_back(op.i1 < (int)obs.size() ? to_string(obs[op.i1]->lastNewStart) : "0");
        } else if (k == "obs_new_end") {
            out.push_back(op.i1 < (int)obs.size() ? to_string(obs[op.i1]->lastNewEnd) : "0");
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
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
