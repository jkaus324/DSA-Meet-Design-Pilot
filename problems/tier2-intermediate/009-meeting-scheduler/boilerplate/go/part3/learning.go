package main

import (
	"math"
	"sort"
)

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type Room struct {
	ID       string
	Name     string
	Capacity int
	HasAV    bool
}

type Meeting struct {
	ID        string
	Title     string
	StartTime int
	EndTime   int
	RoomID    string
}

type Attendee struct {
	ID    string
	Name  string
	Email string
}

// ─── Observer Interface ───────────────────────────────────────────────────────

type MeetingObserver interface {
	OnMeetingBooked(meeting Meeting)
	OnMeetingCancelled(meeting Meeting)
	OnMeetingRescheduled(oldMeeting, newMeeting Meeting)
}

// ─── Strategy Interface ───────────────────────────────────────────────────────

type AllocationStrategy interface {
	SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

// ─── Scheduler ───────────────────────────────────────────────────────────────

type MeetingScheduler struct {
	rooms        map[string]Room
	schedule     map[string][]Meeting
	meetingsById map[string]Meeting
	observers    map[string][]MeetingObserver
	strategy     AllocationStrategy
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{
		rooms:        make(map[string]Room),
		schedule:     make(map[string][]Meeting),
		meetingsById: make(map[string]Meeting),
		observers:    make(map[string][]MeetingObserver),
	}
}

func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy) {
	// TODO: assign strategy to s.strategy
}

func (s *MeetingScheduler) AddRoom(room Room) {
	// TODO: store room in s.rooms keyed by room.ID
}

func (s *MeetingScheduler) GetAllRooms() []Room {
	result := make([]Room, 0, len(s.rooms))
	for _, r := range s.rooms {
		result = append(result, r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool {
	for _, m := range s.schedule[roomID] {
		// TODO: return false if overlap: startTime < m.EndTime && m.StartTime < endTime
		_ = m
	}
	return true
}

func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool {
	if _, ok := s.rooms[meeting.RoomID]; !ok {
		return false
	}
	// TODO: check availability
	// TODO: append to s.schedule[meeting.RoomID]
	// TODO: store in s.meetingsById[meeting.ID]
	// TODO: notify s.observers[meeting.ID] with OnMeetingBooked(meeting)
	// HINT: for _, obs := range s.observers[meeting.ID] { obs.OnMeetingBooked(meeting) }
	return false
}

func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting {
	meetings := s.schedule[roomID]
	result := make([]Meeting, len(meetings))
	copy(result, meetings)
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime < result[j].StartTime
	})
	return result
}

func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string {
	if s.strategy == nil {
		return ""
	}
	allRooms := s.GetAllRooms()
	roomID := s.strategy.SelectRoom(allRooms, s, startTime, endTime, attendeeCount)
	if roomID == "" {
		return ""
	}
	m := Meeting{ID: meetingID, Title: title, StartTime: startTime, EndTime: endTime, RoomID: roomID}
	if s.BookMeeting(m) {
		return roomID
	}
	return ""
}

func (s *MeetingScheduler) SubscribeAttendee(meetingID string, observer MeetingObserver) {
	// TODO: append observer to s.observers[meetingID]
}

func (s *MeetingScheduler) CancelMeeting(meetingID string) bool {
	meeting, ok := s.meetingsById[meetingID]
	if !ok {
		return false
	}
	// TODO: remove meeting from s.schedule[meeting.RoomID]
	//   HINT: iterate slice, rebuild without matching ID
	// TODO: delete s.meetingsById[meetingID]
	// TODO: notify s.observers[meetingID] with OnMeetingCancelled(meeting)
	_ = meeting
	return false
}

func (s *MeetingScheduler) RescheduleMeeting(meetingID string, newStart, newEnd int) bool {
	oldMeeting, ok := s.meetingsById[meetingID]
	if !ok {
		return false
	}
	// TODO: temporarily remove old meeting from s.schedule[oldMeeting.RoomID]
	// TODO: check if s.IsAvailable(oldMeeting.RoomID, newStart, newEnd)
	// TODO: if available:
	//   newMeeting := oldMeeting with StartTime=newStart, EndTime=newEnd
	//   append newMeeting to schedule, update meetingsById
	//   notify observers with OnMeetingRescheduled(oldMeeting, newMeeting)
	//   return true
	// TODO: if NOT available: restore old meeting in schedule, return false
	_ = oldMeeting
	return false
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────

type FirstAvailable struct{}

func (f *FirstAvailable) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	for _, room := range rooms {
		if room.Capacity >= attendeeCount && scheduler.IsAvailable(room.ID, startTime, endTime) {
			return room.ID
		}
	}
	return ""
}

type BestFit struct{}

func (b *BestFit) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	bestID := ""
	bestCap := math.MaxInt32
	for _, room := range rooms {
		// TODO: if fits and available and room.Capacity < bestCap, update bestID and bestCap
		_ = room
		_ = bestCap
	}
	return bestID
}

type PriorityBased struct{}

func (p *PriorityBased) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	bestAV, bestNonAV := "", ""
	capAV, capNonAV := math.MaxInt32, math.MaxInt32
	for _, room := range rooms {
		// TODO: if room fits and is available:
		//   if room.HasAV and room.Capacity < capAV, update bestAV and capAV
		//   else if !room.HasAV and room.Capacity < capNonAV, update bestNonAV and capNonAV
		_ = room
		_ = capAV
		_ = capNonAV
	}
	if bestAV != "" {
		return bestAV
	}
	return bestNonAV
}
