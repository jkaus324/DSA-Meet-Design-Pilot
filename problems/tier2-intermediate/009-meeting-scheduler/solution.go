// Meeting scheduler — rooms, observers, allocation strategies (Go).
package main

import (
	"sort"
	"strconv"
)

type Op struct {
	kind string
	s1   string
	s2   string
	s3   string
	i1   int
	i2   int
	i3   int
}

type Room struct {
	id       string
	name     string
	capacity int
	hasAV    bool
}

type Meeting struct {
	id        string
	title     string
	startTime int
	endTime   int
	roomId    string
}

type CountingObserver struct {
	booked       int
	cancelled    int
	rescheduled  int
	lastNewStart int
	lastNewEnd   int
}

func (o *CountingObserver) onMeetingBooked(m *Meeting)    { o.booked++ }
func (o *CountingObserver) onMeetingCancelled(m *Meeting) { o.cancelled++ }
func (o *CountingObserver) onMeetingRescheduled(oldM, newM *Meeting) {
	o.rescheduled++
	o.lastNewStart = newM.startTime
	o.lastNewEnd = newM.endTime
}

type Strategy interface {
	selectRoom(rooms []*Room, s *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

type FirstAvailable struct{}

func (FirstAvailable) selectRoom(rooms []*Room, s *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	for _, r := range rooms {
		if r.capacity >= attendeeCount && s.isAvailable(r.id, startTime, endTime) {
			return r.id
		}
	}
	return ""
}

type BestFit struct{}

func (BestFit) selectRoom(rooms []*Room, s *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	bestId := ""
	bestCap := int(^uint(0) >> 1)
	for _, r := range rooms {
		if r.capacity >= attendeeCount && s.isAvailable(r.id, startTime, endTime) {
			if r.capacity < bestCap {
				bestCap = r.capacity
				bestId = r.id
			}
		}
	}
	return bestId
}

type PriorityBased struct{}

func (PriorityBased) selectRoom(rooms []*Room, s *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	bestAV, bestNonAV := "", ""
	maxInt := int(^uint(0) >> 1)
	bestAVCap, bestNonAVCap := maxInt, maxInt
	for _, r := range rooms {
		if r.capacity >= attendeeCount && s.isAvailable(r.id, startTime, endTime) {
			if r.hasAV {
				if r.capacity < bestAVCap {
					bestAVCap = r.capacity
					bestAV = r.id
				}
			} else {
				if r.capacity < bestNonAVCap {
					bestNonAVCap = r.capacity
					bestNonAV = r.id
				}
			}
		}
	}
	if bestAV == "" {
		return bestNonAV
	}
	return bestAV
}

type MeetingScheduler struct {
	rooms        map[string]*Room
	schedule     map[string][]*Meeting
	meetingsById map[string]*Meeting
	observers    map[string][]*CountingObserver
	strategy     Strategy
}

func newMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{
		rooms:        map[string]*Room{},
		schedule:     map[string][]*Meeting{},
		meetingsById: map[string]*Meeting{},
		observers:    map[string][]*CountingObserver{},
	}
}

func (s *MeetingScheduler) setStrategy(st Strategy) { s.strategy = st }

func (s *MeetingScheduler) addRoom(r *Room) { s.rooms[r.id] = r }

func (s *MeetingScheduler) getAllRooms() []*Room {
	out := make([]*Room, 0, len(s.rooms))
	for _, r := range s.rooms {
		out = append(out, r)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].id < out[j].id })
	return out
}

func (s *MeetingScheduler) isAvailable(roomId string, startTime, endTime int) bool {
	for _, m := range s.schedule[roomId] {
		if startTime < m.endTime && m.startTime < endTime {
			return false
		}
	}
	return true
}

func (s *MeetingScheduler) bookMeeting(m *Meeting) bool {
	if _, ok := s.rooms[m.roomId]; !ok {
		return false
	}
	if !s.isAvailable(m.roomId, m.startTime, m.endTime) {
		return false
	}
	s.schedule[m.roomId] = append(s.schedule[m.roomId], m)
	s.meetingsById[m.id] = m
	for _, obs := range s.observers[m.id] {
		obs.onMeetingBooked(m)
	}
	return true
}

func (s *MeetingScheduler) getRoomSchedule(roomId string) []*Meeting {
	meets := make([]*Meeting, len(s.schedule[roomId]))
	copy(meets, s.schedule[roomId])
	sort.SliceStable(meets, func(i, j int) bool { return meets[i].startTime < meets[j].startTime })
	return meets
}

func (s *MeetingScheduler) bookWithStrategy(meetingId, title string, startTime, endTime, attendeeCount int) string {
	if s.strategy == nil {
		return ""
	}
	roomId := s.strategy.selectRoom(s.getAllRooms(), s, startTime, endTime, attendeeCount)
	if roomId == "" {
		return ""
	}
	m := &Meeting{meetingId, title, startTime, endTime, roomId}
	if s.bookMeeting(m) {
		return roomId
	}
	return ""
}

