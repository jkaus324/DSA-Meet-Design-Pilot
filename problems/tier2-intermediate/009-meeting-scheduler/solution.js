// Meeting scheduler — rooms, observers, allocation strategies (JavaScript).

class Op {
  constructor(kind, s1 = '', s2 = '', s3 = '', i1 = 0, i2 = 0, i3 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
  }
}

class Room {
  constructor(id, name, capacity, hasAV) {
    this.id = id;
    this.name = name;
    this.capacity = capacity;
    this.hasAV = hasAV;
  }
}

class Meeting {
  constructor(id, title, startTime, endTime, roomId) {
    this.id = id;
    this.title = title;
    this.startTime = startTime;
    this.endTime = endTime;
    this.roomId = roomId;
  }
}

class CountingObserver {
  constructor() {
    this.booked = 0;
    this.cancelled = 0;
    this.rescheduled = 0;
    this.lastNewStart = 0;
    this.lastNewEnd = 0;
  }

  onMeetingBooked(meeting) { this.booked += 1; }

  onMeetingCancelled(meeting) { this.cancelled += 1; }

  onMeetingRescheduled(oldM, newM) {
    this.rescheduled += 1;
    this.lastNewStart = newM.startTime;
    this.lastNewEnd = newM.endTime;
  }
}

class FirstAvailable {
  selectRoom(rooms, scheduler, startTime, endTime, attendeeCount) {
    for (const r of rooms) {
      if (r.capacity >= attendeeCount && scheduler.isAvailable(r.id, startTime, endTime)) {
        return r.id;
      }
    }
    return '';
  }
}

class BestFit {
  selectRoom(rooms, scheduler, startTime, endTime, attendeeCount) {
    let bestId = '';
    let bestCap = Infinity;
    for (const r of rooms) {
      if (r.capacity >= attendeeCount && scheduler.isAvailable(r.id, startTime, endTime)) {
        if (r.capacity < bestCap) {
          bestCap = r.capacity;
          bestId = r.id;
        }
      }
    }
    return bestId;
  }
}

class PriorityBased {
  selectRoom(rooms, scheduler, startTime, endTime, attendeeCount) {
    let bestAV = '';
    let bestNonAV = '';
    let bestAVCap = Infinity;
    let bestNonAVCap = Infinity;
    for (const r of rooms) {
      if (r.capacity >= attendeeCount && scheduler.isAvailable(r.id, startTime, endTime)) {
        if (r.hasAV) {
          if (r.capacity < bestAVCap) {
            bestAVCap = r.capacity;
            bestAV = r.id;
          }
        } else {
          if (r.capacity < bestNonAVCap) {
            bestNonAVCap = r.capacity;
            bestNonAV = r.id;
          }
        }
      }
    }
    return bestAV === '' ? bestNonAV : bestAV;
  }
}

class MeetingScheduler {
  constructor() {
    this.rooms = new Map(); // id -> Room
    this.schedule = new Map(); // roomId -> [Meeting]
    this.meetingsById = new Map();
    this.observers = new Map(); // meetingId -> [observer]
    this.strategy = null;
  }

  setStrategy(s) { this.strategy = s; }

  addRoom(room) { this.rooms.set(room.id, room); }

  getAllRooms() {
    return [...this.rooms.values()].sort((a, b) => (a.id < b.id ? -1 : a.id > b.id ? 1 : 0));
  }

  isAvailable(roomId, startTime, endTime) {
    const meets = this.schedule.has(roomId) ? this.schedule.get(roomId) : [];
    for (const m of meets) {
      if (startTime < m.endTime && m.startTime < endTime) return false;
    }
    return true;
  }

  bookMeeting(meeting) {
    if (!this.rooms.has(meeting.roomId)) return false;
    if (!this.isAvailable(meeting.roomId, meeting.startTime, meeting.endTime)) return false;
    if (!this.schedule.has(meeting.roomId)) this.schedule.set(meeting.roomId, []);
    this.schedule.get(meeting.roomId).push(meeting);
    this.meetingsById.set(meeting.id, meeting);
    for (const obs of (this.observers.has(meeting.id) ? this.observers.get(meeting.id) : [])) {
      obs.onMeetingBooked(meeting);
    }
    return true;
  }

  getRoomSchedule(roomId) {
    const meets = [...(this.schedule.has(roomId) ? this.schedule.get(roomId) : [])];
    return meets.sort((a, b) => a.startTime - b.startTime);
  }

  bookWithStrategy(meetingId, title, startTime, endTime, attendeeCount) {
    if (!this.strategy) return '';
    const allRooms = this.getAllRooms();
    const roomId = this.strategy.selectRoom(allRooms, this, startTime, endTime, attendeeCount);
    if (roomId === '') return '';
    const m = new Meeting(meetingId, title, startTime, endTime, roomId);
    if (this.bookMeeting(m)) return roomId;
    return '';
  }

