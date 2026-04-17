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

// ─── NEW in Extension 1 ───────────────────────────────────────────────────────
//
// The office manager wants AUTOMATIC room allocation using different
// strategies: first-available, best-fit, and priority-based (prefer AV rooms).
//
// Think about:
//   - How do you swap allocation algorithms without modifying the scheduler?
//   - What interface lets all strategies work interchangeably?
//   - How does the Strategy pattern apply here?
//
// Entry points (must exist for tests):
//   func (s *MeetingScheduler) AddRoom(room Room)
//   func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool
//   func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting
//   func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool
//   func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy)
//   func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string
//
// ─────────────────────────────────────────────────────────────────────────────

type AllocationStrategy interface {
	SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

type MeetingScheduler struct {
	// TODO: add your fields here
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
	return nil
}

func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy) {}

func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string {
	return ""
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────

type FirstAvailable struct{}

func (f *FirstAvailable) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	return ""
}

type BestFit struct{}

func (b *BestFit) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	return ""
}

type PriorityBased struct{}

func (p *PriorityBased) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	return ""
}
