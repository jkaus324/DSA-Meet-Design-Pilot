// Meeting Scheduler — Solution (Java, Strategy + Observer)
import java.util.*;

class Op {
    public String kind;
    public String s1;
    public String s2;
    public String s3;
    public int i1;
    public int i2;
    public int i3;

    public Op(String kind) { this(kind, "", "", "", 0, 0, 0); }
    public Op(String kind, String s1) { this(kind, s1, "", "", 0, 0, 0); }
    public Op(String kind, String s1, String s2) { this(kind, s1, s2, "", 0, 0, 0); }
    public Op(String kind, String s1, String s2, String s3) { this(kind, s1, s2, s3, 0, 0, 0); }
    public Op(String kind, String s1, String s2, String s3, int i1) { this(kind, s1, s2, s3, i1, 0, 0); }
    public Op(String kind, String s1, String s2, String s3, int i1, int i2) { this(kind, s1, s2, s3, i1, i2, 0); }
    public Op(String kind, String s1, String s2, String s3, int i1, int i2, int i3) {
        this.kind = kind; this.s1 = s1; this.s2 = s2; this.s3 = s3;
        this.i1 = i1; this.i2 = i2; this.i3 = i3;
    }
}

class MSRoom {
    public String id;
    public String name;
    public int capacity;
    public boolean hasAV;
    public MSRoom(String id, String name, int capacity, boolean hasAV) {
        this.id = id; this.name = name; this.capacity = capacity; this.hasAV = hasAV;
    }
}

class Meeting {
    public String id;
    public String title;
    public int startTime;
    public int endTime;
    public String roomId;
    public Meeting(String id, String title, int startTime, int endTime, String roomId) {
        this.id = id; this.title = title; this.startTime = startTime;
        this.endTime = endTime; this.roomId = roomId;
    }
}

interface MeetingObserver {
    void onMeetingBooked(Meeting meeting);
    void onMeetingCancelled(Meeting meeting);
    void onMeetingRescheduled(Meeting oldMeeting, Meeting newMeeting);
}

class MeetingScheduler {
    public Map<String, MSRoom> rooms = new HashMap<>();
    public Map<String, List<Meeting>> schedule = new HashMap<>();
    public Map<String, Meeting> meetingsById = new HashMap<>();
    public Map<String, List<MeetingObserver>> observers = new HashMap<>();
    public AllocationStrategy strategy = null;

    public void setStrategy(AllocationStrategy s) { this.strategy = s; }
    public void addRoom(MSRoom room) { rooms.put(room.id, room); }

    public List<MSRoom> getAllRooms() {
        List<MSRoom> result = new ArrayList<>(rooms.values());
        result.sort(Comparator.comparing(r -> r.id));
        return result;
    }

    public boolean isAvailable(String roomId, int startTime, int endTime) {
        List<Meeting> ms = schedule.get(roomId);
        if (ms == null) return true;
        for (Meeting m : ms) {
            if (startTime < m.endTime && m.startTime < endTime) return false;
        }
        return true;
    }

    public boolean bookMeeting(Meeting meeting) {
        if (!rooms.containsKey(meeting.roomId)) return false;
        if (!isAvailable(meeting.roomId, meeting.startTime, meeting.endTime)) return false;
        schedule.computeIfAbsent(meeting.roomId, k -> new ArrayList<>()).add(meeting);
        meetingsById.put(meeting.id, meeting);
        List<MeetingObserver> obs = observers.get(meeting.id);
        if (obs != null) for (MeetingObserver o : obs) o.onMeetingBooked(meeting);
        return true;
    }

    public List<Meeting> getRoomSchedule(String roomId) {
        List<Meeting> ms = schedule.get(roomId);
        if (ms == null) return new ArrayList<>();
        List<Meeting> result = new ArrayList<>(ms);
        result.sort(Comparator.comparingInt(m -> m.startTime));
        return result;
    }

