package main

import "sort"

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

type MeetingScheduler struct {
	rooms    map[string]Room
	schedule map[string][]Meeting // roomID -> meetings
}

func NewMeetingScheduler() *MeetingScheduler {
	return &MeetingScheduler{
		rooms:    make(map[string]Room),
		schedule: make(map[string][]Meeting),
	}
}

func (s *MeetingScheduler) AddRoom(room Room) {
	// TODO: store room in s.rooms keyed by room.ID
}

func (s *MeetingScheduler) IsAvailable(roomID string, startTime, endTime int) bool {
	meetings := s.schedule[roomID]
	for _, m := range meetings {
		// TODO: return false if [startTime, endTime) overlaps [m.StartTime, m.EndTime)
		// HINT: overlap condition is startTime < m.EndTime && m.StartTime < endTime
		_ = m
	}
	return true
}

func (s *MeetingScheduler) BookMeeting(meeting Meeting) bool {
	if _, ok := s.rooms[meeting.RoomID]; !ok {
		return false
	}
	// TODO: check availability, then append meeting to s.schedule[meeting.RoomID]
	// HINT: if !s.IsAvailable(...) return false; then append
	return false
}

func (s *MeetingScheduler) GetRoomSchedule(roomID string) []Meeting {
	meetings := s.schedule[roomID]
	result := make([]Meeting, len(meetings))
	copy(result, meetings)
	// TODO: sort result by StartTime using sort.Slice
	sort.Slice(result, func(i, j int) bool {
		// TODO: return result[i].StartTime < result[j].StartTime
		return false
	})
	return result
}
