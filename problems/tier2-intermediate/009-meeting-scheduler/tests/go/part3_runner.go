package main

import "fmt"

// ─── Test Observer ────────────────────────────────────────────────────────────

type TestObserver struct {
	BookedCount              int
	CancelledCount           int
	RescheduledCount         int
	LastBookedMeetingID      string
	LastCancelledMeetingID   string
	LastRescheduledMeetingID string
	LastNewStart             int
	LastNewEnd               int
}

func (o *TestObserver) OnMeetingBooked(meeting Meeting) {
	o.BookedCount++
	o.LastBookedMeetingID = meeting.ID
}

func (o *TestObserver) OnMeetingCancelled(meeting Meeting) {
	o.CancelledCount++
	o.LastCancelledMeetingID = meeting.ID
}

func (o *TestObserver) OnMeetingRescheduled(oldMeeting, newMeeting Meeting) {
	o.RescheduledCount++
	o.LastRescheduledMeetingID = newMeeting.ID
	o.LastNewStart = newMeeting.StartTime
	o.LastNewEnd = newMeeting.EndTime
}

func part3Tests() int {
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

	// Test 1: observer notified on booking
	test("test_observer_on_book", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		obs := &TestObserver{}
		s.SubscribeAttendee("M1", obs)
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		if obs.BookedCount != 1 {
			panic("expected bookedCount == 1")
		}
		if obs.LastBookedMeetingID != "M1" {
			panic("expected lastBookedMeetingID == M1")
		}
	})

	// Test 2: observer notified on cancellation
	test("test_observer_on_cancel", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		obs := &TestObserver{}
		s.SubscribeAttendee("M1", obs)
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		s.CancelMeeting("M1")
		if obs.CancelledCount != 1 {
			panic("expected cancelledCount == 1")
		}
		if obs.LastCancelledMeetingID != "M1" {
			panic("expected lastCancelledMeetingID == M1")
		}
		if !s.IsAvailable("R1", 540, 570) {
			panic("expected room to be free after cancel")
		}
	})

	// Test 3: observer notified on reschedule
	test("test_observer_on_reschedule", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		obs := &TestObserver{}
		s.SubscribeAttendee("M1", obs)
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		ok := s.RescheduleMeeting("M1", 600, 630)
		if !ok {
			panic("expected reschedule to succeed")
		}
		if obs.RescheduledCount != 1 {
			panic("expected rescheduledCount == 1")
		}
		if obs.LastNewStart != 600 {
			panic("expected lastNewStart == 600")
		}
		if obs.LastNewEnd != 630 {
			panic("expected lastNewEnd == 630")
		}
		if !s.IsAvailable("R1", 540, 570) {
			panic("expected old slot to be free")
		}
		if s.IsAvailable("R1", 600, 630) {
			panic("expected new slot to be occupied")
		}
	})

	// Test 4: reschedule fails if new time conflicts
	test("test_reschedule_conflict", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		s.BookMeeting(Meeting{ID: "M2", Title: "Planning", StartTime: 600, EndTime: 660, RoomID: "R1"})
		ok := s.RescheduleMeeting("M1", 610, 650)
		if ok {
			panic("expected reschedule to fail due to conflict")
		}
		if s.IsAvailable("R1", 540, 570) {
			panic("expected M1 to still be in original slot")
		}
	})

	// Test 5: multiple observers on same meeting
	test("test_multiple_observers", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		obs1 := &TestObserver{}
		obs2 := &TestObserver{}
		s.SubscribeAttendee("M1", obs1)
		s.SubscribeAttendee("M1", obs2)
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		if obs1.BookedCount != 1 {
			panic("expected obs1.bookedCount == 1")
		}
		if obs2.BookedCount != 1 {
			panic("expected obs2.bookedCount == 1")
		}
	})

	// Test 6: cancel nonexistent meeting returns false
	test("test_cancel_nonexistent", func() {
		s := NewMeetingScheduler()
		if s.CancelMeeting("M99") {
			panic("expected cancel of nonexistent meeting to return false")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
