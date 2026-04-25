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

// ─── Observer Interface ───────────────────────────────────────────────────────
// HINT: The scheduler notifies all subscribed observers on book/cancel/reschedule.
// The scheduler does NOT know about email, SMS, or Slack — only this interface.

type MeetingObserver interface {
	OnMeetingBooked(meeting Meeting)
	OnMeetingCancelled(meeting Meeting)
	OnMeetingRescheduled(oldMeeting, newMeeting Meeting)
}

// ─── Strategy Interface ───────────────────────────────────────────────────────

type AllocationStrategy interface {
	SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

// ─── Scheduler (extend your Part 2 solution) ──────────────────────────────────
// HINT: Add:
//   - observers    map[string][]MeetingObserver (meetingID -> observers)
//   - meetingsById map[string]Meeting (for cancel/reschedule lookup)
//   - SubscribeAttendee(), CancelMeeting(), RescheduleMeeting()
//   - Notify observers inside BookMeeting(), CancelMeeting(), RescheduleMeeting()

type MeetingScheduler struct {
	// HINT: rooms         map[string]Room
	// HINT: schedule      map[string][]Meeting
	// HINT: meetingsById  map[string]Meeting
	// HINT: observers     map[string][]MeetingObserver
	// HINT: strategy      AllocationStrategy
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{}
}

func (s *MeetingScheduler) SetStrategy(strategy AllocationStrategy) {}

func (s *MeetingScheduler) AddRoom(room Room) {}

func (s *MeetingScheduler) GetAllRooms() []Room { return nil }

func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool {
	// HINT: check s.schedule[roomID] for overlaps
	return false
}

func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool {
	// HINT: validate room, check availability, append to schedule, store in meetingsById
	// HINT: notify observers[meeting.ID] with OnMeetingBooked
	return false
}

func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting { return nil }

func (s *MeetingScheduler) BookWithStrategy(meetingID, title string, startTime, endTime, attendeeCount int) string {
	return ""
}

func (s *MeetingScheduler) SubscribeAttendee(meetingID string, observer MeetingObserver) {
	// HINT: append observer to s.observers[meetingID]
}

func (s *MeetingScheduler) CancelMeeting(meetingID string) bool {
	// HINT: look up meeting in meetingsById; return false if not found
	// HINT: remove from schedule[meeting.RoomID]
	// HINT: delete from meetingsById
	// HINT: notify observers with OnMeetingCancelled
	return false
}

func (s *MeetingScheduler) RescheduleMeeting(meetingID string, newStart, newEnd int) bool {
	// HINT: look up oldMeeting in meetingsById; return false if not found
	// HINT: temporarily remove old meeting from schedule
	// HINT: check if new time is available on same room
	// HINT: if available: add new meeting, notify observers with OnMeetingRescheduled
	// HINT: if NOT available: restore old meeting, return false
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
