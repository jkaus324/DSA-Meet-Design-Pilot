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
	StartTime int // minutes since midnight
	EndTime   int
	RoomID    string
}

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement a MeetingScheduler that:
//   1. Manages rooms and their meeting schedules
//   2. Detects conflicts when booking (two meetings cannot overlap on the
//      same room)
//   3. Returns a room's schedule sorted by start time
//
// Think about:
//   - How do you check if two time intervals overlap?
//   - What data structure maps roomID -> meetings efficiently?
//   - How will this extend to support automatic room allocation later?
//
// Entry points (must exist for tests):
//   func (s *MeetingScheduler) AddRoom(room Room)
//   func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool
//   func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting
//   func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool
//
// ─────────────────────────────────────────────────────────────────────────────

type MeetingScheduler struct {
	// TODO: add your fields here
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{}
}

func (s *MeetingScheduler) AddRoom(room Room) {
}

func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool {
	return false
}

func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool {
	return false
}

func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting {
	return nil
}
