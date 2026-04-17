package main

import "fmt"

func part1Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	// Test 1: book a meeting successfully
	test("test_book_meeting", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		m := Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"}
		if !s.BookMeeting(m) {
			panic("expected BookMeeting to return true")
		}
	})

	// Test 2: conflict detection — overlapping meetings on same room
	test("test_conflict_detection", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		if s.BookMeeting(Meeting{ID: "M2", Title: "Planning", StartTime: 550, EndTime: 600, RoomID: "R1"}) {
			panic("expected conflict to be detected")
		}
	})

	// Test 3: adjacent meetings do NOT conflict
	test("test_adjacent_no_conflict", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		if !s.BookMeeting(Meeting{ID: "M2", Title: "Planning", StartTime: 570, EndTime: 630, RoomID: "R1"}) {
			panic("adjacent meetings should not conflict")
		}
	})

	// Test 4: different rooms don't conflict
	test("test_different_rooms_no_conflict", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Large Room", Capacity: 20, HasAV: true})
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		if !s.BookMeeting(Meeting{ID: "M2", Title: "Planning", StartTime: 540, EndTime: 570, RoomID: "R2"}) {
			panic("different rooms should not conflict")
		}
	})

	// Test 5: IsAvailable returns true for free slot
	test("test_is_available_free", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		if !s.IsAvailable("R1", 540, 570) {
			panic("expected room to be available")
		}
	})

	// Test 6: IsAvailable returns false for occupied slot
	test("test_is_available_occupied", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		if s.IsAvailable("R1", 550, 580) {
			panic("expected room to be unavailable")
		}
	})

	// Test 7: GetRoomSchedule returns meetings sorted by start time
	test("test_schedule_sorted", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		s.BookMeeting(Meeting{ID: "M2", Title: "Planning", StartTime: 600, EndTime: 660, RoomID: "R1"})
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		sched := s.GetRoomSchedule("R1")
		if len(sched) != 2 {
			panic("expected 2 meetings")
		}
		if sched[0].ID != "M1" {
			panic("expected M1 first")
		}
		if sched[1].ID != "M2" {
			panic("expected M2 second")
		}
	})

	// Test 8: booking to nonexistent room fails
	test("test_nonexistent_room", func() {
		s := NewMeetingScheduler()
		if s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R99"}) {
			panic("expected booking to nonexistent room to fail")
		}
	})

	// Test 9: empty schedule
	test("test_empty_schedule", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small Room", Capacity: 4, HasAV: false})
		sched := s.GetRoomSchedule("R1")
		if len(sched) != 0 {
			panic("expected empty schedule")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