    public String bookWithStrategy(String meetingId, String title, int startTime, int endTime, int attendeeCount) {
        if (strategy == null) return "";
        List<MSRoom> all = getAllRooms();
        String roomId = strategy.selectRoom(all, this, startTime, endTime, attendeeCount);
        if (roomId.isEmpty()) return "";
        Meeting m = new Meeting(meetingId, title, startTime, endTime, roomId);
        if (bookMeeting(m)) return roomId;
        return "";
    }

    public void subscribeAttendee(String meetingId, MeetingObserver obs) {
        observers.computeIfAbsent(meetingId, k -> new ArrayList<>()).add(obs);
    }

    public boolean cancelMeeting(String meetingId) {
        Meeting meeting = meetingsById.get(meetingId);
        if (meeting == null) return false;
        List<Meeting> roomMeetings = schedule.get(meeting.roomId);
        if (roomMeetings != null) roomMeetings.removeIf(m -> m.id.equals(meetingId));
        meetingsById.remove(meetingId);
        List<MeetingObserver> obs = observers.get(meetingId);
        if (obs != null) for (MeetingObserver o : obs) o.onMeetingCancelled(meeting);
        return true;
    }

    public boolean rescheduleMeeting(String meetingId, int newStart, int newEnd) {
        Meeting oldMeeting = meetingsById.get(meetingId);
        if (oldMeeting == null) return false;
        List<Meeting> roomMeetings = schedule.get(oldMeeting.roomId);
        if (roomMeetings != null) roomMeetings.removeIf(m -> m.id.equals(meetingId));
        if (!isAvailable(oldMeeting.roomId, newStart, newEnd)) {
            if (roomMeetings != null) roomMeetings.add(oldMeeting);
            return false;
        }
        Meeting newMeeting = new Meeting(oldMeeting.id, oldMeeting.title, newStart, newEnd, oldMeeting.roomId);
        if (roomMeetings == null) {
            roomMeetings = new ArrayList<>();
            schedule.put(oldMeeting.roomId, roomMeetings);
        }
        roomMeetings.add(newMeeting);
        meetingsById.put(meetingId, newMeeting);
        List<MeetingObserver> obs = observers.get(meetingId);
        if (obs != null) for (MeetingObserver o : obs) o.onMeetingRescheduled(oldMeeting, newMeeting);
        return true;
    }
}

interface AllocationStrategy {
    String selectRoom(List<MSRoom> rooms, MeetingScheduler scheduler,
                      int startTime, int endTime, int attendeeCount);
}

class FirstAvailable implements AllocationStrategy {
    public String selectRoom(List<MSRoom> rooms, MeetingScheduler s,
                             int startTime, int endTime, int attendeeCount) {
        for (MSRoom room : rooms) {
            if (room.capacity >= attendeeCount && s.isAvailable(room.id, startTime, endTime)) {
                return room.id;
            }
        }
        return "";
    }
}

class BestFit implements AllocationStrategy {
    public String selectRoom(List<MSRoom> rooms, MeetingScheduler s,
                             int startTime, int endTime, int attendeeCount) {
        String bestId = "";
        int bestCapacity = Integer.MAX_VALUE;
        for (MSRoom room : rooms) {
            if (room.capacity >= attendeeCount && s.isAvailable(room.id, startTime, endTime)) {
                if (room.capacity < bestCapacity) {
                    bestCapacity = room.capacity;
                    bestId = room.id;
                }
            }
        }
        return bestId;
    }
}

