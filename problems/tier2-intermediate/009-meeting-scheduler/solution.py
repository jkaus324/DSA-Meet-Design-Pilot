"""Meeting scheduler — rooms, observers, allocation strategies."""


class Op:
    def __init__(self, kind, s1="", s2="", s3="", i1=0, i2=0, i3=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.i1 = i1
        self.i2 = i2
        self.i3 = i3


class Room:
    def __init__(self, id, name, capacity, hasAV):
        self.id = id
        self.name = name
        self.capacity = capacity
        self.hasAV = hasAV


class Meeting:
    def __init__(self, id, title, startTime, endTime, roomId):
        self.id = id
        self.title = title
        self.startTime = startTime
        self.endTime = endTime
        self.roomId = roomId


class CountingObserver:
    def __init__(self):
        self.booked = 0
        self.cancelled = 0
        self.rescheduled = 0
        self.lastNewStart = 0
        self.lastNewEnd = 0

    def onMeetingBooked(self, meeting):
        self.booked += 1

    def onMeetingCancelled(self, meeting):
        self.cancelled += 1

    def onMeetingRescheduled(self, oldM, newM):
        self.rescheduled += 1
        self.lastNewStart = newM.startTime
        self.lastNewEnd = newM.endTime


class FirstAvailable:
    def selectRoom(self, rooms, scheduler, startTime, endTime, attendeeCount):
        for r in rooms:
            if r.capacity >= attendeeCount and scheduler.isAvailable(r.id, startTime, endTime):
                return r.id
        return ""


class BestFit:
    def selectRoom(self, rooms, scheduler, startTime, endTime, attendeeCount):
        bestId = ""
        bestCap = float("inf")
        for r in rooms:
            if r.capacity >= attendeeCount and scheduler.isAvailable(r.id, startTime, endTime):
                if r.capacity < bestCap:
                    bestCap = r.capacity
                    bestId = r.id
        return bestId


class PriorityBased:
    def selectRoom(self, rooms, scheduler, startTime, endTime, attendeeCount):
        bestAV, bestNonAV = "", ""
        bestAVCap, bestNonAVCap = float("inf"), float("inf")
        for r in rooms:
            if r.capacity >= attendeeCount and scheduler.isAvailable(r.id, startTime, endTime):
                if r.hasAV:
                    if r.capacity < bestAVCap:
                        bestAVCap = r.capacity
                        bestAV = r.id
                else:
                    if r.capacity < bestNonAVCap:
                        bestNonAVCap = r.capacity
                        bestNonAV = r.id
        return bestNonAV if bestAV == "" else bestAV


class MeetingScheduler:
    def __init__(self):
        self.rooms = {}  # id -> Room
        self.schedule = {}  # roomId -> [Meeting]
        self.meetingsById = {}
        self.observers = {}  # meetingId -> [observer]
        self.strategy = None

    def setStrategy(self, s):
        self.strategy = s

    def addRoom(self, room):
        self.rooms[room.id] = room

    def getAllRooms(self):
        return sorted(self.rooms.values(), key=lambda r: r.id)

    def isAvailable(self, roomId, startTime, endTime):
        meets = self.schedule.get(roomId, [])
        for m in meets:
            if startTime < m.endTime and m.startTime < endTime:
                return False
        return True

    def bookMeeting(self, meeting):
        if meeting.roomId not in self.rooms:
            return False
        if not self.isAvailable(meeting.roomId, meeting.startTime, meeting.endTime):
            return False
        self.schedule.setdefault(meeting.roomId, []).append(meeting)
        self.meetingsById[meeting.id] = meeting
        for obs in self.observers.get(meeting.id, []):
            obs.onMeetingBooked(meeting)
        return True

    def getRoomSchedule(self, roomId):
        meets = list(self.schedule.get(roomId, []))
        return sorted(meets, key=lambda m: m.startTime)

    def bookWithStrategy(self, meetingId, title, startTime, endTime, attendeeCount):
        if not self.strategy:
            return ""
        all_rooms = self.getAllRooms()
        roomId = self.strategy.selectRoom(all_rooms, self, startTime, endTime, attendeeCount)
        if roomId == "":
            return ""
        m = Meeting(meetingId, title, startTime, endTime, roomId)
        if self.bookMeeting(m):
            return roomId
        return ""

    def subscribeAttendee(self, meetingId, obs):
        self.observers.setdefault(meetingId, []).append(obs)

    def cancelMeeting(self, meetingId):
        if meetingId not in self.meetingsById:
            return False
        meeting = self.meetingsById[meetingId]
        rm = self.schedule.get(meeting.roomId, [])
        self.schedule[meeting.roomId] = [m for m in rm if m.id != meetingId]
        del self.meetingsById[meetingId]
        for obs in self.observers.get(meetingId, []):
            obs.onMeetingCancelled(meeting)
        return True

    def rescheduleMeeting(self, meetingId, newStart, newEnd):
        if meetingId not in self.meetingsById:
            return False
        old = self.meetingsById[meetingId]
        rm = self.schedule.get(old.roomId, [])
        self.schedule[old.roomId] = [m for m in rm if m.id != meetingId]
        if not self.isAvailable(old.roomId, newStart, newEnd):
            self.schedule[old.roomId].append(old)
            return False
        new = Meeting(old.id, old.title, newStart, newEnd, old.roomId)
        self.schedule[old.roomId].append(new)
        self.meetingsById[meetingId] = new
        for obs in self.observers.get(meetingId, []):
            obs.onMeetingRescheduled(old, new)
        return True


def meeting_simulate(ops):
    out = []
    s = MeetingScheduler()
    fa, bf, pb = FirstAvailable(), BestFit(), PriorityBased()
    obs = []

    def ensure_obs(idx):
        while len(obs) <= idx:
            obs.append(CountingObserver())

    for op in ops:
        k = op.kind
        if k == "reset":
            s = MeetingScheduler()
            obs = []
            out.append("ok")
        elif k == "add_room":
            s.addRoom(Room(op.s1, op.s2, op.i1, op.i2 != 0))
            out.append("ok")
        elif k == "book":
            m = Meeting(op.s1, op.s1, op.i1, op.i2, op.s2)
            out.append("ok" if s.bookMeeting(m) else "fail")
        elif k == "is_available":
            out.append("yes" if s.isAvailable(op.s1, op.i1, op.i2) else "no")
        elif k == "sched_size":
            out.append(str(len(s.getRoomSchedule(op.s1))))
        elif k == "sched_at":
            sched = s.getRoomSchedule(op.s1)
            out.append(sched[op.i1].id if 0 <= op.i1 < len(sched) else "")
        elif k == "cancel":
            out.append("ok" if s.cancelMeeting(op.s1) else "fail")
        elif k == "reschedule":
            out.append("ok" if s.rescheduleMeeting(op.s1, op.i1, op.i2) else "fail")
        elif k == "set_strategy":
            if op.s1 == "first_available":
                s.setStrategy(fa)
            elif op.s1 == "best_fit":
                s.setStrategy(bf)
            elif op.s1 == "priority":
                s.setStrategy(pb)
            out.append("ok")
        elif k == "book_strategy":
            r = s.bookWithStrategy(op.s1, op.s2, op.i1, op.i2, op.i3)
            out.append(r)
        elif k == "sub_obs":
            ensure_obs(op.i1)
            s.subscribeAttendee(op.s1, obs[op.i1])
            out.append("ok")
        elif k == "obs_booked":
            out.append(str(obs[op.i1].booked) if op.i1 < len(obs) else "0")
        elif k == "obs_cancelled":
            out.append(str(obs[op.i1].cancelled) if op.i1 < len(obs) else "0")
        elif k == "obs_rescheduled":
            out.append(str(obs[op.i1].rescheduled) if op.i1 < len(obs) else "0")
        elif k == "obs_new_start":
            out.append(str(obs[op.i1].lastNewStart) if op.i1 < len(obs) else "0")
        elif k == "obs_new_end":
            out.append(str(obs[op.i1].lastNewEnd) if op.i1 < len(obs) else "0")
        else:
            out.append("unknown:" + k)
    return out
