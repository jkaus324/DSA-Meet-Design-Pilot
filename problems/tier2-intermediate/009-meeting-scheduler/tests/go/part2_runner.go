package main

import "fmt"

func part2Tests() int {
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

	// Test 1: FirstAvailable picks first room by ID order
	test("test_first_available", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Medium", Capacity: 10, HasAV: false})
		s.AddRoom(Room{ID: "R3", Name: "Large", Capacity: 20, HasAV: true})
		fa := &FirstAvailable{}
		s.SetStrategy(fa)
		roomID := s.BookWithStrategy("M1", "Standup", 540, 570, 3)
		if roomID != "R1" {
			panic("expected R1, got " + roomID)
		}
	})

	// Test 2: FirstAvailable skips occupied rooms
	test("test_first_available_skip_occupied", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 4, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Medium", Capacity: 10, HasAV: false})
		fa := &FirstAvailable{}
		s.SetStrategy(fa)
		s.BookMeeting(Meeting{ID: "M1", Title: "Standup", StartTime: 540, EndTime: 570, RoomID: "R1"})
		roomID := s.BookWithStrategy("M2", "Planning", 540, 570, 3)
		if roomID != "R2" {
			panic("expected R2, got " + roomID)
		}
	})

	// Test 3: BestFit picks smallest room that fits
	test("test_best_fit", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Tiny", Capacity: 2, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Small", Capacity: 6, HasAV: false})
		s.AddRoom(Room{ID: "R3", Name: "Large", Capacity: 20, HasAV: true})
		bf := &BestFit{}
		s.SetStrategy(bf)
		roomID := s.BookWithStrategy("M1", "Meeting", 540, 570, 5)
		if roomID != "R2" {
			panic("expected R2, got " + roomID)
		}
	})

	// Test 4: BestFit returns empty when no room fits
	test("test_best_fit_no_room", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Tiny", Capacity: 2, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Small", Capacity: 4, HasAV: false})
		bf := &BestFit{}
		s.SetStrategy(bf)
		roomID := s.BookWithStrategy("M1", "Big Meeting", 540, 570, 10)
		if roomID != "" {
			panic("expected empty roomID")
		}
	})

	// Test 5: PriorityBased prefers AV rooms
	test("test_priority_prefers_av", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 6, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Medium AV", Capacity: 10, HasAV: true})
		s.AddRoom(Room{ID: "R3", Name: "Large AV", Capacity: 20, HasAV: true})
		pb := &PriorityBased{}
		s.SetStrategy(pb)
		roomID := s.BookWithStrategy("M1", "Presentation", 540, 570, 5)
		if roomID != "R2" {
			panic("expected R2 (smallest AV), got " + roomID)
		}
	})

	// Test 6: PriorityBased falls back to non-AV when no AV available
	test("test_priority_fallback_non_av", func() {
		s := NewMeetingScheduler()
		s.AddRoom(Room{ID: "R1", Name: "Small", Capacity: 6, HasAV: false})
		s.AddRoom(Room{ID: "R2", Name: "Medium AV", Capacity: 10, HasAV: true})
		pb := &PriorityBased{}
		s.SetStrategy(pb)
		s.BookMeeting(Meeting{ID: "M0", Title: "Existing", StartTime: 540, EndTime: 570, RoomID: "R2"})
		roomID := s.BookWithStrategy("M1", "Meeting", 540, 570, 5)
		if roomID != "R1" {
			panic("expected R1 (non-AV fallback), got " + roomID)
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
