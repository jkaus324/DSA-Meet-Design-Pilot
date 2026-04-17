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

// ─── Strategy Interface ───────────────────────────────────────────────────────

type AllocationStrategy interface {
	SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string
}

// ─── Scheduler ───────────────────────────────────────────────────────────────

type MeetingScheduler struct {
	rooms    map[string]Room
	schedule map[string][]Meeting
	strategy AllocationStrategy
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{
		rooms:    make(map[string]Room),
		schedule: make(map[string][]Meeting),
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
	// TODO: sort result by ID so strategies get a deterministic order
	sort.Slice(result, func(i, j int) bool {
		// TODO: return result[i].ID < result[j].ID
		return false
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
	// TODO: check availability, append to schedule, return true
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
	// TODO: roomID := s.strategy.SelectRoom(allRooms, s, startTime, endTime, attendeeCount)
	// TODO: if roomID == "" return ""
	// TODO: bookMeeting and return roomID
	_ = allRooms
	return ""
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────

type FirstAvailable struct{}

func (f *FirstAvailable) SelectRoom(rooms []Room, scheduler *MeetingScheduler, startTime, endTime, attendeeCount int) string {
	for _, room := range rooms {
		// TODO: return room.ID if room.Capacity >= attendeeCount and scheduler.IsAvailable(...)
		_ = room
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