class PriorityBased implements AllocationStrategy {
    public String selectRoom(List<MSRoom> rooms, MeetingScheduler s,
                             int startTime, int endTime, int attendeeCount) {
        String bestAV = "", bestNonAV = "";
        int bestAVCap = Integer.MAX_VALUE, bestNonAVCap = Integer.MAX_VALUE;
        for (MSRoom room : rooms) {
            if (room.capacity >= attendeeCount && s.isAvailable(room.id, startTime, endTime)) {
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
        return !bestAV.isEmpty() ? bestAV : bestNonAV;
    }
}

class CountingMeetingObserver implements MeetingObserver {
    public int booked = 0;
    public int cancelled = 0;
    public int rescheduled = 0;
    public int lastNewStart = 0;
    public int lastNewEnd = 0;
    public void onMeetingBooked(Meeting m) { booked++; }
    public void onMeetingCancelled(Meeting m) { cancelled++; }
    public void onMeetingRescheduled(Meeting oldM, Meeting nm) {
        rescheduled++;
        lastNewStart = nm.startTime;
        lastNewEnd = nm.endTime;
    }
}

public class Solution {
    public static List<String> meeting_simulate(List<Op> ops) {
        List<String> out = new ArrayList<>();
        MeetingScheduler s = new MeetingScheduler();
        FirstAvailable fa = new FirstAvailable();
        BestFit bf = new BestFit();
        PriorityBased pb = new PriorityBased();
        List<CountingMeetingObserver> obs = new ArrayList<>();

        for (Op op : ops) {
            String k = op.kind;
            switch (k) {
                case "reset":
                    s = new MeetingScheduler();
                    obs.clear();
                    out.add("ok");
                    break;
                case "add_room":
                    s.addRoom(new MSRoom(op.s1, op.s2, op.i1, op.i2 != 0));
                    out.add("ok");
                    break;
                case "book": {
                    Meeting m = new Meeting(op.s1, op.s1, op.i1, op.i2, op.s2);
                    out.add(s.bookMeeting(m) ? "ok" : "fail");
                    break;
                }
                case "is_available":
                    out.add(s.isAvailable(op.s1, op.i1, op.i2) ? "yes" : "no");
                    break;
                case "sched_size":
                    out.add(String.valueOf(s.getRoomSchedule(op.s1).size()));
                    break;
                case "sched_at": {
                    List<Meeting> sched = s.getRoomSchedule(op.s1);
                    out.add(op.i1 >= 0 && op.i1 < sched.size() ? sched.get(op.i1).id : "");
                    break;
                }
                case "cancel":
                    out.add(s.cancelMeeting(op.s1) ? "ok" : "fail");
                    break;
                case "reschedule":
                    out.add(s.rescheduleMeeting(op.s1, op.i1, op.i2) ? "ok" : "fail");
                    break;
                case "set_strategy":
                    if ("first_available".equals(op.s1)) s.setStrategy(fa);
                    else if ("best_fit".equals(op.s1)) s.setStrategy(bf);
                    else if ("priority".equals(op.s1)) s.setStrategy(pb);
                    out.add("ok");
                    break;
                case "book_strategy": {
                    String r = s.bookWithStrategy(op.s1, op.s2, op.i1, op.i2, op.i3);
                    out.add(r);
                    break;
                }
                case "sub_obs": {
                    while (obs.size() <= op.i1) obs.add(new CountingMeetingObserver());
                    s.subscribeAttendee(op.s1, obs.get(op.i1));
                    out.add("ok");
                    break;
                }
                case "obs_booked":
                    out.add(op.i1 < obs.size() ? String.valueOf(obs.get(op.i1).booked) : "0");
                    break;
                case "obs_cancelled":
                    out.add(op.i1 < obs.size() ? String.valueOf(obs.get(op.i1).cancelled) : "0");
                    break;
                case "obs_rescheduled":
                    out.add(op.i1 < obs.size() ? String.valueOf(obs.get(op.i1).rescheduled) : "0");
                    break;
                case "obs_new_start":
                    out.add(op.i1 < obs.size() ? String.valueOf(obs.get(op.i1).lastNewStart) : "0");
                    break;
                case "obs_new_end":
                    out.add(op.i1 < obs.size() ? String.valueOf(obs.get(op.i1).lastNewEnd) : "0");
                    break;
                default:
                    out.add("unknown:" + k);
                    break;
            }
        }
        return out;
    }
}
