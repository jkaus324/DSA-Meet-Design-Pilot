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

// ─── Scheduler ───────────────────────────────────────────────────────────────
// HINT: Use a map to associate roomID -> []Meeting.
// Two intervals [s1,e1) and [s2,e2) overlap if s1 < e2 && s2 < e1.

type MeetingScheduler struct {
	// HINT: rooms  map[string]Room
	// HINT: schedule map[string][]Meeting
}

func NewMeetingScheduler() *MeetingScheduler {
	// HINT: initialise both maps
	return &MeetingScheduler{}
}

func (s *MeetingScheduler) AddRoom(room Room) {
	// HINT: store room in rooms map keyed by room.ID
}

func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool {
	// HINT: iterate schedule[roomID]; return false on overlap
	// HINT: overlap condition: startTime < m.EndTime && m.StartTime < endTime
	return false
}

func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool {
	// HINT: return false if room does not exist
	// HINT: return false if !IsAvailable(...)
	// HINT: append meeting to schedule[meeting.RoomID]
	return false
}

func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting {
	// HINT: copy the slice, sort by StartTime, return it
	return nil
}
