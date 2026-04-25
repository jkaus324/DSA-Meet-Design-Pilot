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

type Attendee struct {
	ID    string
	Name  string
	Email string
}

// ─── NEW in Extension 2 ───────────────────────────────────────────────────────
//
// Attendees must be notified when a meeting is booked, cancelled, or
// rescheduled. The scheduler should NOT know about specific notification
// channels (email, SMS, Slack).
//
// Think about:
//   - What interface lets you decouple the scheduler from notification logic?
//   - How does the Observer pattern apply here?
//   - Should every Attendee be an observer, or should there be an adapter?
//
// Entry points (must exist for tests):
//   func (s *MeetingScheduler) AddRoom(room Room)
//   func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool
//   func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting
//   func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool
//   func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy)
//   func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string
//   func (s *MeetingScheduler) SubscribeAttendee(meetingID string, observer MeetingObserver)
//   func (s *MeetingScheduler) CancelMeeting(meetingID string) bool
//   func (s *MeetingScheduler) RescheduleMeeting(meetingID string, newStart, newEnd int) bool
//
// ─────────────────────────────────────────────────────────────────────────────

type MeetingObserver interface {
	OnMeetingBooked(meeting Meeting)
	OnMeetingCancelled(meeting Meeting)
	OnMeetingRescheduled(oldMeeting, newMeeting Meeting)
}

type AllocationStrategy interface {
	SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

type MeetingScheduler struct {
	// TODO: add your fields here
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{}
}

func (s *MeetingScheduler) AddRoom(room Room)                                {}
func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool { return false }
func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool                  { return false }
func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting           { return nil }
func (s *MeetingScheduler) GetAllRooms() []Room                               { return nil }
func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy)           {}
func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string {
	return ""
}
func (s *MeetingScheduler) SubscribeAttendee(meetingID string, observer MeetingObserver) {}
func (s *MeetingScheduler) CancelMeeting(meetingID string) bool                          { return false }
func (s *MeetingScheduler) RescheduleMeeting(meetingID string, newStart, newEnd int) bool {
	return false
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
