package main

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

// ─── Strategy Interface ───────────────────────────────────────────────────────
// HINT: Each strategy picks a room from a list of available rooms.
// The strategy needs: available rooms, the scheduler (to call IsAvailable),
// time slot, and attendee count.

type AllocationStrategy interface {
	SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

// ─── Scheduler (extend your Part 1 solution) ──────────────────────────────────
// HINT: Add a strategy field and a BookWithStrategy() method that delegates
// room selection to the current strategy.

type MeetingScheduler struct {
	// HINT: rooms    map[string]Room
	// HINT: schedule map[string][]Meeting
	// HINT: strategy AllocationStrategy
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{}
}

func (s *MeetingScheduler) AddRoom(room Room) {}

func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool {
	return false
}

func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool {
	return false
}

func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting {
	return nil
}

func (s *MeetingScheduler) GetAllRooms() []Room {
	// HINT: collect all values from s.rooms, sort by ID, return slice
	return nil
}

func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy) {
	// HINT: s.strategy = strategy
}

func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string {
	// HINT: if strategy is nil, return ""
	// HINT: allRooms := s.GetAllRooms()
	// HINT: roomID := s.strategy.SelectRoom(allRooms, s, startTime, endTime, attendeeCount)
	// HINT: if roomID == "" return ""
	// HINT: bookMeeting and return roomID
	return ""
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────
// TODO: Implement each strategy:
//   - FirstAvailable: first room (by ID order) that fits capacity and is available
//   - BestFit: smallest room (by capacity) that fits capacity and is available
//   - PriorityBased: prefer AV rooms; among those, pick smallest that fits

type FirstAvailable struct{}

func (f *FirstAvailable) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	// HINT: iterate rooms; return first room.ID where room.Capacity >= attendeeCount
	//       and scheduler.IsAvailable(room.ID, startTime, endTime)
	return ""
}

type BestFit struct{}

func (b *BestFit) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	// HINT: track bestID and bestCapacity (start at MaxInt)
	// HINT: update when you find a room with smaller capacity that still fits
	return ""
}

type PriorityBased struct{}

func (p *PriorityBased) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	// HINT: find smallest AV room that fits; if none, fall back to smallest non-AV room
	return ""
}