func (s *MeetingScheduler) subscribeAttendee(meetingId string, obs *CountingObserver) {
	s.observers[meetingId] = append(s.observers[meetingId], obs)
}

func (s *MeetingScheduler) cancelMeeting(meetingId string) bool {
	meeting, ok := s.meetingsById[meetingId]
	if !ok {
		return false
	}
	rm := s.schedule[meeting.roomId]
	kept := rm[:0:0]
	for _, m := range rm {
		if m.id != meetingId {
			kept = append(kept, m)
		}
	}
	s.schedule[meeting.roomId] = kept
	delete(s.meetingsById, meetingId)
	for _, obs := range s.observers[meetingId] {
		obs.onMeetingCancelled(meeting)
	}
	return true
}

func (s *MeetingScheduler) rescheduleMeeting(meetingId string, newStart, newEnd int) bool {
	old, ok := s.meetingsById[meetingId]
	if !ok {
		return false
	}
	rm := s.schedule[old.roomId]
	kept := rm[:0:0]
	for _, m := range rm {
		if m.id != meetingId {
			kept = append(kept, m)
		}
	}
	s.schedule[old.roomId] = kept
	if !s.isAvailable(old.roomId, newStart, newEnd) {
		s.schedule[old.roomId] = append(s.schedule[old.roomId], old)
		return false
	}
	newM := &Meeting{old.id, old.title, newStart, newEnd, old.roomId}
	s.schedule[old.roomId] = append(s.schedule[old.roomId], newM)
	s.meetingsById[meetingId] = newM
	for _, obs := range s.observers[meetingId] {
		obs.onMeetingRescheduled(old, newM)
	}
	return true
}

func meeting_simulate(ops []Op) []string {
	out := []string{}
	s := newMeetingScheduler()
	fa, bf, pb := FirstAvailable{}, BestFit{}, PriorityBased{}
	obs := []*CountingObserver{}

	ensureObs := func(idx int) {
		for len(obs) <= idx {
			obs = append(obs, &CountingObserver{})
		}
	}

	for _, op := range ops {
		switch op.kind {
		case "reset":
			s = newMeetingScheduler()
			obs = []*CountingObserver{}
			out = append(out, "ok")
		case "add_room":
			s.addRoom(&Room{op.s1, op.s2, op.i1, op.i2 != 0})
			out = append(out, "ok")
		case "book":
			m := &Meeting{op.s1, op.s1, op.i1, op.i2, op.s2}
			if s.bookMeeting(m) {
				out = append(out, "ok")
			} else {
				out = append(out, "fail")
			}
		case "is_available":
			if s.isAvailable(op.s1, op.i1, op.i2) {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		case "sched_size":
			out = append(out, strconv.Itoa(len(s.getRoomSchedule(op.s1))))
		case "sched_at":
			sched := s.getRoomSchedule(op.s1)
			if op.i1 >= 0 && op.i1 < len(sched) {
				out = append(out, sched[op.i1].id)
			} else {
				out = append(out, "")
			}
		case "cancel":
			if s.cancelMeeting(op.s1) {
				out = append(out, "ok")
			} else {
				out = append(out, "fail")
			}
		case "reschedule":
			if s.rescheduleMeeting(op.s1, op.i1, op.i2) {
				out = append(out, "ok")
			} else {
				out = append(out, "fail")
			}
		case "set_strategy":
			switch op.s1 {
			case "first_available":
				s.setStrategy(fa)
			case "best_fit":
				s.setStrategy(bf)
			case "priority":
				s.setStrategy(pb)
			}
			out = append(out, "ok")
		case "book_strategy":
			out = append(out, s.bookWithStrategy(op.s1, op.s2, op.i1, op.i2, op.i3))
		case "sub_obs":
			ensureObs(op.i1)
			s.subscribeAttendee(op.s1, obs[op.i1])
			out = append(out, "ok")
		case "obs_booked":
			if op.i1 < len(obs) {
				out = append(out, strconv.Itoa(obs[op.i1].booked))
			} else {
				out = append(out, "0")
			}
		case "obs_cancelled":
			if op.i1 < len(obs) {
				out = append(out, strconv.Itoa(obs[op.i1].cancelled))
			} else {
				out = append(out, "0")
			}
		case "obs_rescheduled":
			if op.i1 < len(obs) {
				out = append(out, strconv.Itoa(obs[op.i1].rescheduled))
			} else {
				out = append(out, "0")
			}
		case "obs_new_start":
			if op.i1 < len(obs) {
				out = append(out, strconv.Itoa(obs[op.i1].lastNewStart))
			} else {
				out = append(out, "0")
			}
		case "obs_new_end":
			if op.i1 < len(obs) {
				out = append(out, strconv.Itoa(obs[op.i1].lastNewEnd))
			} else {
				out = append(out, "0")
			}
		default:
			out = append(out, "unknown:"+op.kind)
		}
	}
	return out
}