  subscribeAttendee(meetingId, obs) {
    if (!this.observers.has(meetingId)) this.observers.set(meetingId, []);
    this.observers.get(meetingId).push(obs);
  }

  cancelMeeting(meetingId) {
    if (!this.meetingsById.has(meetingId)) return false;
    const meeting = this.meetingsById.get(meetingId);
    const rm = this.schedule.has(meeting.roomId) ? this.schedule.get(meeting.roomId) : [];
    this.schedule.set(meeting.roomId, rm.filter(m => m.id !== meetingId));
    this.meetingsById.delete(meetingId);
    for (const obs of (this.observers.has(meetingId) ? this.observers.get(meetingId) : [])) {
      obs.onMeetingCancelled(meeting);
    }
    return true;
  }

  rescheduleMeeting(meetingId, newStart, newEnd) {
    if (!this.meetingsById.has(meetingId)) return false;
    const old = this.meetingsById.get(meetingId);
    const rm = this.schedule.has(old.roomId) ? this.schedule.get(old.roomId) : [];
    this.schedule.set(old.roomId, rm.filter(m => m.id !== meetingId));
    if (!this.isAvailable(old.roomId, newStart, newEnd)) {
      this.schedule.get(old.roomId).push(old);
      return false;
    }
    const nw = new Meeting(old.id, old.title, newStart, newEnd, old.roomId);
    this.schedule.get(old.roomId).push(nw);
    this.meetingsById.set(meetingId, nw);
    for (const obs of (this.observers.has(meetingId) ? this.observers.get(meetingId) : [])) {
      obs.onMeetingRescheduled(old, nw);
    }
    return true;
  }
}

function meeting_simulate(ops) {
  const out = [];
  let s = new MeetingScheduler();
  const fa = new FirstAvailable();
  const bf = new BestFit();
  const pb = new PriorityBased();
  let obs = [];

  function ensureObs(idx) {
    while (obs.length <= idx) obs.push(new CountingObserver());
  }

  for (const op of ops) {
    const k = op.kind;
    if (k === 'reset') {
      s = new MeetingScheduler();
      obs = [];
      out.push('ok');
    } else if (k === 'add_room') {
      s.addRoom(new Room(op.s1, op.s2, op.i1, op.i2 !== 0));
      out.push('ok');
    } else if (k === 'book') {
      const m = new Meeting(op.s1, op.s1, op.i1, op.i2, op.s2);
      out.push(s.bookMeeting(m) ? 'ok' : 'fail');
    } else if (k === 'is_available') {
      out.push(s.isAvailable(op.s1, op.i1, op.i2) ? 'yes' : 'no');
    } else if (k === 'sched_size') {
      out.push(String(s.getRoomSchedule(op.s1).length));
    } else if (k === 'sched_at') {
      const sched = s.getRoomSchedule(op.s1);
      out.push(0 <= op.i1 && op.i1 < sched.length ? sched[op.i1].id : '');
    } else if (k === 'cancel') {
      out.push(s.cancelMeeting(op.s1) ? 'ok' : 'fail');
    } else if (k === 'reschedule') {
      out.push(s.rescheduleMeeting(op.s1, op.i1, op.i2) ? 'ok' : 'fail');
    } else if (k === 'set_strategy') {
      if (op.s1 === 'first_available') s.setStrategy(fa);
      else if (op.s1 === 'best_fit') s.setStrategy(bf);
      else if (op.s1 === 'priority') s.setStrategy(pb);
      out.push('ok');
    } else if (k === 'book_strategy') {
      const r = s.bookWithStrategy(op.s1, op.s2, op.i1, op.i2, op.i3);
      out.push(r);
    } else if (k === 'sub_obs') {
      ensureObs(op.i1);
      s.subscribeAttendee(op.s1, obs[op.i1]);
      out.push('ok');
    } else if (k === 'obs_booked') {
      out.push(op.i1 < obs.length ? String(obs[op.i1].booked) : '0');
    } else if (k === 'obs_cancelled') {
      out.push(op.i1 < obs.length ? String(obs[op.i1].cancelled) : '0');
    } else if (k === 'obs_rescheduled') {
      out.push(op.i1 < obs.length ? String(obs[op.i1].rescheduled) : '0');
    } else if (k === 'obs_new_start') {
      out.push(op.i1 < obs.length ? String(obs[op.i1].lastNewStart) : '0');
    } else if (k === 'obs_new_end') {
      out.push(op.i1 < obs.length ? String(obs[op.i1].lastNewEnd) : '0');
    } else {
      out.push('unknown:' + k);
    }
  }
  return out;
}

module.exports = {
  Op, Room, Meeting, CountingObserver,
  FirstAvailable, BestFit, PriorityBased, MeetingScheduler,
  meeting_simulate,
};
